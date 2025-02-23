// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	local_validator "github.com/aganisatria/terraform-provider-lark/internal/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &departmentResource{}

func NewDepartmentResource() resource.Resource {
	return &departmentResource{}
}

// departmentResource defines the resource implementation.
type departmentResource struct {
	client *common.LarkClient
}

type Leaders struct {
	LeaderType types.Int64  `tfsdk:"leader_type"`
	LeaderID   types.String `tfsdk:"leader_id"`
}

// departmentResourceModel describes the resource data model.
// fields that need to be configured by user.
type departmentResourceModel struct {
	BaseResourceModel
	Name                   types.String   `tfsdk:"name"`
	I18nName               *I18nName      `tfsdk:"i18n_name"`
	ParentDepartmentId     types.String   `tfsdk:"parent_department_id"`
	DepartmentId           types.String   `tfsdk:"department_id"`
	OpenDepartmentId       types.String   `tfsdk:"open_department_id"`
	LeaderUserID           types.String   `tfsdk:"leader_user_id"`
	Order                  types.String   `tfsdk:"order"`
	UnitIDs                []types.String `tfsdk:"unit_ids"`
	CreateGroupChat        types.Bool     `tfsdk:"create_group_chat"`
	ChatID                 types.String   `tfsdk:"chat_id"`
	Leaders                []Leaders      `tfsdk:"leaders"`
	GroupChatEmployeeTypes []types.Int64  `tfsdk:"group_chat_employee_types"`
	MemberCount            types.Int64    `tfsdk:"member_count"`
}

func (r *departmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_department"
}

