package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func readBody[T any](r *http.Request, input *T) error {
	body, err := io.ReadAll(r.Body)
	r.Body.Close()

	if err != nil {
		return fmt.Errorf("can't read body: %w", err)
	}

	err = json.Unmarshal(body, input)
	if err != nil {
		return fmt.Errorf("can't unmarshal body: %w", err)
	}
	return nil
}
