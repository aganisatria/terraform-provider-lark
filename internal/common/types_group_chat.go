// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

type RestrictedModeSetting struct {
	Status                         bool   `json:"status,omitempty"`
	ScreenshotHasPermissionSetting string `json:"screenshot_has_permission_setting,omitempty"`
	DownloadHasPermissionSetting   string `json:"download_has_permission_setting,omitempty"`
	MessageHasPermissionSetting    string `json:"message_has_permission_setting,omitempty"`
}

type GroupChatCreateRequest struct {
	Avatar                 string                 `json:"avatar,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	Description            string                 `json:"description,omitempty"`
	I18nNames              I18nName               `json:"i18n_names"`
	OwnerID                string                 `json:"owner_id,omitempty"`
	UserIDList             []string               `json:"user_id_list,omitempty"`
	BotIDList              []string               `json:"bot_id_list,omitempty"`
	GroupMessageType       string                 `json:"group_message_type,omitempty"`
	ChatMode               string                 `json:"chat_mode,omitempty"`
	ChatType               string                 `json:"chat_type,omitempty"`
	JoinMessageVisibility  string                 `json:"join_message_visibility,omitempty"`
	LeaveMessageVisibility string                 `json:"leave_message_visibility,omitempty"`
	MembershipApproval     string                 `json:"membership_approval,omitempty"`
	RestrictedModeSetting  *RestrictedModeSetting `json:"restricted_mode_setting,omitempty"`
	UrgentSetting          string                 `json:"urgent_setting,omitempty"`
	VideoConferenceSetting string                 `json:"video_conference_setting,omitempty"`
	EditPermission         string                 `json:"edit_permission,omitempty"`
	HideMemberCountSetting string                 `json:"hide_member_count_setting,omitempty"`
}

type GroupChatCreateResponse struct {
	BaseResponse
	Data struct {
		ChatID                 string                `json:"chat_id"`
		Avatar                 string                `json:"avatar"`
		Name                   string                `json:"name"`
		Description            string                `json:"description"`
		I18nNames              I18nName              `json:"i18n_names"`
		OwnerID                string                `json:"owner_id"`
		OwnerIDType            string                `json:"owner_id_type"`
		UrgentSetting          string                `json:"urgent_setting"`
		VideoConferenceSetting string                `json:"video_conference_setting"`
		AddMemberPermission    string                `json:"add_member_permission"`
		ShareCardPermission    string                `json:"share_card_permission"`
		AtAllPermission        string                `json:"at_all_permission"`
		EditPermission         string                `json:"edit_permission"`
		GroupMessageType       string                `json:"group_message_type"`
		ChatMode               string                `json:"chat_mode"`
		ChatType               string                `json:"chat_type"`
		ChatTag                string                `json:"chat_tag"`
		External               bool                  `json:"external"`
		TenantKey              string                `json:"tenant_key"`
		JoinMessageVisibility  string                `json:"join_message_visibility"`
		LeaveMessageVisibility string                `json:"leave_message_visibility"`
		MembershipApproval     string                `json:"membership_approval"`
		ModerationPermission   string                `json:"moderation_permission"`
		RestrictedModeSetting  RestrictedModeSetting `json:"restricted_mode_setting"`
		HideMemberCountSetting string                `json:"hide_member_count_setting"`
	} `json:"data,omitempty"`
}

type GroupChatUpdateRequest struct {
	Avatar                 string                 `json:"avatar,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	Description            string                 `json:"description,omitempty"`
	I18nNames              I18nName               `json:"i18n_names,omitempty"`
	AddMemberPermission    string                 `json:"add_member_permission,omitempty"`
	ShareCardPermission    string                 `json:"share_card_permission,omitempty"`
	AtAllPermission        string                 `json:"at_all_permission,omitempty"`
	EditPermission         string                 `json:"edit_permission,omitempty"`
	OwnerID                string                 `json:"owner_id,omitempty"`
	JoinMessageVisibility  string                 `json:"join_message_visibility,omitempty"`
	LeaveMessageVisibility string                 `json:"leave_message_visibility,omitempty"`
	MembershipApproval     string                 `json:"membership_approval,omitempty"`
	RestrictedModeSetting  *RestrictedModeSetting `json:"restricted_mode_setting,omitempty"`
	ChatType               string                 `json:"chat_type,omitempty"`
	GroupMessageType       string                 `json:"group_message_type,omitempty"`
	UrgentSetting          string                 `json:"urgent_setting,omitempty"`
	VideoConferenceSetting string                 `json:"video_conference_setting,omitempty"`
	HideMemberCountSetting string                 `json:"hide_member_count_setting,omitempty"`
}

type GroupChatSpeechScopesUpdateRequest struct {
	ModerationSetting    string   `json:"moderation_setting,omitempty"`
	ModeratorAddedList   []string `json:"moderator_added_list,omitempty"`
	ModeratorRemovedList []string `json:"moderator_removed_list,omitempty"`
}

type GroupChatGetResponse struct {
	BaseResponse
	Data struct {
		Avatar                 string                `json:"avatar"`
		Name                   string                `json:"name"`
		Description            string                `json:"description"`
		I18nNames              I18nName              `json:"i18n_names"`
		AddMemberPermission    string                `json:"add_member_permission"`
		ShareCardPermission    string                `json:"share_card_permission"`
		AtAllPermission        string                `json:"at_all_permission"`
		EditPermission         string                `json:"edit_permission"`
		OwnerIDType            string                `json:"owner_id_type"`
		OwnerID                string                `json:"owner_id"`
		UserManagerIDList      []string              `json:"user_manager_id_list"`
		BotManagerIDList       []string              `json:"bot_manager_id_list"`
		GroupMessageType       string                `json:"group_message_type"`
		ChatMode               string                `json:"chat_mode"`
		ChatType               string                `json:"chat_type"`
		ChatTag                string                `json:"chat_tag"`
		JoinMessageVisibility  string                `json:"join_message_visibility"`
		LeaveMessageVisibility string                `json:"leave_message_visibility"`
		MembershipApproval     string                `json:"membership_approval"`
		External               bool                  `json:"external"`
		TenantKey              string                `json:"tenant_key"`
		UserCount              string                `json:"user_count"`
		BotCount               string                `json:"bot_count"`
		RestrictedModeSetting  RestrictedModeSetting `json:"restricted_mode_setting"`
		UrgentSetting          string                `json:"urgent_setting"`
		VideoConferenceSetting string                `json:"video_conference_setting"`
		HideMemberCountSetting string                `json:"hide_member_count_setting"`
		ChatStatus             string                `json:"chat_status"`
	} `json:"data"`
}

type GroupChatAdministratorRequest struct {
	ManagerIDs []string `json:"manager_ids"`
}

type GroupChatAdministratorResponse struct {
	BaseResponse
	Data struct {
		ChatManagers    []string `json:"chat_managers"`
		ChatBotManagers []string `json:"chat_bot_managers"`
	} `json:"data"`
}
