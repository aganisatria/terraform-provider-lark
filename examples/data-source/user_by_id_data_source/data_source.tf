data "lark_user_by_id" "example" {
  users = [
    {
      user_id = "1234567890"
    }
  ]
  key_id = "user_id"
}
