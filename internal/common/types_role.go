// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

type RoleRequest struct {
	RoleName string `json:"role_name"`
}

type DataRoleCreateResponse struct {
	RoleID string `json:"role_id"`
}

type RoleCreateResponse struct {
	BaseResponse
	Data DataRoleCreateResponse `json:"data"`
}

type FunctionalRoleMemberResult struct {
	UserID string `json:"user_id"`
	Reason int    `json:"reason"`
}

type RoleMemberCreateRequest struct {
	Members []string `json:"members"`
}

type DataRoleMemberCreateDeleteResponse struct {
	Results []FunctionalRoleMemberResult `json:"results"`
}

type RoleMemberCreateResponse struct {
	BaseResponse
	Data DataRoleMemberCreateDeleteResponse `json:"data"`
}

type RoleMemberUpdateRequest struct {
	Members     []string `json:"members"`
	Departments []string `json:"departments"`
}

type RoleMember struct {
	UserID        string   `json:"user_id"`
	ScopeID       string   `json:"scope_id"`
	DepartmentIDs []string `json:"department_ids"`
}

type DataRoleMemberGetResponse struct {
	Members   []RoleMember `json:"members"`
	PageToken string       `json:"page_token"`
	HasMore   bool         `json:"has_more"`
}

type RoleMemberGetResponse struct {
	BaseResponse
	Data DataRoleMemberGetResponse `json:"data"`
}

type RoleMemberDeleteRequest struct {
	Members []string `json:"members"`
}

type RoleMemberDeleteResponse struct {
	BaseResponse
	Data DataRoleMemberCreateDeleteResponse `json:"data"`
}
