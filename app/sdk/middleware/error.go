package middleware

import (
	"context"
	"errors"
	"net/http"
	"path"

	"foundation/app/sdk/errs"
	"foundation/logger"
	"foundation/web"
)

func Errors(log *logger.Logger) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		f := func(ctx context.Context, r *http.Request) web.Encoder {
			resp := next(ctx, r)

			// Check if response is an error
			err := checkIsError(resp)
			if err == nil {
				return resp
			}

			// Convert error to app error
			var appErr *errs.Error
			if !errors.As(err, &appErr) {
				appErr = errs.New(errs.Internal, err)
			}

			// Log error
			log.Error(ctx, "handled error",
				"err", err,
				"source_err_file", path.Base(appErr.FileName),
				"source_err_func", path.Base(appErr.FuncName))

			// Return app error to web framework to send as response
			return appErr
		}
		return f
	}
	return m
}
