// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &roleResource{}

func NewRoleResource() resource.Resource {
	return &roleResource{}
}

// roleResource defines the resource implementation.
type roleResource struct {
	client *common.LarkClient
}

// roleResourceModel describes the resource data model.
// fields that need to be configured by user.
type roleResourceModel struct {
	BaseResourceModel
	RoleID   types.String `tfsdk:"role_id"`
	RoleName types.String `tfsdk:"role_name"`
}

func (r *roleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *roleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"role_id": schema.StringAttribute{
			Description:         "Unique identity of the role, unique under a single tenant",
			MarkdownDescription: "Unique identity of the role, unique under a single tenant",
			Computed:            true,
		},
		"role_name": schema.StringAttribute{
			Description:         "Role name, unique under single tenant",
			MarkdownDescription: "Role name, unique under single tenant",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 50),
			},
		},
	}

	for k, v := range baseAttributes {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		Description:         "Manages role in Lark",
		MarkdownDescription: "Manages role in Lark",
		Attributes:          attributes,
	}
}

func (r *roleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *roleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleRequest := common.RoleRequest{
		RoleName: data.RoleName.ValueString(),
	}

	roleResponse, err := common.RoleCreateAPI(ctx, r.client, roleRequest)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Role", err.Error())
		return
	}

	data.RoleID = types.StringValue(roleResponse.Data.RoleID)
	data.Id = types.StringValue(common.ConstructID(common.RESOURCE, common.ROLE, roleResponse.Data.RoleID))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *roleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data roleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Implement Read, There is no get API Currently

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *roleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan roleResourceModel
	var state roleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleRequest := common.RoleRequest{
		RoleName: plan.RoleName.ValueString(),
	}

	_, err := common.RoleUpdateAPI(ctx, r.client, state.RoleID.ValueString(), roleRequest)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Role", err.Error())
		return
	}

	plan.Id = state.Id
	plan.RoleID = types.StringValue(state.RoleID.ValueString())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *roleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan roleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := common.RoleDeleteAPI(ctx, r.client, plan.RoleID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Role", err.Error())
		return
	}
}
