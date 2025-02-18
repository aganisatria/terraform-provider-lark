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

func TestAccRoleMemberResource(t *testing.T) {
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
	Mock(common.RoleMemberAddAPI).To(func(ctx context.Context, client *common.LarkClient, roleID string, req common.RoleMemberCreateRequest) (*common.RoleMemberCreateResponse, error) {
		return &common.RoleMemberCreateResponse{
			Data: common.DataRoleMemberCreateDeleteResponse{
				Results: []common.FunctionalRoleMemberResult{
					{UserID: "0", Reason: 1},
				},
			},
		}, nil
	}).Build()

	Mock(common.RoleMemberDeleteAPI).To(func(ctx context.Context, client *common.LarkClient, roleID string, req common.RoleMemberDeleteRequest) (*common.RoleMemberDeleteResponse, error) {
		return &common.RoleMemberDeleteResponse{
			Data: common.DataRoleMemberCreateDeleteResponse{
				Results: []common.FunctionalRoleMemberResult{
					{UserID: "0", Reason: 1},
				},
			},
		}, nil
	}).Build()

	var isUpdated bool
	Mock(common.RoleMemberGetAPI).To(func(ctx context.Context, client *common.LarkClient, roleID string) (*common.RoleMemberGetResponse, error) {
		memberList := []common.RoleMember{{UserID: "ou_0"}}
		if isUpdated {
			memberList = append(memberList, common.RoleMember{UserID: "ou_1"})
		}

		return &common.RoleMemberGetResponse{
			Data: common.DataRoleMemberGetResponse{
				Members:   memberList,
				PageToken: "",
				HasMore:   false,
			},
		}, nil
	}).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Testing
			{
				Config: providerConfig + `resource "lark_role_member" "test" {
					role_id = "role_test"
					member_ids = [
						"ou_0"
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_role_member.test", "role_id", "role_test"),
					resource.TestCheckResourceAttr("lark_role_member.test", "member_ids.#", "1"),
					resource.TestCheckResourceAttr("lark_role_member.test", "member_ids.0", "ou_0"),
				),
			},

			// Update and Read Testing
			{
				PreConfig: func() { isUpdated = true },
				Config: providerConfig + `resource "lark_role_member" "test" {
					role_id = "role_test"
					member_ids = [
						"ou_0",
						"ou_1"
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_role_member.test", "role_id", "role_test"),
					resource.TestCheckResourceAttr("lark_role_member.test", "member_ids.#", "2"),
					resource.TestCheckResourceAttr("lark_role_member.test", "member_ids.0", "ou_0"),
					resource.TestCheckResourceAttr("lark_role_member.test", "member_ids.1", "ou_1"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
