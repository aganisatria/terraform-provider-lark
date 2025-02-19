// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

type I18nContent struct {
	Locale string `json:"locale"`
	Value  string `json:"value"`
}

type WorkforceTypeRequest struct {
	Content     string        `json:"content"`
	EnumType    int           `json:"enum_type"`
	EnumStatus  int           `json:"enum_status"`
	I18nContent []I18nContent `json:"i18n_content"`
}

type EmployeeTypeEnum struct {
	EnumID      string        `json:"enum_id"`
	EnumValue   string        `json:"enum_value"`
	Content     string        `json:"content"`
	EnumType    int           `json:"enum_type"`
	EnumStatus  int           `json:"enum_status"`
	I18nContent []I18nContent `json:"i18n_content"`
}

type WorkforceTypeResponse struct {
	BaseResponse
	Data struct {
		EmployeeTypeEnum EmployeeTypeEnum `json:"employee_type_enum"`
	} `json:"data"`
}

type WorkforceTypeGetResponse struct {
	BaseResponse
	Data struct {
		Items     []EmployeeTypeEnum `json:"items"`
		HasMore   bool               `json:"has_more"`
		PageToken string             `json:"page_token"`
	} `json:"data"`
}
