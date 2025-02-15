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

type UserInfoBatchGetResponse struct {
	BaseResponse
	Data struct {
		UserList []UserInfo `json:"user_list"`
	} `json:"data"`
}

type UserInfoBatchGetRequest struct {
	Emails []string `json:"emails"`
}
