data "lark_user_based_on_email" "example" {
  users = [
    {
      email = "example@gmail.com"
    }
  ]
}