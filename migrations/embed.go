package migrations

import "embed"

// Files contains append-only schema revisions.
//
//go:embed *.sql
var Files embed.FS
