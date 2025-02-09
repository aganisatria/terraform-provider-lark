package common

import "time"

// URL Things.
const (
	BASE_URL = "https://open.larksuite.com/open-apis"
	AUTH_API = BASE_URL + "/auth/v3/tenant_access_token/internal"
)

// HTTP Call Helpers.
const (
	MAX_RETRIES = 3
	BASE_DELAY  = 10 * time.Second
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
	APP_ACCESS_TOKEN   AuthorizationHeader = "app_access_token"
)
