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

func TestAccGroupChatResource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()
	Mock(common.GroupChatCreateAPI).Return(&common.GroupChatCreateResponse{
		Data: struct {
			ChatID                 string                       `json:"chat_id"`
			Avatar                 string                       `json:"avatar"`
			Name                   string                       `json:"name"`
			Description            string                       `json:"description"`
			I18nNames              common.I18nName              `json:"i18n_names"`
			OwnerID                string                       `json:"owner_id"`
			OwnerIDType            string                       `json:"owner_id_type"`
			UrgentSetting          string                       `json:"urgent_setting"`
			VideoConferenceSetting string                       `json:"video_conference_setting"`
			AddMemberPermission    string                       `json:"add_member_permission"`
			ShareCardPermission    string                       `json:"share_card_permission"`
			AtAllPermission        string                       `json:"at_all_permission"`
			EditPermission         string                       `json:"edit_permission"`
			GroupMessageType       string                       `json:"group_message_type"`
			ChatMode               string                       `json:"chat_mode"`
			ChatType               string                       `json:"chat_type"`
			ChatTag                string                       `json:"chat_tag"`
			External               bool                         `json:"external"`
			TenantKey              string                       `json:"tenant_key"`
			JoinMessageVisibility  string                       `json:"join_message_visibility"`
			LeaveMessageVisibility string                       `json:"leave_message_visibility"`
			MembershipApproval     string                       `json:"membership_approval"`
			ModerationPermission   string                       `json:"moderation_permission"`
			RestrictedModeSetting  common.RestrictedModeSetting `json:"restricted_mode_setting"`
			HideMemberCountSetting string                       `json:"hide_member_count_setting"`
		}{
			ChatID:      "test_chat_id",
			Avatar:      "xxxxxx",
			Name:        "ini contoh",
			Description: "ini description",
			I18nNames: common.I18nName{
				ZhCn: "中文",
				JaJp: "日本語",
				EnUs: "English",
			},
			OwnerID:                "test_owner_id",
			OwnerIDType:            "user_id",
			UrgentSetting:          "all_members",
			VideoConferenceSetting: "all_members",
			AddMemberPermission:    "all_members",
			ShareCardPermission:    "allowed",
			AtAllPermission:        "all_members",
			EditPermission:         "all_members",
			GroupMessageType:       "chat",
			ChatMode:               "group",
			ChatType:               "public",
			ChatTag:                "tag",
			External:               false,
			TenantKey:              "test_tenant_key",
			JoinMessageVisibility:  "all_members",
			LeaveMessageVisibility: "all_members",
			MembershipApproval:     "no_approval_required",
			ModerationPermission:   "all_members",
			RestrictedModeSetting: common.RestrictedModeSetting{
				Status:                         true,
				ScreenshotHasPermissionSetting: "not_anyone",
				DownloadHasPermissionSetting:   "all_members",
				MessageHasPermissionSetting:    "all_members",
			},
			HideMemberCountSetting: "all_members",
		},
	}, nil).Build()

	Mock(common.GroupChatGetAPI).Return(&common.GroupChatGetResponse{
		Data: struct {
			Avatar                 string                       `json:"avatar"`
			Name                   string                       `json:"name"`
			Description            string                       `json:"description"`
			I18nNames              common.I18nName              `json:"i18n_names"`
			AddMemberPermission    string                       `json:"add_member_permission"`
			ShareCardPermission    string                       `json:"share_card_permission"`
			AtAllPermission        string                       `json:"at_all_permission"`
			EditPermission         string                       `json:"edit_permission"`
			OwnerIDType            string                       `json:"owner_id_type"`
			OwnerID                string                       `json:"owner_id"`
			UserManagerIDList      []string                     `json:"user_manager_id_list"`
			BotManagerIDList       []string                     `json:"bot_manager_id_list"`
			GroupMessageType       string                       `json:"group_message_type"`
			ChatMode               string                       `json:"chat_mode"`
			ChatType               string                       `json:"chat_type"`
			ChatTag                string                       `json:"chat_tag"`
			JoinMessageVisibility  string                       `json:"join_message_visibility"`
			LeaveMessageVisibility string                       `json:"leave_message_visibility"`
			MembershipApproval     string                       `json:"membership_approval"`
			External               bool                         `json:"external"`
			TenantKey              string                       `json:"tenant_key"`
			UserCount              string                       `json:"user_count"`
			BotCount               string                       `json:"bot_count"`
			RestrictedModeSetting  common.RestrictedModeSetting `json:"restricted_mode_setting"`
			UrgentSetting          string                       `json:"urgent_setting"`
			VideoConferenceSetting string                       `json:"video_conference_setting"`
			HideMemberCountSetting string                       `json:"hide_member_count_setting"`
			ChatStatus             string                       `json:"chat_status"`
		}{
			Avatar:      "xxxxxx",
			Name:        "ini contoh",
			Description: "ini description",
			I18nNames: common.I18nName{
				ZhCn: "中文",
				JaJp: "日本語",
				EnUs: "English",
			},
			AddMemberPermission:    "all_members",
			ShareCardPermission:    "allowed",
			AtAllPermission:        "all_members",
			EditPermission:         "all_members",
			GroupMessageType:       "chat",
			ChatMode:               "group",
			ChatType:               "public",
			JoinMessageVisibility:  "all_members",
			LeaveMessageVisibility: "all_members",
			MembershipApproval:     "no_approval_required",
			RestrictedModeSetting: common.RestrictedModeSetting{
				Status:                         true,
				ScreenshotHasPermissionSetting: "not_anyone",
				DownloadHasPermissionSetting:   "all_members",
				MessageHasPermissionSetting:    "all_members",
			},
			UrgentSetting:          "all_members",
			VideoConferenceSetting: "all_members",
			HideMemberCountSetting: "all_members",
			ChatStatus:             "active",
			UserCount:              "100",
			BotCount:               "100",
			External:               false,
			TenantKey:              "test_tenant_key",
			ChatTag:                "tag",
			UserManagerIDList:      []string{"test_user_id_1", "test_user_id_2"},
			BotManagerIDList:       []string{"test_bot_id_1", "test_bot_id_2"},
		},
	}, nil).Build()

	// Update mock to change state
	Mock(common.GroupChatUpdateAPI).To(func(ctx context.Context, client *common.LarkClient, groupID string, req common.GroupChatUpdateRequest) (*common.BaseResponse, error) {
		if groupID != "test_chat_id" {
			return nil, fmt.Errorf("unexpected group_id: %s", groupID)
		}
		return &common.BaseResponse{
			Code: 0,
		}, nil
	}).Build()

	Mock(common.GroupChatDeleteAPI).Return(&common.BaseResponse{
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
				resource "lark_group_chat" "example" {
				avatar      = "xxxxxx"
				name        = "ini contoh"
				description = "ini description"
				i18n_names = {
					"zh_cn" = "中文"
					"ja_jp" = "日本語"
					"en_us" = "English"
				}
				group_message_type       = "chat"
				chat_mode                = "group"
				chat_type                = "public"
				join_message_visibility  = "all_members"
				leave_message_visibility = "all_members"
				membership_approval      = "no_approval_required"
				restricted_mode_setting = {
					status                            = true
					screenshot_has_permission_setting = "not_anyone"
					download_has_permission_setting   = "all_members"
					message_has_permission_setting    = "all_members"
				}
				urgent_setting            = "all_members"
				video_conference_setting  = "all_members"
				edit_permission           = "all_members"
				hide_member_count_setting = "all_members"
				add_member_permission     = "all_members"
				share_card_permission     = "allowed"
				at_all_permission         = "all_members"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_group_chat.example", "chat_id", "test_chat_id"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "name", "ini contoh"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "description", "ini description"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "i18n_names.zh_cn", "中文"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "i18n_names.ja_jp", "日本語"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "i18n_names.en_us", "English"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "group_message_type", "chat"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "chat_mode", "group"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "chat_type", "public"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "join_message_visibility", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "leave_message_visibility", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "membership_approval", "no_approval_required"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "restricted_mode_setting.status", "true"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "restricted_mode_setting.screenshot_has_permission_setting", "not_anyone"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "restricted_mode_setting.download_has_permission_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "restricted_mode_setting.message_has_permission_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "urgent_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "video_conference_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "edit_permission", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "hide_member_count_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "add_member_permission", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "share_card_permission", "allowed"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "at_all_permission", "all_members"),
				),
			},
			// Update and Read Testing
			{
				Config: providerConfig + `
				resource "lark_group_chat" "example" {
				avatar      = "xxxxxx"
				name        = "ini contoh"
				description = "ini update description"
				i18n_names = {
					"zh_cn" = "中文"
					"ja_jp" = "日本語"
					"en_us" = "English"
				}
				group_message_type       = "chat"
				chat_mode                = "group"
				chat_type                = "public"
				join_message_visibility  = "all_members"
				leave_message_visibility = "all_members"
				membership_approval      = "no_approval_required"
				restricted_mode_setting = {
					status                            = true
					screenshot_has_permission_setting = "not_anyone"
					download_has_permission_setting   = "all_members"
					message_has_permission_setting    = "all_members"
				}
				urgent_setting            = "all_members"
				video_conference_setting  = "all_members"
				edit_permission           = "all_members"
				hide_member_count_setting = "all_members"
				add_member_permission     = "all_members"
				share_card_permission     = "allowed"
				at_all_permission         = "all_members"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_group_chat.example", "chat_id", "test_chat_id"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "name", "ini contoh"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "description", "ini update description"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "i18n_names.zh_cn", "中文"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "i18n_names.ja_jp", "日本語"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "i18n_names.en_us", "English"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "group_message_type", "chat"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "chat_mode", "group"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "chat_type", "public"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "join_message_visibility", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "leave_message_visibility", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "membership_approval", "no_approval_required"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "restricted_mode_setting.status", "true"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "restricted_mode_setting.screenshot_has_permission_setting", "not_anyone"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "restricted_mode_setting.download_has_permission_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "restricted_mode_setting.message_has_permission_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "urgent_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "video_conference_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "edit_permission", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "hide_member_count_setting", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "add_member_permission", "all_members"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "share_card_permission", "allowed"),
					resource.TestCheckResourceAttr("lark_group_chat.example", "at_all_permission", "all_members"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}
