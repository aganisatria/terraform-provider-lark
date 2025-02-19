// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAccessTokenAPI(t *testing.T) {
	genericResponse := AccessTokenResponse{
		TenantAccessToken: "test_tenant_token",
		AppAccessToken:    "test_app_token",
	}

	tests := []struct {
		name          string
		appID         string
		appSecret     string
		mockResponse  AccessTokenResponse
		mockError     error
		wantTenant    string
		wantApp       string
		wantErr       bool
		expectedError string
	}{
		{
			name:          "error when calling DoInitializeRequest",
			appID:         "test_app_id",
			appSecret:     "test_secret",
			mockResponse:  genericResponse,
			mockError:     fmt.Errorf("invalid credentials"),
			wantTenant:    "",
			wantApp:       "",
			wantErr:       true,
			expectedError: "failed to get access token: invalid credentials",
		},
		{
			name:         "success",
			appID:        "test_app_id",
			appSecret:    "test_secret",
			mockResponse: genericResponse,
			mockError:    nil,
			wantTenant:   "test_tenant_token",
			wantApp:      "test_app_token",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoInitializeRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, requestBody interface{}, response interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := response.(*AccessTokenResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			gotTenant, gotApp, err := GetAccessTokenAPI(tt.appID, tt.appSecret, BASE_DELAY, BASE_RETRY_COUNT)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(gotTenant, ShouldBeEmpty)
				So(gotApp, ShouldBeEmpty)
			} else {
				So(err, ShouldBeNil)
				So(gotTenant, ShouldEqual, tt.wantTenant)
				So(gotApp, ShouldEqual, tt.wantApp)
			}
			UnPatchAll()
		})
	}
}

func TestUsergroupCreateAPI(t *testing.T) {
	tests := []struct {
		name             string
		req              UsergroupCreateRequest
		mockError        error
		expectedResponse UsergroupCreateResponse
		wantErr          bool
	}{
		{
			name:             "error on create",
			req:              UsergroupCreateRequest{},
			mockError:        fmt.Errorf("create failed"),
			expectedResponse: UsergroupCreateResponse{},
			wantErr:          true,
		},
		{
			name:             "success create",
			req:              UsergroupCreateRequest{},
			expectedResponse: UsergroupCreateResponse{},
			mockError:        nil,
			wantErr:          false,
		},
	}

	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			cleanup := SetupDoTenantRequest(tt.mockError, tt.expectedResponse)
			defer cleanup()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, 1)
			got, err := UsergroupCreateAPI(context.Background(), client, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				if got != nil {
					t.Errorf("expected nil response but got: %v", got)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(got, &tt.expectedResponse) {
					t.Errorf("expected response %v, got %v", tt.expectedResponse, got)
				}
			}
			UnPatchAll()
		})
	}
}

func TestUsergroupGetAPI(t *testing.T) {
	tests := []struct {
		name         string
		groupID      string
		mockResponse UsergroupGetResponse
		mockError    error
		wantErr      bool
	}{
		{
			name:         "error on get",
			groupID:      "group1",
			mockResponse: UsergroupGetResponse{},
			mockError:    fmt.Errorf("get failed"),
			wantErr:      true,
		},
		{
			name:         "success get",
			groupID:      "group1",
			mockResponse: UsergroupGetResponse{},
			mockError:    nil,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*UsergroupGetResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupGetAPI(context.Background(), client, tt.groupID)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
			UnPatchAll()
		})
	}
}

func TestUsergroupUpdateAPI(t *testing.T) {
	tests := []struct {
		name         string
		groupID      string
		req          UsergroupUpdateRequest
		mockResponse BaseResponse
		mockError    error
		wantErr      bool
	}{
		{
			name:         "error on update",
			groupID:      "group1",
			req:          UsergroupUpdateRequest{},
			mockResponse: BaseResponse{},
			mockError:    fmt.Errorf("update failed"),
			wantErr:      true,
		},
		{
			name:         "success update",
			groupID:      "group1",
			req:          UsergroupUpdateRequest{},
			mockResponse: BaseResponse{},
			mockError:    nil,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*BaseResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupUpdateAPI(context.Background(), client, tt.groupID, tt.req)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
			UnPatchAll()
		})
	}
}

