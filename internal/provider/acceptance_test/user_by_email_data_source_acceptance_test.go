// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_acceptance_test

import (
	"testing"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserByEmailDataSource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()
	Mock(common.GetUserIdByEmailsAPI).Return(&common.UserInfoByEmailOrMobileBatchGetResponse{
		Data: struct {
			UserList []common.UserInfo `json:"user_list"`
		}{
			UserList: []common.UserInfo{{Email: "example@gmail.com", UserID: "user-0"}},
		},
	}, nil).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "lark_user_by_email" "test" {
					users = [
						{
							email = "example@gmail.com"
						}
					]
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verified that the number of users returned is 1
					resource.TestCheckResourceAttr("data.lark_user_by_email.test", "users.#", "1"),
					// Verified that the first user attribute is consistent with the fake response
					resource.TestCheckResourceAttr("data.lark_user_by_email.test", "users.0.email", "example@gmail.com"),
					resource.TestCheckResourceAttr("data.lark_user_by_email.test", "users.0.user_id", "user-0"),
				),
			},
		},
	})
}
