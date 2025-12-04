package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"foundation/app/sdk/errs"
	"foundation/logger"
	"foundation/web"
)

func Logging(log *logger.Logger) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		f := func(ctx context.Context, r *http.Request) web.Encoder {
			now := time.Now()
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info(ctx, "request received",
				"method", r.Method,
				"path", path,
				"remote_addr", r.RemoteAddr)

			resp := next(ctx, r)

			status := errs.None
			if err := checkIsError(resp); err != nil {
				status = errs.Internal
				var appError *errs.Error
				if errors.As(err, &appError) {
					status = appError.Code
				}
			}
			log.Info(ctx, "request completed",
				"method", r.Method,
				"path", path,
				"remote_addr", r.RemoteAddr,
				"status", status,
				"duration", time.Since(now).String())

			return resp
		}
		return f
	}
	return m
}
