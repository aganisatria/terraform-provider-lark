terraform {
  required_providers {
    lark = {
      source  = "custom/aganisatria/lark"
      version = "0.1.0"
    }
  }
}

provider "lark" {
  app_id      = "cli_a718cd690138d02f"
  app_secret  = "J6m7yQiJ5MF0u4MT4q9AVb7nZQPNSRLd"
  delay       = 1000
  retry_count = 3
}

resource "lark_user_group" "example" {
  group_id    = "cekiceki"
  name        = "examplelelealay"
  description = "example"
}

resource "lark_user_group_member" "example" {
  depends_on    = [lark_user_group.example]
  user_group_id = lark_user_group.example.group_id
  member_ids = [
    "ou_8fc0c1843c33c130462669327fb2113c"
  ]
}
