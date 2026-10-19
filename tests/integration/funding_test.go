//go:build integration

package integration

import (
	"adweb3.local/platform/internal/campaigns"
	"adweb3.local/platform/internal/store"
	"bytes"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

const token = "0x0000000000000000000000000000000000000001"
const payer = "0x0000000000000000000000000000000000000002"

func database(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Fatal("DATABASE_URL required")
	}
	admin, e := pgxpool.New(ctx, url)
	if e != nil {
		t.Fatal(e)
	}
	schema := "test_" + strings.ReplaceAll(store.UUID(), "-", "")
	if _, e = admin.Exec(ctx, "CREATE SCHEMA "+schema); e != nil {
		t.Fatal(e)
	}
	cfg, e := pgxpool.ParseConfig(url)
	if e != nil {
		t.Fatal(e)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	p, e := pgxpool.NewWithConfig(ctx, cfg)
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { p.Close(); admin.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE"); admin.Close() })
	if e = store.Migrate(ctx, p); e != nil {
		t.Fatal(e)
	}
	return p
}
func TestFundingIdempotencyAndBacking(t *testing.T) {
	ctx := context.Background()
	key := store.UUID()
	s := campaigns.Service{DB: database(t), ChainID: 31337, Token: token}
	in := campaigns.Order{Payer: payer, Amount: 25_000_000, ChainID: 31337, Token: token, Start: 100, End: 200}
	o, e := s.CreateOrder(ctx, "advertiser-a", key, in)
	if e != nil {
		t.Fatal("create valid order:", e)
	}
	again, e := s.CreateOrder(ctx, "advertiser-a", key, in)
	if e != nil || again.ID != o.ID {
		t.Fatal("create retry not idempotent", e)
	}
	in.Amount++
	if _, e = s.CreateOrder(ctx, "advertiser-a", key, in); e == nil {
		t.Fatal("changed idempotency body accepted")
	}
	if _, e = s.Get(ctx, "advertiser-b", o.ID); e == nil {
		t.Fatal("tenant isolation")
	}
	event := campaigns.ChainEvent{ID: "event-" + key, OrderID: o.ID, Kind: "fund", Payer: payer, Token: token, ChainID: 31337, Amount: 25_000_000, Final: true, TxHash: "tx-a", Block: 1}
	wrong := event
	wrong.ID = "wrong"
	wrong.Token = payer
	if s.ApplyFinalizedEvent(ctx, wrong) == nil {
		t.Fatal("wrong token accepted")
	}
	pending := event
	pending.Final = false
	if s.ApplyFinalizedEvent(ctx, pending) == nil {
		t.Fatal("unfinalized event accepted")
	}
	if e = s.ApplyFinalizedEvent(ctx, event); e != nil {
		t.Fatal(e)
	}
	if e = s.ApplyFinalizedEvent(ctx, event); e != nil {
		t.Fatal("replay", e)
	}
	o, e = s.Get(ctx, "advertiser-a", o.ID)
	if e != nil || o.Funded != 25_000_000 {
		t.Fatal("funding doubled/missing", o, e)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	count := 0
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, e := s.RequestRefund(ctx, "advertiser-a", o.ID, string(rune('a'+i)), 20_000_000)
			if e == nil {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	if count != 1 {
		t.Fatalf("concurrent refund reservations: %d", count)
	}
	var imbalance int
	if e = s.DB.QueryRow(ctx, "SELECT count(*) FROM (SELECT journal_id FROM ledger_entries GROUP BY journal_id HAVING sum(amount) <> 0) x").Scan(&imbalance); e != nil || imbalance != 0 {
		t.Fatal("unbalanced ledger", e)
	}
}

func TestMigrationAndImmutableHistory(t *testing.T) {
	ctx := context.Background()
	db := database(t)
	if e := store.Migrate(ctx, db); e != nil {
		t.Fatal("repeat migration", e)
	}
	s := campaigns.Service{DB: db, ChainID: 31337, Token: token}
	key := store.UUID()
	o, e := s.CreateOrder(ctx, "a", key, campaigns.Order{Payer: payer, Amount: 100, ChainID: 31337, Token: token, Start: 1, End: 2})
	if e != nil {
		t.Fatal(e)
	}
	event := campaigns.ChainEvent{ID: key, OrderID: o.ID, Kind: "fund", Payer: payer, Token: token, ChainID: 31337, Amount: 100, Block: 1, TxHash: "tx", Final: true}
	for _, changed := range []campaigns.ChainEvent{func() campaigns.ChainEvent { x := event; x.ChainID = 43113; return x }(), func() campaigns.ChainEvent { x := event; x.Amount = 99; return x }()} {
		if e = s.ApplyFinalizedEvent(ctx, changed); e == nil {
			t.Fatal("invalid funding accepted")
		}
	}
	if e = s.ApplyFinalizedEvent(ctx, event); e != nil {
		t.Fatal(e)
	}
	if _, e = db.Exec(ctx, "UPDATE ledger_entries SET amount=1"); e == nil {
		t.Fatal("history mutable")
	}
	if e = s.ApplyFinalizedEvent(ctx, event); e != nil {
		t.Fatal("restart replay", e)
	}
	// A new service instance after replay uses the durable funding and original terms.
	restarted := campaigns.Service{DB: db, ChainID: 31337, Token: token}
	got, e := restarted.Get(ctx, "a", o.ID)
	if e != nil || got.Funded != 100 || got.Terms != o.Terms {
		t.Fatal("durable state/terms changed", e)
	}
}

func TestProtectedHTTPOrderIsolation(t *testing.T) {
	ctx := context.Background()
	db := database(t)
	s := campaigns.Service{DB: db, ChainID: 31337, Token: token}
	o, e := s.CreateOrder(ctx, "a", store.UUID(), campaigns.Order{Payer: payer, Amount: 100, ChainID: 31337, Token: token, Start: 1, End: 2})
	if e != nil {
		t.Fatal(e)
	}
	h := campaigns.Handler(s, map[string]string{"token-a": "a", "token-b": "b"})
	r := httptest.NewRequest("GET", "/v1/orders/"+o.ID, nil)
	r.Header.Set("Authorization", "Bearer token-b")
	r.Header.Set("X-Actor", "a")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 404 {
		t.Fatal("cross-actor UI request authorized", w.Code)
	}
}

func TestBackupRestorePreservesFunding(t *testing.T) {
	ctx := context.Background()
	db := database(t)
	s := campaigns.Service{DB: db, ChainID: 31337, Token: token}
	o, e := s.CreateOrder(ctx, "a", store.UUID(), campaigns.Order{Payer: payer, Amount: 100, ChainID: 31337, Token: token, Start: 1, End: 2})
	if e != nil {
		t.Fatal(e)
	}
	if e = s.ApplyFinalizedEvent(ctx, campaigns.ChainEvent{ID: store.UUID(), OrderID: o.ID, Kind: "fund", Payer: payer, Token: token, ChainID: 31337, Amount: 100, Block: 1, TxHash: "fixture", Final: true}); e != nil {
		t.Fatal(e)
	}
	var schema string
	if e = db.QueryRow(ctx, "SELECT current_schema()").Scan(&schema); e != nil {
		t.Fatal(e)
	}
	container := os.Getenv("POSTGRES_CONTAINER")
	if container == "" {
		container = "adweb3-postgres"
	}
	backup, e := exec.CommandContext(ctx, "docker", "exec", container, "pg_dump", "-U", "postgres", "-d", "adweb3", "--no-owner", "--no-acl", "--schema", schema).Output()
	if e != nil {
		t.Fatal("backup required", e)
	}
	if _, e = db.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE"); e != nil {
		t.Fatal(e)
	}
	command := exec.CommandContext(ctx, "docker", "exec", "-i", container, "psql", "-U", "postgres", "-d", "adweb3", "-v", "ON_ERROR_STOP=1")
	command.Stdin = bytes.NewReader(backup)
	if out, e := command.CombinedOutput(); e != nil {
		t.Fatalf("restore failed: %v %s", e, out)
	}
	got, e := s.Get(ctx, "a", o.ID)
	if e != nil || got.Funded != 100 || got.Terms != o.Terms {
		t.Fatal("restored funding/terms mismatch", got, e)
	}
}
