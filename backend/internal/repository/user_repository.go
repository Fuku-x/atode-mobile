package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) UpsertFirebaseUser(ctx context.Context, firebaseUID string, email string) (uuid.UUID, error) {
	id := uuid.New()

	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (id, firebase_uid, email)
VALUES ($1, $2, $3)
ON CONFLICT (firebase_uid) DO UPDATE
SET email = EXCLUDED.email
RETURNING id;`,
		id,
		firebaseUID,
		email,
	)

	var out uuid.UUID
	if err := row.Scan(&out); err != nil {
		return uuid.UUID{}, err
	}

	return out, nil
}
