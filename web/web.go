package web

import (
	"context"
	"fmt"
	"net/http"
)

// Encoder is an interface for encoding data into a response(Eq. to JSON.stringify in JS).
type Encoder interface {
	Encode() ([]byte, string, error)
}

type Logger func(ctx context.Context, msg string, args ...any)

type HandlerFunc func(ctx context.Context, r *http.Request) Encoder

type App struct {
	logger Logger
	mux    *http.ServeMux
	mw     []MidFunc
	cors   []string
}

func New(logger Logger, mw ...MidFunc) *App {
	return &App{
		logger: logger,
		mux:    http.NewServeMux(),
		mw:     mw,
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(a.cors) > 0 {
		reqOrigin := r.Header.Get("Origin")
		for _, origin := range a.cors {
			if reqOrigin == origin || reqOrigin == "*" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

	a.mux.ServeHTTP(w, r)
}

func (a *App) EnableCORS(origins []string) {
	a.cors = origins
}

func (a *App) HandlerFunc(method string, group string, path string, handler HandlerFunc, mw ...MidFunc) {
	mwChain := NewChain(append(a.mw, mw...)...)
	handler = mwChain.Then(handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		// TODO: Set tracing span for request

		resp := handler(r.Context(), r)
		if err := Respond(r.Context(), w, resp); err != nil {
			a.logger(r.Context(), "web.respond", "ERROR", err)
		}
	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	finalPath = fmt.Sprintf("/%s %s", method, finalPath)

	a.mux.HandleFunc(finalPath, h)
}
