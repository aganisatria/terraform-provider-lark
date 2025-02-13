// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ACCESS TOKEN API.
// https://open.larksuite.com/document/server-docs/getting-started/api-access-token/auth-v3/tenant_access_token_internal.
func GetAccessTokenAPI(appID, appSecret string, baseDelay int, retryCount int) (string, string, error) {
	tflog.Info(context.Background(), "Getting access token from Lark API")

	client := NewLarkClient("", "", appID, baseDelay, retryCount)

	requestBody := AccessTokenRequest{
		AppID:     appID,
		AppSecret: appSecret,
	}

	var response AccessTokenResponse

	err := client.DoInitializeRequest(context.Background(), POST, AUTH_API, requestBody, &response)

	if err != nil {
		fmt.Println("err", err)
		return "", "", fmt.Errorf("failed to get access token: %w", err)
	}
	fmt.Println("response", response)

	tflog.Info(context.Background(), "Access token retrieved successfully")

	return response.TenantAccessToken, response.AppAccessToken, nil
}

// USERGROUP API.
// https://open.larksuite.com/document/server-docs/contact-v3/group/create.
func UsergroupCreateAPI(ctx context.Context, client *LarkClient, request UsergroupCreateRequest) (*UsergroupCreateResponse, error) {
	response := &UsergroupCreateResponse{}
	err := client.DoTenantRequest(ctx, POST, USERGROUP_API, request, response)
	tflog.Info(ctx, "Creating User Group")

	if err != nil {
		tflog.Error(ctx, "Failed to create user group", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "User Group Created")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/group/get.
func UsergroupGetAPI(ctx context.Context, client *LarkClient, groupID string) (*UsergroupGetResponse, error) {
	response := &UsergroupGetResponse{}
	path := fmt.Sprintf("%s/%s", USERGROUP_API, groupID)
	tflog.Info(ctx, "Getting User Group")

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil {
		tflog.Error(ctx, "Failed to get user group", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "User Group Retrieved")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/group/patch.
func UsergroupUpdateAPI(ctx context.Context, client *LarkClient, groupID string, request UsergroupUpdateRequest) (*BaseResponse, error) {
	response := &BaseResponse{}
	path := fmt.Sprintf("%s/%s", USERGROUP_API, groupID)
	tflog.Info(ctx, "Updating User Group")

	err := client.DoTenantRequest(ctx, PATCH, path, request, response)
	if err != nil {
		tflog.Error(ctx, "Failed to update user group", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "User Group Updated")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/group/delete.
func UsergroupDeleteAPI(ctx context.Context, client *LarkClient, groupID string) (*BaseResponse, error) {
	response := &BaseResponse{}
	path := fmt.Sprintf("%s/%s", USERGROUP_API, groupID)
	tflog.Info(ctx, "Deleting User Group")

	err := client.DoTenantRequest(ctx, DELETE, path, nil, response)
	if err != nil {
		tflog.Error(ctx, "Failed to delete user group", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "User Group Deleted")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/group/simplelist.
func UsergroupListAPI(ctx context.Context, client *LarkClient) (*UsergroupListResponse, error) {
	var allGroups []Group
	pageSize := 100
	pageToken := ""

	for {
		response := &UsergroupListResponse{}
		path := fmt.Sprintf("%s/simplelist?page_size=%d", USERGROUP_API, pageSize)

		if pageToken != "" {
			path += fmt.Sprintf("&page_token=%s", pageToken)
		}

		tflog.Info(ctx, "Getting User Groups", map[string]interface{}{
			"page_size":  pageSize,
			"page_token": pageToken,
		})

		err := client.DoTenantRequest(ctx, GET, path, nil, response)
		if err != nil {
			tflog.Error(ctx, "Failed to get user groups", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		allGroups = append(allGroups, response.Data.GroupList...)

		if response.Data.PageToken == "" {
			break
		}

		pageToken = response.Data.PageToken
	}

	finalResponse := &UsergroupListResponse{
		BaseResponse: BaseResponse{
			Code: 0,
			Msg:  "success",
		},
		Data: struct {
			GroupList []Group `json:"grouplist"`
			PageToken string  `json:"page_token"`
			HasMore   bool    `json:"has_more"`
		}{
			GroupList: allGroups,
			PageToken: "",
			HasMore:   false,
		},
	}

	tflog.Info(ctx, "All User Groups Retrieved", map[string]interface{}{
		"total_groups": len(allGroups),
	})

	return finalResponse, nil
}

// USERGROUP MEMBER API.
// https://open.larksuite.com/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/group-member/batch_add.
func UsergroupMemberAddAPI(ctx context.Context, client *LarkClient, groupID string, request UsergroupMemberAddRequest) (*UsergroupMemberAddResponse, error) {
	response := &UsergroupMemberAddResponse{}
	path := fmt.Sprintf("%s/%s/member/batch_add", USERGROUP_API, groupID)
	tflog.Info(ctx, "Adding User Group Member")

	memberIDs := []string{}
	for _, member := range request.Members {
		memberIDs = append(memberIDs, member.MemberID)
	}

	// Each request will be limited to 100 members.
	for i := 0; i < len(memberIDs); i += 100 {
		batchMemberIDs := memberIDs[i:min(i+100, len(memberIDs))]

		currentTurnMemberIDs := []UsergroupMember{}
		for _, memberID := range batchMemberIDs {
			currentTurnMemberIDs = append(currentTurnMemberIDs, UsergroupMember{
				MemberID:     memberID,
				MemberType:   "user",
				MemberIDType: "open_id",
			})
		}

		request := UsergroupMemberAddRequest{
			Members: currentTurnMemberIDs,
		}
		err := client.DoTenantRequest(ctx, POST, path, request, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add user group member", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	}

	tflog.Info(ctx, "User Group Member Added")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/group/group-member/simplelist.
func UsergroupMemberGetByMemberTypeAPI(ctx context.Context, client *LarkClient, groupID string, memberType string) (*UsergroupMemberGetResponse, error) {
	response := &UsergroupMemberGetResponse{}
	path := fmt.Sprintf("%s/%s/member/simplelist?member_type=%s", USERGROUP_API, groupID, memberType)
	tflog.Info(ctx, "Getting User Group Member by Member Type")

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil {
		tflog.Error(ctx, "Failed to get user group member", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "User Group Member Retrieved")
	return response, nil
}

// https://open.larksuite.com/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/group-member/batch_remove.
func UsergroupMemberRemoveAPI(ctx context.Context, client *LarkClient, groupID string, request UsergroupMemberRemoveRequest) (*BaseResponse, error) {
	response := &BaseResponse{}
	path := fmt.Sprintf("%s/%s/member/batch_remove", USERGROUP_API, groupID)
	tflog.Info(ctx, "Removing User Group Member")

	memberIDs := []string{}
	for _, member := range request.Members {
		memberIDs = append(memberIDs, member.MemberID)
	}

	// Each request will be limited to 100 members.
	for i := 0; i < len(memberIDs); i += 100 {
		batchMemberIDs := memberIDs[i:min(i+100, len(memberIDs))]
		currentTurnMemberIDs := []UsergroupMember{}
		for _, memberID := range batchMemberIDs {
			currentTurnMemberIDs = append(currentTurnMemberIDs, UsergroupMember{
				MemberID:     memberID,
				MemberType:   "user",
				MemberIDType: "open_id",
			})
		}
		request := UsergroupMemberRemoveRequest{
			Members: currentTurnMemberIDs,
		}
		err := client.DoTenantRequest(ctx, POST, path, request, response)
		if err != nil {
			tflog.Error(ctx, "Failed to remove user group member", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	}

	tflog.Info(ctx, "User Group Member Removed")
	return response, nil
}

// USER API.
// https://open.larksuite.com/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/user/batch?appId=cli_a718cd690138d02f.
func GetUsersByOpenIDAPI(ctx context.Context, client *LarkClient, userIds []string) (*UserInfoBatchGetResponse, error) {
	response := &UserInfoBatchGetResponse{}
	params := url.Values{}
	for _, id := range userIds {
		params.Add("user_ids", id)
	}
	path := fmt.Sprintf("%s/batch?%s", USER_API, params.Encode())
	tflog.Info(ctx, "Getting Users by OpenID")

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil {
		tflog.Error(ctx, "Failed to get users by OpenID", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Users by OpenID Retrieved")
	return response, nil
}

// GROUP CHAT API.
// https://open.larksuite.com/document/server-docs/group/chat/create.
func GroupChatCreateAPI(ctx context.Context, client *LarkClient, request GroupChatCreateRequest) (*GroupChatCreateResponse, error) {
	response := &GroupChatCreateResponse{}
	tflog.Info(ctx, "Creating Group Chat")

	err := client.DoTenantRequest(ctx, POST, GROUP_CHAT_API, request, response)
	if err != nil {
		tflog.Error(ctx, "Failed to create group chat", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Group Chat Created")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/group/chat/delete.
func GroupChatDeleteAPI(ctx context.Context, client *LarkClient, chatID string) (*BaseResponse, error) {
	response := &BaseResponse{}
	tflog.Info(ctx, "Deleting Group Chat")
	path := fmt.Sprintf("%s/%s", GROUP_CHAT_API, chatID)

	err := client.DoTenantRequest(ctx, DELETE, path, nil, response)
	tflog.Info(ctx, "Deleting Group Chat", map[string]interface{}{
		"path": path,
	})
	if err != nil {
		tflog.Error(ctx, "Failed to delete group chat", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Group Chat Deleted")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/group/chat/update.
func GroupChatUpdateAPI(ctx context.Context, client *LarkClient, chatID string, request GroupChatUpdateRequest) (*BaseResponse, error) {
	response := &BaseResponse{}
	tflog.Info(ctx, "Updating Group Chat")
	path := fmt.Sprintf("%s/%s", GROUP_CHAT_API, chatID)

	err := client.DoTenantRequest(ctx, PUT, path, request, response)
	if err != nil {
		tflog.Error(ctx, "Failed to update group chat", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Group Chat Updated")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/group/chat/update-2.
func GroupChatSpeechScopesUpdateAPI(ctx context.Context, client *LarkClient, chatID string, request GroupChatSpeechScopesUpdateRequest) (*BaseResponse, error) {
	response := &BaseResponse{}
	tflog.Info(ctx, "Updating Group Chat Speech Scopes")
	path := fmt.Sprintf("%s/%s/speech_scopes", GROUP_CHAT_API, chatID)

	err := client.DoTenantRequest(ctx, PATCH, path, request, response)
	if err != nil {
		tflog.Error(ctx, "Failed to update group chat speech scopes", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Group Chat Speech Scopes Updated")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/group/chat/get.
func GroupChatGetAPI(ctx context.Context, client *LarkClient, chatID string) (*GroupChatGetResponse, error) {
	response := &GroupChatGetResponse{}
	tflog.Info(ctx, "Getting Group Chat")
	path := fmt.Sprintf("%s/%s", GROUP_CHAT_API, chatID)

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil {
		tflog.Error(ctx, "Failed to get group chat", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Group Chat Retrieved")
	return response, nil
}

// GROUP MEMBER API.
// https://open.larksuite.com/document/server-docs/group/chat-member/get.
func GroupMemberGetAPI(ctx context.Context, client *LarkClient, chatID string) (*GetUserGroupMemberResponse, error) {
	response := &GetUserGroupMemberResponse{}
	tflog.Info(ctx, "Getting Group Member")
	path := fmt.Sprintf("%s/%s/members", GROUP_CHAT_API, chatID)

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil {
		tflog.Error(ctx, "Failed to get group member", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Group Member Retrieved")
	return response, nil
}

// GROUP ADMINISTRATOR API.
// https://open.larksuite.com/document/server-docs/group/chat-member/add_managers.
func GroupChatAdministratorAddAPI(ctx context.Context, client *LarkClient, chatID string, GroupChatType GroupChatType, request GroupChatAdministratorRequest) (*GroupChatAdministratorResponse, error) {
	response := &GroupChatAdministratorResponse{}
	tflog.Info(ctx, "Adding Group Administrator")
	path := fmt.Sprintf("%s/%s/managers/add_managers", GROUP_CHAT_API, chatID)

	botList, personList, err := splitUserAndBotList(request)
	if err != nil {
		return nil, err
	}

	// For Common Groups, up to 10 administrators can be specified.
	// For Super Lark Groups, up to 20 administrators can be specified.

	if GroupChatType == GroupChatTypeCommon && len(personList)+len(botList) > 10 {
		return nil, errors.New("invalid administrator count, max 10 administrators for common group")
	}

	if GroupChatType == GroupChatTypeSuperLarge && len(personList)+len(botList) > 20 {
		return nil, errors.New("invalid administrator count, max 20 administrators for super large group")
	}

	tflog.Info(ctx, "Adding Group Administrator", map[string]interface{}{
		"bot_count":  len(botList),
		"user_count": len(personList),
	})

	// Up to 50 users or 5 bots can be specified for each request.
	for i := 0; i < len(botList); i += 5 {
		end := i + 5
		if end > len(botList) {
			end = len(botList)
		}

		batchRequest := GroupChatAdministratorRequest{
			ManagerIDs: botList[i:end],
		}

		path = fmt.Sprintf("%s?member_id_type=app_id", path)

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add bot administrators", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	}

	for i := 0; i < len(personList); i += 50 {
		end := i + 50
		if end > len(personList) {
			end = len(personList)
		}

		batchRequest := GroupChatAdministratorRequest{
			ManagerIDs: personList[i:end],
		}

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add user administrators", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	}

	tflog.Info(ctx, "Group Administrator Added")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/group/chat-member/delete_managers.
func GroupChatAdministratorDeleteAPI(ctx context.Context, client *LarkClient, chatID string, request GroupChatAdministratorRequest) (*GroupChatAdministratorResponse, error) {
	response := &GroupChatAdministratorResponse{}
	tflog.Info(ctx, "Deleting Group Administrator")
	path := fmt.Sprintf("%s/%s/managers/delete_managers", GROUP_CHAT_API, chatID)

	botList, personList, err := splitUserAndBotList(request)
	if err != nil {
		return nil, err
	}

	// Up to 50 users or 5 bots can be specified for each request.
	for i := 0; i < len(botList); i += 5 {
		end := i + 5
		if end > len(botList) {
			end = len(botList)
		}

		batchRequest := GroupChatAdministratorRequest{
			ManagerIDs: botList[i:end],
		}

		path = fmt.Sprintf("%s?member_id_type=app_id", path)

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add bot administrators", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	}

	for i := 0; i < len(personList); i += 50 {
		end := i + 50
		if end > len(personList) {
			end = len(personList)
		}

		batchRequest := GroupChatAdministratorRequest{
			ManagerIDs: personList[i:end],
		}

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add user administrators", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	}

	tflog.Info(ctx, "Group Administrator Deleted")
	return response, nil
}
