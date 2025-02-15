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
    screenshot_has_permission_setting = "all_members"
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