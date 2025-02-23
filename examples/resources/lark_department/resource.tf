resource "lark_department" "example" {
  name = "testooo"
  i18n_name = {
    ja_jp = "テストooo"
    en_us = "Testooo"
    zh_cn = "测试ooo"
  }
  parent_department_id = "0"
  department_id        = "digidawi"
  leader_user_id       = "ou_8fc0c1843c33c130462669327fb2113c"
  order                = "1"
  unit_ids             = ["ou_8fc0c1843c33c130462669327fb2113c"]
  create_group_chat    = false
  leaders = [
    {
      leader_type = 1
      leader_id   = "ou_8fc0c1843c33c130462669327fb2113c"
    }
  ]
  group_chat_employee_types = [1, 2]
}
