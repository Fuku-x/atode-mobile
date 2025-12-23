package auth

import "context"

type uidContextKey struct{}

type UID string

func WithUID(ctx context.Context, uid UID) context.Context {
	return context.WithValue(ctx, uidContextKey{}, uid)
}

func UIDFromContext(ctx context.Context) (UID, bool) {
	v := ctx.Value(uidContextKey{})
	uid, ok := v.(UID)
	return uid, ok
}
