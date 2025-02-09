package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

type testResponse struct {
	Message string `json:"message"`
}

type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewLarkClient(t *testing.T) {
	cases := []struct {
		name              string
		tenantAccessToken string
		appAccessToken    string
		want             *LarkClient
	}{
		{
			name:              "success create new client",
			tenantAccessToken: "tenant-token",
			appAccessToken:    "app-token",
			want: &LarkClient{
				httpClient:        &http.Client{},
				TenantAccessToken: "tenant-token",
				AppAccessToken:    "app-token",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NewLarkClient(tc.tenantAccessToken, tc.appAccessToken, BASE_DELAY)
			if got.TenantAccessToken != tc.want.TenantAccessToken {
				t.Errorf("TenantAccessToken = %v, want %v", got.TenantAccessToken, tc.want.TenantAccessToken)
			}
			if got.AppAccessToken != tc.want.AppAccessToken {
				t.Errorf("AppAccessToken = %v, want %v", got.AppAccessToken, tc.want.AppAccessToken)
			}
		})
	}
}

func TestLarkClient_DoRequest(t *testing.T) {
	cases := []struct {
		name              string
		method            HTTPMethod
		path              string
		requestBody       interface{}
		mockResponseCode  int
		mockResponseBody  interface{}
		authHeader        AuthorizationHeader
		expectedError     bool
		expectedResponse  *testResponse
		mockHTTPError     error
		cancelContext     bool
	}{
		{
			name:              "successful request",
			method:           GET,
			path:             "/test",
			mockResponseCode: 200,
			mockResponseBody: testResponse{Message: "success"},
			authHeader:       TENANT_ACCESS_TOKEN,
			expectedResponse: &testResponse{Message: "success"},
		},
		{
			name:              "error response",
			method:           POST,
			path:             "/test",
			mockResponseCode: 400,
			mockResponseBody: map[string]interface{}{
				"code": 400,
				"msg":  "bad request",
			},
			authHeader:       APP_ACCESS_TOKEN,
			expectedError:    true,
		},
		{
			name:           "connection error",
			method:        GET,
			path:          "/test",
			mockHTTPError: fmt.Errorf("connection refused"),
			authHeader:    TENANT_ACCESS_TOKEN,
			expectedError: true,
		},
		{
			name:          "context cancelled",
			method:        GET,
			path:          "/test",
			authHeader:    TENANT_ACCESS_TOKEN,
			cancelContext: true,
			expectedError: true,
		},
		{
			name:           "invalid request body",
			method:        POST,
			path:          "/test",
			requestBody:   make(chan int),
			authHeader:    TENANT_ACCESS_TOKEN,
			expectedError: true,
		},
		{
			name:              "invalid response body",
			method:           GET,
			path:             "/test",
			mockResponseCode: 200,
			mockResponseBody: "invalid json}",
			authHeader:       TENANT_ACCESS_TOKEN,
			expectedError:    true,
		},
		{
			name:           "invalid authorization header",
			method:        GET,
			path:          "/test",
			authHeader:    "invalid-header",
			expectedError: true,
		},
		{
			name:              "error response with invalid json",
			method:           GET,
			path:             "/test",
			mockResponseCode: 400,
			mockResponseBody: "invalid error json}",
			authHeader:       TENANT_ACCESS_TOKEN,
			expectedError:    true,
		},
		{
			name:           "retry exhausted",
			method:        GET,
			path:          "/test",
			mockHTTPError: fmt.Errorf("connection refused"),
			authHeader:    TENANT_ACCESS_TOKEN,
			expectedError: true,
		},
		{
			name:           "invalid URL",
			method:        GET,
			path:          string([]byte{0x7f}), // Invalid URL character
			authHeader:    TENANT_ACCESS_TOKEN,
			expectedError: true,
		},
		{
			name:           "empty authorization header",
			method:        GET,
			path:          "/test",
			authHeader:    "",
			mockResponseCode: 200,
			mockResponseBody: testResponse{Message: "success"},
			expectedResponse: &testResponse{Message: "success"},
		},
		{
			name:           "nil response",
			method:        GET,
			path:          "/test",
			authHeader:    TENANT_ACCESS_TOKEN,
			mockResponseCode: 200,
			mockResponseBody: testResponse{Message: "success"},
		},
		{
			name:           "empty auth header with nil response",
			method:        GET,
			path:          "/test",
			authHeader:    "",
			mockResponseCode: 200,
			mockResponseBody: map[string]interface{}{
				"message": "success",
			},
			expectedError: false,
		},
		{
			name:           "nil response with success status",
			method:        GET,
			path:          "/test",
			authHeader:    TENANT_ACCESS_TOKEN,
			mockResponseCode: 200,
			mockResponseBody: nil,
			expectedError: false,
		},
		{
			name:           "request with body",
			method:        POST,
			path:          "/test",
			requestBody:   map[string]interface{}{
				"key": "value",
				"number": 123,
				"nested": map[string]string{
					"inner": "data",
				},
			},
			mockResponseCode: 200,
			mockResponseBody: testResponse{Message: "success"},
			authHeader:    TENANT_ACCESS_TOKEN,
			expectedResponse: &testResponse{Message: "success"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if tc.cancelContext {
						return nil, context.Canceled
					}

					if tc.mockHTTPError != nil {
						return nil, tc.mockHTTPError
					}

					if string(tc.method) != req.Method {
						t.Errorf("expected method %s, got %s", tc.method, req.Method)
					}

					if !strings.HasSuffix(req.URL.Path, tc.path) {
						t.Errorf("expected path to end with %s, got %s", tc.path, req.URL.Path)
					}

					if tc.authHeader != "" {
						authHeader := req.Header.Get("Authorization")
						var expectedToken string
						if tc.authHeader == APP_ACCESS_TOKEN {
							expectedToken = "Bearer test-app-token"
						} else if tc.authHeader == TENANT_ACCESS_TOKEN {
							expectedToken = "Bearer test-tenant-token"
						}
						if authHeader != expectedToken {
							t.Errorf("expected Authorization header %s, got %s", expectedToken, authHeader)
						}
					} else {
						if auth := req.Header.Get("Authorization"); auth != "" {
							t.Errorf("expected no Authorization header, got %s", auth)
						}
					}

					if tc.requestBody != nil {
						body, err := io.ReadAll(req.Body)
						if err != nil {
							t.Fatalf("failed to read request body: %v", err)
						}
						req.Body = io.NopCloser(bytes.NewBuffer(body))

						expectedBody, err := json.Marshal(tc.requestBody)
						if err != nil {
							t.Fatalf("failed to marshal expected body: %v", err)
						}

						var got, want interface{}
						if err := json.Unmarshal(body, &got); err != nil {
							t.Fatalf("failed to unmarshal actual body: %v", err)
						}
						if err := json.Unmarshal(expectedBody, &want); err != nil {
							t.Fatalf("failed to unmarshal expected body: %v", err)
						}

						if !reflect.DeepEqual(got, want) {
							t.Errorf("request body = %v, want %v", got, want)
						}
					}

					responseBody, err := json.Marshal(tc.mockResponseBody)
					if err != nil {
						t.Fatalf("failed to marshal mock response: %v", err)
					}

					return &http.Response{
						StatusCode: tc.mockResponseCode,
						Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
						Header:     make(http.Header),
					}, nil
				},
			}

			client := &LarkClient{
				httpClient:        mockClient,
				TenantAccessToken: "test-tenant-token",
				AppAccessToken:    "test-app-token",
				BaseDelay:         1 * time.Millisecond,
			}

			ctx := context.Background()
			if tc.cancelContext {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			var response testResponse
			err := client.DoRequest(ctx, tc.method, tc.path, tc.requestBody, &response, tc.authHeader)

			if tc.expectedError && err == nil {
				t.Error("expected error but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tc.expectedResponse != nil && response.Message != tc.expectedResponse.Message {
				t.Errorf("response = %v, want %v", response, tc.expectedResponse)
			}
		})
	}
}

