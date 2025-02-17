// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validator

import (
	"context"
	"fmt"
	"strings"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UserIDValidator is a validator for a list of user_id.
type UserIDValidator struct {
	Path          path.Path
	doesAList     bool
	doesSkipAppID bool
	idType        common.UserIDType
	client        *common.LarkClient
}

func NewUserIDValidator(pathToField string, doesAList bool, doesSkipAppID bool, idType common.UserIDType, client *common.LarkClient) resource.ConfigValidator {
	return &UserIDValidator{
		Path:          path.Root(pathToField),
		doesAList:     doesAList,
		doesSkipAppID: doesSkipAppID,
		idType:        idType,
		client:        client,
	}
}

func (v UserIDValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("validates that %s is a valid list of %s", v.Path, v.idType)
}

func (v UserIDValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("validates that `%s` is a valid list of %s", v.Path, v.idType)
}

func (v UserIDValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	availableIDs := []string{}

	if !v.doesAList {
		var id types.String

		diags := req.Config.GetAttribute(ctx, v.Path, &id)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() || id.IsNull() || id.IsUnknown() {
			return
		}

		if id.ValueString() == "" {
			resp.Diagnostics.AddAttributeError(
				v.Path,
				"Invalid ID",
				"ID cannot be empty",
			)
			return
		}

		availableIDs = append(availableIDs, id.ValueString())
	} else {
		var ids []types.String

		diags := req.Config.GetAttribute(ctx, v.Path, &ids)
		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() || len(ids) == 0 {
			return
		}

		for _, element := range ids {
			if !element.IsNull() && !element.IsUnknown() {
				if element.ValueString() == "" {
					resp.Diagnostics.AddAttributeError(
						v.Path,
						"Invalid ID",
						"ID cannot be empty",
					)
					return
				}

				if v.doesSkipAppID && strings.HasPrefix(element.ValueString(), "cli_") {
					continue
				}

				availableIDs = append(availableIDs, element.ValueString())
			}
		}
	}

	if len(availableIDs) > 0 {
		if err := isListOfUserIDsExist(ctx, v.client, availableIDs, v.idType); err != nil {
			resp.Diagnostics.AddAttributeError(
				v.Path,
				"Invalid ID",
				fmt.Sprintf("Error validating IDs: %s", err.Error()),
			)
		}
	}
}

func isListOfUserIDsExist(ctx context.Context, client *common.LarkClient, ids []string, idType common.UserIDType) error {
	_, err := common.GetUsersByIDAPI(ctx, client, ids, idType)
	if err != nil {
		return err
	}

	return nil
}
