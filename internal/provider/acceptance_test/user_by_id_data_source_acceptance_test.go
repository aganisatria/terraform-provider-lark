// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_acceptance_test

import (
	"testing"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserByIdDataSource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()
	Mock(common.GetUsersByIDAPI).Return(&common.UserInfoBatchGetResponse{
		Data: struct {
			Items []common.User `json:"items"`
		}{
			Items: []common.User{
				{
					UserID:  "0",
					OpenID:  "ou_0",
					UnionID: "on_0",
				},
			},
		},
	}, nil).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "lark_user_by_id" "test" {
					users = [
						{
							user_id = "0"
						}
					]
					key_id = "user_id"
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verified that the number of users returned is 1
					resource.TestCheckResourceAttr("data.lark_user_by_id.test", "users.#", "1"),
					// Verified that the first user attribute is consistent with the fake response
					resource.TestCheckResourceAttr("data.lark_user_by_id.test", "users.0.user_id", "0"),
					resource.TestCheckResourceAttr("data.lark_user_by_id.test", "users.0.open_id", "ou_0"),
					resource.TestCheckResourceAttr("data.lark_user_by_id.test", "users.0.union_id", "on_0"),
				),
			},
		},
	})
}
