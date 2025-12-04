package mux

import (
	"embed"
	"net/http"

	"foundation/app/sdk/middleware"
	"foundation/logger"
	"foundation/web"
)

type StaticSite struct {
	react      bool
	static     embed.FS
	staticPath string
	staticDir  string
}

type Options struct {
	corsOrigins []string
	sites       []StaticSite
}

func WithCORS(origins []string) func(*Options) {
	return func(o *Options) {
		o.corsOrigins = origins
	}
}

func WithFileServer(react bool, static embed.FS, dir string, path string) func(*Options) {
	return func(o *Options) {
		o.sites = append(o.sites, StaticSite{
			react:      react,
			static:     static,
			staticPath: path,
			staticDir:  dir,
		})
	}
}

type Config struct {
	Log *logger.Logger
}

func New(config Config, opts ...func(*Options)) http.Handler {
	app := web.New(
		config.Log.Info,
		middleware.Logging(config.Log),
		middleware.Errors(config.Log))

	var options Options
	for _, opt := range opts {
		opt(&options)
	}

	if len(options.corsOrigins) > 0 {
		app.EnableCORS(options.corsOrigins)
	}

	return app
}
