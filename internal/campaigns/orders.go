package campaigns

import (
	"adweb3.local/platform/internal/ledger"
	"adweb3.local/platform/internal/store"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

var ErrInvalid = errors.New("invalid request")
var ErrForbidden = errors.New("forbidden")
var ErrConflict = errors.New("conflict")

type Order struct {
	ID         string `json:"id"`
	Advertiser string `json:"advertiser"`
	Payer      string `json:"payer"`
	Amount     int64  `json:"amount,string"`
	Funded     int64  `json:"funded,string"`
	Spent      int64  `json:"spent,string"`
	ChainID    int64  `json:"chainId"`
	Token      string `json:"token"`
	Terms      string `json:"terms"`
	Start      int64  `json:"start"`
	End        int64  `json:"end"`
}
type ChainEvent struct {
	ID, OrderID, Kind, Payer, Token, TxHash string
	ChainID, Amount, Block                  int64
	Final                                   bool
}
type Service struct {
	DB      *pgxpool.Pool
	ChainID int64
	Token   string
}

const columns = "id::text,advertiser,payer,amount,funded,spent,chain_id,token,terms,start_at,end_at"

func scan(row pgx.Row) (o Order, e error) {
	e = row.Scan(&o.ID, &o.Advertiser, &o.Payer, &o.Amount, &o.Funded, &o.Spent, &o.ChainID, &o.Token, &o.Terms, &o.Start, &o.End)
	return
}
func fingerprint(v any) string {
	b, _ := json.Marshal(v)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func (s Service) CreateOrder(ctx context.Context, actor, key string, in Order) (Order, error) {
	if actor == "" || len(actor) > 128 || len(key) < 1 || len(key) > 128 || in.Amount <= 0 || in.Amount > 9_000_000_000_000_000 || in.Start < 0 || in.End <= in.Start || !common.IsHexAddress(in.Payer) || common.HexToAddress(in.Payer) == (common.Address{}) || in.ChainID != s.ChainID || (s.ChainID != 31337 && s.ChainID != 43113) || !strings.EqualFold(in.Token, s.Token) {
		return Order{}, ErrInvalid
	}
	o := Order{Payer: strings.ToLower(in.Payer), Amount: in.Amount, ChainID: in.ChainID, Token: strings.ToLower(in.Token), Start: in.Start, End: in.End}
	fp := fingerprint(o)
	o.ID = store.UUID()
	o.Advertiser = actor
	o.Terms = crypto.Keccak256Hash([]byte(fp + ":3000:6000:1000:v1")).Hex()
	tx, e := s.DB.Begin(ctx)
	if e != nil {
		return Order{}, e
	}
	defer tx.Rollback(ctx)
	_, e = tx.Exec(ctx, "INSERT INTO orders(id,advertiser,idempotency_key,fingerprint,payer,amount,chain_id,token,terms,start_at,end_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT(advertiser,idempotency_key) DO NOTHING", o.ID, actor, key, fp, o.Payer, o.Amount, o.ChainID, o.Token, o.Terms, o.Start, o.End)
	if e != nil {
		return Order{}, e
	}
	var oldFP string
	e = tx.QueryRow(ctx, "SELECT fingerprint FROM orders WHERE advertiser=$1 AND idempotency_key=$2", actor, key).Scan(&oldFP)
	if e != nil {
		return Order{}, e
	}
	if oldFP != fp {
		return Order{}, ErrConflict
	}
	o, e = scan(tx.QueryRow(ctx, "SELECT "+columns+" FROM orders WHERE advertiser=$1 AND idempotency_key=$2", actor, key))
	if e != nil {
		return Order{}, e
	}
	if _, e = tx.Exec(ctx, "UPDATE orders SET chain_order_id=$2 WHERE id=$1 AND chain_order_id IS NULL", o.ID, crypto.Keccak256Hash([]byte(o.ID)).Hex()); e != nil {
		return Order{}, e
	}
	return o, tx.Commit(ctx)
}
func (s Service) Get(ctx context.Context, actor, id string) (Order, error) {
	o, e := scan(s.DB.QueryRow(ctx, "SELECT "+columns+" FROM orders WHERE id::text=$1 AND advertiser=$2", id, actor))
	if errors.Is(e, pgx.ErrNoRows) {
		return Order{}, ErrForbidden
	}
	return o, e
}
func (s Service) ApplyFinalizedEvent(ctx context.Context, event ChainEvent) error {
	tx, e := s.DB.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	if e = s.ApplyEventTx(ctx, tx, event); e != nil {
		return e
	}
	return tx.Commit(ctx)
}
func (s Service) ApplyEventTx(ctx context.Context, tx pgx.Tx, event ChainEvent) error {
	if !event.Final || event.ID == "" || event.TxHash == "" || event.Block < 0 || event.Amount <= 0 || event.ChainID != s.ChainID || !strings.EqualFold(event.Token, s.Token) {
		return ErrInvalid
	}
	o, e := scan(tx.QueryRow(ctx, "SELECT "+columns+" FROM orders WHERE id::text=$1 FOR UPDATE", event.OrderID))
	if e != nil {
		return e
	}
	fp := fingerprint(event)
	var previous string
	e = tx.QueryRow(ctx, "SELECT fingerprint FROM chain_events WHERE id=$1", event.ID).Scan(&previous)
	if e == nil {
		if previous != fp {
			return ErrConflict
		}
		return nil
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return e
	}
	if event.Kind != "fund" || !strings.EqualFold(event.Payer, o.Payer) || event.Amount != o.Amount || o.Funded != 0 {
		return ErrInvalid
	}
	_, e = tx.Exec(ctx, "UPDATE orders SET funded=$2 WHERE id=$1", o.ID, event.Amount)
	if e != nil {
		return e
	}
	if e = ledger.Post(ctx, tx, "chain:"+event.ID, "confirmed campaign funding", ledger.Entry{Account: "cash:" + o.ID, Amount: event.Amount}, ledger.Entry{Account: "escrow:" + o.ID, Amount: -event.Amount}); e != nil {
		return e
	}
	_, e = tx.Exec(ctx, "INSERT INTO chain_events VALUES($1,$2,$3,$4,$5,$6)", event.ID, fp, o.ID, event.Kind, event.Block, event.TxHash)
	if e != nil {
		return e
	}
	return nil
}
func (s Service) RequestRefund(ctx context.Context, actor, orderID, key string, amount int64) (string, error) {
	if amount <= 0 || key == "" || len(key) > 128 {
		return "", ErrInvalid
	}
	tx, e := s.DB.Begin(ctx)
	if e != nil {
		return "", e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(hashtextextended($1,0))", actor+":"+key); e != nil {
		return "", e
	}
	if e = CheckChain(ctx, tx, s.ChainID); e != nil {
		return "", e
	}
	o, e := scan(tx.QueryRow(ctx, "SELECT "+columns+" FROM orders WHERE id::text=$1 AND advertiser=$2 FOR UPDATE", orderID, actor))
	if errors.Is(e, pgx.ErrNoRows) {
		return "", ErrForbidden
	}
	if e != nil {
		return "", e
	}
	var id, oldOrder string
	var oldAmount int64
	e = tx.QueryRow(ctx, "SELECT id::text,order_id::text,amount FROM refund_intents WHERE actor=$1 AND idempotency_key=$2", actor, key).Scan(&id, &oldOrder, &oldAmount)
	if e == nil {
		if oldOrder != o.ID || oldAmount != amount {
			return "", ErrConflict
		}
		return id, tx.Commit(ctx)
	}
	if !errors.Is(e, pgx.ErrNoRows) {
		return "", e
	}
	var reserved int64
	if e = tx.QueryRow(ctx, "SELECT COALESCE(sum(amount),0) FROM refund_intents WHERE order_id=$1 AND status IN ('pending','submitted')", o.ID).Scan(&reserved); e != nil {
		return "", e
	}
	if amount > o.Funded-o.Spent-reserved {
		return "", ErrConflict
	}
	id = store.UUID()
	_, e = tx.Exec(ctx, "INSERT INTO refund_intents(id,order_id,actor,idempotency_key,amount,chain_refund_id) VALUES($1,$2,$3,$4,$5,$6)", id, o.ID, actor, key, amount, crypto.Keccak256Hash([]byte(id)).Hex())
	if e != nil {
		return "", e
	}
	payload := fmt.Sprintf(`{"payer":%q,"amount":%q}`, o.Payer, fmt.Sprint(amount))
	_, e = tx.Exec(ctx, "INSERT INTO outbox(kind,aggregate_id,payload) VALUES('refund',$1,$2::jsonb)", id, payload)
	if e != nil {
		return "", e
	}
	return id, tx.Commit(ctx)
}

// CheckChain holds a shared health lock for the transaction. Reorg quarantine blocks new spending.
func CheckChain(ctx context.Context, tx pgx.Tx, chainID int64) error {
	if _, e := tx.Exec(ctx, "INSERT INTO chain_health(chain_id) VALUES($1) ON CONFLICT DO NOTHING", chainID); e != nil {
		return e
	}
	var quarantined bool
	if e := tx.QueryRow(ctx, "SELECT quarantined FROM chain_health WHERE chain_id=$1 FOR SHARE", chainID).Scan(&quarantined); e != nil {
		return e
	}
	if quarantined {
		return ErrConflict
	}
	return nil
}
