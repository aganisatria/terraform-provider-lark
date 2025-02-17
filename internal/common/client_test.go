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
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewLarkClient(t *testing.T) {
	tests := []struct {
		name string
		want *LarkClient
	}{
		{
			name: "success create new client",
			want: &LarkClient{
				httpClient:        &http.Client{},
				TenantAccessToken: "tenant-token",
				AppAccessToken:    "app-token",
			},
		},
	}

	for _, tc := range tests {
		PatchConvey(tc.name, t, func() {
			got := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			So(got.TenantAccessToken, ShouldEqual, tc.want.TenantAccessToken)
			So(got.AppAccessToken, ShouldEqual, tc.want.AppAccessToken)
		})
	}
}

func TestLarkClient_DoInitializeRequest(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "success do initialize request",
		},
	}

	for _, tc := range tests {
		PatchConvey(tc.name, t, func() {
			Mock((*LarkClient).DoRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, requestBody interface{}, response interface{}, authorizationHeader AuthorizationHeader) error {
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			err := client.DoInitializeRequest(context.Background(), GET, "/test", nil, nil)
			if err != nil {
				t.Errorf("DoInitializeRequest = %v", err)
			}
		})
	}
}

func TestLarkClient_DoTenantRequest(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "success do tenant request",
		},
	}

	for _, tc := range tests {
		PatchConvey(tc.name, t, func() {
			Mock((*LarkClient).DoRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, requestBody interface{}, response interface{}, authorizationHeader AuthorizationHeader) error {
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			err := client.DoTenantRequest(context.Background(), GET, "/test", nil, nil)
			if err != nil {
				t.Errorf("DoTenantRequest = %v", err)
			}
		})
	}
}

func TestLarkClient_DoAppRequest(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "success do app request",
		},
	}

	for _, tc := range tests {
		PatchConvey(tc.name, t, func() {
			Mock((*LarkClient).DoRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, requestBody interface{}, response interface{}, authorizationHeader AuthorizationHeader) error {
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			err := client.DoAppRequest(context.Background(), GET, "/test", nil, nil)
			if err != nil {
				t.Errorf("DoAppRequest = %v", err)
			}
		})
	}
}

func TestLarkClient_IsConnectionError(t *testing.T) {
	tests := []struct {
		name   string
		params error
		want   bool
	}{
		{
			name:   "success is connection error",
			params: fmt.Errorf("connection refused"),
			want:   true,
		},
		{
			name:   "success is not connection error",
			params: nil,
			want:   false,
		},
	}

	for _, tc := range tests {
		PatchConvey(tc.name, t, func() {
			err := isConnectionError(tc.params)
			So(err, ShouldEqual, tc.want)
		})
	}
}

