// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

type UserStatus struct {
	IsFrozen    bool `json:"is_frozen"`
	IsResigned  bool `json:"is_resigned"`
	IsActivated bool `json:"is_activated"`
	IsExited    bool `json:"is_exited"`
	IsUnjoin    bool `json:"is_unjoin"`
}

type UserInfo struct {
	UserID string     `json:"user_id"`
	Mobile string     `json:"mobile"`
	Email  string     `json:"email"`
	Status UserStatus `json:"status"`
}

type UserInfoByEmailOrMobileBatchGetResponse struct {
	BaseResponse
	Data struct {
		UserList []UserInfo `json:"user_list"`
	} `json:"data"`
}

type UserInfoBatchGetRequest struct {
	Emails []string `json:"emails"`
}

type AvatarInfo struct {
	Avatar72     string `json:"avatar_72"`
	Avatar240    string `json:"avatar_240"`
	Avatar640    string `json:"avatar_640"`
	AvatarOrigin string `json:"avatar_origin"`
}

type Order struct {
	DepartmentID    string `json:"department_id"`
	UserOrder       int    `json:"user_order"`
	DepartmentOrder int    `json:"department_order"`
	IsPrimaryDept   bool   `json:"is_primary_dept"`
}

type CustomAttrGenericUser struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type CustomAttrValue struct {
	Text        string                `json:"text"`
	Url         string                `json:"url"`
	PCUrl       string                `json:"pc_url"`
	OptionID    string                `json:"option_id"`
	OptionValue string                `json:"option_value"`
	Name        string                `json:"name"`
	PictureUrl  string                `json:"picture_url"`
	GenericUser CustomAttrGenericUser `json:"generic_user"`
}

type CustomAttr struct {
	Type  string          `json:"type"`
	ID    string          `json:"id"`
	Value CustomAttrValue `json:"value"`
}

type AssignInfo struct {
	SubscriptionID string   `json:"subscription_id"`
	LicensePlanKey string   `json:"license_plan_key"`
	ProductName    string   `json:"product_name"`
	I18nName       I18nName `json:"i18n_name"`
	StartTime      string   `json:"start_time"`
	EndTime        string   `json:"end_time"`
}

type DepartmentPathName struct {
	Name     string   `json:"name"`
	I18nName I18nName `json:"i18n_name"`
}

type DepartmentPath struct {
	DepartmentIDs      []string           `json:"department_ids"`
	DepartmentPathName DepartmentPathName `json:"department_path_name"`
}

type DepartmentDetail struct {
	DepartmentID   string             `json:"department_id"`
	DepartmentName DepartmentPathName `json:"department_name"`
	DepartmentPath DepartmentPath     `json:"department_path"`
}

type User struct {
	UnionID                 string           `json:"union_id"`
	UserID                  string           `json:"user_id"`
	OpenID                  string           `json:"open_id"`
	Name                    string           `json:"name"`
	EnName                  string           `json:"en_name"`
	NickName                string           `json:"nick_name"`
	Email                   string           `json:"email"`
	Mobile                  string           `json:"mobile"`
	MobileVisible           bool             `json:"mobile_visible"`
	Gender                  int              `json:"gender"`
	AvatarKey               string           `json:"avatar_key"`
	Avatar                  AvatarInfo       `json:"avatar"`
	Status                  UserStatus       `json:"status"`
	DepartmentIDs           []string         `json:"department_ids"`
	LeaderUserID            string           `json:"leader_user_id"`
	City                    string           `json:"city"`
	Country                 string           `json:"country"`
	WorkStation             string           `json:"work_station"`
	JoinTime                int64            `json:"join_time"`
	IsTenantManager         bool             `json:"is_tenant_manager"`
	EmployeeNo              string           `json:"employee_no"`
	EmployeeType            int              `json:"employee_type"`
	Orders                  []Order          `json:"orders"`
	CustomAttrs             []CustomAttr     `json:"custom_attrs"`
	EnterpriseEmail         string           `json:"enterprise_email"`
	JobTitle                string           `json:"job_title"`
	Geo                     string           `json:"geo"`
	JobLevelID              string           `json:"job_level_id"`
	JobFamilyID             string           `json:"job_family_id"`
	SubscriptionIDs         []string         `json:"subscription_ids"`
	AssignInfo              AssignInfo       `json:"assign_info"`
	DepartmentPath          DepartmentDetail `json:"department_path"`
	DottedLineLeaderUserIDs []string         `json:"dotted_line_leader_user_ids"`
}

type UserInfoBatchGetResponse struct {
	BaseResponse
	Data struct {
		Items []User `json:"items"`
	} `json:"data"`
}
