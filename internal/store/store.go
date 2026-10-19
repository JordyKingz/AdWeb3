package store

import (
	"adweb3.local/platform/migrations"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"sort"
)

func UUID() string {
	var b [16]byte
	if _, e := rand.Read(b[:]); e != nil {
		panic(e)
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	s := hex.EncodeToString(b[:])
	return s[:8] + "-" + s[8:12] + "-" + s[12:16] + "-" + s[16:20] + "-" + s[20:]
}
func Migrate(ctx context.Context, p *pgxpool.Pool) error {
	tx, e := p.Begin(ctx)
	if e != nil {
		return e
	}
	defer tx.Rollback(ctx)
	if _, e = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(8293142)"); e != nil {
		return e
	}
	if _, e = tx.Exec(ctx, "CREATE TABLE IF NOT EXISTS schema_migrations(version text PRIMARY KEY)"); e != nil {
		return e
	}
	names, e := migrations.Files.ReadDir(".")
	if e != nil {
		return e
	}
	sort.Slice(names, func(i, j int) bool { return names[i].Name() < names[j].Name() })
	for _, n := range names {
		var exists bool
		if e = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)", n.Name()).Scan(&exists); e != nil {
			return e
		}
		if exists {
			continue
		}
		data, e := migrations.Files.ReadFile(n.Name())
		if e != nil {
			return e
		}
		if _, e = tx.Exec(ctx, string(data)); e != nil {
			return fmt.Errorf("migration %s: %w", n.Name(), e)
		}
		if _, e = tx.Exec(ctx, "INSERT INTO schema_migrations VALUES($1)", n.Name()); e != nil {
			return e
		}
	}
	return tx.Commit(ctx)
}
