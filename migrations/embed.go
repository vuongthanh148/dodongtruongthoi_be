package migrations

import "embed"

// Files holds all SQL migration files, embedded into the binary so the
// server can apply schema changes on startup without external tooling.
//
//go:embed *.sql
var Files embed.FS