func TestUsergroupDeleteAPI(t *testing.T) {
	tests := []struct {
		name         string
		groupID      string
		mockResponse BaseResponse
		mockError    error
		wantErr      bool
	}{
		{
			name:         "error on delete",
			groupID:      "group1",
			mockResponse: BaseResponse{},
			mockError:    fmt.Errorf("delete failed"),
			wantErr:      true,
		},
		{
			name:         "success delete",
			groupID:      "group1",
			mockResponse: BaseResponse{},
			mockError:    nil,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*BaseResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupDeleteAPI(context.Background(), client, tt.groupID)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
			UnPatchAll()
		})
	}
}

func TestUsergroupListAPI(t *testing.T) {
	PatchConvey("error on list", t, func() {
		Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
			return fmt.Errorf("list failed")
		}).Build()
		client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
		got, err := UsergroupListAPI(context.Background(), client)
		So(err, ShouldNotBeNil)
		So(got, ShouldBeNil)
		UnPatchAll()
	})

	PatchConvey("success list with pagination", t, func() {
		callCount := 0
		Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
			callCount++
			if r, ok := resp.(*UsergroupListResponse); ok {
				if callCount == 1 {
					r.Data.GroupList = []Group{{}}
					r.Data.PageToken = "token123"
				} else {
					r.Data.GroupList = []Group{{}}
					r.Data.PageToken = ""
				}
			} else {
				return fmt.Errorf("unexpected type for response")
			}
			return nil
		}).Build()

		client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
		got, err := UsergroupListAPI(context.Background(), client)
		So(err, ShouldBeNil)
		So(len(got.Data.GroupList), ShouldEqual, 2)
		UnPatchAll()
	})
}

