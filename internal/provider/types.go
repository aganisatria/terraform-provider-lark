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
