// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

const (
	OpenID          string = "open_id"
	UnionID         string = "union_id"
	UserIDInATenant string = "user_id"
)

const (
	DepartmentIDInATenant string = "department_id"
	OpenDepartmentID      string = "open_department_id"
)

type DepartmentLeader struct {
	LeaderID   string `json:"leaderID"`
	LeaderType int64  `json:"leaderType"`
}

type BaseDepartment struct {
	Name                   string             `json:"name"`
	I18nName               I18nName           `json:"i18n_name,omitempty"`
	ParentDepartmentID     string             `json:"parent_department_id"`
	LeaderUserID           string             `json:"leader_user_id,omitempty"`
	Order                  string             `json:"order,omitempty"`
	UnitIDs                []string           `json:"unit_ids,omitempty"`
	Leaders                []DepartmentLeader `json:"leaders,omitempty"`
	GroupChatEmployeeTypes []int64            `json:"group_chat_employee_types,omitempty"`
	CreateGroupChat        bool               `json:"create_group_chat,omitempty"`
}

type DepartmentCreateRequest struct {
	BaseDepartment
	DepartmentID string `json:"department_id,omitempty"`
}

type DepartmentUpdateRequest struct {
	BaseDepartment
}

type DepartmentStatus struct {
	IsDeleted bool `json:"is_deleted"`
}

type Department struct {
	BaseDepartment
	DepartmentID     string           `json:"department_id"`
	OpenDepartmentID string           `json:"open_department_id"`
	ChatID           string           `json:"chat_id"`
	MemberCount      int              `json:"member_count"`
	Status           DepartmentStatus `json:"status"`
}

type DepartmentGetResponse struct {
	BaseResponse
	Data struct {
		Department Department `json:"department"`
	} `json:"data"`
}

type DepartmentDeleteResponse struct {
	BaseResponse
	Data struct{} `json:"data"`
}
