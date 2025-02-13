// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"reflect"

	. "github.com/bytedance/mockey"
)

// SetupDoTenantRequest helper that arrange mock for DoTenantRequest.
// return cleanup function to unpatch the mock.
// mockError: error that want to simulate.
// mockResponse: response that want to assign to response variable.
func SetupDoTenantRequest(mockError error, mockResponse interface{}) func() {
	patch := Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
		if mockError != nil {
			return mockError
		}
		reflect.ValueOf(resp).Elem().Set(reflect.ValueOf(mockResponse))
		return nil
	}).Build()
	return func() {
		patch.UnPatch()
	}
}
