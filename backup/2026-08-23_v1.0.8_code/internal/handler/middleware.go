package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"khan/internal/models"
)

// contextKey for storing the current user in request context
type ctxKey int

const userKey ctxKey = 1

// userFrom extracts the current user from request context
func userFrom(r *http.Request) *models.User {
	if u, ok := r.Context().Value(userKey).(*models.User); ok {
		return u
	}
	return nil
}

// contextWithUser stores the user in request context
func contextWithUser(r *http.Request, u *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userKey, u)
	return r.WithContext(ctx)
}

// RequireAuth is the exported auth middleware
func RequireAuth(auth interface {
	ValidateToken(token string) (*models.User, error)
}) func(http.Handler) http.Handler {
	return requireAuth(auth)
}

// requireAuth ensures a valid session token
func requireAuth(auth interface {
	ValidateToken(token string) (*models.User, error)
}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				writeErr(w, http.StatusUnauthorized, "وارد نشده‌اید")
				return
			}
			u, err := auth.ValidateToken(token)
			if err != nil {
				writeErr(w, http.StatusUnauthorized, err.Error())
				return
			}
			next.ServeHTTP(w, contextWithUser(r, u))
		})
	}
}

// RequireRole ensures the user has at least the given role (exported)
func RequireRole(role string) func(http.Handler) http.Handler {
	return requireRole(role)
}

// requireRole ensures the user has at least the given role
func requireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := userFrom(r)
			if u == nil {
				writeErr(w, http.StatusUnauthorized, "وارد نشده‌اید")
				return
			}
			if !hasRole(u, role) {
				writeErr(w, http.StatusForbidden, "دسترسی غیرمجاز")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// role rank for comparisons
func roleRank(u *models.User) int {
	switch u.Role {
	case models.RoleSuperAdmin:
		return 4
	case models.RoleAdmin:
		return 3
	case models.RoleSupervisor:
		return 2
	default:
		return 1
	}
}

func hasRole(u *models.User, required string) bool {
	return roleRank(u) >= roleRank(&models.User{Role: required})
}

// pathInt64 parses an int64 path parameter
func pathInt64(r *http.Request, key string) (int64, error) {
	val := r.PathValue(key)
	return strconv.ParseInt(val, 10, 64)
}

// parseBool reads a query param as bool
func parseBool(r *http.Request, key string) bool {
	return strings.EqualFold(r.URL.Query().Get(key), "true")
}
