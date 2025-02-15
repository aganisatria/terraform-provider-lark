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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &groupChatMemberResource{}
var _ resource.ResourceWithConfigValidators = &groupChatMemberResource{}

func NewGroupChatMemberResource() resource.Resource {
	return &groupChatMemberResource{}
}

// groupChatMemberResource defines the resource implementation.
type groupChatMemberResource struct {
	client *common.LarkClient
}

// groupChatMemberResourceModel describes the resource data model.
// fields that need to be configured by user.
type groupChatMemberResourceModel struct {
	BaseResourceModel
	GroupChatID      types.String   `tfsdk:"group_chat_id"`
	MemberIDs        []types.String `tfsdk:"member_ids"`
	AdministratorIDs []types.String `tfsdk:"administrator_ids"`
}

func (r *groupChatMemberResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_chat_member"
}

func (r *groupChatMemberResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"group_chat_id": schema.StringAttribute{
			Description:         "Unique identity of the group chat, unique under a single tenant",
			MarkdownDescription: "Unique identity of the group chat, unique under a single tenant",
			Required:            true,
		},
		"member_ids": schema.ListAttribute{
			Description:         "List of members added by the group chat. Can be UserID (starts with ou) or BotID (starts with cli)",
			MarkdownDescription: "List of members added by the group chat. Can be UserID (starts with ou) or BotID (starts with cli)",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"administrator_ids": schema.ListAttribute{
			Description:         "List of administrator added by the group chat. Can be UserID (starts with ou) or BotID (starts with cli)",
			MarkdownDescription: "List of administrator added by the group chat. Can be UserID (starts with ou) or BotID (starts with cli)",
			Optional:            true,
			ElementType:         types.StringType,
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

func (r *groupChatMemberResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupChatMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data groupChatMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addedMembers := []string{}
	for _, member := range data.MemberIDs {
		addedMembers = append(addedMembers, member.ValueString())
	}
	addedAdministrators := []string{}
	for _, member := range data.AdministratorIDs {
		addedAdministrators = append(addedAdministrators, member.ValueString())
	}

	err := r.AddHelper(ctx, data, addedMembers, addedAdministrators, data.GroupChatID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	data.Id = types.StringValue(data.GroupChatID.ValueString())
	data.GroupChatID = types.StringValue(data.GroupChatID.ValueString())
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupChatMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupChatMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members := []string{}
	for _, member := range state.MemberIDs {
		members = append(members, member.ValueString())
	}
	administrators := []string{}
	for _, member := range state.AdministratorIDs {
		administrators = append(administrators, member.ValueString())
	}

	if len(members) > 0 || len(administrators) > 0 {
		groupChat, err := common.GroupChatMemberGetAPI(ctx, r.client, state.GroupChatID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(err.Error(), err.Error())
			return
		}

		groupChatMembersInServer := []string{}
		for _, member := range groupChat.Data.Items {
			groupChatMembersInServer = append(groupChatMembersInServer, member.MemberID)
		}

		members = []string{}
		for _, member := range state.MemberIDs {
			if !common.Contains(member.ValueString(), groupChatMembersInServer...) {
				resp.Diagnostics.AddError(fmt.Sprintf(
					"Member %s not found in group chat %s, usually this is because the member force deleted from group chat without using this resource",
					member.ValueString(),
					state.GroupChatID.ValueString(),
				),
					"Member not found in group chat",
				)

			}
		}

		// INFO: Since lark doesn't support get group chat member, we need to add it again.
		errAdd := r.AddHelper(ctx, state, members, administrators, state.GroupChatID.ValueString())
		if errAdd != nil {
			resp.Diagnostics.AddError(errAdd.Summary(), errAdd.Detail())
			return
		}
	}

	state.Id = types.StringValue(state.GroupChatID.ValueString())
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	state.GroupChatID = types.StringValue(state.GroupChatID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *groupChatMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupChatMemberResourceModel
	var state groupChatMemberResourceModel

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

	planAdministrators := []string{}
	for _, member := range plan.AdministratorIDs {
		planAdministrators = append(planAdministrators, member.ValueString())
	}

	stateAdministrators := []string{}
	for _, member := range state.AdministratorIDs {
		stateAdministrators = append(stateAdministrators, member.ValueString())
	}

	addedAdministrators := []string{}
	for _, member := range planAdministrators {
		if !slices.Contains(stateAdministrators, member) {
			addedAdministrators = append(addedAdministrators, member)
		}
	}

	removedAdministrators := []string{}
	for _, member := range stateAdministrators {
		if !slices.Contains(planAdministrators, member) {
			removedAdministrators = append(removedAdministrators, member)
		}
	}

	err := r.AddHelper(ctx, plan, addedMembers, addedAdministrators, plan.GroupChatID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	err = r.DeleteHelper(ctx, plan, removedMembers, removedAdministrators, plan.GroupChatID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}

	plan.Id = types.StringValue(state.Id.ValueString())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupChatMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan groupChatMemberResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	members := []string{}
	for _, member := range plan.MemberIDs {
		members = append(members, member.ValueString())
	}
	administrators := []string{}
	for _, member := range plan.AdministratorIDs {
		administrators = append(administrators, member.ValueString())
	}

	err := r.DeleteHelper(ctx, plan, members, administrators, plan.GroupChatID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(err.Summary(), err.Detail())
		return
	}
}

func (r *groupChatMemberResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if r.client == nil {
		return []resource.ConfigValidator{}
	}
	return []resource.ConfigValidator{
		local_validator.NewUserIDValidator("member_ids", true, true, r.client),
		local_validator.NewUserIDValidator("administrator_ids", true, true, r.client),
		local_validator.NewListShouldBeMemberOfAnotherListValidator("administrator_ids", "member_ids", r.client),
	}
}

func (r *groupChatMemberResource) AddHelper(ctx context.Context, plan groupChatMemberResourceModel, addedMembers []string, addedAdministrators []string, groupChatID string) *diag.ErrorDiagnostic {
	members := common.GroupChatMemberRequest{}
	if len(addedMembers) > 0 {
		members.IDList = append(members.IDList, addedMembers...)
	}

	administrators := common.GroupChatAdministratorRequest{}
	if len(addedAdministrators) > 0 {
		administrators.ManagerIDs = append(administrators.ManagerIDs, addedAdministrators...)
	}

	if len(members.IDList) > 0 {
		_, err := common.GroupChatMemberAddAPI(ctx, r.client, groupChatID, members)
		if err != nil {
			errorDiag := diag.NewErrorDiagnostic(
				"API Error Adding User Group Member",
				err.Error(),
			)
			return &errorDiag
		}
	}

	if len(administrators.ManagerIDs) > 0 {
		_, err := common.GroupChatAdministratorAddAPI(ctx, r.client, groupChatID, administrators)
		if err != nil {
			errorDiag := diag.NewErrorDiagnostic(
				"API Error Adding User Group Administrator",
				err.Error(),
			)
			return &errorDiag
		}
	}

	return nil
}

func (r *groupChatMemberResource) DeleteHelper(ctx context.Context, plan groupChatMemberResourceModel, removedMembers []string, removedAdministrators []string, groupChatID string) *diag.ErrorDiagnostic {
	members := common.GroupChatMemberRequest{}
	if len(removedMembers) > 0 {
		for _, member := range removedMembers {
			members.IDList = append(members.IDList, member)
		}
	}

	administrators := common.GroupChatAdministratorRequest{}
	if len(removedAdministrators) > 0 {
		administrators.ManagerIDs = append(administrators.ManagerIDs, removedAdministrators...)
	}

	if len(administrators.ManagerIDs) > 0 {
		_, err := common.GroupChatAdministratorDeleteAPI(ctx, r.client, groupChatID, administrators)
		if err != nil {
			errorDiag := diag.NewErrorDiagnostic(
				"API Error Removing User Group Administrator",
				err.Error(),
			)
			return &errorDiag
		}
	}

	if len(members.IDList) > 0 {
		_, err := common.GroupChatMemberDeleteAPI(ctx, r.client, groupChatID, members)
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
