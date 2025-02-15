// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

// URL Things.
const (
	BASE_URL       = "https://open.larksuite.com/open-apis"
	AUTH_API       = "/auth/v3/tenant_access_token/internal"
	GROUP_CHAT_API = "/im/v1/chats"
	USERGROUP_API  = "/contact/v3/group"
	USER_API       = "/contact/v3/users"
)

// HTTP Call Helpers.
const (
	BASE_RETRY_COUNT = 3
	BASE_DELAY       = 10 //
)

type HTTPMethod string

// HTTP Method.
const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PATCH  HTTPMethod = "PATCH"
	DELETE HTTPMethod = "DELETE"
	PUT    HTTPMethod = "PUT"
)

type AuthorizationHeader string

// Authorization Header.
const (
	TENANT_ACCESS_TOKEN AuthorizationHeader = "tenant_access_token"
	APP_ACCESS_TOKEN    AuthorizationHeader = "app_access_token"
)

type UserIDType string

// User ID Type.
const (
	USER_ID  UserIDType = "user_id"
	OPEN_ID  UserIDType = "open_id"
	UNION_ID UserIDType = "union_id"
)
