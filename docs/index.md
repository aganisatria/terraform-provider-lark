---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lark Provider"
description: |-
  
---

# lark Provider



## Example Usage

```terraform
provider "lark" {
  app_id      = "app_id"
  app_secret  = "app_secret"
  delay       = 1000
  retry_count = 3
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_id` (String, Sensitive) The App ID for authenticating with Lark API
- `app_secret` (String, Sensitive) The App Secret for authenticating with Lark API

### Optional

- `delay` (Number) The delay for retrying the request
- `retry_count` (Number) The retry count for retrying the request
