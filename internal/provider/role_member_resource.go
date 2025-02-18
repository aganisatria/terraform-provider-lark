// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	local_validator "github.com/aganisatria/terraform-provider-lark/internal/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &roleMemberResource{}
var _ resource.ResourceWithConfigValidators = &roleMemberResource{}

func NewRoleMemberResource() resource.Resource {
	return &roleMemberResource{}
}

// roleMemberResource defines the resource implementation.
type roleMemberResource struct {
	client *common.LarkClient
}

// roleMemberResourceModel describes the resource data model.
// fields that need to be configured by user.
type roleMemberResourceModel struct {
	BaseResourceModel
	RoleID    types.String   `tfsdk:"role_id"`
	MemberIDs []types.String `tfsdk:"member_ids"`
}

func (r *roleMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_member"
}

func (r *roleMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"role_id": schema.StringAttribute{
			Description:         "Unique identity of the role, unique under a single tenant",
			MarkdownDescription: "Unique identity of the role, unique under a single tenant",
			Required:            true,
		},
		"member_ids": schema.ListAttribute{
			Description:         "List of role members added by the role (UserID list of a batch of users)",
			MarkdownDescription: "List of role members added by the role (UserID list of a batch of users)",
			Required:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeBetween(1, 100),
			},
		},
	}

	for k, v := range baseAttributes {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		Description:         "Manages role member in Lark",
		MarkdownDescription: "Manages role member in Lark",
		Attributes:          attributes,
	}
}

func (r *roleMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*common.LarkClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *LarkClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *roleMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data roleMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members := []string{}
	memberIDs := []types.String{}
	for _, member := range data.MemberIDs {
		members = append(members, member.ValueString())
		memberIDs = append(memberIDs, types.StringValue(member.ValueString()))
	}

	err := r.AddHelper(ctx, data, members, data.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	data.Id = types.StringValue(data.RoleID.ValueString())
	data.RoleID = types.StringValue(data.RoleID.ValueString())
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	data.MemberIDs = memberIDs

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *roleMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data roleMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members := []string{}
	memberIDs := []types.String{}
	if data.MemberIDs != nil {
		for _, member := range data.MemberIDs {
			members = append(members, member.ValueString())
			memberIDs = append(memberIDs, types.StringValue(member.ValueString()))
		}
	}

	if len(members) > 0 {
		response, err := common.RoleMemberGetAPI(ctx, r.client, data.RoleID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("API Error Getting Role Member", err.Error())
			return
		}

		dbMembersIDs := []string{}
		for _, member := range response.Data.Members {
			dbMembersIDs = append(dbMembersIDs, member.UserID)
		}

		for _, member := range members {
			if !slices.Contains(dbMembersIDs, member) {
				resp.Diagnostics.AddError("API Error Getting Role Member", fmt.Sprintf("Member %s not found in role %s", member, data.RoleID.ValueString()))
				return
			}
		}
	}

	data.Id = types.StringValue(data.RoleID.ValueString())
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	data.RoleID = types.StringValue(data.RoleID.ValueString())
	data.MemberIDs = memberIDs

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *roleMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleMemberResourceModel
	var state roleMemberResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planMembers := []string{}
	for _, member := range plan.MemberIDs {
		planMembers = append(planMembers, member.ValueString())
	}
	stateMembers := []string{}
	for _, member := range state.MemberIDs {
		stateMembers = append(stateMembers, member.ValueString())
	}

	addedMembers := []string{}
	for _, member := range planMembers {
		if !slices.Contains(stateMembers, member) {
			addedMembers = append(addedMembers, member)
		}
	}
	removedMembers := []string{}
	for _, member := range stateMembers {
		if !slices.Contains(planMembers, member) {
			removedMembers = append(removedMembers, member)
		}
	}

	err := r.AddHelper(ctx, plan, addedMembers, state.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	err = r.DeleteHelper(ctx, plan, removedMembers, state.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	plan.Id = types.StringValue(state.Id.ValueString())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *roleMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan roleMemberResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	members := []string{}
	for _, member := range plan.MemberIDs {
		members = append(members, member.ValueString())
	}

	err := r.DeleteHelper(ctx, plan, members, plan.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}
}

func (r *roleMemberResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	// There is a mini gap that client not yet initialized before terraform plan is executed, so we need to check if the client is nil
	if r.client == nil {
		return []resource.ConfigValidator{}
	}
	return []resource.ConfigValidator{
		local_validator.NewUserIDValidator("member_ids", true, false, common.OPEN_ID, r.client),
	}
}

func (r *roleMemberResource) AddHelper(ctx context.Context, plan roleMemberResourceModel, addedMembers []string, roleID string) *diag.ErrorDiagnostic {
	if len(addedMembers) > 0 {
		addRoleMemberRequest := common.RoleMemberCreateRequest{
			Members: addedMembers,
		}

		response, err := common.RoleMemberAddAPI(ctx, r.client, plan.RoleID.ValueString(), addRoleMemberRequest)
		if err != nil {
			errorDiag := diag.NewErrorDiagnostic(
				"API Error Adding Role Member",
				err.Error(),
			)
			return &errorDiag
		}

		for _, member := range response.Data.Results {
			if member.Reason != 1 && member.Reason != 4 {
				errorDiag := diag.NewErrorDiagnostic(
					"API Error Adding Role Member",
					fmt.Sprintf("Failed to add member %s to role %s: %d", member.UserID, plan.RoleID.ValueString(), member.Reason),
				)
				return &errorDiag
			}
		}
	}

	return nil
}

func (r *roleMemberResource) DeleteHelper(ctx context.Context, plan roleMemberResourceModel, removedMembers []string, roleID string) *diag.ErrorDiagnostic {
	if len(removedMembers) > 0 {
		removeRoleMemberRequest := common.RoleMemberDeleteRequest{
			Members: removedMembers,
		}

		response, err := common.RoleMemberDeleteAPI(ctx, r.client, plan.RoleID.ValueString(), removeRoleMemberRequest)
		if err != nil {
			errorDiag := diag.NewErrorDiagnostic(
				"API Error Removing Role Member",
				err.Error(),
			)
			return &errorDiag
		}

		for _, member := range response.Data.Results {
			if member.Reason != 1 && member.Reason != 5 {
				errorDiag := diag.NewErrorDiagnostic(
					"API Error Removing Role Member",
					fmt.Sprintf("Failed to remove member %s from role %s: %d", member.UserID, plan.RoleID.ValueString(), member.Reason),
				)
				return &errorDiag
			}
		}
	}

	return nil
}
