// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_acceptance_test

import (
	"testing"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkforceTypeResource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()

	Mock(common.WorkforceTypeCreateAPI).Return(&common.WorkforceTypeResponse{
		Data: struct {
			EmployeeTypeEnum common.EmployeeTypeEnum `json:"employee_type_enum"`
		}{
			EmployeeTypeEnum: common.EmployeeTypeEnum{
				EnumID:     "test_enum_id",
				EnumType:   1,
				EnumStatus: 1,
				EnumValue:  "test_enum_value",
			},
		},
	}, nil).Build()

	// Update mock to change state
	Mock(common.WorkforceTypeUpdateAPI).Return(&common.WorkforceTypeResponse{
		Data: struct {
			EmployeeTypeEnum common.EmployeeTypeEnum `json:"employee_type_enum"`
		}{
			EmployeeTypeEnum: common.EmployeeTypeEnum{
				EnumID:     "test_enum_id",
				EnumType:   1,
				EnumStatus: 1,
				EnumValue:  "test_enum_value",
			},
		},
	}, nil).Build()

	Mock(common.WorkforceTypeDeleteAPI).Return(&common.BaseResponse{
		Code: 0,
	}, nil).Build()

	Mock(common.WorkforceTypeGetAllAPI).Return(&common.WorkforceTypeGetResponse{
		Data: struct {
			Items     []common.EmployeeTypeEnum `json:"items"`
			HasMore   bool                      `json:"has_more"`
			PageToken string                    `json:"page_token"`
		}{
			Items: []common.EmployeeTypeEnum{
				{
					EnumID:     "test_enum_id",
					EnumType:   2,
					EnumStatus: 1,
					EnumValue:  "test_enum_value",
				},
			},
			HasMore:   false,
			PageToken: "",
		},
	}, nil).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Testing
			{
				Config: providerConfig + `
				resource "lark_workforce_type" "test" {
					content     = "Test Content"
					enum_type   = 2
					enum_status = 1
					i18n_content = [
						{
							locale = "en"
							content = "Test Content"
						}
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_workforce_type.test", "content", "Test Content"),
					resource.TestCheckResourceAttr("lark_workforce_type.test", "enum_type", "2"),
					resource.TestCheckResourceAttr("lark_workforce_type.test", "enum_status", "1"),
					resource.TestCheckResourceAttrSet("lark_workforce_type.test", "enum_id"),
				),
			},
			// Update and Read Testing
			{
				Config: providerConfig + `
				resource "lark_workforce_type" "test" {
					content     = "Updated Test Content"
					enum_type   = 2
					enum_status = 1
					i18n_content = [
						{
							locale = "en"
							content = "Updated Test Content"
						}
					]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_workforce_type.test", "content", "Updated Test Content"),
					resource.TestCheckResourceAttr("lark_workforce_type.test", "enum_type", "2"),
					resource.TestCheckResourceAttr("lark_workforce_type.test", "enum_status", "1"),
					resource.TestCheckResourceAttrSet("lark_workforce_type.test", "enum_id"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
