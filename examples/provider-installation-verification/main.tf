terraform {
  required_providers {
    lark = {
      source  = "custom/aganisatria/lark"
      version = "0.1.0"
    }
  }
}

provider "lark" {
  app_id      = ""
  app_secret  = ""
  delay       = 1000
  retry_count = 3
}