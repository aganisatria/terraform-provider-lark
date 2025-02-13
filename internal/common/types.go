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
	AppID     string `json:"app_id,omitempty"`
	AppSecret string `json:"app_secret,omitempty"`
}

// Access Token Response.
type AccessTokenResponse struct {
	BaseResponse
	AppAccessToken    string `json:"app_access_token"`
	TenantAccessToken string `json:"tenant_access_token"`
	Expire            int    `json:"expire"`
}

// I18nName is the internationalized name of the group chat.
type I18nName struct {
	ZhCn string `json:"zh_cn,omitempty"`
	JaJp string `json:"ja_jp,omitempty"`
	EnUs string `json:"en_us,omitempty"`
}

type TypeOfUsersInAListValidator string

const (
	TypeOfUsersInAListValidatorUser TypeOfUsersInAListValidator = "USER_ID_IN_GROUP_CHAT"
)
