// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

// FolderMetaData represents folder metadata.
type FolderMetaData struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Token     string `json:"token"`
	CreateUid string `json:"create_uid"`
	EditUid   string `json:"edit_uid"`
	ParentID  string `json:"parent_id"`
	OwnUid    string `json:"own_uid"`
}

// FolderMetaGetResponse represents the response for getting folder metadata.
type FolderMetaGetResponse struct {
	BaseResponse
	Data FolderMetaData `json:"data"`
}

// RootFolderMetaData represents root folder metadata.
type RootFolderMetaData struct {
	Token  string `json:"token"`
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

// RootFolderMetaGetResponse represents the response for getting the root folder metadata.
// Based on: https://open.larksuite.com/document/server-docs/docs/drive-v1/folder/get-root-folder-meta
type RootFolderMetaGetResponse struct {
	BaseResponse
	Data RootFolderMetaData `json:"data"`
}

// FolderCreateRequest represents the request body for creating a folder.
// Based on: https://open.larksuite.com/document/server-docs/docs/drive-v1/folder/create_folder
type FolderCreateRequest struct {
	Name        string `json:"name"`
	FolderToken string `json:"folder_token"`
}

// FolderCreateResponseData represents the data for a folder creation response.
type FolderCreateResponseData struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// FolderCreateResponse represents the response for creating a folder.
type FolderCreateResponse struct {
	BaseResponse
	Data FolderCreateResponseData `json:"data"`
}

// FileMoveRequest represents the request to move a file.
// Based on: https://open.larksuite.com/document/server-docs/docs/drive-v1/file/move
type FileMoveRequest struct {
	Type        string `json:"type"`
	FolderToken string `json:"folder_token"`
}

// FileTaskResponseData represents the data for a file task response.
type FileTaskResponseData struct {
	TaskID string `json:"task_id"`
}

// FileTaskResponse represents a generic response for file operations that return a task ID.
// Based on: https://open.larksuite.com/document/server-docs/docs/drive-v1/file/move and https://open.larksuite.com/document/server-docs/docs/drive-v1/file/delete
type FileTaskResponse struct {
	BaseResponse
	Data FileTaskResponseData `json:"data"`
}

// FileChild represents a file or folder in a folder's children list.
type FileChild struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

// FolderChildrenListData represents the data for a folder children list response.
type FolderChildrenListData struct {
	Files         []FileChild `json:"files"`
	NextPageToken string      `json:"next_page_token"`
	HasMore       bool        `json:"has_more"`
}

// FolderChildrenListResponse represents the response for listing children of a folder.
// Based on: https://open.larksuite.com/document/server-docs/docs/drive-v1/folder/list
type FolderChildrenListResponse struct {
	BaseResponse
	Data FolderChildrenListData `json:"data"`
}
