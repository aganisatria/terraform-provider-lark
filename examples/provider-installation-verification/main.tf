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
  app_secret  = "032XJ3TJXKj5fjWX6RGuWfV0LW8DaIee"
  delay       = 10
  retry_count = 3
}
