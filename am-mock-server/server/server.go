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

// Package server is the AM Automation mock HTTP server and in-memory StrictServer implementation.
package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gravitee-io-labs/gravitee-automation-tools/common/pkg/auth"
)

// BasePath is the default API prefix ("/automation").
const BasePath = "/automation"

// New mounts the API at BasePath with no auth.
func New(impl *MockAM) http.Handler {
	return NewWithPath(impl, BasePath, nil)
}

// NewWithPath mounts the API at basePath. nil Registry leaves the API open. nil impl serves Unimplemented.
func NewWithPath(mockAM *MockAM, basePath string, reg *auth.Registry) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	var middlewares []MiddlewareFunc
	if reg != nil {
		r.Use(auth.AuthnMiddleware(reg))
		middlewares = append(middlewares, auth.AuthzMiddleware(reg, chiRouteInfoExtractor()))
	}
	middlewares = append(middlewares, DomainParentCheck(chiRouteInfoExtractor(), func(org, env, key string) bool {
		_, exists := mockAM.getTenant(orgEnv{org: org, env: env}).Domains.Get(key)
		return exists
	}))

	var si ServerInterface = Unimplemented{}
	if mockAM != nil {
		si = NewStrictHandler(mockAM, nil)
	}

	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:     basePath,
		BaseRouter:  r,
		Middlewares: middlewares,
	})
}

func chiRouteInfoExtractor() auth.RouteInfoExtractor {
	return func(r *http.Request) auth.RouteInfo {
		context := chi.RouteContext(r.Context())
		return auth.RouteInfo{
			Method:       r.Method,
			RoutePattern: context.RoutePattern(),
			RouteParams: func(params chi.RouteParams) map[string]string {
				r := make(map[string]string)
				for i, key := range params.Keys {
					r[key] = params.Values[i]
				}
				return r
			}(context.URLParams),
		}
	}

}
