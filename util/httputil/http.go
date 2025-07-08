package httputil

import (
	"app/log"
	"app/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func RequestSimple(ctx context.Context, method, url string, body any, header map[string]string) (string, error) {
	return RequestBase(ctx, method, url, body, header, nil, 0)
}

func Request(ctx context.Context, method, url string, body any, header map[string]string, transport *http.Transport) (string, error) {
	return RequestBase(ctx, method, url, body, header, transport, 0)
}

func RequestBase(ctx context.Context, method, url string, body any, header map[string]string, transport *http.Transport, timeout time.Duration) (string, error) {
	log.T(ctx).Infof("[%s] %s", method, url)
	var bodyByte []byte
	var err error
	if body != nil {
		bodyByte, err = json.Marshal(body)
		if err != nil {
			return "", err
		}
		log.T(ctx).Infof("request body: %s", string(bodyByte))
	}

	client := &http.Client{}
	if transport != nil {
		client.Transport = transport
	}

	ctx2 := context.Background()
	var cancel context.CancelFunc
	if timeout != 0 {
		client.Timeout = timeout
		ctx2, cancel = context.WithTimeout(context.Background(), timeout)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx2, method, url, bytes.NewBuffer(bodyByte))
	if err != nil {
		return "", err
	}

	if header != nil {
		log.T(ctx).Infof("request header: %s", util.ToJson(header))
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.T(ctx).Warn("request failed: ", err)
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.T(ctx).Warn("read response body failed: ", err)
		return "", err
	}
	if len(bodyBytes) > 4096 {
		log.T(ctx).Infof("response status: %s, response body (truncated): %s...", resp.Status, string(bodyBytes[:1024]))
	} else {
		log.T(ctx).Infof("response status: %s, response body: %s", resp.Status, string(bodyBytes))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return "", fmt.Errorf("status code failed: %s", resp.Status)
	}
	return string(bodyBytes), nil
}

func IsNetworkFailed(err error) bool {
	if err != nil {
		if strings.HasPrefix(err.Error(), "request failed: ") {
			return true
		}
	}
	return false
}

func IsStatusFailed(err error) bool {
	if err != nil {
		if strings.HasPrefix(err.Error(), "status code failed: ") {
			return true
		}
	}
	return false
}
