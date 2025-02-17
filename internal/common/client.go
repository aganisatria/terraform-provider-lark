// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LarkClient struct {
	httpClient        HTTPClient
	TenantAccessToken string
	AppAccessToken    string
	BaseDelay         time.Duration
	RetryCount        int
	AppID             string
}

func NewLarkClient(tenantAccessToken, appAccessToken, appID string, baseDelay int, retryCount int) *LarkClient {
	return &LarkClient{
		httpClient:        &http.Client{},
		TenantAccessToken: tenantAccessToken,
		AppAccessToken:    appAccessToken,
		BaseDelay:         time.Duration(baseDelay) * time.Second,
		RetryCount:        retryCount,
		AppID:             appID,
	}
}

func (c *LarkClient) DoInitializeRequest(
	ctx context.Context,
	method HTTPMethod,
	path string,
	requestBody interface{},
	response interface{},
) error {
	return c.DoRequest(ctx, method, path, requestBody, response, "")
}

func (c *LarkClient) DoTenantRequest(
	ctx context.Context,
	method HTTPMethod,
	path string,
	requestBody interface{},
	response interface{},
) error {
	return c.DoRequest(ctx, method, path, requestBody, response, TENANT_ACCESS_TOKEN)
}

func (c *LarkClient) DoAppRequest(
	ctx context.Context,
	method HTTPMethod,
	path string,
	requestBody interface{},
	response interface{},
) error {
	return c.DoRequest(ctx, method, path, requestBody, response, APP_ACCESS_TOKEN)
}

func (c *LarkClient) DoRequest(
	ctx context.Context,
	method HTTPMethod,
	path string,
	requestBody interface{},
	response interface{},
	authorizationHeader AuthorizationHeader,
) error {
	var lastErr error

	for attempt := 0; attempt < c.RetryCount; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			delay := c.BaseDelay * time.Duration(1<<uint(attempt-1))
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := c.doSingleRequest(ctx, method, path, requestBody, response, authorizationHeader)
		if err == nil {
			return nil
		}

		// Only retry for connection errors
		if isConnectionError(err) {
			lastErr = err
			continue
		}

		// For other errors, return immediately
		return err
	}

	return fmt.Errorf("failed after %d retries. Last error: %w", c.RetryCount, lastErr)
}

func (c *LarkClient) doSingleRequest(
	ctx context.Context,
	method HTTPMethod,
	path string,
	requestBody interface{},
	response interface{},
	authorizationHeader AuthorizationHeader,
) error {
	url := fmt.Sprintf("%s%s", BASE_URL, path)

	var bodyReader io.Reader
	if requestBody != nil {
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("error marshaling request: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, string(method), url, bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	if authorizationHeader != "" {
		var token string
		if authorizationHeader == APP_ACCESS_TOKEN {
			token = c.AppAccessToken
		} else if authorizationHeader == TENANT_ACCESS_TOKEN {
			token = c.TenantAccessToken
		} else {
			return fmt.Errorf("invalid authorization header: %s", authorizationHeader)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp BaseResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("error response with status code %d", resp.StatusCode)
		}
		return fmt.Errorf("API error: code=%d, message=%s", errResp.Code, errResp.Msg)
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
	}

	return nil
}

// isConnectionError is a helper function to check if the error is a connection error.
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	return Contains(errStr,
		"connection refused",
		"no such host",
		"timeout",
		"connection reset",
		"EOF",
		"broken pipe",
		"TLS handshake timeout",
	)
}
