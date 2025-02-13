// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.String = &groupNameValidator{}

type groupNameValidator struct{}

func GroupNameValidator() validator.String {
	return &groupNameValidator{}
}

func (v groupNameValidator) Description(ctx context.Context) string {
	return "validates group name based on chat type"
}

func (v groupNameValidator) MarkdownDescription(ctx context.Context) string {
	return "validates group name based on chat type"
}

func (v groupNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		// If name is empty for private chat, it's fine because it will default to "(no title)".
		return
	}

	var chatType types.String
	diags := req.Config.GetAttribute(ctx, path.Root("chat_type"), &chatType)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	name := req.ConfigValue.ValueString()

	switch chatType.ValueString() {
	case "public":
		if len(name) < 2 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Group Name",
				"Public group name must be at least 2 characters long",
			)
		}
	case "private":
		// If name is not empty, it must be at least 2 characters long.
		if name != "" && len(name) < 2 {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Group Name",
				"If provided, private group name must be at least 2 characters long",
			)
		}
	}
}