func TestUsergroupMemberAddAPI(t *testing.T) {
	tests := []struct {
		name      string
		groupID   string
		req       UsergroupMemberAddRequest
		mockError error
		wantErr   bool
	}{
		{
			name:    "error on add member",
			groupID: "group1",
			req: UsergroupMemberAddRequest{
				Members: []UsergroupMember{
					{MemberID: "user1", MemberType: "user", MemberIDType: "open_id"},
				},
			},
			mockError: fmt.Errorf("add member failed"),
			wantErr:   true,
		},
		{
			name:    "success add member",
			groupID: "group1",
			req: UsergroupMemberAddRequest{
				Members: []UsergroupMember{
					{MemberID: "user1", MemberType: "user", MemberIDType: "open_id"},
					{MemberID: "user2", MemberType: "user", MemberIDType: "open_id"},
				},
			},
			mockError: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			callCount := 0
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				callCount++
				if tt.mockError != nil && callCount == 1 {
					return tt.mockError
				}
				if _, ok := resp.(*UsergroupMemberAddResponse); !ok {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupMemberAddAPI(context.Background(), client, tt.groupID, tt.req)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldNotBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestUsergroupMemberGetByMemberTypeAPI(t *testing.T) {
	tests := []struct {
		name         string
		groupID      string
		memberType   string
		mockResponse UsergroupMemberGetResponse
		mockError    error
		wantErr      bool
	}{
		{
			name:         "error on get by member type",
			groupID:      "group1",
			memberType:   "user",
			mockResponse: UsergroupMemberGetResponse{},
			mockError:    fmt.Errorf("get member failed"),
			wantErr:      true,
		},
		{
			name:         "success get by member type",
			groupID:      "group1",
			memberType:   "user",
			mockResponse: UsergroupMemberGetResponse{},
			mockError:    nil,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*UsergroupMemberGetResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupMemberGetByMemberTypeAPI(context.Background(), client, tt.groupID, tt.memberType)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
			UnPatchAll()
		})
	}
}

func TestUsergroupMemberRemoveAPI(t *testing.T) {
	tests := []struct {
		name      string
		groupID   string
		req       UsergroupMemberRemoveRequest
		mockError error
		wantErr   bool
	}{
		{
			name:    "error on remove member",
			groupID: "group1",
			req: UsergroupMemberRemoveRequest{
				Members: []UsergroupMember{
					{MemberID: "user1", MemberType: "user", MemberIDType: "open_id"},
				},
			},
			mockError: fmt.Errorf("remove member failed"),
			wantErr:   true,
		},
		{
			name:    "success remove member",
			groupID: "group1",
			req: UsergroupMemberRemoveRequest{
				Members: []UsergroupMember{
					{MemberID: "user1", MemberType: "user", MemberIDType: "open_id"},
					{MemberID: "user2", MemberType: "user", MemberIDType: "open_id"},
				},
			},
			mockError: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			callCount := 0
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				callCount++
				if tt.mockError != nil && callCount == 1 {
					return tt.mockError
				}
				if _, ok := resp.(*BaseResponse); !ok {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupMemberRemoveAPI(context.Background(), client, tt.groupID, tt.req)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldNotBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestGetUsersByIDAPI(t *testing.T) {
	tests := []struct {
		name         string
		userIDs      []string
		mockResponse UserInfoBatchGetResponse
		mockError    error
		wantErr      bool
	}{
		{
			name:         "error on get users",
			userIDs:      []string{"uid1", "uid2"},
			mockResponse: UserInfoBatchGetResponse{},
			mockError:    fmt.Errorf("get users failed"),
			wantErr:      true,
		},
		{
			name:         "success get users",
			userIDs:      []string{"uid1", "uid2"},
			mockResponse: UserInfoBatchGetResponse{},
			mockError:    nil,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				for _, uid := range tt.userIDs {
					if !Contains(path, uid) {
						return fmt.Errorf("missing user id in path")
					}
				}
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*UserInfoBatchGetResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GetUsersByIDAPI(context.Background(), client, tt.userIDs, OPEN_ID)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
			UnPatchAll()
		})
	}
}

func TestGetUserIdByEmailsAPI(t *testing.T) {
	tests := []struct {
		name         string
		request      UserInfoBatchGetRequest
		mockResponse UserInfoByEmailOrMobileBatchGetResponse
		mockError    error
		wantErr      bool
	}{
		{
			name:         "error on get user id by emails",
			request:      UserInfoBatchGetRequest{Emails: []string{"test@test.com", "test2@test.com"}},
			mockResponse: UserInfoByEmailOrMobileBatchGetResponse{},
			mockError:    fmt.Errorf("get user id by emails failed"),
			wantErr:      true,
		},
		{
			name:    "success get user id by emails",
			request: UserInfoBatchGetRequest{Emails: []string{"test@test.com", "test2@test.com"}},
			mockResponse: UserInfoByEmailOrMobileBatchGetResponse{
				Data: struct {
					UserList []UserInfo `json:"user_list"`
				}{
					UserList: []UserInfo{
						{UserID: "uid1", Email: "test@test.com"},
						{UserID: "uid2", Email: "test2@test.com"},
					},
				},
			},
			mockError: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*UserInfoByEmailOrMobileBatchGetResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GetUserIdByEmailsAPI(context.Background(), client, tt.request)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatCreateAPI(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse GroupChatCreateResponse
		mockError    error
		wantErr      bool
	}{
		{
			name:         "error on create",
			mockResponse: GroupChatCreateResponse{},
			mockError:    fmt.Errorf("create failed"),
			wantErr:      true,
		},
		{
			name:         "success create",
			mockResponse: GroupChatCreateResponse{},
			mockError:    nil,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*GroupChatCreateResponse); ok {
					*r = tt.mockResponse
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatCreateAPI(context.Background(), client, GroupChatCreateRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatDeleteAPI(t *testing.T) {
	tests := []struct {
		name      string
		chatID    string
		mockError error
		wantErr   bool
	}{
		{
			name:      "error on delete",
			chatID:    "chat1",
			mockError: fmt.Errorf("delete failed"),
			wantErr:   true,
		},
		{
			name:      "success delete",
			chatID:    "chat1",
			mockError: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*BaseResponse); ok {
					*r = BaseResponse{}
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatDeleteAPI(context.Background(), client, tt.chatID)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &BaseResponse{})
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatUpdateAPI(t *testing.T) {
	tests := []struct {
		name      string
		chatID    string
		mockError error
		wantErr   bool
	}{
		{
			name:      "error on update",
			chatID:    "chat1",
			mockError: fmt.Errorf("update failed"),
			wantErr:   true,
		},
		{
			name:      "success update",
			chatID:    "chat1",
			mockError: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*BaseResponse); ok {
					*r = BaseResponse{}
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatUpdateAPI(context.Background(), client, tt.chatID, GroupChatUpdateRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &BaseResponse{})
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatGetAPI(t *testing.T) {
	tests := []struct {
		name      string
		chatID    string
		mockError error
		wantErr   bool
	}{
		{
			name:      "error on get",
			chatID:    "chat1",
			mockError: fmt.Errorf("get failed"),
			wantErr:   true,
		},
		{
			name:      "success get",
			chatID:    "chat1",
			mockError: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				if tt.mockError != nil {
					return tt.mockError
				}
				if r, ok := resp.(*GroupChatGetResponse); ok {
					*r = GroupChatGetResponse{}
				}
				return nil
			}).Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatGetAPI(context.Background(), client, tt.chatID)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &GroupChatGetResponse{})
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatMemberGetAPI(t *testing.T) {
	tests := []struct {
		name         string
		chatID       string
		mockFn       func() *MockBuilder
		wantErr      bool
		expectedResp *GroupChatMemberGetResponse
	}{
		{
			name:   "success with multiple pages",
			chatID: "chat1",
			mockFn: func() *MockBuilder {
				responses := []GroupChatMemberGetResponse{
					{
						BaseResponse: BaseResponse{Code: 0, Msg: "success"},
						Data: struct {
							Items       []ListMember `json:"items"`
							PageToken   string       `json:"page_token"`
							HasMore     bool         `json:"has_more"`
							MemberTotal int64        `json:"member_total"`
						}{
							Items:     []ListMember{{Name: "User1"}},
							PageToken: "next_page",
							HasMore:   true,
						},
					},
					{
						BaseResponse: BaseResponse{Code: 0, Msg: "success"},
						Data: struct {
							Items       []ListMember `json:"items"`
							PageToken   string       `json:"page_token"`
							HasMore     bool         `json:"has_more"`
							MemberTotal int64        `json:"member_total"`
						}{
							Items:     []ListMember{{Name: "User2"}},
							PageToken: "",
							HasMore:   false,
						},
					},
				}

				callCount := 0
				return Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
					if callCount >= len(responses) {
						return fmt.Errorf("unexpected API call")
					}
					if r, ok := resp.(*GroupChatMemberGetResponse); ok {
						*r = responses[callCount]
					}
					callCount++
					return nil
				})
			},
			wantErr: false,
			expectedResp: &GroupChatMemberGetResponse{
				BaseResponse: BaseResponse{Code: 0, Msg: "success"},
				Data: struct {
					Items       []ListMember `json:"items"`
					PageToken   string       `json:"page_token"`
					HasMore     bool         `json:"has_more"`
					MemberTotal int64        `json:"member_total"`
				}{
					Items:       []ListMember{{Name: "User1"}, {Name: "User2"}},
					PageToken:   "",
					HasMore:     false,
					MemberTotal: 2,
				},
			},
		},
		{
			name:   "error on second page",
			chatID: "chat1",
			mockFn: func() *MockBuilder {
				responses := []GroupChatMemberGetResponse{
					{
						BaseResponse: BaseResponse{Code: 0, Msg: "success"},
						Data: struct {
							Items       []ListMember `json:"items"`
							PageToken   string       `json:"page_token"`
							HasMore     bool         `json:"has_more"`
							MemberTotal int64        `json:"member_total"`
						}{
							Items:     []ListMember{{Name: "User1"}},
							PageToken: "error_token",
							HasMore:   true,
						},
					},
				}

				callCount := 0
				return Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
					if callCount > 0 {
						return fmt.Errorf("error on second page")
					}
					if r, ok := resp.(*GroupChatMemberGetResponse); ok {
						*r = responses[callCount]
					}
					callCount++
					return nil
				})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			tt.mockFn().Build()

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatMemberGetAPI(context.Background(), client, tt.chatID)

			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, tt.expectedResp)
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatMemberAddAPI(t *testing.T) {
	tests := []struct {
		name    string
		chatID  string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name:   "error on split user and bot list",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return nil, nil, fmt.Errorf("split failed")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error on add bot",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("add failed")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error on add person",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				callCount := 0
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						callCount++
						if callCount == 1 {
							return nil
						}
						return fmt.Errorf("failed to add person members")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "success add",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatMemberAddAPI(context.Background(), client, tt.chatID, GroupChatMemberRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldNotBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatMemberDeleteAPI(t *testing.T) {
	tests := []struct {
		name    string
		chatID  string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name:   "error on split user and bot list",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return nil, nil, fmt.Errorf("split failed")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error on delete bot",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("delete failed")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error on delete person",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				callCount := 0
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						callCount++
						if callCount == 1 {
							return nil
						}
						return fmt.Errorf("failed to delete person members")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "there is error code on success delete",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						mockResponse := GroupChatMemberRemoveResponse{
							Data: struct {
								InvalidIDList []string `json:"invalid_id_list"`
							}{
								InvalidIDList: []string{"user1", "user2"},
							},
						}
						reflect.ValueOf(resp).Elem().Set(reflect.ValueOf(mockResponse))
						return nil
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "success delete",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatMemberDeleteAPI(context.Background(), client, tt.chatID, GroupChatMemberRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldNotBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatAdministratorAddAPI(t *testing.T) {
	tests := []struct {
		name    string
		chatID  string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name:   "error on split user and bot list",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return nil, nil, fmt.Errorf("split failed")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error because too many administrators",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2", "bot3", "bot4", "bot5", "bot6", "bot7", "bot8", "bot9", "bot10"}, []string{"user1", "user2"}, nil
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error on add administrator",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("add administrator failed")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error on add administrator",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				callCount := 0
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						callCount++
						if callCount == 1 {
							return nil
						}
						return fmt.Errorf("failed to add administrator")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "success add administrator",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatAdministratorAddAPI(context.Background(), client, tt.chatID, GroupChatAdministratorRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldNotBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestGroupChatAdministratorDeleteAPI(t *testing.T) {
	tests := []struct {
		name    string
		chatID  string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name:   "error on split user and bot list",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return nil, nil, fmt.Errorf("split failed")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error on delete administrator",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("delete administrator failed")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "error on delete administrator",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				callCount := 0
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						callCount++
						if callCount == 1 {
							return nil
						}
						return fmt.Errorf("failed to delete administrator")
					}),
				}
			},
			wantErr: true,
		},
		{
			name:   "success delete administrator",
			chatID: "chat1",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock(splitUserAndBotList).To(func(ids []string) (botList []string, personList []string, err error) {
						return []string{"bot1", "bot2"}, []string{"user1", "user2"}, nil
					}),
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GroupChatAdministratorDeleteAPI(context.Background(), client, tt.chatID, GroupChatAdministratorRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldNotBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestRoleCreateAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success create",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
		},
		{
			name: "error on create",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("create failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := RoleCreateAPI(context.Background(), client, RoleRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestRoleUpdateAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success update",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
		},
		{
			name: "error on update",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("update failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := RoleUpdateAPI(context.Background(), client, "role1", RoleRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestRoleDeleteAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success delete",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
		},
		{
			name: "error on delete",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("delete failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := RoleDeleteAPI(context.Background(), client, "role1")
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestRoleMemberAddAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success add",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
		},
		{
			name: "error on add",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("add failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := RoleMemberAddAPI(context.Background(), client, "role1", RoleMemberCreateRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestRoleMemberGetAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success get",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
		},
		{
			name: "error on get",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("get failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := RoleMemberGetAPI(context.Background(), client, "role1")
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestRoleMemberDeleteAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success delete",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on delete",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("delete failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := RoleMemberDeleteAPI(context.Background(), client, "role1", RoleMemberDeleteRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestDepartmentCreateAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success create",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on create",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("create failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := DepartmentCreateAPI(context.Background(), client, DepartmentCreateRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestDepartmentUpdateAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success update",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on update",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("update failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := DepartmentUpdateAPI(context.Background(), client, "department1", DepartmentUpdateRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestDepartmentDeleteAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success delete",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on delete",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("delete failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := DepartmentDeleteAPI(context.Background(), client, "department1")
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestDepartmentGetAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success get",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on get",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("get failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := DepartmentGetAPI(context.Background(), client, "department1")
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestWorkforceTypeCreateAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success create",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on create",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("create failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := WorkforceTypeCreateAPI(context.Background(), client, WorkforceTypeRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}

}

func TestWorkforceTypeUpdateAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success update",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on update",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("update failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := WorkforceTypeUpdateAPI(context.Background(), client, "enum1", WorkforceTypeRequest{})
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestWorkforceTypeDeleteAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success delete",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on delete",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("delete failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := WorkforceTypeDeleteAPI(context.Background(), client, "enum1")
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}

func TestWorkforceTypeGetAPI(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func() []*MockBuilder
		wantErr bool
	}{
		{
			name: "success get",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return nil
					}),
				}
			},
			wantErr: false,
		},
		{
			name: "error on get",
			mockFn: func() []*MockBuilder {
				return []*MockBuilder{
					Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
						return fmt.Errorf("get failed")
					}),
				}
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			for _, mockBuilder := range tt.mockFn() {
				mockBuilder.Build()
			}

			client := NewLarkClient("tenant-token", "app-token", "app-id", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := WorkforceTypeGetAllAPI(context.Background(), client)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			}
			UnPatchAll()
		})
	}
}
