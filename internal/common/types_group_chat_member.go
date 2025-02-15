// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

type ListMember struct {
	MemberID     string `json:"member_id"`
	MemberIDType string `json:"member_id_type"`
	Name         string `json:"name"`
	TenantKey    string `json:"tenant_key"`
}

type GroupChatMemberGetResponse struct {
	BaseResponse
	Data struct {
		Items       []ListMember `json:"items"`
		PageToken   string       `json:"page_token"`
		HasMore     bool         `json:"has_more"`
		MemberTotal int64        `json:"member_total"`
	} `json:"data"`
}

type GroupChatMemberRequest struct {
	IDList []string `json:"id_list"`
}

type GroupChatMemberAddResponse struct {
	BaseResponse
	Data struct {
		InvalidIDList         []string `json:"invalid_id_list"`
		NotExistedIDList      []string `json:"not_existed_id_list"`
		PendingApprovalIDList []string `json:"pending_approval_id_list"`
	} `json:"data"`
}

type GroupChatMemberRemoveResponse struct {
	BaseResponse
	Data struct {
		InvalidIDList []string `json:"invalid_id_list"`
	} `json:"data"`
}

type GroupChatAdministratorRequest struct {
	ManagerIDs []string `json:"manager_ids"`
}

type GroupChatAdministratorResponse struct {
	BaseResponse
	Data struct {
		ChatManagers    []string `json:"chat_managers"`
		ChatBotManagers []string `json:"chat_bot_managers"`
	} `json:"data"`
}
