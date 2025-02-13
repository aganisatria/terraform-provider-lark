// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaseResourceModel struct {
	Id          types.String `tfsdk:"id"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func BaseSchemaResourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Resource ID.",
			MarkdownDescription: "Resource ID.",
			Computed:            true,
		},
		"last_updated": schema.StringAttribute{
			Description:         "Timestamp of the last update.",
			MarkdownDescription: "Timestamp of the last update.",
			Computed:            true,
		},
	}
}

type I18nName struct {
	ZhCn types.String `tfsdk:"zh_cn"`
	JaJp types.String `tfsdk:"ja_jp"`
	EnUs types.String `tfsdk:"en_us"`
}

type RestrictedModeSetting struct {
	Status                         types.Bool   `tfsdk:"status"`
	ScreenshotHasPermissionSetting types.String `tfsdk:"screenshot_has_permission_setting"`
	DownloadHasPermissionSetting   types.String `tfsdk:"download_has_permission_setting"`
	MessageHasPermissionSetting    types.String `tfsdk:"message_has_permission_setting"`
}
