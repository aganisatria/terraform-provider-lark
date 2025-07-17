// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider_acceptance_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	. "github.com/bytedance/mockey"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDocsSpaceFolderResource(t *testing.T) {
	Mock(common.GetAccessTokenAPI).Return("test_tenant_access_token", "test_app_access_token", nil).Build()

	Mock(common.RootFolderMetaGetAPI).Return(&common.RootFolderMetaGetResponse{
		Data: common.RootFolderMetaData{
			Token: "root_folder_token",
			ID:    "root_id",
		},
	}, nil).Build()

	createCallCount := 0
	Mock(common.FolderCreateAPI).To(func(ctx context.Context, client *common.LarkClient, req common.FolderCreateRequest) (*common.FolderCreateResponse, error) {
		createCallCount++
		token := ""
		switch req.Name {
		case "Test Folder":
			token = "test_folder_token"
		case "Updated Folder Name":
			token = "renamed_folder_token"
		case "Another Folder":
			token = "another_folder_token"
		default:
			token = fmt.Sprintf("new_folder_token_%d", createCallCount)
		}

		return &common.FolderCreateResponse{
			Data: common.FolderCreateResponseData{
				Token: token,
				ID:    token,
				Name:  req.Name,
				URL:   "https://example.com/" + token,
			},
		}, nil
	}).Build()

	folderStates := make(map[string]struct{ Name, ParentID string })
	folderStates["root_folder_token"] = struct{ Name, ParentID string }{"Root", ""}
	folderStates["test_folder_token"] = struct{ Name, ParentID string }{"Test Folder", "root_folder_token"}
	folderStates["renamed_folder_token"] = struct{ Name, ParentID string }{"Updated Folder Name", "root_folder_token"}
	folderStates["another_folder_token"] = struct{ Name, ParentID string }{"Another Folder", "root_folder_token"}

	Mock(common.FolderMetaGetAPI).To(func(ctx context.Context, client *common.LarkClient, folderToken string) (*common.FolderMetaGetResponse, error) {
		var name, parentID string
		switch folderToken {
		case "test_folder_token":
			name = "Test Folder"
			parentID = "root_folder_token"
		case "renamed_folder_token":
			name = "Updated Folder Name"
			parentID = "root_folder_token"
			if state, ok := folderStates[folderToken]; ok {
				parentID = state.ParentID
			}
		case "another_folder_token":
			name = "Another Folder"
			parentID = "root_folder_token"
		default:
			return nil, fmt.Errorf("folder not found: %s", folderToken)
		}

		return &common.FolderMetaGetResponse{
			Data: common.FolderMetaData{
				ID:       folderToken,
				Name:     name,
				Token:    folderToken,
				ParentID: parentID,
			},
		}, nil
	}).Build()

	Mock(common.FolderChildrenListAPI).Return(&common.FolderChildrenListResponse{
		Data: common.FolderChildrenListData{
			Files:         []common.FileChild{},
			NextPageToken: "",
			HasMore:       false,
		},
	}, nil).Build()

	Mock(common.FileMoveAPI).To(func(ctx context.Context, client *common.LarkClient, fileToken string, req common.FileMoveRequest) (*common.FileTaskResponse, error) {
		if state, ok := folderStates[fileToken]; ok {
			state.ParentID = req.FolderToken
			folderStates[fileToken] = state
		}

		return &common.FileTaskResponse{
			Data: common.FileTaskResponseData{
				TaskID: "move_task_123",
			},
		}, nil
	}).Build()

	Mock(common.FileDeleteAPI).Return(&common.FileTaskResponse{
		Data: common.FileTaskResponseData{
			TaskID: "delete_task_123",
		},
	}, nil).Build()

	defer UnPatchAll()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
				resource "lark_docs_space_folder" "test" {
					name = "Test Folder"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_docs_space_folder.test", "name", "Test Folder"),
					resource.TestCheckResourceAttr("lark_docs_space_folder.test", "parent_folder_token", "root_folder_token"),
					resource.TestCheckResourceAttr("lark_docs_space_folder.test", "token", "test_folder_token"),
					resource.TestCheckResourceAttrSet("lark_docs_space_folder.test", "id"),
				),
			},
			{
				Config: providerConfig + `
				resource "lark_docs_space_folder" "test" {
					name = "Updated Folder Name"
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_docs_space_folder.test", "name", "Updated Folder Name"),
					resource.TestCheckResourceAttr("lark_docs_space_folder.test", "token", "renamed_folder_token"),
				),
			},
			{
				Config: providerConfig + `
				resource "lark_docs_space_folder" "another_folder" {
				  name = "Another Folder"
				}

				resource "lark_docs_space_folder" "test" {
					name                = "Updated Folder Name"
					parent_folder_token = lark_docs_space_folder.another_folder.token
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("lark_docs_space_folder.test", "name", "Updated Folder Name"),
					resource.TestCheckResourceAttr("lark_docs_space_folder.test", "parent_folder_token", "another_folder_token"),
				),
			},
		},
	})
}
