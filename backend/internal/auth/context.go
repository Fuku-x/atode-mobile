package auth

import (
	"context"

	"github.com/google/uuid"
)

type uidContextKey struct{}

type emailContextKey struct{}

type userIDContextKey struct{}

type UID string

func WithUID(ctx context.Context, uid UID) context.Context {
	return context.WithValue(ctx, uidContextKey{}, uid)
}

func UIDFromContext(ctx context.Context) (UID, bool) {
	v := ctx.Value(uidContextKey{})
	uid, ok := v.(UID)
	return uid, ok
}

func WithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, emailContextKey{}, email)
}

func EmailFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(emailContextKey{})
	email, ok := v.(string)
	return email, ok
}

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey{}, id)
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(userIDContextKey{})
	id, ok := v.(uuid.UUID)
	return id, ok
}
