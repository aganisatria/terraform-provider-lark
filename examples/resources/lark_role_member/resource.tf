resource "lark_role_member" "example" {
  role_id    = lark_role.example.id
  member_ids = ["ou_8fc0c1843c33c130462669327fb2113c"]
}