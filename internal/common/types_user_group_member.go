// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

type Items struct {
	MemberIDType string `json:"member_id_type"`
	MemberID     string `json:"member_id"`
	Name         string `json:"name"`
	TenantKey    string `json:"tenant_key"`
}

type GetUserGroupMemberResponse struct {
	BaseResponse
	Data struct {
		Items       []Items `json:"items"`
		PageToken   string  `json:"page_token"`
		HasMore     bool    `json:"has_more"`
		MemberTotal int     `json:"member_total"`
	} `json:"data"`
}
