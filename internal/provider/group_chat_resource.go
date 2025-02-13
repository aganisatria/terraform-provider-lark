// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	. "github.com/aganisatria/terraform-provider-lark/internal/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &groupChatResource{}

func NewGroupChatResource() resource.Resource {
	return &groupChatResource{}
}

// groupChatResource defines the resource implementation.
type groupChatResource struct {
	client *common.LarkClient
}

// groupChatResourceModel describes the resource data model.
// fields that need to be configured by user.
type groupChatResourceModel struct {
	BaseResourceModel
	ChatID                 types.String           `tfsdk:"chat_id"`
	Avatar                 types.String           `tfsdk:"avatar"`
	Name                   types.String           `tfsdk:"name"`
	Description            types.String           `tfsdk:"description"`
	I18nNames              *I18nName              `tfsdk:"i18n_names"`
	UserIDList             []types.String         `tfsdk:"user_id_list"`
	BotIDList              []types.String         `tfsdk:"bot_id_list"`
	GroupMessageType       types.String           `tfsdk:"group_message_type"`
	ChatMode               types.String           `tfsdk:"chat_mode"`
	ChatType               types.String           `tfsdk:"chat_type"`
	JoinMessageVisibility  types.String           `tfsdk:"join_message_visibility"`
	LeaveMessageVisibility types.String           `tfsdk:"leave_message_visibility"`
	MembershipApproval     types.String           `tfsdk:"membership_approval"`
	RestrictedModeSetting  *RestrictedModeSetting `tfsdk:"restricted_mode_setting"`
	UrgentSetting          types.String           `tfsdk:"urgent_setting"`
	VideoConferenceSetting types.String           `tfsdk:"video_conference_setting"`
	EditPermission         types.String           `tfsdk:"edit_permission"`
	HideMemberCountSetting types.String           `tfsdk:"hide_member_count_setting"`
	AddMemberPermission    types.String           `tfsdk:"add_member_permission"`
	ShareCardPermission    types.String           `tfsdk:"share_card_permission"`
	AtAllPermission        types.String           `tfsdk:"at_all_permission"`
	GroupType              types.String           `tfsdk:"group_type"` // To validate how many members and administrators can be added.
}

func (r *groupChatResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_chat"
}

