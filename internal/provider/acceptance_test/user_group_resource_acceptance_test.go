// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_acceptance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupResource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()

	var isUpdated bool = false

	mocker := Mock(common.UsergroupGetAPI)

	// Condition 1: When create (check existing)
	mocker.When(func(ctx context.Context, client *common.LarkClient, groupID string) bool {
		return groupID == ""
	}).Return(nil, fmt.Errorf("group not found"))

	// Condition 2: When read (both initial and after update)
	mocker.When(func(ctx context.Context, client *common.LarkClient, groupID string) bool {
		return groupID == "test_group_id"
	}).To(func(ctx context.Context, client *common.LarkClient, groupID string) (*common.UsergroupGetResponse, error) {
		name := "Test Group"
		description := "Test Description"
		if isUpdated {
			name = "Updated Test Group"
			description = "Updated Test Description"
		}

		return &common.UsergroupGetResponse{
			Data: struct {
				Group struct {
					ID                    string `json:"id"`
					Name                  string `json:"name"`
					Description           string `json:"description"`
					MemberUserCount       int64  `json:"member_user_count"`
					MemberDepartmentCount int64  `json:"member_department_count"`
				} `json:"group"`
			}{
				Group: struct {
					ID                    string `json:"id"`
					Name                  string `json:"name"`
					Description           string `json:"description"`
					MemberUserCount       int64  `json:"member_user_count"`
					MemberDepartmentCount int64  `json:"member_department_count"`
				}{
					ID:                    groupID,
					Name:                  name,
					Description:           description,
					MemberUserCount:       1,
					MemberDepartmentCount: 1,
				},
			},
		}, nil
	})

	mocker.Build()

	Mock(common.UsergroupCreateAPI).Return(&common.UsergroupCreateResponse{
		Data: struct {
			GroupID string `json:"group_id"`
		}{
			GroupID: "test_group_id",
		},
	}, nil).Build()

	// Update mock to change state
	Mock(common.UsergroupUpdateAPI).To(func(ctx context.Context, client *common.LarkClient, groupID string, req common.UsergroupUpdateRequest) (*common.BaseResponse, error) {
		isUpdated = true
		return &common.BaseResponse{
			Code: 0,
		}, nil
	}).Build()

	Mock(common.UsergroupDeleteAPI).Return(&common.BaseResponse{
		Code: 0,
	}, nil).Build()

	Mock(common.UsergroupListAPI).Return(&common.UsergroupListResponse{
		Data: struct {
			GroupList []common.Group `json:"grouplist"`
			PageToken string         `json:"page_token"`
			HasMore   bool           `json:"has_more"`
		}{
			GroupList: []common.Group{},
			PageToken: "",
			HasMore:   false,
		},
	}, nil).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Testing
			{
				Config: providerConfig + `
				resource "lark_user_group" "test" {
					name        = "Test Group"
					description = "Test Description"
					type        = "1"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_user_group.test", "name", "Test Group"),
					resource.TestCheckResourceAttr("lark_user_group.test", "description", "Test Description"),
					resource.TestCheckResourceAttr("lark_user_group.test", "type", "1"),
					resource.TestCheckResourceAttrSet("lark_user_group.test", "group_id"),
				),
			},
			// Update and Read Testing
			{
				Config: providerConfig + `
				resource "lark_user_group" "test" {
					name        = "Updated Test Group"
					description = "Updated Test Description"
					type        = "1"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_user_group.test", "name", "Updated Test Group"),
					resource.TestCheckResourceAttr("lark_user_group.test", "description", "Updated Test Description"),
					resource.TestCheckResourceAttr("lark_user_group.test", "type", "1"),
					resource.TestCheckResourceAttrSet("lark_user_group.test", "group_id"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
