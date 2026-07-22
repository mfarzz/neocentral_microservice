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
const UserRolesKey contextKey = "userRoles"

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
			ctx = context.WithValue(ctx, UserRolesKey, claims.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(UserIDKey).(string)
	return v
}

func GetUserRolesFromContext(ctx context.Context) []string {
	v, _ := ctx.Value(UserRolesKey).([]string)
	return v
}

// RequireRoleMiddleware ensures the user has at least one of the allowed roles.
func RequireRoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRoles := GetUserRolesFromContext(r.Context())
			
			hasAccess := false
			for _, userRole := range userRoles {
				for _, allowedRole := range allowedRoles {
					// Compare roles case-insensitively just in case
					if strings.EqualFold(userRole, allowedRole) {
						hasAccess = true
						break
					}
				}
				if hasAccess {
					break
				}
			}

			if !hasAccess {
				response.Error(w, http.StatusForbidden, "Forbidden: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
