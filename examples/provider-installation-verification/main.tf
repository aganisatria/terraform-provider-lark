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

data "lark_user_by_id" "user" {
  users = [
    {
      union_id = "on_f96eca3d3bacf1f3dd54136083c33faa"
      open_id  = "ou_f96eca3d3bacf1f3dd54136083c33faa"
    },
  ]
  key_id = "union_id"
}

output "user" {
  value = data.lark_user_by_id.user
}
