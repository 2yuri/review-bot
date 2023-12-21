package http

import (
	"bytes"
	"encoding/json"
	"errors"
	l "github.com/2yuri/review-bot/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

var timeout = 5 * time.Second
var client = &http.Client{
	Timeout: timeout,
}

func Get(url string, headers map[string]string, dest interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("invalid status code")
	}

	if err := json.NewDecoder(resp.Body).Decode(&dest); err != nil {
		return err
	}

	return nil
}

func Post(url string, headers map[string]string, body interface{}, dest interface{}) error {
	b := new(bytes.Buffer)

	if body != nil {
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(http.MethodPost, url, b)
	if err != nil {
		return err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		l.Logger.Debug("Response", zap.Any("resp", string(b)))

		return errors.New("invalid status code")
	}

	if dest == nil {
		if err := json.NewDecoder(resp.Body).Decode(&dest); err != nil {
			return err
		}
	}

	return nil
}
