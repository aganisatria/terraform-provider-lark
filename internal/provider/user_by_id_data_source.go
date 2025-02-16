// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &UserByIDDataSource{}
var _ datasource.DataSourceWithConfigValidators = &UserByIDDataSource{}

func NewUserByIDDataSource() datasource.DataSource {
	return &UserByIDDataSource{}
}

// UserByIDDataSource defines the data source implementation.
type UserByIDDataSource struct {
	client *common.LarkClient
}

type UserBasedOnUserID struct {
	UserID  types.String `tfsdk:"user_id"`
	OpenID  types.String `tfsdk:"open_id"`
	UnionID types.String `tfsdk:"union_id"`
}

// UserByIDDataSourceModel describes the data source data model.
type UserByIDDataSourceModel struct {
	BaseResourceModel
	Users []UserBasedOnUserID `tfsdk:"users"`
	KeyID types.String        `tfsdk:"key_id"`
}

func (d *UserByIDDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_by_id"
}

func (d *UserByIDDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"users": schema.ListNestedAttribute{
			Description:         "List of users.",
			MarkdownDescription: "List of users.",
			Required:            true,
			Validators: []validator.List{
				listvalidator.UniqueValues(),
				listvalidator.SizeAtLeast(1),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"open_id": schema.StringAttribute{
						Description:         "Unique identifier of the user in the tenant, namely the user's open_id.",
						MarkdownDescription: "Unique identifier of the user in the tenant, namely the user's open_id.",
						Optional:            true,
						Computed:            true,
					},
					"user_id": schema.StringAttribute{
						Description:         "Unique identifier of the user in the tenant, namely the user's user_id.",
						MarkdownDescription: "Unique identifier of the user in the tenant, namely the user's user_id.",
						Optional:            true,
						Computed:            true,
					},
					"union_id": schema.StringAttribute{
						Description:         "Unique identifier of the user in the tenant, namely the user's union_id.",
						MarkdownDescription: "Unique identifier of the user in the tenant, namely the user's union_id.",
						Optional:            true,
						Computed:            true,
					},
				},
			},
		},
		"key_id": schema.StringAttribute{
			Description:         "Unique identifier of the user in the tenant, namely the user's key_id.",
			MarkdownDescription: "Unique identifier of the user in the tenant, namely the user's key_id.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf("user_id", "union_id", "open_id"),
			},
		},
	}

	for k, v := range baseAttributes {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		Description:         "Manages user data in Lark",
		MarkdownDescription: "Manages user data in Lark",
		Attributes:          attributes,
	}
}

func (d *UserByIDDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*common.LarkClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *LarkClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *UserByIDDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserByIDDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := getIDsFromUsers(ctx, data.Users, data.KeyID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid User ID", err.Error())
		return
	}

	if len(ids) > 0 {
		response, err := common.GetUsersByIDAPI(ctx, d.client, ids, common.UserIDType(data.KeyID.ValueString()))
		if err != nil {
			resp.Diagnostics.AddError("API Error Getting User ID by User ID", err.Error())
			return
		}

		users := make([]UserBasedOnUserID, 0, len(response.Data.Items))
		for _, user := range response.Data.Items {
			if user.UserID == "" {
				resp.Diagnostics.AddError("API Error Getting User ID by User ID", "User ID is not found for user ID: "+user.UserID)
				return
			}

			users = append(users, UserBasedOnUserID{
				UserID:  types.StringValue(user.UserID),
				OpenID:  types.StringValue(user.OpenID),
				UnionID: types.StringValue(user.UnionID),
			})
		}
		data.Users = users
	}

	data.Id = types.StringValue("user-data-source-user-id" + "_" + time.Now().Format("20060102150405"))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// keyIDInUsersValidator is a validator to ensure that the value of key_id
// (user_id / union_id / open_id) is present in each element of the users list.
type keyIDInUsersValidator struct {
	client *common.LarkClient
}

func (v *keyIDInUsersValidator) Description(ctx context.Context) string {
	return "Ensure that each element in the 'users' list has a non-empty value for the attribute specified by 'key_id'."
}

func (v *keyIDInUsersValidator) MarkdownDescription(ctx context.Context) string {
	return "Ensure that each element in the 'users' list has a non-empty value for the attribute specified by 'key_id'."
}

func (v *keyIDInUsersValidator) ValidateDataSource(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var config UserByIDDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := getIDsFromUsers(ctx, config.Users, config.KeyID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid User ID", err.Error())
		return
	}

	// We can't do API call here because the client is not initialized yet
	// So we need to do it in the Read function
}

func (d *UserByIDDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		&keyIDInUsersValidator{client: d.client},
	}
}

func getIDsFromUsers(ctx context.Context, users []UserBasedOnUserID, keyID string) ([]string, error) {
	ids := []string{}
	err_ids := []string{}
	for _, user := range users {
		var value string
		switch keyID {
		case "user_id":
			value = user.UserID.ValueString()

			if (!user.UnionID.IsUnknown() && !user.UnionID.IsNull()) || (!user.OpenID.IsUnknown() && !user.OpenID.IsNull()) {
				err_ids = append(err_ids, user.UserID.ValueString())
			}

		case "union_id":
			value = user.UnionID.ValueString()

			if (!user.UserID.IsUnknown() && !user.UserID.IsNull()) || (!user.OpenID.IsUnknown() && !user.OpenID.IsNull()) {
				err_ids = append(err_ids, user.UnionID.ValueString())
			}

		case "open_id":
			value = user.OpenID.ValueString()

			if (!user.UserID.IsUnknown() && !user.UserID.IsNull()) || (!user.UnionID.IsUnknown() && !user.UnionID.IsNull()) {
				err_ids = append(err_ids, user.OpenID.ValueString())
			}
		}

		if value == "" {
			err_ids = append(err_ids, user.UserID.ValueString())
		}

		ids = append(ids, value)
	}

	var errMsg string
	if len(err_ids) > 0 {
		errMsg = fmt.Sprintf("user %s is missing a non-empty value for the specified key_id: %s", strings.Join(err_ids, ", "), keyID)
	}

	if errMsg != "" {
		return ids, fmt.Errorf(errMsg)
	}

	return ids, nil
}
