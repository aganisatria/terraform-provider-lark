// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &userGroupResource{}
var _ resource.ResourceWithModifyPlan = &userGroupResource{}

func NewUserGroupResource() resource.Resource {
	return &userGroupResource{}
}

// userGroupResource defines the resource implementation.
type userGroupResource struct {
	client *common.LarkClient
}

// userGroupResourceModel describes the resource data model.
// fields that need to be configured by user.
type userGroupResourceModel struct {
	BaseResourceModel
	GroupId     types.String `tfsdk:"group_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Type        types.String `tfsdk:"type"`
}

func (r *userGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (r *userGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"group_id": schema.StringAttribute{
			Description:         "User group ID.",
			MarkdownDescription: "User group ID.",
			Computed:            true,
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9]+$`), "must be a valid group ID"),
			},
		},
		"name": schema.StringAttribute{
			Description:         "User group name.",
			MarkdownDescription: "User group name.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.LengthAtMost(64),
			},
		},
		"description": schema.StringAttribute{
			Description:         "User group description.",
			MarkdownDescription: "User group description.",
			Optional:            true,
			Computed:            true,
		},
		"type": schema.StringAttribute{
			Description:         "Type of group: public or private.",
			MarkdownDescription: "Type of group: 1 for common user group or 2 for dynamic user group.",
			Computed:            true,
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("1", "2"),
			},
		},
	}

	for k, v := range baseAttributes {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		Description:         "Manages user groups in Lark",
		MarkdownDescription: "Manages user groups in Lark",
		Attributes:          attributes,
	}
}

func (r *userGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *userGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data userGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userGroupCreateRequestBody := common.UsergroupCreateRequest{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		GroupID:     data.GroupId.ValueString(),
		Type:        data.Type.ValueString(),
	}

	_, err := common.UsergroupGetAPI(ctx, r.client, userGroupCreateRequestBody.GroupID)
	if err == nil {
		resp.Diagnostics.AddError("User Group Already Exists", "User Group already exists")
		return
	}

	userGroupCreateResponse, err := common.UsergroupCreateAPI(ctx, r.client, userGroupCreateRequestBody)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating User Group", err.Error())
		return
	}

	data.Id = types.StringValue(userGroupCreateResponse.Data.GroupID)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	data.GroupId = types.StringValue(userGroupCreateResponse.Data.GroupID)
	data.Name = types.StringValue(userGroupCreateRequestBody.Name)
	data.Description = types.StringValue(userGroupCreateRequestBody.Description)
	data.Type = types.StringValue(userGroupCreateRequestBody.Type)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data userGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userGroupGetResponse, err := common.UsergroupGetAPI(ctx, r.client, data.GroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Getting User Group", err.Error())
		return
	}

	data.Id = types.StringValue(userGroupGetResponse.Data.Group.ID)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	data.GroupId = types.StringValue(userGroupGetResponse.Data.Group.ID)
	data.Name = types.StringValue(userGroupGetResponse.Data.Group.Name)
	data.Description = types.StringValue(userGroupGetResponse.Data.Group.Description)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *userGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state userGroupResourceModel
	var plan userGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userGroupUpdateRequestBody := common.UsergroupUpdateRequest{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	_, err := common.UsergroupUpdateAPI(ctx, r.client, state.GroupId.ValueString(), userGroupUpdateRequestBody)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating User Group", err.Error())
		return
	}

	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	state.Name = plan.Name
	state.Description = plan.Description

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *userGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan userGroupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := common.UsergroupDeleteAPI(ctx, r.client, plan.GroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting User Group", err.Error())
		return
	}
}

// We use modify plan when we need both plan and state when validating.
func (r *userGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *userGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// plan null means resource is being deleted.
	if req.Plan.Raw.IsNull() {
		return
	}

	// State null means resource is being created.
	if req.State.Raw.IsNull() {

		// Checking if the group name is available.
		userGroupListResponse, err := common.UsergroupListAPI(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError("API Error Getting User Group", err.Error())
			return
		}

		for _, group := range userGroupListResponse.Data.GroupList {
			if group.Name == plan.Name.ValueString() {
				resp.Diagnostics.AddError("Group Name Already Exists", "Group name already exists")
				return
			}
		}

		return
	}

	if !state.GroupId.IsNull() && !plan.GroupId.IsUnknown() && plan.GroupId.ValueString() != state.GroupId.ValueString() {
		resp.Diagnostics.AddError(
			"Group ID Cannot Be Modified",
			"Group ID cannot be changed after resource creation",
		)
	}

	if plan.Name.ValueString() != state.Name.ValueString() {
		// Checking if the group name is available.
		userGroupListResponse, err := common.UsergroupListAPI(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError("API Error Getting User Group", err.Error())
			return
		}

		for _, group := range userGroupListResponse.Data.GroupList {
			if group.Name == plan.Name.ValueString() {
				resp.Diagnostics.AddError("Group Name Already Exists", "Group name already exists")
				return
			}
		}
	}
}