func TestIsConnectionError(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "connection refused",
			err:      fmt.Errorf("connection refused"),
			expected: true,
		},
		{
			name:     "timeout error",
			err:      fmt.Errorf("timeout occurred"),
			expected: true,
		},
		{
			name:     "non-connection error",
			err:      fmt.Errorf("some other error"),
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := isConnectionError(tc.err)
			if got != tc.expected {
				t.Errorf("isConnectionError() = %v, want %v", got, tc.expected)
			}
		})
	}
}

func TestLarkClient_DoTenantRequest(t *testing.T) {
	mockClient := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if auth := req.Header.Get("Authorization"); auth != "Bearer test-tenant-token" {
				t.Errorf("expected tenant token, got %s", auth)
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"success"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	client := &LarkClient{
		httpClient:        mockClient,
		TenantAccessToken: "test-tenant-token",
		AppAccessToken:    "test-app-token",
		BaseDelay:         1 * time.Millisecond,
	}

	var response struct{ Message string }
	err := client.DoTenantRequest(
		context.Background(),
		GET,
		"/test",
		nil,
		&response,
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if response.Message != "success" {
		t.Errorf("expected success message, got %s", response.Message)
	}
}

func TestLarkClient_DoAppRequest(t *testing.T) {
	mockClient := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if auth := req.Header.Get("Authorization"); auth != "Bearer test-app-token" {
				t.Errorf("expected app token, got %s", auth)
			}
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"success"}`)),
				Header:     make(http.Header),
			}, nil
		},
	}

	client := &LarkClient{
		httpClient:        mockClient,
		TenantAccessToken: "test-tenant-token",
		AppAccessToken:    "test-app-token",
		BaseDelay:         1 * time.Millisecond,
	}

	var response struct{ Message string }
	err := client.DoAppRequest(
		context.Background(),
		GET,
		"/test",
		nil,
		&response,
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if response.Message != "success" {
		t.Errorf("expected success message, got %s", response.Message)
	}
}

func TestNewLarkClient_WithVariousParams(t *testing.T) {
	cases := []struct {
		name              string
		tenantAccessToken string
		appAccessToken    string
		baseDelay         time.Duration
	}{
		{
			name:              "with empty tokens",
			tenantAccessToken: "",
			appAccessToken:    "",
			baseDelay:         time.Second,
		},
		{
			name:              "with zero delay",
			tenantAccessToken: "tenant",
			appAccessToken:    "app",
			baseDelay:         0,
		},
		{
			name:              "with all valid params",
			tenantAccessToken: "tenant",
			appAccessToken:    "app",
			baseDelay:         time.Second,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewLarkClient(tc.tenantAccessToken, tc.appAccessToken, tc.baseDelay)
			if client.TenantAccessToken != tc.tenantAccessToken {
				t.Errorf("TenantAccessToken = %v, want %v", client.TenantAccessToken, tc.tenantAccessToken)
			}
			if client.AppAccessToken != tc.appAccessToken {
				t.Errorf("AppAccessToken = %v, want %v", client.AppAccessToken, tc.appAccessToken)
			}
			if client.BaseDelay != tc.baseDelay {
				t.Errorf("BaseDelay = %v, want %v", client.BaseDelay, tc.baseDelay)
			}
			if client.httpClient == nil {
				t.Error("httpClient should not be nil")
			}
		})
	}
}

func TestLarkClient_DoRequest_RetryWithContext(t *testing.T) {
	attemptCount := 0
	mockClient := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			attemptCount++
			time.Sleep(2 * time.Millisecond)
			return nil, fmt.Errorf("connection refused")
		},
	}

	client := &LarkClient{
		httpClient:        mockClient,
		TenantAccessToken: "test-tenant-token",
		AppAccessToken:    "test-app-token",
		BaseDelay:         1 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()

	var response testResponse
	err := client.DoRequest(ctx, GET, "/test", nil, &response, TENANT_ACCESS_TOKEN)

	if err == nil {
		t.Error("expected error but got none")
	}
	if !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Errorf("expected deadline exceeded error, got: %v", err)
	}
	if attemptCount == 0 {
		t.Error("expected at least one attempt")
	}
}

