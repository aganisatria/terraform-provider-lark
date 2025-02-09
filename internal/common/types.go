// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import "net/http"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Base Response that all Lark API responses should implement.
type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Access Token Request.
type AccessTokenRequest struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

// Access Token Response.
type AccessTokenResponse struct {
	BaseResponse
	AppAccessToken    string `json:"app_access_token"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}
