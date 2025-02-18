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

// ListShouldBeMemberOfAnotherListValidator is a validator does a list is member of another list.
type ListShouldBeMemberOfAnotherListValidator struct {
	Path          path.Path
	ValidatorPath path.Path
	client        *common.LarkClient
}

func NewListShouldBeMemberOfAnotherListValidator(pathToField string, validatorPath string, client *common.LarkClient) resource.ConfigValidator {
	return &ListShouldBeMemberOfAnotherListValidator{
		Path:          path.Root(pathToField),
		ValidatorPath: path.Root(validatorPath),
		client:        client,
	}
}

func (v ListShouldBeMemberOfAnotherListValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("validates that %s is a member of %s", v.Path, v.ValidatorPath)
}

func (v ListShouldBeMemberOfAnotherListValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("validates that `%s` is a member of `%s`", v.Path, v.ValidatorPath)
}

func (v ListShouldBeMemberOfAnotherListValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var list types.List
	var validatorList types.List

	diags := req.Config.GetAttribute(ctx, v.Path, &list)
	resp.Diagnostics.Append(diags...)

	diags = req.Config.GetAttribute(ctx, v.ValidatorPath, &validatorList)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() || list.IsNull() || list.IsUnknown() || list.Elements() == nil {
		return
	}

	if !list.IsNull() && !list.IsUnknown() && list.Elements() != nil && (validatorList.IsNull() || validatorList.IsUnknown() || validatorList.Elements() == nil) {
		resp.Diagnostics.AddAttributeError(
			v.Path,
			fmt.Sprintf("If %s is not empty, %s must not be empty", v.ValidatorPath, v.Path),
			fmt.Sprintf("If %s is not empty, %s must not be empty", v.ValidatorPath, v.Path),
		)
	}

	listValue := []string{}

	for _, element := range list.Elements() {
		if !element.IsNull() && !element.IsUnknown() {
			listValue = append(listValue, element.String())
		}
	}

	validatorListValue := []string{}

	for _, element := range validatorList.Elements() {
		if !element.IsNull() && !element.IsUnknown() {
			validatorListValue = append(validatorListValue, element.String())
		}
	}

	notExistOnValidatorList := []string{}

	for _, element := range listValue {
		if !common.Contains(element, validatorListValue...) {
			notExistOnValidatorList = append(notExistOnValidatorList, element)
		}
	}

	if len(notExistOnValidatorList) > 0 {
		resp.Diagnostics.AddAttributeError(
			v.Path,
			fmt.Sprintf("There is some %s element that is not exist on %s element", v.Path, v.ValidatorPath),
			fmt.Sprintf("Error element: %s", notExistOnValidatorList),
		)
	}
}
