// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"context"
	"fmt"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserIDValidator is a validator for a list of user_id.
type UserIDValidator struct {
	Path      path.Path
	doesAList bool
	client    *common.LarkClient
}

func NewUserIDValidator(pathToField string, doesAList bool, client *common.LarkClient) resource.ConfigValidator {
	return &UserIDValidator{
		Path:      path.Root(pathToField),
		doesAList: doesAList,
		client:    client,
	}
}

func (v UserIDValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("validates that %s is a valid list of user_id", v.Path)
}

func (v UserIDValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("validates that `%s` is a valid list of user_id", v.Path)
}

func (v UserIDValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	availableUserIDs := []string{}

	if !v.doesAList {
		var userID types.String

		diags := req.Config.GetAttribute(ctx, v.Path, &userID)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() || userID.IsNull() || userID.IsUnknown() {
			return
		}

		availableUserIDs = append(availableUserIDs, userID.ValueString())
	} else {
		var userIDs []types.String

		diags := req.Config.GetAttribute(ctx, v.Path, &userIDs)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() || len(userIDs) == 0 {
			return
		}

		for _, element := range userIDs {
			if !element.IsNull() && !element.IsUnknown() {
				availableUserIDs = append(availableUserIDs, element.ValueString())
			}
		}
	}

	if len(availableUserIDs) > 0 {
		if err := isListOfUserIDsExist(ctx, v.client, availableUserIDs); err != nil {
			resp.Diagnostics.AddAttributeError(
				v.Path,
				"Invalid User ID",
				fmt.Sprintf("Error validating user IDs: %s", err.Error()),
			)
		}
	}
}

func isListOfUserIDsExist(ctx context.Context, client *common.LarkClient, userIDs []string) error {
	_, err := common.GetUsersByOpenIDAPI(ctx, client, userIDs)

	if err != nil {
		return err
	}

	return nil
}
