// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"
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
		name string
	}{
		{
			name: "success is connection error",
		},
	}

	for _, tc := range tests {
		PatchConvey(tc.name, t, func() {
			err := isConnectionError(fmt.Errorf("connection refused"))
			So(err, ShouldBeTrue)
		})
	}
}
