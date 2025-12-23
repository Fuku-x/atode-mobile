package repository

import (
	"context"

	"github.com/google/uuid"
)

type UserUpserter interface {
	UpsertFirebaseUser(ctx context.Context, firebaseUID string, email string) (uuid.UUID, error)
}
