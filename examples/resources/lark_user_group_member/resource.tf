resource "lark_user_group_member" "example" {
  user_group_id = lark_user_group.example.user_group_id
  member_ids = [
    "ou_8fc0c1843c33c130462669327fb2113c"
  ]
}
