package handler

import (
	"context"
	"net/http"
	"strings"

	"neocentral-go/pkg/auth"
	"neocentral-go/pkg/response"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserEmailKey contextKey = "userEmail"

// AuthGuardMiddleware validates the JWT from the Authorization header
// and injects the user ID into the request context.
func AuthGuardMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := ""
			if authHeader := r.Header.Get("Authorization"); authHeader != "" {
				if strings.HasPrefix(authHeader, "Bearer ") {
					tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
				} else {
					tokenStr = authHeader
				}
			} else if q := r.URL.Query().Get("token"); q != "" {
				tokenStr = q
			}

			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			claims, err := auth.VerifyAccessToken(tokenStr, jwtSecret)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
			ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(UserIDKey).(string)
	return v
}