func (r *groupChatResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"chat_id": schema.StringAttribute{
			Description:         "Group chat ID.",
			MarkdownDescription: "Group chat ID.",
			Computed:            true,
		},
		"avatar": schema.StringAttribute{
			Description:         "URL group of photo profile. You can check on this first https://open.larksuite.com/document/server-docs/im-v1/image/create",
			MarkdownDescription: "URL group of photo profile. You can check on this first https://open.larksuite.com/document/server-docs/im-v1/image/create",
			Optional:            true,
		},
		"name": schema.StringAttribute{
			Description:         "Group chat name. The length of the public group name must be at least 2 characters, and if the private group does not fill in the group name, the group name defaults to \"(no title)\".",
			MarkdownDescription: "Group chat name. The length of the public group name must be at least 2 characters, and if the private group does not fill in the group name, the group name defaults to \"(no title)\".",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				GroupNameValidator(),
			},
		},
		"description": schema.StringAttribute{
			Description:         "Group chat description.",
			MarkdownDescription: "Group chat description.",
			Optional:            true,
		},
		"i18n_names": schema.SingleNestedAttribute{
			Description:         "Internationalized group chat name.",
			MarkdownDescription: "Internationalized group chat name.",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"zh_cn": schema.StringAttribute{
					Description:         "Group chat's Chinese name.",
					MarkdownDescription: "Group chat's Chinese name.",
					Optional:            true,
				},
				"ja_jp": schema.StringAttribute{
					Description:         "Group chat's Japanese name.",
					MarkdownDescription: "Group chat's Japanese name.",
					Optional:            true,
				},
				"en_us": schema.StringAttribute{
					Description:         "Group chat's English name.",
					MarkdownDescription: "Group chat's English name.",
					Optional:            true,
				},
			},
		},
		// We remove owner id, because we dont want bot ability to delete group chat is deleted.
		// PLACEHOLDER OWNER ID.
		"user_id_list": schema.ListAttribute{
			Description:         "Group chat user ID list.",
			MarkdownDescription: "Group chat user ID list.",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"bot_id_list": schema.ListAttribute{
			Description:         "Group chat bot ID list.",
			MarkdownDescription: "Group chat bot ID list.",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"group_message_type": schema.StringAttribute{
			Description:         "Group chat message type.",
			MarkdownDescription: "Group chat message type.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("chat"),
			Validators: []validator.String{
				stringvalidator.OneOf("chat", "thread"),
			},
		},
		"chat_mode": schema.StringAttribute{
			Description:         "Group chat mode.",
			MarkdownDescription: "Group chat mode.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("group"),
			Validators: []validator.String{
				// Based on Lark API, only group is supported
				stringvalidator.OneOf("group"),
			},
		},
		"chat_type": schema.StringAttribute{
			Description:         "Group chat type.",
			MarkdownDescription: "Group chat type.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("private"),
			Validators: []validator.String{
				stringvalidator.OneOf("private", "public"),
			},
		},
		"join_message_visibility": schema.StringAttribute{
			Description:         "Group chat join message visibility.",
			MarkdownDescription: "Group chat join message visibility.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("all_members"),
			Validators: []validator.String{
				stringvalidator.OneOf("all_members", "only_owner", "not_anyone"),
			},
		},
		"leave_message_visibility": schema.StringAttribute{
			Description:         "Group chat leave message visibility.",
			MarkdownDescription: "Group chat leave message visibility.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("all_members"),
			Validators: []validator.String{
				stringvalidator.OneOf("all_members", "only_owner", "not_anyone"),
			},
		},
		"membership_approval": schema.StringAttribute{
			Description:         "Group chat membership approval.",
			MarkdownDescription: "Group chat membership approval.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("no_approval_required"),
			Validators: []validator.String{
				stringvalidator.OneOf("no_approval_required", "approval_required"),
			},
		},
		"restricted_mode_setting": schema.SingleNestedAttribute{
			Description:         "Group chat restricted mode setting.",
			MarkdownDescription: "Group chat restricted mode setting.",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"status": schema.BoolAttribute{
					Description:         "Whether the restricted mode is enabled.",
					MarkdownDescription: "Whether the restricted mode is enabled.",
					Optional:            true,
				},
				"screenshot_has_permission_setting": schema.StringAttribute{
					Description:         "Whether the screenshot permission is enabled.",
					MarkdownDescription: "Whether the screenshot permission is enabled.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("all_members", "not_anyone"),
					},
				},
				"download_has_permission_setting": schema.StringAttribute{
					Description:         "Whether the download permission is enabled.",
					MarkdownDescription: "Whether the download permission is enabled.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("all_members", "not_anyone"),
					},
				},
				"message_has_permission_setting": schema.StringAttribute{
					Description:         "Whether the message permission is enabled.",
					MarkdownDescription: "Whether the message permission is enabled.",
					Optional:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("all_members", "not_anyone"),
					},
				},
			},
		},
		"urgent_setting": schema.StringAttribute{
			Description:         "Group chat urgent setting.",
			MarkdownDescription: "Group chat urgent setting.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("all_members", "not_anyone"),
			},
		},
		"video_conference_setting": schema.StringAttribute{
			Description:         "Group chat video conference setting.",
			MarkdownDescription: "Group chat video conference setting.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("all_members", "not_anyone"),
			},
		},
		"edit_permission": schema.StringAttribute{
			Description:         "Group chat edit permission.",
			MarkdownDescription: "Group chat edit permission.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("all_members", "only_owner"),
			},
		},
		"hide_member_count_setting": schema.StringAttribute{
			Description:         "Group chat hide member count setting.",
			MarkdownDescription: "Group chat hide member count setting.",
			Optional:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("all_members", "only_owner"),
			},
		},
		"add_member_permission": schema.StringAttribute{
			Description:         "Group chat add member permission.",
			MarkdownDescription: "Group chat add member permission.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("all_members"),
			Validators: []validator.String{
				stringvalidator.OneOf("all_members", "only_owner"),
			},
		},
		"share_card_permission": schema.StringAttribute{
			Description:         "Group chat share card permission.",
			MarkdownDescription: "Group chat share card permission.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("allowed"),
			Validators: []validator.String{
				stringvalidator.OneOf("allowed", "not_allowed"),
			},
		},
		"at_all_permission": schema.StringAttribute{
			Description:         "Group chat at all permission.",
			MarkdownDescription: "Group chat at all permission.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("all_members"),
			Validators: []validator.String{
				stringvalidator.OneOf("all_members", "only_owner"),
			},
		},
		"group_type": schema.StringAttribute{
			Description:         "Group chat type. Whether it is a common group or a super large group. Only 5 super large groups can be created for each enterprise tenant.",
			MarkdownDescription: "Group chat type. Whether it is a common group or a super large group. Only 5 super large groups can be created for each enterprise tenant.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString("common"),
			Validators: []validator.String{
				stringvalidator.OneOf("common", "super_large"),
			},
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

func (r *groupChatResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *groupChatResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data groupChatResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var i18nNames common.I18nName
	if data.I18nNames != nil {
		i18nNames = common.I18nName{
			ZhCn: data.I18nNames.ZhCn.ValueString(),
			JaJp: data.I18nNames.JaJp.ValueString(),
			EnUs: data.I18nNames.EnUs.ValueString(),
		}
	}

	userIDList := []string{}
	for _, userID := range data.UserIDList {
		userIDList = append(userIDList, userID.ValueString())
	}

	botIDList := []string{}
	for _, botID := range data.BotIDList {
		botIDList = append(botIDList, botID.ValueString())
	}

	requestBody := common.GroupChatCreateRequest{
		Avatar:                 data.Avatar.ValueString(),
		Name:                   data.Name.ValueString(),
		Description:            data.Description.ValueString(),
		I18nNames:              i18nNames,
		UserIDList:             userIDList,
		BotIDList:              append(botIDList, r.client.AppID),
		GroupMessageType:       data.GroupMessageType.ValueString(),
		ChatMode:               data.ChatMode.ValueString(),
		ChatType:               data.ChatType.ValueString(),
		JoinMessageVisibility:  data.JoinMessageVisibility.ValueString(),
		LeaveMessageVisibility: data.LeaveMessageVisibility.ValueString(),
		MembershipApproval:     data.MembershipApproval.ValueString(),
		UrgentSetting:          data.UrgentSetting.ValueString(),
		VideoConferenceSetting: data.VideoConferenceSetting.ValueString(),
		EditPermission:         data.EditPermission.ValueString(),
		HideMemberCountSetting: data.HideMemberCountSetting.ValueString(),
	}

	if data.RestrictedModeSetting != nil && data.RestrictedModeSetting.Status.ValueBool() {
		requestBody.RestrictedModeSetting = &common.RestrictedModeSetting{
			Status:                         data.RestrictedModeSetting.Status.ValueBool(),
			ScreenshotHasPermissionSetting: data.RestrictedModeSetting.ScreenshotHasPermissionSetting.ValueString(),
			DownloadHasPermissionSetting:   data.RestrictedModeSetting.DownloadHasPermissionSetting.ValueString(),
			MessageHasPermissionSetting:    data.RestrictedModeSetting.MessageHasPermissionSetting.ValueString(),
		}
	}

	groupChatCreateResponse, err := common.GroupChatCreateAPI(ctx, r.client, requestBody)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Group Chat", err.Error())
		return
	}

	data.Id = types.StringValue(groupChatCreateResponse.Data.ChatID)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	data.ChatID = types.StringValue(groupChatCreateResponse.Data.ChatID)
	data.Name = types.StringValue(groupChatCreateResponse.Data.Name)
	data.AddMemberPermission = types.StringValue(groupChatCreateResponse.Data.AddMemberPermission)
	data.ShareCardPermission = types.StringValue(groupChatCreateResponse.Data.ShareCardPermission)
	data.AtAllPermission = types.StringValue(groupChatCreateResponse.Data.AtAllPermission)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupChatResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data groupChatResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupChatGetResponse, err := common.GroupChatGetAPI(ctx, r.client, data.ChatID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Reading Group Chat", err.Error())
		return
	}

	data.AddMemberPermission = types.StringValue(groupChatGetResponse.Data.AddMemberPermission)
	data.ShareCardPermission = types.StringValue(groupChatGetResponse.Data.ShareCardPermission)
	data.AtAllPermission = types.StringValue(groupChatGetResponse.Data.AtAllPermission)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupChatResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupChatResourceModel
	var state groupChatResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var i18nNames common.I18nName
	if plan.I18nNames != nil {
		i18nNames = common.I18nName{
			ZhCn: plan.I18nNames.ZhCn.ValueString(),
			JaJp: plan.I18nNames.JaJp.ValueString(),
			EnUs: plan.I18nNames.EnUs.ValueString(),
		}
	}

	requestBody := common.GroupChatUpdateRequest{
		Avatar:                 plan.Avatar.ValueString(),
		Name:                   plan.Name.ValueString(),
		Description:            plan.Description.ValueString(),
		I18nNames:              i18nNames,
		AddMemberPermission:    plan.AddMemberPermission.ValueString(),
		ShareCardPermission:    plan.ShareCardPermission.ValueString(),
		AtAllPermission:        plan.AtAllPermission.ValueString(),
		EditPermission:         plan.EditPermission.ValueString(),
		JoinMessageVisibility:  plan.JoinMessageVisibility.ValueString(),
		LeaveMessageVisibility: plan.LeaveMessageVisibility.ValueString(),
		MembershipApproval:     plan.MembershipApproval.ValueString(),
		ChatType:               plan.ChatType.ValueString(),
		GroupMessageType:       plan.GroupMessageType.ValueString(),
		UrgentSetting:          plan.UrgentSetting.ValueString(),
		VideoConferenceSetting: plan.VideoConferenceSetting.ValueString(),
		HideMemberCountSetting: plan.HideMemberCountSetting.ValueString(),
	}

	if plan.RestrictedModeSetting != nil && plan.RestrictedModeSetting.Status.ValueBool() {
		requestBody.RestrictedModeSetting = &common.RestrictedModeSetting{
			Status:                         plan.RestrictedModeSetting.Status.ValueBool(),
			ScreenshotHasPermissionSetting: plan.RestrictedModeSetting.ScreenshotHasPermissionSetting.ValueString(),
			DownloadHasPermissionSetting:   plan.RestrictedModeSetting.DownloadHasPermissionSetting.ValueString(),
			MessageHasPermissionSetting:    plan.RestrictedModeSetting.MessageHasPermissionSetting.ValueString(),
		}
	}

	_, err := common.GroupChatUpdateAPI(ctx, r.client, state.ChatID.ValueString(), requestBody)
	if err != nil {
		resp.Diagnostics.AddError("API Error Updating Group Chat", err.Error())
		return
	}

	plan.Id = types.StringValue(state.Id.ValueString())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	plan.ChatID = types.StringValue(state.ChatID.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *groupChatResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan groupChatResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := common.GroupChatDeleteAPI(ctx, r.client, plan.ChatID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Group Chat", err.Error())
		return
	}
}

func (r *groupChatResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	if r.client == nil {
		return []resource.ConfigValidator{}
	}

	return []resource.ConfigValidator{
		NewUserIDValidator("user_id_list", true, r.client),
	}
}

// We use modify plan when we need state when validating.
func (r *groupChatResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan, state *groupChatResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// plan null means resource is being deleted.
	if req.Plan.Raw.IsNull() {
		return
	}

	// State null means resource is being created.
	if req.State.Raw.IsNull() {
		if plan.AddMemberPermission.ValueString() != "all_members" {
			resp.Diagnostics.AddError("Add Member Permission Must Be All Members when creating terraform resource", "Add Member Permission Must Be All Members when creating group chat")
			return
		}

		if plan.ShareCardPermission.ValueString() != "allowed" {
			resp.Diagnostics.AddError("Share Card Permission Must Be Allowed when creating terraform resource", "Share Card Permission Must Be Allowed when creating group chat")
			return
		}
	}

	if (plan.AddMemberPermission.ValueString() == "only_owner") && (plan.ShareCardPermission.ValueString() != "not_allowed") {
		resp.Diagnostics.AddError("If Add Member Permission is Only Owner, Share Card Permission is must be allowed", "If Add Member Permission is Only Owner, Share Card Permission is must be allowed")
		return
	}

	if (plan.AddMemberPermission.ValueString() == "all_members") && (plan.ShareCardPermission.ValueString() != "allowed") {
		resp.Diagnostics.AddError("If Add Member Permission is All Members, Share Card Permission is must be not allowed", "If Add Member Permission is All Members, Share Card Permission is must be not allowed")
		return
	}

	if plan.RestrictedModeSetting.Status.ValueBool() {
		if (plan.RestrictedModeSetting.ScreenshotHasPermissionSetting.ValueString() == "all_members") &&
			(plan.RestrictedModeSetting.DownloadHasPermissionSetting.ValueString() == "all_members") &&
			(plan.RestrictedModeSetting.MessageHasPermissionSetting.ValueString() == "all_members") {
			resp.Diagnostics.AddError(
				"If Restricted Mode Setting is Enabled, Screenshot Has Permission Setting, Download Has Permission Setting, and Message Has Permission Setting must at least 1 not be all members",
				"If Restricted Mode Setting is Enabled, Screenshot Has Permission Setting, Download Has Permission Setting, and Message Has Permission Setting must at least 1 not be all members",
			)
			return
		}
	}
}
