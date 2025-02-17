// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

type UsergroupCreateRequest struct {
	GroupID     string `json:"group_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

type UsergroupCreateResponse struct {
	BaseResponse
	Data struct {
		GroupID string `json:"group_id"`
	} `json:"data"`
}

type UsergroupGetResponse struct {
	BaseResponse
	Data struct {
		Group struct {
			ID                    string `json:"id"`
			Name                  string `json:"name"`
			Description           string `json:"description"`
			MemberUserCount       int64  `json:"member_user_count"`
			MemberDepartmentCount int64  `json:"member_department_count"`
		} `json:"group"`
	} `json:"data"`
}

type UsergroupUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

const (
	UsergroupMemberTypeUser       string = "user"
	UsergroupMemberTypeDepartment string = "department" // available soon
)

type UsergroupMember struct {
	MemberID     string `json:"member_id"`
	MemberType   string `json:"member_type"`
	MemberIDType string `json:"member_id_type"`
}

type UsergroupMemberAddRequest struct {
	Members []UsergroupMember `json:"members"`
}

type UsergroupMemberAddResponse struct {
	BaseResponse
	Data struct {
		Results []struct {
			MemberID string `json:"member_id"`
			Code     int    `json:"code"`
		} `json:"results"`
	} `json:"data"`
}

type UsergroupMemberGetResponse struct {
	BaseResponse
	Data struct {
		MemberList []UsergroupMember `json:"memberlist"`
		PageToken  string            `json:"page_token"`
		HasMore    bool              `json:"has_more"`
	} `json:"data"`
}

type UsergroupMemberRemoveRequest struct {
	Members []UsergroupMember `json:"members"`
}

type Group struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	MemberUserCount       int64  `json:"member_user_count"`
	MemberDepartmentCount int64  `json:"member_department_count"`
}

type UsergroupListResponse struct {
	BaseResponse
	Data struct {
		GroupList []Group `json:"grouplist"`
		PageToken string  `json:"page_token"`
		HasMore   bool    `json:"has_more"`
	} `json:"data"`
}
