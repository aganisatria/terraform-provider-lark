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

func TestAccUserGroupMemberResource(t *testing.T) {
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
	Mock(common.UsergroupMemberAddAPI).To(func(ctx context.Context, client *common.LarkClient, userGroupID string, req common.UsergroupMemberAddRequest) (*common.UsergroupMemberAddResponse, error) {
		results := make([]struct {
			MemberID string `json:"member_id"`
			Code     int    `json:"code"`
		}, len(req.Members))

		for i, member := range req.Members {
			results[i] = struct {
				MemberID string `json:"member_id"`
				Code     int    `json:"code"`
			}{
				MemberID: member.MemberID,
				Code:     0,
			}
		}

		return &common.UsergroupMemberAddResponse{
			Data: struct {
				Results []struct {
					MemberID string `json:"member_id"`
					Code     int    `json:"code"`
				} `json:"results"`
			}{
				Results: results,
			},
		}, nil
	}).Build()

	Mock(common.UsergroupMemberRemoveAPI).Return(&common.BaseResponse{
		Code: 0,
	}, nil).Build()

	Mock(common.UsergroupMemberGetByMemberTypeAPI).To(func(ctx context.Context, client *common.LarkClient, userGroupID string, pageToken string) (*common.UsergroupMemberGetResponse, error) {
		memberList := []common.UsergroupMember{{MemberID: "ou_0"}}

		if len(memberList) == 1 && pageToken == "" {
			memberList = append(memberList, common.UsergroupMember{MemberID: "ou_1"})
		}

		return &common.UsergroupMemberGetResponse{
			Data: struct {
				MemberList []common.UsergroupMember `json:"memberlist"`
				PageToken  string                   `json:"page_token"`
				HasMore    bool                     `json:"has_more"`
			}{
				MemberList: memberList,
				PageToken:  "0",
				HasMore:    false,
			},
		}, nil
	}).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Testing
			{
				Config: providerConfig + `resource "lark_user_group_member" "test" {
					user_group_id = "ug_test"
					member_ids = [
						"ou_0"
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_user_group_member.test", "user_group_id", "ug_test"),
					resource.TestCheckResourceAttr("lark_user_group_member.test", "member_ids.#", "1"),
					resource.TestCheckResourceAttr("lark_user_group_member.test", "member_ids.0", "ou_0"),
				),
			},

			// Update and Read Testing
			{
				Config: providerConfig + `resource "lark_user_group_member" "test" {
					user_group_id = "ug_test"
					member_ids = [
						"ou_0",
						"ou_1"
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_user_group_member.test", "user_group_id", "ug_test"),
					resource.TestCheckResourceAttr("lark_user_group_member.test", "member_ids.#", "2"),
					resource.TestCheckResourceAttr("lark_user_group_member.test", "member_ids.0", "ou_0"),
					resource.TestCheckResourceAttr("lark_user_group_member.test", "member_ids.1", "ou_1"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
