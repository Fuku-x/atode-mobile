package db

import (
	"context"
	"database/sql"
)

func EnsureUsersTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (
	id uuid PRIMARY KEY,
	firebase_uid text UNIQUE NOT NULL,
	email text NOT NULL
);`)
	return err
}
