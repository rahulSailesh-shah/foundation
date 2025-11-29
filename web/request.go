package web

import (
	"fmt"
	"io"
	"net/http"
)

type Decoder interface {
	Decode(data []byte) error
}

type Validator interface {
	Validate() error
}

func Decode(r *http.Request, v Decoder) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := v.Decode(body); err != nil {
		return fmt.Errorf("request: decode: %w", err)
	}

	if v, ok := v.(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
