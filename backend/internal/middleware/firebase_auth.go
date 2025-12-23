package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"example.com/atode/backend/internal/auth"
)

type authResponse struct {
	Error string `json:"error"`
}

func RequireFirebaseAuth(verifier auth.Verifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := bearerTokenFromHeader(r.Header.Get("Authorization"))
			if !ok {
				writeUnauthorized(w)
				return
			}

			uid, err := verifier.VerifyIDToken(r.Context(), token)
			if err != nil {
				writeUnauthorized(w)
				return
			}

			ctx := auth.WithUID(r.Context(), auth.UID(uid))
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

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(authResponse{Error: "unauthorized"})
}
