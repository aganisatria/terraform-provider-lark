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

	if err != nil || response.Code != 0 {
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

	if err != nil || response.Code != 0 {
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
	if err != nil || response.Code != 0 {
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
	if err != nil || response.Code != 0 {
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
	if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
	if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
func GetUsersByIDAPI(ctx context.Context, client *LarkClient, userIds []string, idType UserIDType) (*UserInfoBatchGetResponse, error) {
	response := &UserInfoBatchGetResponse{}
	params := url.Values{}
	for _, id := range userIds {
		params.Add("user_ids", id)
	}
	path := fmt.Sprintf("%s/batch?%s&user_id_type=%s", USER_API, params.Encode(), string(idType))
	tflog.Info(ctx, "Getting Users by OpenID")

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to get users by OpenID", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Users by OpenID Retrieved", map[string]interface{}{
		"response": response.Data.Items,
	})
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/user/batch_get_id.
func GetUserIdByEmailsAPI(ctx context.Context, client *LarkClient, request UserInfoBatchGetRequest) (*UserInfoByEmailOrMobileBatchGetResponse, error) {
	batchResponse := &UserInfoByEmailOrMobileBatchGetResponse{}
	path := fmt.Sprintf("%s/batch_get_id", USER_API)
	tflog.Info(ctx, "Getting User ID by Emails")

	// Every request can only contain up to 50 emails.
	for i := 0; i < len(request.Emails); i += 50 {
		response := &UserInfoByEmailOrMobileBatchGetResponse{}
		batchEmails := request.Emails[i:min(i+50, len(request.Emails))]
		request := UserInfoBatchGetRequest{
			Emails: batchEmails,
		}

		err := client.DoTenantRequest(ctx, POST, path, request, response)
		if err != nil || response.Code != 0 {
			tflog.Error(ctx, "Failed to get user ID by emails", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
		batchResponse.Data.UserList = append(batchResponse.Data.UserList, response.Data.UserList...)
	}
	tflog.Info(ctx, "User ID Retrieved")
	return batchResponse, nil
}

// GROUP CHAT API.
// https://open.larksuite.com/document/server-docs/group/chat/create.
func GroupChatCreateAPI(ctx context.Context, client *LarkClient, request GroupChatCreateRequest) (*GroupChatCreateResponse, error) {
	response := &GroupChatCreateResponse{}
	tflog.Info(ctx, "Creating Group Chat")

	err := client.DoTenantRequest(ctx, POST, GROUP_CHAT_API, request, response)
	if err != nil || response.Code != 0 {
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
	if err != nil || response.Code != 0 {
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
	if err != nil || response.Code != 0 {
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
	if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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
		if err != nil || response.Code != 0 {
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

// ROLE API.
// https://open.larksuite.com/document/server-docs/contact-v3/functional_role/create.
func RoleCreateAPI(ctx context.Context, client *LarkClient, request RoleRequest) (*RoleCreateResponse, error) {
	response := &RoleCreateResponse{}
	tflog.Info(ctx, "Creating Role")

	err := client.DoTenantRequest(ctx, POST, ROLE_API, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to create role", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Role Created")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/functional_role/update.
func RoleUpdateAPI(ctx context.Context, client *LarkClient, roleID string, request RoleRequest) (*BaseResponse, error) {
	response := &BaseResponse{}
	tflog.Info(ctx, "Updating Role")
	path := fmt.Sprintf("%s/%s", ROLE_API, roleID)

	err := client.DoTenantRequest(ctx, PUT, path, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to update role", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Role Updated")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/functional_role/delete.
func RoleDeleteAPI(ctx context.Context, client *LarkClient, roleID string) (*BaseResponse, error) {
	response := &BaseResponse{}
	tflog.Info(ctx, "Deleting Role")
	path := fmt.Sprintf("%s/%s", ROLE_API, roleID)

	err := client.DoTenantRequest(ctx, DELETE, path, nil, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to delete role", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Role Deleted")
	return response, nil
}

// ROLE MEMBER API.
// https://open.larksuite.com/document/server-docs/contact-v3/functional_role-member/batch_create
func RoleMemberAddAPI(ctx context.Context, client *LarkClient, roleID string, request RoleMemberCreateRequest) (*RoleMemberCreateResponse, error) {
	response := &RoleMemberCreateResponse{}
	tflog.Info(ctx, "Adding Role Member")
	path := fmt.Sprintf("%s/%s/members/batch_create", ROLE_API, roleID)

	err := client.DoTenantRequest(ctx, POST, path, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to add role member", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Role Member Added")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/functional_role-member/list
func RoleMemberGetAPI(ctx context.Context, client *LarkClient, roleID string) (*RoleMemberGetResponse, error) {
	response := &RoleMemberGetResponse{}
	tflog.Info(ctx, "Getting Role Member")
	path := fmt.Sprintf("%s/%s/members", ROLE_API, roleID)

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to get role member", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Role Member Retrieved")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/functional_role-member/batch_delete
func RoleMemberDeleteAPI(ctx context.Context, client *LarkClient, roleID string, request RoleMemberDeleteRequest) (*RoleMemberDeleteResponse, error) {
	response := &RoleMemberDeleteResponse{}
	tflog.Info(ctx, "Deleting Role Member")
	path := fmt.Sprintf("%s/%s/members/batch_delete", ROLE_API, roleID)

	err := client.DoTenantRequest(ctx, PATCH, path, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to delete role member", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Role Member Deleted")
	return response, nil
}

// DEPARTMENT API.
// https://open.larksuite.com/document/server-docs/contact-v3/department/create.
func DepartmentCreateAPI(ctx context.Context, client *LarkClient, request DepartmentCreateRequest) (*DepartmentGetResponse, error) {
	response := &DepartmentGetResponse{}
	tflog.Info(ctx, "Sending Department Create Request")
	path := fmt.Sprintf("%s?department_id_type=open_department_id", DEPARTMENT_API)

	err := client.DoTenantRequest(ctx, POST, path, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to create department", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	tflog.Info(ctx, "Department Create API Response Received")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/department/update
func DepartmentUpdateAPI(ctx context.Context, client *LarkClient, departmentID string, request DepartmentUpdateRequest) (*DepartmentGetResponse, error) {
	response := &DepartmentGetResponse{}
	tflog.Info(ctx, "Updating Department")
	path := fmt.Sprintf("%s/%s?department_id_type=open_department_id", DEPARTMENT_API, departmentID)

	err := client.DoTenantRequest(ctx, PUT, path, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to update department", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Department Updated")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/department/get
func DepartmentGetByDepartmentIDAPI(ctx context.Context, client *LarkClient, departmentID string) (*DepartmentGetResponse, error) {
	return DepartmentGetAPI(ctx, client, departmentID, DEPARTMENT_ID)
}

func DepartmentGetByOpenDepartmentIDAPI(ctx context.Context, client *LarkClient, departmentID string) (*DepartmentGetResponse, error) {
	return DepartmentGetAPI(ctx, client, departmentID, OPEN_DEPARTMENT_ID)
}

func DepartmentGetAPI(ctx context.Context, client *LarkClient, departmentID string, departmentIDType DepartmentIDType) (*DepartmentGetResponse, error) {
	response := &DepartmentGetResponse{}
	tflog.Info(ctx, "Getting Department")

	path := fmt.Sprintf("%s/%s?department_id_type=%s", DEPARTMENT_API, departmentID, departmentIDType)

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to get department", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Department Retrieved")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/department/delete
func DepartmentDeleteAPI(ctx context.Context, client *LarkClient, departmentID string) (*DepartmentDeleteResponse, error) {
	response := &DepartmentDeleteResponse{}
	tflog.Info(ctx, "Deleting Department")

	path := fmt.Sprintf("%s/%s", DEPARTMENT_API, departmentID)
	err := client.DoTenantRequest(ctx, DELETE, path, nil, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to delete department", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Department Deleted")
	return response, nil
}

// https://open.larksuite.com/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/department/update_department_id?appId=cli_a718cd690138d02f.
func DepartmentUpdateIDAPI(ctx context.Context, client *LarkClient, parentDepartmentID string, request DepartmentUpdateIDRequest) (*BaseResponse, error) {
	response := &BaseResponse{}
	tflog.Info(ctx, "Updating Department ID")
	path := fmt.Sprintf("%s/%s/update_department_id?department_id_type=open_department_id", DEPARTMENT_API, parentDepartmentID)

	err := client.DoTenantRequest(ctx, PATCH, path, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to update department ID", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Department ID Updated")
	return response, nil
}

// WORKFORCE TYPE API.
// https://open.larksuite.com/document/server-docs/contact-v3/employee_type_enum/create.
func WorkforceTypeCreateAPI(ctx context.Context, client *LarkClient, request WorkforceTypeRequest) (*WorkforceTypeResponse, error) {
	response := &WorkforceTypeResponse{}
	tflog.Info(ctx, "Creating Workforce Type")

	err := client.DoTenantRequest(ctx, POST, WORKFORCE_TYPE_API, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to create workforce type", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Workforce Type Created")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/employee_type_enum/update.
func WorkforceTypeUpdateAPI(ctx context.Context, client *LarkClient, enumID string, request WorkforceTypeRequest) (*WorkforceTypeResponse, error) {
	response := &WorkforceTypeResponse{}
	tflog.Info(ctx, "Updating Workforce Type")
	path := fmt.Sprintf("%s/%s", WORKFORCE_TYPE_API, enumID)

	err := client.DoTenantRequest(ctx, PUT, path, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to update workforce type", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Workforce Type Updated")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/employee_type_enum/list.
func WorkforceTypeGetAllAPI(ctx context.Context, client *LarkClient) (*WorkforceTypeGetResponse, error) {
	response := &WorkforceTypeGetResponse{}
	tflog.Info(ctx, "Getting Workforce Type")

	err := client.DoTenantRequest(ctx, GET, WORKFORCE_TYPE_API, nil, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to get workforce type", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Workforce Type Retrieved")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/contact-v3/employee_type_enum/delete.
func WorkforceTypeDeleteAPI(ctx context.Context, client *LarkClient, enumID string) (*BaseResponse, error) {
	response := &BaseResponse{}
	tflog.Info(ctx, "Deleting Workforce Type")
	path := fmt.Sprintf("%s/%s", WORKFORCE_TYPE_API, enumID)

	err := client.DoTenantRequest(ctx, DELETE, path, nil, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to delete workforce type", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	tflog.Info(ctx, "Workforce Type Deleted")
	return response, nil
}

// DOCS SPACE API.

// DOCS SPACE FOLDER API.
// https://open.larksuite.com/document/server-docs/docs/drive-v1/folder/get-root-folder-meta.
func RootFolderMetaGetAPI(ctx context.Context, client *LarkClient) (*RootFolderMetaGetResponse, error) {
	response := &RootFolderMetaGetResponse{}
	tflog.Info(ctx, "Getting Root Folder Meta")
	path := fmt.Sprintf("%s/meta", EXPLORER_ROOT_FOLDER_API)

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		tflog.Error(ctx, "Failed to get root folder meta", map[string]interface{}{"response": response})
		return nil, fmt.Errorf("failed to get root folder meta: %s", response.Msg)
	}
	tflog.Info(ctx, "Root Folder Meta Retrieved")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/docs/drive-v1/folder/get-folder-meta.
func FolderMetaGetAPI(ctx context.Context, client *LarkClient, folderToken string) (*FolderMetaGetResponse, error) {
	response := &FolderMetaGetResponse{}
	tflog.Info(ctx, "Getting Folder Meta")
	path := fmt.Sprintf("%s/%s/meta", EXPLORER_FOLDER_API, folderToken)

	err := client.DoTenantRequest(ctx, GET, path, nil, response)
	if err != nil {
		return nil, err
	}

	if response.Code != 0 {
		tflog.Error(ctx, "Failed to get folder meta", map[string]interface{}{
			"response": response,
		})
		return nil, fmt.Errorf("failed to get folder meta: %s", response.Msg)
	}
	tflog.Info(ctx, "Folder Meta Retrieved")
	return response, nil
}

// https://open.larksuite.com/document/server-docs/docs/drive-v1/folder/list
func FolderChildrenListAPI(ctx context.Context, client *LarkClient, folderToken string) (*FolderChildrenListResponse, error) {
	var allChildren []FileChild
	pageSize := 200
	pageToken := ""

	tflog.Info(ctx, "Listing folder children", map[string]interface{}{
		"folder_token": folderToken,
	})

	for {
		response := &FolderChildrenListResponse{}
		path := fmt.Sprintf("%s?folder_token=%s&page_size=%d", DOCS_FILE_API, folderToken, pageSize)

		if pageToken != "" {
			path += fmt.Sprintf("&page_token=%s", pageToken)
		}

		err := client.DoTenantRequest(ctx, GET, path, nil, response)
		if err != nil {
			tflog.Error(ctx, "Failed to list folder children page", map[string]interface{}{"error": err})
			return nil, err
		}
		if response.Code != 0 {
			err = fmt.Errorf("API error when listing folder children: %s", response.Msg)
			tflog.Error(ctx, err.Error(), map[string]interface{}{"response": response})
			return nil, err
		}

		allChildren = append(allChildren, response.Data.Files...)

		if !response.Data.HasMore || response.Data.NextPageToken == "" {
			break
		}

		pageToken = response.Data.NextPageToken
	}

	finalResponse := &FolderChildrenListResponse{
		BaseResponse: BaseResponse{
			Code: 0,
			Msg:  "success",
		},
	}
	finalResponse.Data.Files = allChildren
	finalResponse.Data.NextPageToken = ""
	finalResponse.Data.HasMore = false

	tflog.Info(ctx, "All folder children retrieved", map[string]interface{}{
		"total_children": len(allChildren),
	})

	return finalResponse, nil
}

// https://open.larksuite.com/document/server-docs/docs/drive-v1/folder/create_folder.
func FolderCreateAPI(ctx context.Context, client *LarkClient, request FolderCreateRequest) (*FolderCreateResponse, error) {
	response := &FolderCreateResponse{}
	tflog.Info(ctx, "Creating Folder", map[string]interface{}{
		"folder_token": request.FolderToken,
		"name":         request.Name,
	})
	path := fmt.Sprintf("%s/create_folder", DOCS_FILE_API)

	err := client.DoTenantRequest(ctx, POST, path, request, response)
	if err != nil {
		tflog.Error(ctx, "Failed to create folder", map[string]interface{}{
			"error": err,
		})
		return nil, err
	}

	if response.Code != 0 {
		tflog.Error(ctx, "API returned an error when creating folder", map[string]interface{}{"response": response})
		return nil, fmt.Errorf("API error when creating folder: %s", response.Msg)
	}
	tflog.Info(ctx, "Folder Created successfully", map[string]interface{}{"new_folder_token": response.Data.Token})
	return response, nil
}

// DOCS SPACE FILE API.
// https://open.larksuite.com/document/server-docs/docs/drive-v1/file/move.
func FileMoveAPI(ctx context.Context, client *LarkClient, fileToken string, request FileMoveRequest) (*FileTaskResponse, error) {
	response := &FileTaskResponse{}
	tflog.Info(ctx, "Moving File", map[string]interface{}{
		"file_token":      fileToken,
		"destination_dir": request.FolderToken,
	})
	path := fmt.Sprintf("%s/%s/move", DOCS_FILE_API, fileToken)

	err := client.DoTenantRequest(ctx, POST, path, request, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "Failed to move file", map[string]interface{}{"error": err, "response": response})
		return nil, fmt.Errorf("API error when moving file: %s", response.Msg)
	}

	tflog.Info(ctx, "File Move task created successfully", map[string]interface{}{"task_id": response.Data.TaskID})
	return response, nil
}

// https://open.larksuite.com/document/server-docs/docs/drive-v1/file/delete.
func FileDeleteAPI(ctx context.Context, client *LarkClient, fileToken, fileType string) (*FileTaskResponse, error) {
	response := &FileTaskResponse{}
	tflog.Info(ctx, "Deleting File (moving to trash)", map[string]interface{}{
		"file_token": fileToken,
		"file_type":  fileType,
	})
	path := fmt.Sprintf("%s/%s?type=%s", DOCS_FILE_API, fileToken, fileType)

	err := client.DoTenantRequest(ctx, DELETE, path, nil, response)
	if err != nil || response.Code != 0 {
		tflog.Error(ctx, "API returned an error when deleting file", map[string]interface{}{"error": err, "response": response})
		return nil, fmt.Errorf("API error when deleting file: %s", response.Msg)
	}
	tflog.Info(ctx, "File Delete task created successfully", map[string]interface{}{"task_id": response.Data.TaskID})
	return response, nil
}