func TestLarkClient_doSingleRequest(t *testing.T) {
	tests := []struct {
		name                string
		requestBody         interface{}
		authorizationHeader AuthorizationHeader
		mockFn              func() []*MockBuilder
		expectedError       error
		wantErr             bool
	}{
		{
			name: "error marshal request",
			requestBody: map[string]string{
				"key": "value",
			},
			authorizationHeader: APP_ACCESS_TOKEN,
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(json.Marshal).To(func(v any) ([]byte, error) {
						return nil, fmt.Errorf("marshal failed")
					}),
				}
			},
			expectedError: fmt.Errorf("error marshaling request: marshal failed"),
			wantErr:       true,
		},
		{
			name: "error NewRequestWithContext",
			requestBody: map[string]string{
				"key": "value",
			},
			authorizationHeader: APP_ACCESS_TOKEN,
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(json.Marshal).To(func(v any) ([]byte, error) {
						return []byte(`{"key": "value"}`), nil
					}),
					Mock(http.NewRequestWithContext).To(func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
						return nil, fmt.Errorf("no such host")
					}),
				}
			},
			expectedError: fmt.Errorf("error creating request: no such host"),
			wantErr:       true,
		},
		{
			name: "invalid authorization header",
			requestBody: map[string]string{
				"key": "value",
			},
			authorizationHeader: "invalid",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(json.Marshal).To(func(v any) ([]byte, error) {
						return []byte(`{"key": "value"}`), nil
					}),
					Mock(http.NewRequestWithContext).To(func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
						return &http.Request{
							Header: http.Header{
								"Authorization": []string{"Bearer tenant-token"},
							},
						}, nil
					}),
				}
			},
			expectedError: fmt.Errorf("invalid authorization header: invalid"),
			wantErr:       true,
		},
		{
			name: "error http client do",
			requestBody: map[string]string{
				"key": "value",
			},
			authorizationHeader: APP_ACCESS_TOKEN,
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(json.Marshal).To(func(v any) ([]byte, error) {
						return []byte(`{"key": "value"}`), nil
					}),
					Mock(http.NewRequestWithContext).To(func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
						return &http.Request{
							Header: http.Header{
								"Authorization": []string{"Bearer tenant-token"},
							},
						}, nil
					}),
					Mock((*http.Client).Do).To(func(req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("do failed")
					}),
				}
			},
			expectedError: fmt.Errorf("error executing request: do failed"),
			wantErr:       true,
		},
		{
			name: "error response with status code",
			requestBody: map[string]string{
				"key": "value",
			},
			authorizationHeader: TENANT_ACCESS_TOKEN,
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(json.Marshal).To(func(v any) ([]byte, error) {
						return []byte(`{"key": "value"}`), nil
					}),
					Mock(http.NewRequestWithContext).To(func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
						return &http.Request{
							Header: http.Header{
								"Authorization": []string{"Bearer tenant-token"},
							},
						}, nil
					}),
					Mock((*http.Client).Do).To(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 400,
							Body:       io.NopCloser(bytes.NewBufferString(`{"code": 400, "msg": "error"}`)),
						}, nil
					}),
					Mock(json.NewDecoder).To(func(r io.Reader) *json.Decoder {
						return &json.Decoder{}
					}),
					Mock((*json.Decoder).Decode).To(func(dec *json.Decoder, v any) error {
						return fmt.Errorf("decode failed")
					}),
				}
			},
			expectedError: fmt.Errorf("error response with status code 400"),
			wantErr:       true,
		},
		{
			name: "error API error",
			requestBody: map[string]string{
				"key": "value",
			},
			authorizationHeader: TENANT_ACCESS_TOKEN,
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(json.Marshal).To(func(v any) ([]byte, error) {
						return []byte(`{"key": "value"}`), nil
					}),
					Mock(http.NewRequestWithContext).To(func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
						return &http.Request{
							Header: http.Header{
								"Authorization": []string{"Bearer tenant-token"},
							},
						}, nil
					}),
					Mock((*http.Client).Do).To(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 400,
							Body:       io.NopCloser(bytes.NewBufferString(`{"code": 400, "msg": "error"}`)),
						}, nil
					}),
				}
			},
			expectedError: fmt.Errorf("API error: code=400, message=error"),
			wantErr:       true,
		},
		{
			name: "success",
			requestBody: map[string]string{
				"key": "value",
			},
			authorizationHeader: TENANT_ACCESS_TOKEN,
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(json.Marshal).To(func(v any) ([]byte, error) {
						return []byte(`{"key": "value"}`), nil
					}),
					Mock(http.NewRequestWithContext).To(func(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
						return &http.Request{
							Header: http.Header{
								"Authorization": []string{"Bearer tenant-token"},
							},
						}, nil
					}),
					Mock((*http.Client).Do).To(func(req *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode: 200,
							Body:       io.NopCloser(bytes.NewBufferString(`{"code": 400, "msg": "error"}`)),
						}, nil
					}),
				}
			},
			expectedError: nil,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			err := client.doSingleRequest(context.Background(), GET, "/test", tt.requestBody, nil, tt.authorizationHeader)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestLarkClient_DoRequest(t *testing.T) {
	tests := []struct {
		name          string
		mockFn        func() []*MockBuilder
		expectedError error
		wantErr       bool
	}{
		{
			name: "error",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).doSingleRequest).To(func(ctx context.Context, method HTTPMethod, path string, requestBody interface{}, response interface{}, authorizationHeader AuthorizationHeader) error {
						return fmt.Errorf("do single request failed")
					}),
				}
			},
			expectedError: fmt.Errorf("do single request failed"),
			wantErr:       true,
		},
		{
			name: "error failed after %d retries. Last error: %w",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).doSingleRequest).To(func(ctx context.Context, method HTTPMethod, path string, requestBody interface{}, response interface{}, authorizationHeader AuthorizationHeader) error {
						return fmt.Errorf("connection error")
					}),
					Mock(isConnectionError).To(func(err error) bool {
						return true
					}),
				}
			},
			expectedError: fmt.Errorf("failed after %d retries. Last error: %w", BASE_RETRY_COUNT, fmt.Errorf("connection error")),
			wantErr:       true,
		},
		{
			name: "success",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).doSingleRequest).To(func(ctx context.Context, method HTTPMethod, path string, requestBody interface{}, response interface{}, authorizationHeader AuthorizationHeader) error {
						return nil
					}),
				}
			},
			expectedError: nil,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			err := client.DoRequest(context.Background(), GET, "/test", nil, nil, APP_ACCESS_TOKEN)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}
