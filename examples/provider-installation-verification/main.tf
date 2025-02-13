terraform {
  required_providers {
    lark = {
      source  = "custom/aganisatria/lark"
      version = "0.1.0"
    }
  }
}

provider "lark" {
  app_id      = "app_id"
  app_secret  = "app_secret"
  delay       = 1000
  retry_count = 3
}

resource "lark_user_group" "example" {
  group_id    = "test"
  name        = "example1"
  description = "example"
}

