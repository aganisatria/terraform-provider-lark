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
			name:      "error when calling DoInitializeRequest",
			appID:     "test_app_id",
			appSecret: "test_secret",
			mockResponse: genericResponse,
			mockError:    fmt.Errorf("invalid credentials"),
			wantTenant:   "",
			wantApp:      "",
			wantErr:      true,
			expectedError: "failed to get access token: invalid credentials",
		},
		{
			name:      "success",
			appID:     "test_app_id",
			appSecret: "test_secret",
			mockResponse: genericResponse,
			mockError:  nil,
			wantTenant: "test_tenant_token",
			wantApp:    "test_app_token",
			wantErr:    false,
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
			name: "success create",
			req:  UsergroupCreateRequest{},
			expectedResponse: UsergroupCreateResponse{},
			mockError: nil,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cleanup := SetupDoTenantRequest(tt.mockError, tt.expectedResponse)
			defer cleanup()

			client := NewLarkClient("", "", BASE_DELAY, 1)
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
			name:    "success get",
			groupID: "group1",
			mockResponse: UsergroupGetResponse{},
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
				if r, ok := resp.(*UsergroupGetResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupGetAPI(context.Background(), client, tt.groupID)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
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
			name:    "success update",
			groupID: "group1",
			req:     UsergroupUpdateRequest{},
			mockResponse: BaseResponse{},
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
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupUpdateAPI(context.Background(), client, tt.groupID, tt.req)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
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
			name:    "success delete",
			groupID: "group1",
			mockResponse: BaseResponse{},
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
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupDeleteAPI(context.Background(), client, tt.groupID)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
		})
	}
}

func TestUsergroupListAPI(t *testing.T) {
	PatchConvey("error on list", t, func() {
		Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
			return fmt.Errorf("list failed")
		}).Build()
		client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
		got, err := UsergroupListAPI(context.Background(), client)
		So(err, ShouldNotBeNil)
		So(got, ShouldBeNil)
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

		client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
		got, err := UsergroupListAPI(context.Background(), client)
		So(err, ShouldBeNil)
		So(len(got.Data.GroupList), ShouldEqual, 2)
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

			client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupMemberAddAPI(context.Background(), client, tt.groupID, tt.req)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldNotBeNil)
			}
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
			name:       "success get by member type",
			groupID:    "group1",
			memberType: "user",
			mockResponse: UsergroupMemberGetResponse{},
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
				if r, ok := resp.(*UsergroupMemberGetResponse); ok {
					*r = tt.mockResponse
				} else {
					return fmt.Errorf("unexpected type for response")
				}
				return nil
			}).Build()

			client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupMemberGetByMemberTypeAPI(context.Background(), client, tt.groupID, tt.memberType)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
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

			client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := UsergroupMemberRemoveAPI(context.Background(), client, tt.groupID, tt.req)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldNotBeNil)
			}
		})
	}
}

func TestGetUsersByOpenIDAPI(t *testing.T) {
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
			name:    "success get users",
			userIDs: []string{"uid1", "uid2"},
			mockResponse: UserInfoBatchGetResponse{},
			mockError: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		PatchConvey(tt.name, t, func() {
			Mock((*LarkClient).DoTenantRequest).To(func(c *LarkClient, ctx context.Context, method HTTPMethod, path string, reqBody interface{}, resp interface{}) error {
				for _, uid := range tt.userIDs {
					if !contains(path, uid) {
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

			client := NewLarkClient("", "", BASE_DELAY, BASE_RETRY_COUNT)
			got, err := GetUsersByOpenIDAPI(context.Background(), client, tt.userIDs)
			if tt.wantErr {
				So(err, ShouldNotBeNil)
				So(got, ShouldBeNil)
			} else {
				So(err, ShouldBeNil)
				So(got, ShouldResemble, &tt.mockResponse)
			}
		})
	}
}
