resource "lark_group_chat_member" "example" {
  group_chat_id     = "test"
  member_ids        = ["ou_test"]
  administrator_ids = ["ou_test"]
}