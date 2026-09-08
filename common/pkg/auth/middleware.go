package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey struct{}

func UserFromContext(ctx context.Context) *User {
	u, _ := ctx.Value(contextKey{}).(*User)
	return u
}

func withUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, contextKey{}, u)
}

type RouteInfo struct {
	Method       string
	RoutePattern string
}

type RouteInfoExtractor func(r *http.Request) RouteInfo

func AuthnMiddleware(reg *Registry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				writeError(w, http.StatusUnauthorized, "Missing Authorization header")
				return
			}

			var user *User
			var ok bool

			if token, found := strings.CutPrefix(header, "Bearer "); found {
				user, ok = reg.AuthenticateBearer(token)
			} else if encoded, found := strings.CutPrefix(header, "Basic "); found {
				username, password, valid := decodeBasic(encoded)
				if !valid {
					writeError(w, http.StatusUnauthorized, "Invalid Basic credentials encoding")
					return
				}
				user, ok = reg.AuthenticateBasic(username, password)
			} else {
				writeError(w, http.StatusUnauthorized, "Unsupported authorization scheme")
				return
			}

			if !ok {
				writeError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			next.ServeHTTP(w, r.WithContext(withUser(r.Context(), user)))
		})
	}
}

func AuthzMiddleware(reg *Registry, extract RouteInfoExtractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user == nil {
				writeError(w, http.StatusUnauthorized, "Not authenticated")
				return
			}

			if user.AllPermissions {
				next.ServeHTTP(w, r)
				return
			}

			info := extract(r)
			perm, ok := reg.RequiredPermission(info.Method, info.RoutePattern)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			if !user.HasPermission(perm) {
				writeError(w, http.StatusForbidden, "Missing permission: "+perm+" for route "+info.RoutePattern)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func decodeBasic(encoded string) (username, password string, ok bool) {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}

type errorBody struct {
	HttpStatus int32  `json:"http_status"`
	Message    string `json:"message"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorBody{
		HttpStatus: int32(status),
		Message:    message,
	})
}
