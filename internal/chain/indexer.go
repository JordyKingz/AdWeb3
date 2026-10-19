package chain

import (
	"adweb3.local/platform/internal/campaigns"
	"adweb3.local/platform/internal/ledger"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/big"
	"strings"
)

const eventABI = `[{"type":"event","name":"Funded","anonymous":false,"inputs":[{"name":"orderId","type":"bytes32","indexed":true},{"name":"payer","type":"address","indexed":true},{"name":"amount","type":"uint256","indexed":false}]},{"type":"event","name":"Refunded","anonymous":false,"inputs":[{"name":"refundId","type":"bytes32","indexed":true},{"name":"orderId","type":"bytes32","indexed":true},{"name":"payer","type":"address","indexed":true},{"name":"amount","type":"uint256","indexed":false}]}]`

type RPC interface {
	ChainID(context.Context) (*big.Int, error)
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	CodeAt(context.Context, common.Address, *big.Int) ([]byte, error)
	CallContract(context.Context, ethereum.CallMsg, *big.Int) ([]byte, error)
}
type Indexer struct {
	Client        RPC
	DB            *pgxpool.Pool
	Campaigns     campaigns.Service
	Escrow        common.Address
	Confirmations uint64
}

func (i Indexer) quarantine(ctx context.Context, reason string) error {
	_, e := i.DB.Exec(ctx, "INSERT INTO chain_health(chain_id,quarantined,reason) VALUES($1,true,$2) ON CONFLICT(chain_id) DO UPDATE SET quarantined=true,reason=EXCLUDED.reason,updated_at=now()", i.Campaigns.ChainID, reason)
	if e != nil {
		return e
	}
	return errors.New(reason)
}
func (i Indexer) Scan(ctx context.Context) error {
	id, e := i.Client.ChainID(ctx)
	if e != nil {
		return e
	}
	if id.Int64() != i.Campaigns.ChainID || id.Int64() != 31337 && id.Int64() != 43113 {
		return errors.New("chain mismatch")
	}
	code, e := i.Client.CodeAt(ctx, i.Escrow, nil)
	if e != nil {
		return e
	}
	if len(code) == 0 {
		return errors.New("missing escrow code")
	}
	output, e := i.Client.CallContract(ctx, ethereum.CallMsg{To: &i.Escrow, Data: common.FromHex("0xfc0c546a")}, nil)
	if e != nil {
		return e
	}
	if len(output) != 32 || common.BytesToAddress(output) != common.HexToAddress(i.Campaigns.Token) {
		return errors.New("escrow token mismatch")
	}
	head, e := i.Client.HeaderByNumber(ctx, nil)
	if e != nil {
		return e
	}
	if head.Number.Uint64() < i.Confirmations {
		return nil
	}
	end := head.Number.Uint64() - i.Confirmations
	tx, e := i.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(hashtextextended($1,0))", fmt.Sprintf("index:%d:%s", id.Int64(), i.Escrow.Hex())); e != nil {
		return e
	}
	if e = campaigns.CheckChain(ctx, tx, id.Int64()); e != nil {
		return e
	}
	var checkpoint int64
	var hash string
	start := uint64(0)
	e = tx.QueryRow(ctx, "SELECT block_number,block_hash FROM chain_checkpoints WHERE chain_id=$1 AND contract=$2", id.Int64(), i.Escrow.Hex()).Scan(&checkpoint, &hash)
	if e == nil {
		h, e := i.Client.HeaderByNumber(ctx, big.NewInt(checkpoint))
		if e != nil {
			return e
		}
		if h.Hash().Hex() != hash {
			tx.Rollback(ctx)
			return i.quarantine(ctx, "finalized checkpoint reorg")
		}
		start = uint64(checkpoint) + 1
	} else if !errors.Is(e, pgx.ErrNoRows) {
		return e
	}
	if start > end {
		return nil
	}
	if end-start > 2000 {
		end = start + 2000
	}
	anchor, e := i.Client.HeaderByNumber(ctx, new(big.Int).SetUint64(end))
	if e != nil {
		return e
	}
	parsed, e := abi.JSON(strings.NewReader(eventABI))
	if e != nil {
		return e
	}
	logs, e := i.Client.FilterLogs(ctx, ethereum.FilterQuery{FromBlock: new(big.Int).SetUint64(start), ToBlock: new(big.Int).SetUint64(end), Addresses: []common.Address{i.Escrow}, Topics: [][]common.Hash{{parsed.Events["Funded"].ID, parsed.Events["Refunded"].ID}}})
	if e != nil {
		return e
	}
	for _, log := range logs {
		canonical, e := i.Client.HeaderByNumber(ctx, new(big.Int).SetUint64(log.BlockNumber))
		if e != nil {
			return e
		}
		if log.Removed || canonical.Hash() != log.BlockHash {
			tx.Rollback(ctx)
			return i.quarantine(ctx, "noncanonical log")
		}
		if e = i.apply(ctx, tx, parsed, log); e != nil {
			return e
		}
	}
	final, e := i.Client.HeaderByNumber(ctx, new(big.Int).SetUint64(end))
	if e != nil {
		return e
	}
	if final.Hash() != anchor.Hash() {
		tx.Rollback(ctx)
		return i.quarantine(ctx, "chain changed during scan")
	}
	_, e = tx.Exec(ctx, "INSERT INTO chain_checkpoints VALUES($1,$2,$3,$4) ON CONFLICT(chain_id,contract) DO UPDATE SET block_number=EXCLUDED.block_number,block_hash=EXCLUDED.block_hash WHERE chain_checkpoints.block_number<EXCLUDED.block_number", id.Int64(), i.Escrow.Hex(), end, anchor.Hash().Hex())
	if e != nil {
		return e
	}
	return tx.Commit(ctx)
}
func (i Indexer) apply(ctx context.Context, tx pgx.Tx, parsed abi.ABI, log types.Log) error {
	if len(log.Topics) < 2 {
		return errors.New("malformed event")
	}
	eventID := fmt.Sprintf("%d:%s:%d", i.Campaigns.ChainID, log.TxHash.Hex(), log.Index)
	if log.Topics[0] == parsed.Events["Funded"].ID {
		if len(log.Topics) != 3 {
			return errors.New("bad funding topics")
		}
		values, e := parsed.Unpack("Funded", log.Data)
		if e != nil || len(values) != 1 {
			return errors.New("bad funding data")
		}
		amount := values[0].(*big.Int)
		if !amount.IsInt64() {
			return errors.New("amount overflow")
		}
		var orderID string
		if e = tx.QueryRow(ctx, "SELECT id::text FROM orders WHERE chain_order_id=$1", log.Topics[1].Hex()).Scan(&orderID); e != nil {
			return e
		}
		return i.Campaigns.ApplyEventTx(ctx, tx, campaigns.ChainEvent{ID: eventID, OrderID: orderID, Kind: "fund", Payer: common.BytesToAddress(log.Topics[2].Bytes()).Hex(), Token: i.Campaigns.Token, ChainID: i.Campaigns.ChainID, Amount: amount.Int64(), Block: int64(log.BlockNumber), TxHash: log.TxHash.Hex(), Final: true})
	}
	if len(log.Topics) != 4 {
		return errors.New("bad refund topics")
	}
	values, e := parsed.Unpack("Refunded", log.Data)
	if e != nil || len(values) != 1 {
		return errors.New("bad refund data")
	}
	amount := values[0].(*big.Int)
	if !amount.IsInt64() {
		return errors.New("amount overflow")
	}
	var orderID, payer string
	var funded, spent int64
	e = tx.QueryRow(ctx, "SELECT id::text,payer,funded,spent FROM orders WHERE chain_order_id=$1 FOR UPDATE", log.Topics[2].Hex()).Scan(&orderID, &payer, &funded, &spent)
	if e != nil {
		return e
	}
	var exists bool
	if e = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM chain_events WHERE id=$1)", eventID).Scan(&exists); e != nil {
		return e
	}
	if exists {
		return nil
	}
	var refundID string
	var expected int64
	e = tx.QueryRow(ctx, "SELECT id::text,amount FROM refund_intents WHERE chain_refund_id=$1 AND order_id=$2 AND status IN ('pending','submitted') FOR UPDATE", log.Topics[1].Hex(), orderID).Scan(&refundID, &expected)
	if e != nil {
		return e
	}
	if expected != amount.Int64() || common.HexToAddress(payer) != common.BytesToAddress(log.Topics[3].Bytes()) || expected > funded-spent {
		return errors.New("refund mismatch")
	}
	if e = ledger.Post(ctx, tx, "chain:"+eventID, "confirmed refund", ledger.Entry{Account: "cash:" + orderID, Amount: -expected}, ledger.Entry{Account: "escrow:" + orderID, Amount: expected}); e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, "UPDATE orders SET spent=spent+$2 WHERE id=$1", orderID, expected); e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, "UPDATE refund_intents SET status='confirmed' WHERE id=$1", refundID); e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, "INSERT INTO chain_events VALUES($1,$2,$3,'refund',$4,$5)", eventID, log.BlockHash.Hex(), orderID, int64(log.BlockNumber), log.TxHash.Hex()); e != nil {
		return e
	}
	return nil
}
