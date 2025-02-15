// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

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

// GROUP CHAT MEMBER API.
// https://open.larksuite.com/document/server-docs/group/chat-member/get.
func GroupChatMemberGetAPI(ctx context.Context, client *LarkClient, chatID string) (*GroupChatMemberGetResponse, error) {
	var allMembers []ListMember
	pageSize := 100
	pageToken := ""

	for {
		response := &GroupChatMemberGetResponse{}
		path := fmt.Sprintf("%s/%s/members?page_size=%d", GROUP_CHAT_API, chatID, pageSize)

		if pageToken != "" {
			path += fmt.Sprintf("&page_token=%s", pageToken)
		}

		err := client.DoTenantRequest(ctx, GET, path, nil, response)
		if err != nil {
			tflog.Error(ctx, "Failed to get user groups", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		allMembers = append(allMembers, response.Data.Items...)

		if response.Data.PageToken == "" {
			break
		}

		pageToken = response.Data.PageToken
	}

	finalResponse := &GroupChatMemberGetResponse{
		BaseResponse: BaseResponse{
			Code: 0,
			Msg:  "success",
		},
		Data: struct {
			Items       []ListMember `json:"items"`
			PageToken   string       `json:"page_token"`
			HasMore     bool         `json:"has_more"`
			MemberTotal int64        `json:"member_total"`
		}{
			Items:       allMembers,
			PageToken:   "",
			HasMore:     false,
			MemberTotal: int64(len(allMembers)),
		},
	}

	tflog.Info(ctx, "All User Groups Retrieved", map[string]interface{}{
		"total_members": len(allMembers),
	})

	return finalResponse, nil
}

// https://open.larksuite.com/document/server-docs/group/chat-member/create.
func GroupChatMemberAddAPI(ctx context.Context, client *LarkClient, chatID string, request GroupChatMemberRequest) (*GroupChatMemberAddResponse, error) {
	fullResponse := GroupChatMemberAddResponse{}
	tflog.Info(ctx, "Adding Group Member")
	// there are 3 succeed_type: 0, 1, 2.
	// 0: When there is a separated ID, other available IDs will be pulled into the group chat, and a successful response will be returned.
	// 1: APull all the available IDs in the parameters into the group chat, return the successful response of pulling the group, and show the remaining unavailable IDs and reasons.
	// 2: As long as there is any unavailable ID in the parameter, the group will fail, an error response will be returned, and the unavailable ID will be displayed.
	path := fmt.Sprintf("%s/%s/members?succeed_type=2", GROUP_CHAT_API, chatID)

	botList, personList, err := splitUserAndBotList(request.IDList)
	if err != nil {
		return nil, err
	}

	// Up to 50 users or 5 bots can be specified for each request.
	for i := 0; i < len(botList); i += 5 {
		end := i + 5
		if end > len(botList) {
			end = len(botList)
		}

		batchRequest := GroupChatMemberRequest{
			IDList: botList[i:end],
		}

		response := &GroupChatMemberAddResponse{}

		path = fmt.Sprintf("%s?member_id_type=app_id", path)

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add bot members", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		fullResponse.Data.InvalidIDList = append(fullResponse.Data.InvalidIDList, response.Data.InvalidIDList...)
	}

	for i := 0; i < len(personList); i += 50 {
		end := i + 50
		if end > len(personList) {
			end = len(personList)
		}

		batchRequest := GroupChatMemberRequest{
			IDList: personList[i:end],
		}

		response := &GroupChatMemberAddResponse{}

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add user members", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		fullResponse.Data.InvalidIDList = append(fullResponse.Data.InvalidIDList, response.Data.InvalidIDList...)
	}

	tflog.Info(ctx, "Group Member Added")
	return &fullResponse, nil
}

// https://open.larksuite.com/document/server-docs/group/chat-member/delete.
func GroupChatMemberDeleteAPI(ctx context.Context, client *LarkClient, chatID string, request GroupChatMemberRequest) (*GroupChatMemberRemoveResponse, error) {
	fullResponse := GroupChatMemberRemoveResponse{}
	tflog.Info(ctx, "Deleting Group Members")
	path := fmt.Sprintf("%s/%s/members", GROUP_CHAT_API, chatID)

	botList, personList, err := splitUserAndBotList(request.IDList)
	if err != nil {
		return nil, err
	}

	// Up to 50 users or 5 bots can be specified for each request.
	for i := 0; i < len(botList); i += 5 {
		end := i + 5
		if end > len(botList) {
			end = len(botList)
		}

		batchRequest := GroupChatMemberRequest{
			IDList: botList[i:end],
		}

		path = fmt.Sprintf("%s?member_id_type=app_id", path)

		response := &GroupChatMemberRemoveResponse{}

		err := client.DoTenantRequest(ctx, DELETE, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to delete bot members", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		fullResponse.Data.InvalidIDList = append(fullResponse.Data.InvalidIDList, response.Data.InvalidIDList...)
	}

	for i := 0; i < len(personList); i += 50 {
		end := i + 50
		if end > len(personList) {
			end = len(personList)
		}

		batchRequest := GroupChatMemberRequest{
			IDList: personList[i:end],
		}

		response := &GroupChatMemberRemoveResponse{}

		err := client.DoTenantRequest(ctx, DELETE, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to delete user members", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		fullResponse.Data.InvalidIDList = append(fullResponse.Data.InvalidIDList, response.Data.InvalidIDList...)
	}

	if len(fullResponse.Data.InvalidIDList) > 0 {
		return nil, errors.New("invalid ID list, " + strings.Join(fullResponse.Data.InvalidIDList, ", "))
	}

	tflog.Info(ctx, "Group Member Deleted")
	return &fullResponse, nil
}

// GROUP ADMINISTRATOR API.
// https://open.larksuite.com/document/server-docs/group/chat-member/add_managers.
func GroupChatAdministratorAddAPI(ctx context.Context, client *LarkClient, chatID string, request GroupChatAdministratorRequest) (*GroupChatAdministratorResponse, error) {
	fullResponse := GroupChatAdministratorResponse{}
	tflog.Info(ctx, "Adding Group Administrator")
	path := fmt.Sprintf("%s/%s/managers/add_managers", GROUP_CHAT_API, chatID)

	botList, personList, err := splitUserAndBotList(request.ManagerIDs)
	if err != nil {
		return nil, err
	}

	// For Common Groups, up to 10 administrators can be specified.
	if len(personList)+len(botList) > 10 {
		return nil, errors.New("invalid administrator count, max 10 administrators for common group")
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

		response := &GroupChatAdministratorResponse{}

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add bot administrators", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		fullResponse.Data.ChatManagers = append(fullResponse.Data.ChatManagers, response.Data.ChatManagers...)
		fullResponse.Data.ChatBotManagers = append(fullResponse.Data.ChatBotManagers, response.Data.ChatBotManagers...)
	}

	for i := 0; i < len(personList); i += 50 {
		end := i + 50
		if end > len(personList) {
			end = len(personList)
		}

		batchRequest := GroupChatAdministratorRequest{
			ManagerIDs: personList[i:end],
		}

		response := &GroupChatAdministratorResponse{}

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add user administrators", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		fullResponse.Data.ChatManagers = append(fullResponse.Data.ChatManagers, response.Data.ChatManagers...)
		fullResponse.Data.ChatBotManagers = append(fullResponse.Data.ChatBotManagers, response.Data.ChatBotManagers...)
	}

	tflog.Info(ctx, "Group Administrator Added")
	return &fullResponse, nil
}

// https://open.larksuite.com/document/server-docs/group/chat-member/delete_managers.
func GroupChatAdministratorDeleteAPI(ctx context.Context, client *LarkClient, chatID string, request GroupChatAdministratorRequest) (*GroupChatAdministratorResponse, error) {
	fullResponse := GroupChatAdministratorResponse{}
	tflog.Info(ctx, "Deleting Group Administrator")
	path := fmt.Sprintf("%s/%s/managers/delete_managers", GROUP_CHAT_API, chatID)

	botList, personList, err := splitUserAndBotList(request.ManagerIDs)
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

		response := &GroupChatAdministratorResponse{}

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add bot administrators", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		fullResponse.Data.ChatManagers = append(fullResponse.Data.ChatManagers, response.Data.ChatManagers...)
		fullResponse.Data.ChatBotManagers = append(fullResponse.Data.ChatBotManagers, response.Data.ChatBotManagers...)
	}

	for i := 0; i < len(personList); i += 50 {
		end := i + 50
		if end > len(personList) {
			end = len(personList)
		}

		batchRequest := GroupChatAdministratorRequest{
			ManagerIDs: personList[i:end],
		}

		response := &GroupChatAdministratorResponse{}

		err := client.DoTenantRequest(ctx, POST, path, batchRequest, response)
		if err != nil {
			tflog.Error(ctx, "Failed to add user administrators", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		fullResponse.Data.ChatManagers = append(fullResponse.Data.ChatManagers, response.Data.ChatManagers...)
		fullResponse.Data.ChatBotManagers = append(fullResponse.Data.ChatBotManagers, response.Data.ChatBotManagers...)
	}

	tflog.Info(ctx, "Group Administrator Deleted")
	return &fullResponse, nil
}
