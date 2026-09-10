// Copyright (C) 2015 The Gravitee team (http://gravitee.io)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

type contextKey struct{}

// UserFromContext returns the User stored by AuthnMiddleware, or nil if the request is unauthenticated.
func UserFromContext(ctx context.Context) *User {
	u, _ := ctx.Value(contextKey{}).(*User)
	return u
}

func withUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, contextKey{}, u)
}

// RouteInfo is the method and chi route pattern used to look up a required permission.
// RoutePattern must match how permissions were registered (basePath + YAML path).
type RouteInfo struct {
	Method       string
	RoutePattern string
	RouteParams  map[string]string
}

// RouteInfoExtractor reads RouteInfo from the current request (typically chi.RouteContext).
type RouteInfoExtractor func(r *http.Request) RouteInfo

// AuthnMiddleware rejects missing/invalid Authorization with 401. On success it stores *User in the request context.
func AuthnMiddleware(reg *Registry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				WriteError(w, http.StatusUnauthorized, "Missing Authorization header")
				return
			}

			var user *User
			var ok bool

			if token, found := strings.CutPrefix(header, "Bearer "); found {
				user, ok = reg.AuthenticateBearer(token)
			} else if encoded, found := strings.CutPrefix(header, "Basic "); found {
				username, password, valid := decodeBasic(encoded)
				if !valid {
					WriteError(w, http.StatusUnauthorized, "Invalid Basic credentials encoding")
					return
				}
				user, ok = reg.AuthenticateBasic(username, password)
			} else {
				WriteError(w, http.StatusUnauthorized, "Unsupported authorization scheme")
				return
			}

			if !ok {
				WriteError(w, http.StatusUnauthorized, "Invalid credentials")
				return
			}

			next.ServeHTTP(w, r.WithContext(withUser(r.Context(), user)))
		})
	}
}

// AuthzMiddleware requires a User in context. Missing user is 401. Missing permission is 403.
// Unconfigured routes and AllPermissions users are allowed through.
func AuthzMiddleware(reg *Registry, extract RouteInfoExtractor) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := UserFromContext(r.Context())
			if user == nil {
				WriteError(w, http.StatusUnauthorized, "Not authenticated")
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
				WriteError(w, http.StatusForbidden, "Missing permission: "+perm+" for route "+info.RoutePattern)
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

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorBody{
		HttpStatus: int32(status),
		Message:    message,
	})
}
