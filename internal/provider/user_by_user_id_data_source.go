// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &UserBasedOnUserIDDataSource{}

func NewUserBasedOnUserIDDataSource() datasource.DataSource {
	return &UserBasedOnUserIDDataSource{}
}

// UserBasedOnUserIDDataSource defines the data source implementation.
type UserBasedOnUserIDDataSource struct {
	client *common.LarkClient
}

type UserBasedOnUserID struct {
	UserID types.String `tfsdk:"user_id"`
	OpenID types.String `tfsdk:"open_id"`
}

// UserDataSourceModel describes the data source data model.
type UserBasedOnUserIDDataSourceModel struct {
	BaseResourceModel
	Users []UserBasedOnUserID `tfsdk:"users"`
}

func (d *UserBasedOnUserIDDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_based_on_user_id"
}

func (d *UserBasedOnUserIDDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
						Computed:            true,
					},
					"user_id": schema.StringAttribute{
						Description:         "Unique identifier of the user in the tenant, namely the user's user_id.",
						MarkdownDescription: "Unique identifier of the user in the tenant, namely the user's user_id.",
						Required:            true,
					},
				},
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

func (d *UserBasedOnUserIDDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserBasedOnUserIDDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserBasedOnUserIDDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	userIDs := []string{}
	for _, user := range data.Users {
		userIDs = append(userIDs, user.UserID.ValueString())
	}

	response, err := common.GetUsersByUserIDAPI(ctx, d.client, userIDs)
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
			UserID: types.StringValue(user.UserID),
			OpenID: types.StringValue(user.OpenID),
		})
	}
	data.Users = users

	data.Id = types.StringValue("user-data-source-user-id" + "_" + time.Now().Format("20060102150405"))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
