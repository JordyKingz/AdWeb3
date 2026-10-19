package ledger

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"math/big"
)

type Entry struct {
	Account string
	Amount  int64
}

func Post(ctx context.Context, tx pgx.Tx, id, description string, entries ...Entry) error {
	sum := new(big.Int)
	seen := map[string]bool{}
	for _, e := range entries {
		if e.Account == "" || e.Amount == 0 || seen[e.Account] {
			return errors.New("invalid journal entry")
		}
		seen[e.Account] = true
		sum.Add(sum, big.NewInt(e.Amount))
	}
	if len(entries) < 2 || sum.Sign() != 0 {
		return errors.New("unbalanced journal")
	}
	if _, e := tx.Exec(ctx, "INSERT INTO journals(id,description) VALUES($1,$2)", id, description); e != nil {
		return e
	}
	for _, v := range entries {
		if _, e := tx.Exec(ctx, "INSERT INTO ledger_entries VALUES($1,$2,$3)", id, v.Account, v.Amount); e != nil {
			return e
		}
	}
	return nil
}
