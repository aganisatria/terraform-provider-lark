// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &workforceTypeResource{}

func NewWorkforceTypeResource() resource.Resource {
	return &workforceTypeResource{}
}

// workforceTypeResource defines the resource implementation.
type workforceTypeResource struct {
	client *common.LarkClient
}

type I18nContent struct {
	Locale types.String `tfsdk:"locale"`
	Value  types.String `tfsdk:"value"`
}

// workforceTypeResourceModel describes the resource data model.
// fields that need to be configured by user.
type workforceTypeResourceModel struct {
	BaseResourceModel
	Content     types.String  `tfsdk:"content"`
	EnumType    types.Int64   `tfsdk:"enum_type"`
	EnumStatus  types.Int64   `tfsdk:"enum_status"`
	I18nContent []I18nContent `tfsdk:"i18n_content"`
	EnumID      types.String  `tfsdk:"enum_id"`
	EnumValue   types.String  `tfsdk:"enum_value"`
}

func (r *workforceTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workforce_type"
}

func (r *workforceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"content": schema.StringAttribute{
			Description:         "Enum content",
			MarkdownDescription: "Enum content",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthBetween(1, 100),
			},
		},
		"enum_type": schema.Int64Attribute{
			Description:         "Type",
			MarkdownDescription: "Type",
			Required:            true,
			Validators: []validator.Int64{
				int64validator.OneOf(2), // Currently only support 2 (custom)
			},
		},
		"enum_status": schema.Int64Attribute{
			Description:         "Status",
			MarkdownDescription: "Status",
			Required:            true,
			Validators: []validator.Int64{
				int64validator.OneOf(1, 2),
			},
		},
		"i18n_content": schema.ListNestedAttribute{
			Description:         "Internationalization definition",
			MarkdownDescription: "Internationalization definition",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"locale": schema.StringAttribute{
						Description:         "Language version",
						MarkdownDescription: "Language version",
						Optional:            true,
					},
					"value": schema.StringAttribute{
						Description:         "Field name",
						MarkdownDescription: "Field name",
						Optional:            true,
					},
				},
			},
		},
		"enum_id": schema.StringAttribute{
			Description:         "Enum ID",
			MarkdownDescription: "Enum ID",
			Computed:            true,
		},
		"enum_value": schema.StringAttribute{
			Description:         "Enum value, which is automatically generated for a newly created workforce type",
			MarkdownDescription: "Enum value, which is automatically generated for a newly created workforce type",
			Computed:            true,
		},
	}

	for k, v := range baseAttributes {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		Description:         "Manages workforce type in Lark",
		MarkdownDescription: "Manages workforce type in Lark",
		Attributes:          attributes,
	}
}

func (r *workforceTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *workforceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data workforceTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var i18nContent []common.I18nContent
	if len(data.I18nContent) > 0 {
		for _, v := range data.I18nContent {
			i18nContent = append(i18nContent, common.I18nContent{
				Locale: v.Locale.ValueString(),
				Value:  v.Value.ValueString(),
			})
		}
	}

	request := common.WorkforceTypeRequest{
		Content:     data.Content.ValueString(),
		EnumType:    int(data.EnumType.ValueInt64()),
		EnumStatus:  int(data.EnumStatus.ValueInt64()),
		I18nContent: i18nContent,
	}

	response, err := common.WorkforceTypeCreateAPI(ctx, r.client, request)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Workforce Type", err.Error())
		return
	}

	data.Id = types.StringValue(common.ConstructID(common.RESOURCE, common.WORKFORCE_TYPE, response.Data.EmployeeTypeEnum.EnumID))
	data.EnumID = types.StringValue(response.Data.EmployeeTypeEnum.EnumID)
	data.EnumValue = types.StringValue(response.Data.EmployeeTypeEnum.EnumValue)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *workforceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data workforceTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	response, err := common.WorkforceTypeGetAllAPI(ctx, r.client)
	if err != nil {
		resp.Diagnostics.AddError("API Error Getting Workforce Type", err.Error())
		return
	}

	found := false
	for _, enum := range response.Data.Items {
		if enum.EnumID == data.EnumID.ValueString() {
			found = true
			break
		}
	}

	if !found {
		resp.Diagnostics.AddError("API Error Getting Workforce Type", "Workforce type not found")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *workforceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan workforceTypeResourceModel
	var state workforceTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var i18nContent []common.I18nContent
	if len(plan.I18nContent) > 0 {
		for _, v := range plan.I18nContent {
			i18nContent = append(i18nContent, common.I18nContent{
				Locale: v.Locale.ValueString(),
				Value:  v.Value.ValueString(),
			})
		}
	}

	request := common.WorkforceTypeRequest{
		Content:     plan.Content.ValueString(),
		EnumType:    int(plan.EnumType.ValueInt64()),
		EnumStatus:  int(plan.EnumStatus.ValueInt64()),
		I18nContent: i18nContent,
	}

	response, err := common.WorkforceTypeUpdateAPI(ctx, r.client, state.EnumID.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Workforce Type", err.Error())
		return
	}

	plan.Id = state.Id
	plan.EnumID = types.StringValue(response.Data.EmployeeTypeEnum.EnumID)
	plan.EnumValue = types.StringValue(response.Data.EmployeeTypeEnum.EnumValue)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *workforceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan workforceTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := common.WorkforceTypeDeleteAPI(ctx, r.client, plan.EnumID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Workforce Type", err.Error())
		return
	}
}
