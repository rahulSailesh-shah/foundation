package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type httpStatus interface {
	HTTPStatus() int
}

type NoResponse struct{}

func NewNoResponse() NoResponse {
	return NoResponse{}
}

func (n NoResponse) HTTPStatus() int {
	return http.StatusNoContent
}

func (n NoResponse) Encode() ([]byte, string, error) {
	return nil, "", nil
}

func Respond(ctx context.Context, w http.ResponseWriter, resp Encoder) error {
	if _, ok := resp.(NoResponse); ok {
		return nil
	}

	if err := ctx.Err(); err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.New("client disconnected, do not send response")
		}
	}

	// TODO: Set tracing span for response

	statusCode := http.StatusOK
	switch v := resp.(type) {
	case httpStatus:
		statusCode = v.HTTPStatus()
	case error:
		statusCode = http.StatusInternalServerError
	default:
		if resp == nil {
			statusCode = http.StatusNoContent
		}
	}

	data, contentType, err := resp.Encode()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("respond: encode: %w", err)
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)

	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("respond: write: %w", err)
	}

	return nil
}
