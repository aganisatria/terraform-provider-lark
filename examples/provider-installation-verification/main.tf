terraform {
  required_providers {
    lark = {
      source  = "custom/aganisatria/lark"
      version = "0.1.0"
    }
  }
}

provider "lark" {
  app_id     = "cli_a718cd690138d02f"
  app_secret = "J6m7yQiJ5MF0u4MT4q9AVb7nZQPNSRLd"
}

resource "lark_user_group" "example" {
  group_id    = "test"
  name        = "example1"
  description = "example"
}

resource "lark_user_group_member" "example2" {
  depends_on    = [lark_user_group.example]
  user_group_id = lark_user_group.example.group_id
  member_ids    = []
}
