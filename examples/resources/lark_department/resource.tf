resource "lark_department" "example" {
  name = "example"
  i18n_name = {
    en_us = "example"
    zh_cn = "示例"
    ja_jp = "例"
  }
  parent_department_id = "0"
  department_id        = "dp_919c0000000000000000000000000000"
  order                = "1"
  unit_ids             = ["unit_v1_919c0000000000000000000000000000"]
  create_group_chat    = true
  leaders = [{
    leader_id   = "user_v1_919c0000000000000000000000000000"
    leader_type = 1
  }]
  group_chat_employee_types = [1, 2]
}