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

func TestAccRoleResource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()

	Mock(common.RoleCreateAPI).Return(&common.RoleCreateResponse{
		Data: common.DataRoleCreateResponse{
			RoleID: "test_role_id",
		},
	}, nil).Build()

	// Update mock to change state
	Mock(common.RoleUpdateAPI).To(func(ctx context.Context, client *common.LarkClient, roleID string, req common.RoleRequest) (*common.BaseResponse, error) {
		return &common.BaseResponse{
			Code: 0,
		}, nil
	}).Build()

	Mock(common.RoleDeleteAPI).Return(&common.BaseResponse{
		Code: 0,
	}, nil).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Testing
			{
				Config: providerConfig + `
				resource "lark_role" "test" {
					role_name        = "Test Role"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_role.test", "role_name", "Test Role"),
					resource.TestCheckResourceAttrSet("lark_role.test", "role_id"),
				),
			},
			// Update and Read Testing
			{
				Config: providerConfig + `
				resource "lark_role" "test" {
					role_name        = "Updated Test Role"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_role.test", "role_name", "Updated Test Role"),
					resource.TestCheckResourceAttrSet("lark_role.test", "role_id"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