func (r *departmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "Department name.",
			MarkdownDescription: "Department name.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"i18n_name": schema.SingleNestedAttribute{
			Description:         "Internationalized department name.",
			MarkdownDescription: "Internationalized department name.",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"zh_cn": schema.StringAttribute{
					Description:         "Department's Chinese name.",
					MarkdownDescription: "Department's Chinese name.",
					Optional:            true,
				},
				"ja_jp": schema.StringAttribute{
					Description:         "Department's Japanese name.",
					MarkdownDescription: "Department's Japanese name.",
					Optional:            true,
				},
				"en_us": schema.StringAttribute{
					Description:         "Department's English name.",
					MarkdownDescription: "Department's English name.",
					Optional:            true,
				},
			},
		},
		"parent_department_id": schema.StringAttribute{
			Description:         "Parent department ID. If you want to create a root department, set it to 0.",
			MarkdownDescription: "Parent department ID. If you want to create a root department, set it to 0.",
			Required:            true,
		},
		"department_id": schema.StringAttribute{
			Description:         "Department's custom department ID.",
			MarkdownDescription: "Department's custom department ID.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.LengthAtMost(64),
				stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_\-@.]{0,63}$`), "must be a valid department ID if you decide to fill it"),
			},
		},
		"open_department_id": schema.StringAttribute{
			Description:         "Department's open department ID.",
			MarkdownDescription: "Department's open department ID.",
			Computed:            true,
		},
		"leader_user_id": schema.StringAttribute{
			Description:         "Department manager's user ID.",
			MarkdownDescription: "Department manager's user ID.",
			Optional:            true,
		},
		"order": schema.StringAttribute{
			Description:         "Department order, i.e. the order in which the department is displayed among the departments at the same level.",
			MarkdownDescription: "Department order, i.e. the order in which the department is displayed among the departments at the same level.",
			Optional:            true,
		},
		"unit_ids": schema.ListAttribute{
			Description:         "List of the department unit's custom IDs. Only one custom ID is supported currently.",
			MarkdownDescription: "List of the department unit's custom IDs. Only one custom ID is supported currently.",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.List{
				// As the docs says, the current maximum number of unit IDs is 1.
				listvalidator.SizeAtMost(1),
			},
		},
		"create_group_chat": schema.BoolAttribute{
			Description:         "Whether to create a group chat for the department.",
			MarkdownDescription: "Whether to create a group chat for the department.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(false),
		},
		"chat_id": schema.StringAttribute{
			Description:         "Department group chat ID.",
			MarkdownDescription: "Department group chat ID.",
			Computed:            true,
		},
		"leaders": schema.ListNestedAttribute{
			Description:         "Head of department.",
			MarkdownDescription: "Head of department.",
			Optional:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"leader_type": schema.Int64Attribute{
						Description:         "Person in charge type.",
						MarkdownDescription: "Person in charge type.",
						Required:            true,
						Validators: []validator.Int64{
							int64validator.OneOf(1, 2),
						},
					},
					"leader_id": schema.StringAttribute{
						Description:         "Person in charge ID.",
						MarkdownDescription: "Person in charge ID.",
						Required:            true,
					},
				},
			},
		},
		"group_chat_employee_types": schema.ListAttribute{
			Description:         "Department group employee type restriction.",
			MarkdownDescription: "Department group employee type restriction.",
			Optional:            true,
			ElementType:         types.Int64Type,
		},
		"member_count": schema.Int64Attribute{
			Description:         "Number of users under the department.",
			MarkdownDescription: "Number of users under the department.",
			Computed:            true,
		},
	}

	for k, v := range baseAttributes {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		Description:         "Manages department in Lark",
		MarkdownDescription: "Manages department in Lark",
		Attributes:          attributes,
	}
}

func (r *departmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *departmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data departmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	departmentId := data.DepartmentId.ValueString()
	if departmentId != "" {
		_, err := common.DepartmentGetByDepartmentIDAPI(ctx, r.client, data.DepartmentId.ValueString())
		if err == nil {
			resp.Diagnostics.AddError("API Error Reading Department", "Department with the same ID already exists")
			return
		}
	}

	tempRequestBody, err := r.modelToRequest(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Department", err.Error())
		return
	}

	requestBody := common.DepartmentCreateRequest{
		BaseDepartment: tempRequestBody,
		DepartmentID:   departmentId,
	}

	departmentCreateResponse, err := common.DepartmentCreateAPI(ctx, r.client, requestBody)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Department", err.Error())
		return
	}

	data.DepartmentId = types.StringValue(departmentCreateResponse.Data.Department.DepartmentID)
	data.OpenDepartmentId = types.StringValue(departmentCreateResponse.Data.Department.OpenDepartmentID)
	data.ChatID = types.StringValue(departmentCreateResponse.Data.Department.ChatID)
	data.MemberCount = types.Int64Value(int64(departmentCreateResponse.Data.Department.MemberCount))
	data.Id = types.StringValue(common.ConstructID(common.RESOURCE, common.DEPARTMENT, departmentCreateResponse.Data.Department.DepartmentID))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *departmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data departmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentDepartmentId := data.ParentDepartmentId.ValueString()
	_, err := common.DepartmentGetByDepartmentIDAPI(ctx, r.client, parentDepartmentId)
	if err != nil && parentDepartmentId != "0" {
		resp.Diagnostics.AddError("API Error Reading Department", "Parent department not found")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *departmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan departmentResourceModel
	var state departmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate the plan
	if state.CreateGroupChat.ValueBool() && !plan.CreateGroupChat.ValueBool() {
		resp.Diagnostics.AddError("API Error Updating User Group", "Cannot disable group chat for department")
		return
	}

	tempRequestBody, err := r.modelToRequest(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Department", err.Error())
		return
	}

	requestBody := common.DepartmentUpdateRequest{
		BaseDepartment: tempRequestBody,
	}

	departmentUpdateResponse, err := common.DepartmentUpdateAPI(ctx, r.client, state.OpenDepartmentId.ValueString(), requestBody)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Department", err.Error())
		return
	}

	if plan.DepartmentId.ValueString() != departmentUpdateResponse.Data.Department.DepartmentID {
		_, err = common.DepartmentUpdateIDAPI(ctx, r.client, state.OpenDepartmentId.ValueString(), common.DepartmentUpdateIDRequest{
			NewDepartmentID: plan.DepartmentId.ValueString(),
		})

		if err != nil {
			resp.Diagnostics.AddError("API Error Updating Department ID", err.Error())
			return
		}
	} else {
		plan.DepartmentId = types.StringValue(departmentUpdateResponse.Data.Department.DepartmentID)
	}

	plan.OpenDepartmentId = types.StringValue(departmentUpdateResponse.Data.Department.OpenDepartmentID)
	plan.ChatID = types.StringValue(departmentUpdateResponse.Data.Department.ChatID)
	plan.MemberCount = types.Int64Value(int64(departmentUpdateResponse.Data.Department.MemberCount))
	plan.Id = state.Id
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *departmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan departmentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := common.DepartmentDeleteAPI(ctx, r.client, plan.OpenDepartmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Department", err.Error())
		return
	}
}

func (r *departmentResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if r.client == nil {
		return []resource.ConfigValidator{}
	}
	return []resource.ConfigValidator{
		local_validator.NewUserIDValidator("leader_user_id", false, true, common.OPEN_ID, r.client),
	}
}

func (r *departmentResource) modelToRequest(ctx context.Context, data *departmentResourceModel) (common.BaseDepartment, error) {

	_, err := common.DepartmentGetByDepartmentIDAPI(ctx, r.client, data.ParentDepartmentId.ValueString())
	if err != nil && data.ParentDepartmentId.ValueString() != "0" {
		return common.BaseDepartment{}, err
	}

	// Special case handle null, need use pointer to avoid nil pointer dereference
	i18nName := common.I18nName{}
	if data.I18nName != nil {
		i18nName = common.I18nName{
			ZhCn: data.I18nName.ZhCn.ValueString(),
			JaJp: data.I18nName.JaJp.ValueString(),
			EnUs: data.I18nName.EnUs.ValueString(),
		}
	}

	leaders := []common.DepartmentLeader{}
	for _, leader := range data.Leaders {
		leaders = append(leaders, common.DepartmentLeader{
			LeaderType: leader.LeaderType.ValueInt64(),
			LeaderID:   leader.LeaderID.ValueString(),
		})
	}

	groupChatEmployeeTypes := []int64{}
	for _, groupChatEmployeeType := range data.GroupChatEmployeeTypes {
		groupChatEmployeeTypes = append(groupChatEmployeeTypes, groupChatEmployeeType.ValueInt64())
	}

	return common.BaseDepartment{
		Name:                   data.Name.ValueString(),
		I18nName:               i18nName,
		ParentDepartmentID:     data.ParentDepartmentId.ValueString(),
		LeaderUserID:           data.LeaderUserID.ValueString(),
		Order:                  data.Order.ValueString(),
		UnitIDs:                common.StringValuesToStrings(data.UnitIDs),
		Leaders:                leaders,
		GroupChatEmployeeTypes: groupChatEmployeeTypes,
		CreateGroupChat:        data.CreateGroupChat.ValueBool(),
	}, nil
}

// We use modify plan when we need both plan and state when validating.
func (r *departmentResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *departmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// plan null means resource is being deleted.
	if req.Plan.Raw.IsNull() {
		return
	}

	// State null means resource is being created.
	if req.State.Raw.IsNull() {

		// Checking if the department id is available.
		if plan.DepartmentId.ValueString() != "" {
			_, err := common.DepartmentGetByDepartmentIDAPI(ctx, r.client, plan.DepartmentId.ValueString())
			if err == nil {
				resp.Diagnostics.AddError("API Error Getting Department", "Department with the same ID already exists")
				return
			}
		}

		return
	}

	if plan.DepartmentId.ValueString() != "" && plan.DepartmentId.ValueString() != state.DepartmentId.ValueString() {
		// Checking if the department id is available.
		_, err := common.DepartmentGetByDepartmentIDAPI(ctx, r.client, plan.DepartmentId.ValueString())
		if err == nil {
			resp.Diagnostics.AddError("API Error Getting Department", "Department with the same ID already exists")
			return
		}
	}
}
