package httputil

import (
	"app/log"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func RequestSimple(method, url string, body string, header map[string]string) (string, error) {
	return RequestBase(method, url, body, header, nil, 0)
}

func Request(method, url string, body string, header map[string]string, transport *http.Transport) (string, error) {
	return RequestBase(method, url, body, header, transport, 0)
}

func RequestBase(method, url string, body string, header map[string]string, transport *http.Transport, timeout time.Duration) (string, error) {
	log.Infof("[%s] %s request body: %s", method, url, body)

	client := &http.Client{}

	if transport != nil {
		client.Transport = transport
	}

	ctx := context.Background()
	var cancel context.CancelFunc
	if timeout != 0 {
		client.Timeout = timeout
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Warn("request failed: ", err)
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("read response body failed: ", err)
		return "", err
	}
	if len(bodyBytes) > 4096 {
		log.Infof("response body (truncated): %s...", string(bodyBytes[:1024]))
	} else {
		log.Infof("response body: %s", string(bodyBytes))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("status code failed：%s", resp.Status)
	}
	return string(bodyBytes), nil
}
