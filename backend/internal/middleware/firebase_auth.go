package middleware

import (
	"context"
	"net/http"
	"strings"

	"example.com/atode/backend/internal/auth"
	"example.com/atode/backend/internal/httpx"
	"github.com/google/uuid"
)

type UserUpserter interface {
	UpsertFirebaseUser(ctx context.Context, firebaseUID string, email string) (uuid.UUID, error)
}

func RequireFirebaseAuth(verifier auth.Verifier, upserter UserUpserter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerTokenFromHeader(r.Header.Get("Authorization"))
			if !ok {
				writeUnauthorized(w, r)
				return
			}

			fbUser, err := verifier.VerifyIDToken(r.Context(), token)
			if err != nil {
				writeUnauthorized(w, r)
				return
			}

			userID, err := upserter.UpsertFirebaseUser(r.Context(), fbUser.UID, fbUser.Email)
			if err != nil {
				writeInternalServerError(w, r)
				return
			}

			ctx := auth.WithUID(r.Context(), auth.UID(fbUser.UID))
			ctx = auth.WithEmail(ctx, fbUser.Email)
			ctx = auth.WithUserID(ctx, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func bearerTokenFromHeader(v string) (string, bool) {
	if v == "" {
		return "", false
	}
	parts := strings.SplitN(v, " ", 2)
	if len(parts) != 2 {
		return "", false
	}
	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	if parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func writeUnauthorized(w http.ResponseWriter, r *http.Request) {
	httpx.WriteError(w, r, http.StatusUnauthorized, "unauthorized", "unauthorized")
}

func writeInternalServerError(w http.ResponseWriter, r *http.Request) {
	httpx.WriteError(w, r, http.StatusInternalServerError, "internal_server_error", "internal server error")
}
