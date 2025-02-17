// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_acceptance_test

import (
	"context"
	"testing"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupChatMemberResource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()
	Mock(common.GetUsersByIDAPI).Return(&common.UserInfoBatchGetResponse{
		Data: struct {
			Items []common.User `json:"items"`
		}{
			Items: []common.User{
				{UserID: "0"},
			},
		},
	}, nil).Build()
	Mock(common.GroupChatMemberAddAPI).To(func(ctx context.Context, client *common.LarkClient, chatID string, req common.GroupChatMemberRequest) (*common.GroupChatMemberAddResponse, error) {
		return &common.GroupChatMemberAddResponse{
			Data: struct {
				InvalidIDList         []string `json:"invalid_id_list"`
				NotExistedIDList      []string `json:"not_existed_id_list"`
				PendingApprovalIDList []string `json:"pending_approval_id_list"`
			}{
				InvalidIDList:         []string{},
				NotExistedIDList:      []string{},
				PendingApprovalIDList: []string{},
			},
		}, nil
	}).Build()

	Mock(common.GroupChatAdministratorAddAPI).To(func(ctx context.Context, client *common.LarkClient, chatID string, req common.GroupChatAdministratorRequest) (*common.GroupChatAdministratorResponse, error) {
		return &common.GroupChatAdministratorResponse{
			Data: struct {
				ChatManagers    []string `json:"chat_managers"`
				ChatBotManagers []string `json:"chat_bot_managers"`
			}{
				ChatManagers:    req.ManagerIDs,
				ChatBotManagers: []string{},
			},
		}, nil
	}).Build()

	Mock(common.GroupChatMemberDeleteAPI).To(func(ctx context.Context, client *common.LarkClient, chatID string, req common.GroupChatMemberRequest) (*common.GroupChatMemberRemoveResponse, error) {
		return &common.GroupChatMemberRemoveResponse{
			Data: struct {
				InvalidIDList []string `json:"invalid_id_list"`
			}{
				InvalidIDList: []string{},
			},
		}, nil
	}).Build()

	Mock(common.GroupChatAdministratorDeleteAPI).To(func(ctx context.Context, client *common.LarkClient, chatID string, req common.GroupChatAdministratorRequest) (*common.GroupChatAdministratorResponse, error) {
		return &common.GroupChatAdministratorResponse{
			Data: struct {
				ChatManagers    []string `json:"chat_managers"`
				ChatBotManagers []string `json:"chat_bot_managers"`
			}{
				ChatManagers:    []string{},
				ChatBotManagers: []string{},
			},
		}, nil
	}).Build()

	var isUpdated bool
	Mock(common.GroupChatMemberGetAPI).To(func(ctx context.Context, client *common.LarkClient, chatID string) (*common.GroupChatMemberGetResponse, error) {
		memberList := []common.ListMember{{MemberID: "ou_0"}}
		if isUpdated {
			memberList = append(memberList, common.ListMember{MemberID: "ou_1"})
		}

		return &common.GroupChatMemberGetResponse{
			Data: struct {
				Items       []common.ListMember `json:"items"`
				PageToken   string              `json:"page_token"`
				HasMore     bool                `json:"has_more"`
				MemberTotal int64               `json:"member_total"`
			}{
				Items:       memberList,
				PageToken:   "",
				HasMore:     false,
				MemberTotal: int64(len(memberList)),
			},
		}, nil
	}).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Testing
			{
				Config: providerConfig + `resource "lark_group_chat_member" "test" {
					group_chat_id = "gc_test"
					member_ids = [
						"ou_0"
					]
					administrator_ids = [
						"ou_0"
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "group_chat_id", "gc_test"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "member_ids.#", "1"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "member_ids.0", "ou_0"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "administrator_ids.#", "1"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "administrator_ids.0", "ou_0"),
				),
			},

			// Update and Read Testing
			{
				PreConfig: func() { isUpdated = true },
				Config: providerConfig + `resource "lark_group_chat_member" "test" {
					group_chat_id = "gc_test"
					member_ids = [
						"ou_0",
						"ou_1"
					]
					administrator_ids = [
						"ou_0",
						"ou_1"
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "group_chat_id", "gc_test"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "member_ids.#", "2"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "member_ids.0", "ou_0"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "member_ids.1", "ou_1"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "administrator_ids.#", "2"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "administrator_ids.0", "ou_0"),
					resource.TestCheckResourceAttr("lark_group_chat_member.test", "administrator_ids.1", "ou_1"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
