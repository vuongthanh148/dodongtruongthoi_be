package database

import (
"context"

"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
if databaseURL == "" {
return nil, nil
}

pool, err := pgxpool.New(ctx, databaseURL)
if err != nil {
return nil, err
}

if err := pool.Ping(ctx); err != nil {
pool.Close()
return nil, err
}

return pool, nil
}

func Close(pool *pgxpool.Pool) {
if pool == nil {
return
}
pool.Close()
}
