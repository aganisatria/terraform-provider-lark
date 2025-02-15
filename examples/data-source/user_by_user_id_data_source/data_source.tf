data "lark_user_based_on_union_id" "example" {
  users = [
    {
      union_id = "on_1234567890"
    }
  ]
}
