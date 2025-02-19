resource "lark_workforce_type" "example" {
  content     = "example3"
  enum_type   = 2
  enum_status = 1
  i18n_content = [
    {
      locale = "en_us"
      value  = "Example"
    },
    {
      locale = "zh_cn"
      value  = "示例"
    }
  ]
}