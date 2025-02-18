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

func TestAccDepartmentResource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()

	Mock(common.GetUsersByIDAPI).Return(&common.UserInfoBatchGetResponse{
		Data: struct {
			Items []common.User `json:"items"`
		}{
			Items: []common.User{
				{
					UserID:  "ou_8fc0c1843c33c130462669327fb2113c",
					OpenID:  "ou_8fc0c1843c33c130462669327fb2113c",
					UnionID: "on_8fc0c1843c33c130462669327fb2113c",
				},
			},
		},
	}, nil).Build()

	// Mock DepartmentGetAPI dengan cara yang benar
	Mock(common.DepartmentGetAPI).When(func(ctx context.Context, client *common.LarkClient, departmentID string) bool {
		return departmentID == "dp_919c0000000000000000000000000000"
	}).Return(nil, fmt.Errorf("department not found")).When(func(ctx context.Context, client *common.LarkClient, departmentID string) bool {
		return departmentID == "0" || departmentID == "test_department_id" || departmentID == "od_test_department_id"
	}).Return(&common.DepartmentGetResponse{
		Data: struct {
			Department common.Department `json:"department"`
		}{
			Department: common.Department{
				DepartmentID:     "test_department_id",
				OpenDepartmentID: "od_test_department_id",
				BaseDepartment: common.BaseDepartment{
					Name: "Test Department",
					I18nName: common.I18nName{
						ZhCn: "测试部门",
						JaJp: "テスト部門",
						EnUs: "Test Department",
					},
					ParentDepartmentID: "0",
					LeaderUserID:       "ou_8fc0c1843c33c130462669327fb2113c",
					Order:              "1",
					UnitIDs:            []string{"unit_v1_919c0000000000000000000000000000"},
					CreateGroupChat:    true,
					Leaders: []common.DepartmentLeader{
						{
							LeaderID:   "user_v1_919c0000000000000000000000000000",
							LeaderType: 1,
						},
					},
					GroupChatEmployeeTypes: []int64{1, 2},
				},
				ChatID:      "oc_test_chat_id",
				MemberCount: 10,
				Status: common.DepartmentStatus{
					IsDeleted: false,
				},
			},
		},
	}, nil).Build()

	Mock(common.DepartmentCreateAPI).Return(&common.DepartmentGetResponse{
		Data: struct {
			Department common.Department `json:"department"`
		}{
			Department: common.Department{
				DepartmentID:     "dp_919c0000000000000000000000000000",
				OpenDepartmentID: "od_test_department_id",
				BaseDepartment: common.BaseDepartment{
					Name: "Test Department",
					I18nName: common.I18nName{
						ZhCn: "测试部门",
						JaJp: "テスト部門",
						EnUs: "Test Department",
					},
					ParentDepartmentID: "0",
					LeaderUserID:       "ou_8fc0c1843c33c130462669327fb2113c",
					Order:              "1",
					UnitIDs:            []string{"unit_v1_919c0000000000000000000000000000"},
					CreateGroupChat:    true,
					Leaders: []common.DepartmentLeader{
						{
							LeaderID:   "user_v1_919c0000000000000000000000000000",
							LeaderType: 1,
						},
					},
					GroupChatEmployeeTypes: []int64{1, 2},
				},
				ChatID:      "oc_test_chat_id",
				MemberCount: 10,
				Status: common.DepartmentStatus{
					IsDeleted: false,
				},
			},
		},
	}, nil).Build()

	// Mock untuk update department
	Mock(common.DepartmentUpdateAPI).Return(&common.DepartmentGetResponse{
		Data: struct {
			Department common.Department `json:"department"`
		}{
			Department: common.Department{
				DepartmentID:     "dp_919c0000000000000000000000000000",
				OpenDepartmentID: "od_test_department_id",
				BaseDepartment: common.BaseDepartment{
					Name: "Updated Test Department",
					I18nName: common.I18nName{
						ZhCn: "测试部门",
						JaJp: "テスト部門",
						EnUs: "Test Department",
					},
					ParentDepartmentID: "0",
					LeaderUserID:       "ou_8fc0c1843c33c130462669327fb2113c",
					Order:              "1",
					UnitIDs:            []string{"unit_v1_919c0000000000000000000000000000"},
					CreateGroupChat:    true,
					Leaders: []common.DepartmentLeader{
						{
							LeaderID:   "user_v1_919c0000000000000000000000000000",
							LeaderType: 1,
						},
					},
					GroupChatEmployeeTypes: []int64{1, 2},
				},
			},
		},
	}, nil).Build()

	Mock(common.DepartmentDeleteAPI).Return(&common.DepartmentDeleteResponse{
		BaseResponse: common.BaseResponse{
			Code: 0,
		},
	}, nil).Build()
	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read Testing
			{
				Config: providerConfig + `
				resource "lark_department" "test" {
					name = "Test Department"
					i18n_name = {
						en_us = "Test Department"
						zh_cn = "测试部门"
						ja_jp = "テスト部門"
					}
					parent_department_id = "0"
					department_id        = "dp_919c0000000000000000000000000000"
					leader_user_id       = "ou_8fc0c1843c33c130462669327fb2113c"
					order                = "1"
					unit_ids             = ["unit_v1_919c0000000000000000000000000000"]
					create_group_chat    = true
					leaders = [{
						leader_id   = "user_v1_919c0000000000000000000000000000"
						leader_type = 1
					}]
					group_chat_employee_types = [1, 2]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_department.test", "name", "Test Department"),
					resource.TestCheckResourceAttr("lark_department.test", "i18n_name.zh_cn", "测试部门"),
					resource.TestCheckResourceAttr("lark_department.test", "i18n_name.ja_jp", "テスト部門"),
					resource.TestCheckResourceAttr("lark_department.test", "i18n_name.en_us", "Test Department"),
					resource.TestCheckResourceAttr("lark_department.test", "parent_department_id", "0"),
					resource.TestCheckResourceAttr("lark_department.test", "department_id", "dp_919c0000000000000000000000000000"),
					resource.TestCheckResourceAttr("lark_department.test", "open_department_id", "od_test_department_id"),
					resource.TestCheckResourceAttr("lark_department.test", "leader_user_id", "ou_8fc0c1843c33c130462669327fb2113c"),
					resource.TestCheckResourceAttr("lark_department.test", "order", "1"),
					resource.TestCheckResourceAttr("lark_department.test", "unit_ids.0", "unit_v1_919c0000000000000000000000000000"),
					resource.TestCheckResourceAttr("lark_department.test", "create_group_chat", "true"),
					resource.TestCheckResourceAttr("lark_department.test", "chat_id", "oc_test_chat_id"),
					resource.TestCheckResourceAttr("lark_department.test", "leaders.0.leader_id", "user_v1_919c0000000000000000000000000000"),
					resource.TestCheckResourceAttr("lark_department.test", "leaders.0.leader_type", "1"),
					resource.TestCheckResourceAttr("lark_department.test", "group_chat_employee_types.0", "1"),
					resource.TestCheckResourceAttr("lark_department.test", "group_chat_employee_types.1", "2"),
					resource.TestCheckResourceAttr("lark_department.test", "member_count", "10"),
				),
			},
			// Update and Read Testing
			{
				Config: providerConfig + `
				resource "lark_department" "test" {
					name = "Updated Test Department"
					i18n_name = {
						en_us = "Test Department"
						zh_cn = "测试部门"
						ja_jp = "テスト部門"
					}
					parent_department_id = "0"
					department_id        = "dp_919c0000000000000000000000000000"
					leader_user_id       = "ou_8fc0c1843c33c130462669327fb2113c"
					order                = "1"
					unit_ids             = ["unit_v1_919c0000000000000000000000000000"]
					create_group_chat    = true
					leaders = [{
						leader_id   = "user_v1_919c0000000000000000000000000000"
						leader_type = 1
					}]
					group_chat_employee_types = [1, 2]
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_department.test", "name", "Updated Test Department"),
					resource.TestCheckResourceAttrSet("lark_department.test", "department_id"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
