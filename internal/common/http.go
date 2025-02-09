// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ACCESS TOKEN API.
// https://open.larksuite.com/document/server-docs/getting-started/api-access-token/auth-v3/tenant_access_token_internal.
func GetAccessTokenAPI(appID, appSecret string) (string, string, error) {
	tflog.Info(context.Background(), "Getting access token from Lark API")

	client := NewLarkClient("", "", BASE_DELAY)

	requestBody := AccessTokenRequest{
		AppID:     appID,
		AppSecret: appSecret,
	}

	var response AccessTokenResponse

	err := client.DoTenantRequest(context.Background(), POST, AUTH_API, requestBody, &response)
	if err != nil {
		return "", "", fmt.Errorf("failed to get access token: %w", err)
	}

	if response.Code != 0 {
		return "", "", fmt.Errorf("failed to get access token: %s", response.Msg)
	}

	tflog.Info(context.Background(), "Access token retrieved successfully")

	return response.TenantAccessToken, response.AppAccessToken, nil
}
