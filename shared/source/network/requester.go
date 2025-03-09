package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"social/shared/models"
)

func MakeRequest[IN any, OUT any](ctx context.Context, method string, url string, input IN, out *OUT) (err error) {
	inBody, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("can't marshal input: %w", err)
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(inBody))
	if err != nil {
		return fmt.Errorf("can't create request: %w", err)
	}

	// TODO add request id from context to headers

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("can't do request: %w", err)
	}

	sdBody, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return fmt.Errorf("can't read Body request: %w", err)
	}
	slog.Debug("sdBody", "body", sdBody)

	sd := models.StatusData{}
	err = json.Unmarshal(sdBody, &sd)
	if err != nil {
		return fmt.Errorf("can't unmarshal body: %w", err)
	}

	*out, err = models.ParseStatusData[OUT](sd)
	return err
}

func ReadBody[T any](r *http.Request, input *T) error {
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
