package chain

import (
	"adweb3.local/platform/internal/campaigns"
	"context"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/big"
	"os"
	"strings"
	"time"
)

type Worker struct {
	DB            *pgxpool.Pool
	Client        *ethclient.Client
	Campaigns     campaigns.Service
	Escrow        common.Address
	Signer        *bind.TransactOpts
	ABI           abi.ABI
	Confirmations uint64
}

func ContractABI(path string) (abi.ABI, error) {
	data, e := os.ReadFile(path)
	if e != nil {
		return abi.ABI{}, e
	}
	var artifact struct {
		ABI json.RawMessage `json:"abi"`
	}
	if e = json.Unmarshal(data, &artifact); e != nil {
		return abi.ABI{}, e
	}
	return abi.JSON(strings.NewReader(string(artifact.ABI)))
}
func (w Worker) ProcessOnce(ctx context.Context) error {
	conn, e := w.DB.Acquire(ctx)
	if e != nil {
		return e
	}
	defer conn.Release()
	var acquired bool
	if e = conn.QueryRow(ctx, "SELECT pg_try_advisory_lock(8293143)").Scan(&acquired); e != nil {
		return e
	}
	if !acquired {
		return nil
	}
	defer conn.Exec(context.Background(), "SELECT pg_advisory_unlock(8293143)")
	chainID, e := w.Client.ChainID(ctx)
	if e != nil {
		return e
	}
	if chainID.Int64() != w.Campaigns.ChainID || chainID.Int64() != 31337 && chainID.Int64() != 43113 {
		return errors.New("unsafe chain")
	}
	indexer := Indexer{Client: w.Client, DB: w.DB, Campaigns: w.Campaigns, Escrow: w.Escrow, Confirmations: w.Confirmations}
	if e = indexer.Scan(ctx); e != nil {
		return e
	}
	if e = w.recoverPending(ctx); e != nil {
		return e
	}
	contract := bind.NewBoundContract(w.Escrow, w.ABI, w.Client, w.Client, w.Client)
	rows, e := w.DB.Query(ctx, "SELECT o.id::text,o.payer,o.amount,o.terms,o.end_at FROM orders o WHERE o.funded=0 AND o.chain_id=$1 AND lower(o.token)=lower($2) AND NOT EXISTS(SELECT 1 FROM chain_registrations r WHERE r.order_id=o.id) ORDER BY o.id LIMIT 100", w.Campaigns.ChainID, w.Campaigns.Token)
	if e != nil {
		return e
	}
	type order struct {
		id, payer, terms string
		amount, end      int64
	}
	var orders []order
	for rows.Next() {
		var o order
		if e = rows.Scan(&o.id, &o.payer, &o.amount, &o.terms, &o.end); e != nil {
			rows.Close()
			return e
		}
		orders = append(orders, o)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return e
	}
	for _, o := range orders {
		var output []any
		id := crypto.Keccak256Hash([]byte(o.id))
		if e = contract.Call(&bind.CallOpts{Context: ctx}, &output, "orders", id); e != nil {
			return e
		}
		if output[0].(common.Address) != (common.Address{}) {
			if _, e = w.DB.Exec(ctx, "INSERT INTO chain_registrations(order_id,chain_id,contract) VALUES($1,$2,$3) ON CONFLICT DO NOTHING", o.id, w.Campaigns.ChainID, w.Escrow.Hex()); e != nil {
				return e
			}
			continue
		}
		if o.end+86400 <= time.Now().Unix() {
			continue
		}
		if e = w.submit(ctx, "register:"+o.id, "registerOrder", id, common.HexToAddress(o.payer), big.NewInt(o.amount), common.HexToHash(o.terms), uint64(o.end+86400)); e != nil {
			return e
		}

		if _, e = w.DB.Exec(ctx, "INSERT INTO chain_registrations(order_id,chain_id,contract) VALUES($1,$2,$3) ON CONFLICT DO NOTHING", o.id, w.Campaigns.ChainID, w.Escrow.Hex()); e != nil {
			return e
		}
	}
	rows, e = w.DB.Query(ctx, "SELECT r.id::text,r.order_id::text,r.amount,o.payer FROM refund_intents r JOIN orders o ON o.id=r.order_id WHERE r.status IN ('pending','submitted') AND o.chain_id=$1 AND lower(o.token)=lower($2) ORDER BY r.id LIMIT 100", w.Campaigns.ChainID, w.Campaigns.Token)
	if e != nil {
		return e
	}
	type refund struct {
		id, order, payer string
		amount           int64
	}
	var refunds []refund
	for rows.Next() {
		var r refund
		if e = rows.Scan(&r.id, &r.order, &r.amount, &r.payer); e != nil {
			rows.Close()
			return e
		}
		refunds = append(refunds, r)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return e
	}
	for _, r := range refunds {
		var out []any
		id := crypto.Keccak256Hash([]byte(r.id))
		if e = contract.Call(&bind.CallOpts{Context: ctx}, &out, "processed", id); e != nil {
			return e
		}
		if out[0].(bool) {
			continue
		}
		if e = w.submit(ctx, "refund:"+r.id, "refundOrder", id, crypto.Keccak256Hash([]byte(r.order)), common.HexToAddress(r.payer), big.NewInt(r.amount)); e != nil {
			return e
		}
		if _, e = w.DB.Exec(ctx, "UPDATE refund_intents SET status='submitted' WHERE id=$1 AND status='pending'", r.id); e != nil {
			return e
		}

	}
	return indexer.Scan(ctx)
}
