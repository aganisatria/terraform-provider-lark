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
var _ resource.Resource = &userGroupMemberResource{}
var _ resource.ResourceWithConfigValidators = &userGroupMemberResource{}

func NewUserGroupMemberResource() resource.Resource {
	return &userGroupMemberResource{}
}

// userGroupMemberResource defines the resource implementation.
type userGroupMemberResource struct {
	client *common.LarkClient
}

// userGroupMemberResourceModel describes the resource data model.
// fields that need to be configured by user.
type userGroupMemberResourceModel struct {
	BaseResourceModel
	UserGroupID types.String   `tfsdk:"user_group_id"`
	MemberIDs   []types.String `tfsdk:"member_ids"`
}

func (r *userGroupMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group_member"
}

func (r *userGroupMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"user_group_id": schema.StringAttribute{
			Description:         "Unique identity of the role, unique under a single tenant",
			MarkdownDescription: "Unique identity of the role, unique under a single tenant",
			Required:            true,
		},
		"member_ids": schema.ListAttribute{
			Description:         "List of user group members added by the user group (OpenID list of a batch of users)",
			MarkdownDescription: "List of user group members added by the user group (OpenID list of a batch of users)",
			Required:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				listvalidator.SizeAtLeast(1),
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

func (r *userGroupMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userGroupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userGroupMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members := []string{}
	for _, member := range data.MemberIDs {
		members = append(members, member.ValueString())
	}

	err := r.AddHelper(ctx, data, members, data.UserGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	data.Id = types.StringValue(common.ConstructID(common.RESOURCE, common.USER_GROUP_MEMBER, data.UserGroupID.ValueString()))
	data.UserGroupID = types.StringValue(data.UserGroupID.ValueString())
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userGroupMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members := []string{}
	if data.MemberIDs != nil {
		for _, member := range data.MemberIDs {
			members = append(members, member.ValueString())
		}
	}

	if len(members) > 0 {
		response, err := common.UsergroupMemberGetByMemberTypeAPI(ctx, r.client, data.UserGroupID.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError("API Error Getting User Group Member", err.Error())
			return
		}

		dbMembersIDs := []string{}
		for _, member := range response.Data.MemberList {
			dbMembersIDs = append(dbMembersIDs, member.MemberID)
		}

		for _, member := range members {
			if !slices.Contains(dbMembersIDs, member) {
				resp.Diagnostics.AddError("API Error Getting User Group Member", fmt.Sprintf("Member %s not found in user group %s", member, data.UserGroupID.ValueString()))
				return
			}
		}
	}

	data.Id = types.StringValue(data.Id.ValueString())
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	data.UserGroupID = types.StringValue(data.UserGroupID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan userGroupMemberResourceModel
	var state userGroupMemberResourceModel

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

	err := r.AddHelper(ctx, plan, addedMembers, state.UserGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	err = r.DeleteHelper(ctx, plan, removedMembers, state.UserGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	plan.Id = types.StringValue(state.Id.ValueString())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *userGroupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan userGroupMemberResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	members := []string{}
	for _, member := range plan.MemberIDs {
		members = append(members, member.ValueString())
	}

	err := r.DeleteHelper(ctx, plan, members, plan.UserGroupID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}
}

func (r *userGroupMemberResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if r.client == nil {
		return []resource.ConfigValidator{}
	}
	return []resource.ConfigValidator{
		local_validator.NewUserIDValidator("member_ids", true, false, common.OPEN_ID, r.client),
	}
}

func (r *userGroupMemberResource) AddHelper(ctx context.Context, plan userGroupMemberResourceModel, addedMembers []string, userGroupID string) *diag.ErrorDiagnostic {
	if len(addedMembers) > 0 {
		members := []common.UsergroupMember{}
		for _, member := range addedMembers {
			members = append(members, common.UsergroupMember{
				MemberID:     member,
				MemberType:   "user",
				MemberIDType: "open_id",
			})
		}
		addUserGroupMemberRequest := common.UsergroupMemberAddRequest{
			Members: members,
		}

		_, err := common.UsergroupMemberAddAPI(ctx, r.client, plan.UserGroupID.ValueString(), addUserGroupMemberRequest)
		if err != nil {
			errorDiag := diag.NewErrorDiagnostic(
				"API Error Adding User Group Member",
				err.Error(),
			)
			return &errorDiag
		}
	}

	return nil
}

func (r *userGroupMemberResource) DeleteHelper(ctx context.Context, plan userGroupMemberResourceModel, removedMembers []string, roleID string) *diag.ErrorDiagnostic {
	if len(removedMembers) > 0 {
		members := []common.UsergroupMember{}
		for _, member := range removedMembers {
			members = append(members, common.UsergroupMember{
				MemberID:     member,
				MemberType:   "user",
				MemberIDType: "open_id",
			})
		}
		removeUserGroupMemberRequest := common.UsergroupMemberRemoveRequest{
			Members: members,
		}

		_, err := common.UsergroupMemberRemoveAPI(ctx, r.client, plan.UserGroupID.ValueString(), removeUserGroupMemberRequest)
		if err != nil {
			errorDiag := diag.NewErrorDiagnostic(
				"API Error Removing User Group Member",
				err.Error(),
			)
			return &errorDiag
		}
	}

	return nil
}
