package middleware

import "foundation/web"

func checkIsError(e web.Encoder) error {
	err, hasError := e.(error)
	if hasError {
		return err
	}
	return nil
}
