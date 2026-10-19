package chain

import (
	"adweb3.local/platform/internal/campaigns"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
	"strings"
)

// Prepare persists the signed transaction before any broadcast. Restart reuses its nonce and bytes.
func (w Worker) Prepare(ctx context.Context, id, method string, args ...any) (*types.Transaction, error) {
	input, e := w.ABI.Pack(method, args...)
	if e != nil {
		return nil, e
	}
	hash := crypto.Keccak256Hash(input).Hex()
	tx, e := w.DB.Begin(ctx)
	if e != nil {
		return nil, e
	}
	defer tx.Rollback(ctx)
	if e = campaigns.CheckChain(ctx, tx, w.Campaigns.ChainID); e != nil {
		return nil, e
	}
	if _, e = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(hashtextextended($1,0))", fmt.Sprintf("nonce:%d:%s", w.Campaigns.ChainID, w.Signer.From.Hex())); e != nil {
		return nil, e
	}
	var raw []byte
	var oldHash, status string
	e = tx.QueryRow(ctx, "SELECT raw,input_hash,status FROM transaction_intents WHERE operation_id=$1", id).Scan(&raw, &oldHash, &status)
	if e == nil {
		if oldHash != hash || status == "failed" {
			return nil, errors.New("conflicting or failed transaction intent")
		}
		var signed types.Transaction
		if e = signed.UnmarshalBinary(raw); e != nil {
			return nil, e
		}
		return &signed, tx.Commit(ctx)
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return nil, e
	}
	contract := bind.NewBoundContract(w.Escrow, w.ABI, w.Client, w.Client, w.Client)
	opts := *w.Signer
	opts.Context = ctx
	opts.NoSend = true
	signed, e := contract.Transact(&opts, method, args...)
	if e != nil {
		return nil, e
	}
	raw, e = signed.MarshalBinary()
	if e != nil {
		return nil, e
	}
	_, e = tx.Exec(ctx, "INSERT INTO transaction_intents VALUES($1,$2,$3,$4,$5,$6,$7,$8,'pending')", id, w.Campaigns.ChainID, w.Signer.From.Hex(), w.Escrow.Hex(), signed.Nonce(), hash, signed.Hash().Hex(), raw)
	if e != nil {
		return nil, e
	}
	return signed, tx.Commit(ctx)
}
func (w Worker) broadcast(ctx context.Context, id string, signed *types.Transaction) error {
	receipt, e := w.Client.TransactionReceipt(ctx, signed.Hash())
	if errors.Is(e, ethereum.NotFound) {
		if e = w.Client.SendTransaction(ctx, signed); e != nil && !strings.Contains(e.Error(), "already known") && !strings.Contains(e.Error(), "nonce too low") {
			return e
		}
		receipt, e = bind.WaitMined(ctx, w.Client, signed)
	}
	if e != nil {
		return e
	}
	status := "confirmed"
	if receipt.Status != 1 {
		status = "failed"
	}
	if _, e = w.DB.Exec(ctx, "UPDATE transaction_intents SET status=$2 WHERE operation_id=$1", id, status); e != nil {
		return e
	}
	if status == "failed" {
		return errors.New("transaction reverted")
	}
	return nil
}
func (w Worker) submit(ctx context.Context, id, method string, args ...any) error {
	signed, e := w.Prepare(ctx, id, method, args...)
	if e != nil {
		return e
	}
	return w.broadcast(ctx, id, signed)
}
func (w Worker) recoverPending(ctx context.Context) error {
	rows, e := w.DB.Query(ctx, "SELECT operation_id,raw FROM transaction_intents WHERE chain_id=$1 AND sender=$2 AND contract=$3 AND status='pending' ORDER BY nonce", w.Campaigns.ChainID, w.Signer.From.Hex(), w.Escrow.Hex())
	if e != nil {
		return e
	}
	type pending struct {
		id  string
		raw []byte
	}
	var items []pending
	for rows.Next() {
		var p pending
		if e = rows.Scan(&p.id, &p.raw); e != nil {
			rows.Close()
			return e
		}
		items = append(items, p)
	}
	e = rows.Err()
	rows.Close()
	if e != nil {
		return e
	}
	for _, p := range items {
		var tx types.Transaction
		if e = tx.UnmarshalBinary(p.raw); e != nil {
			return e
		}
		if e = w.broadcast(ctx, p.id, &tx); e != nil {
			return e
		}
	}
	return nil
}
