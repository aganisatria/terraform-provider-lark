// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

// URL Things.
const (
	BASE_URL                 = "https://open.larksuite.com/open-apis"
	AUTH_API                 = "/auth/v3/tenant_access_token/internal"
	DEPARTMENT_API           = "/contact/v3/departments"
	GROUP_CHAT_API           = "/im/v1/chats"
	USERGROUP_API            = "/contact/v3/group"
	USER_API                 = "/contact/v3/users"
	ROLE_API                 = "/contact/v3/functional_roles"
	UNIT_API                 = "/contact/v3/units"
	EXPLORER_ROOT_FOLDER_API = "/drive/explorer/v2/root_folder"
	EXPLORER_FOLDER_API      = "/drive/explorer/v2/folder"
	DOCS_FILE_API            = "/drive/v1/files"
	WORKFORCE_TYPE_API       = "/contact/v3/employee_type_enums"
)

// HTTP Call Helpers.
const (
	BASE_RETRY_COUNT = 2
	BASE_DELAY       = 1
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

type TerraformType string

// Terraform Type.
const (
	RESOURCE           TerraformType = "resource"
	DATA_SOURCE        TerraformType = "data_source"
	EPHEMERAL_RESOURCE TerraformType = "ephemeral_resource"
	FUNCTION           TerraformType = "function"
)

type TerraformName string

// Terraform Name.
const (
	DEPARTMENT        TerraformName = "department"
	GROUP_CHAT        TerraformName = "group_chat"
	GROUP_CHAT_MEMBER TerraformName = "group_chat_member"
	ROLE              TerraformName = "role"
	ROLE_MEMBER       TerraformName = "role_member"
	USER_GROUP        TerraformName = "user_group"
	USER_GROUP_MEMBER TerraformName = "user_group_member"
	USER_BY_EMAIL     TerraformName = "user_by_email"
	USER_BY_ID        TerraformName = "user_by_id"
	WORKFORCE_TYPE    TerraformName = "workforce_type"
)

type DepartmentIDType string

// Department ID Type.
const (
	DEPARTMENT_ID      DepartmentIDType = "department_id"
	OPEN_DEPARTMENT_ID DepartmentIDType = "open_department_id"
)
