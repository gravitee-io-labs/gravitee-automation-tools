package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const BasePath = "/automation"

func New(impl StrictServerInterface) http.Handler {
	return NewWithPath(impl, BasePath)
}

func NewWithPath(impl StrictServerInterface, basePath string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	var si ServerInterface = Unimplemented{}
	if impl != nil {
		si = NewStrictHandler(impl, nil)
	}

	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    basePath,
		BaseRouter: r,
	})
}
