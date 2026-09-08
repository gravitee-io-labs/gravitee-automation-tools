package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gravitee-io-labs/gravitee-automation-sdks/common/pkg/auth"
)

const BasePath = "/automation"

func New(impl StrictServerInterface) http.Handler {
	return NewWithPath(impl, BasePath, nil)
}

func NewWithPath(impl StrictServerInterface, basePath string, reg *auth.Registry) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	var middlewares []MiddlewareFunc
	if reg != nil {
		r.Use(auth.AuthnMiddleware(reg))
		middlewares = append(middlewares, auth.AuthzMiddleware(reg, chiRouteInfoExtractor()))
	}

	var si ServerInterface = Unimplemented{}
	if impl != nil {
		si = NewStrictHandler(impl, nil)
	}

	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:     basePath,
		BaseRouter:  r,
		Middlewares: middlewares,
	})
}

func chiRouteInfoExtractor() auth.RouteInfoExtractor {
	return func(r *http.Request) auth.RouteInfo {
		return auth.RouteInfo{
			Method:       r.Method,
			RoutePattern: chi.RouteContext(r.Context()).RoutePattern(),
		}
	}
}
