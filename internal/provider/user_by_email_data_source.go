// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &UserByEmailDataSource{}

func NewUserByEmailDataSource() datasource.DataSource {
	return &UserByEmailDataSource{}
}

// UserByEmailDataSource defines the data source implementation.
type UserByEmailDataSource struct {
	client *common.LarkClient
}

type UserByEmail struct {
	UserID types.String `tfsdk:"user_id"`
	Email  types.String `tfsdk:"email"`
}

// UserDataSourceModel describes the data source data model.
type UserByEmailDataSourceModel struct {
	BaseResourceModel
	Users []UserByEmail `tfsdk:"users"`
}

func (d *UserByEmailDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_by_email"
}

func (d *UserByEmailDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
					"user_id": schema.StringAttribute{
						Description:         "Unique identifier of the user in the tenant, namely the user's user_id.",
						MarkdownDescription: "Unique identifier of the user in the tenant, namely the user's user_id.",
						Computed:            true,
					},
					"email": schema.StringAttribute{
						Description:         "User's email.",
						MarkdownDescription: "User's email.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`), "Email must be a valid email address"),
						},
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

func (d *UserByEmailDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserByEmailDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserByEmailDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	emails := []string{}
	for _, user := range data.Users {
		emails = append(emails, user.Email.ValueString())
	}

	response, err := common.GetUserIdByEmailsAPI(ctx, d.client, common.UserInfoBatchGetRequest{
		Emails: emails,
	})
	if err != nil {
		resp.Diagnostics.AddError("API Error Getting User ID by Emails", err.Error())
		return
	}

	users := make([]UserByEmail, 0, len(response.Data.UserList))
	for _, user := range response.Data.UserList {
		if user.UserID == "" {
			resp.Diagnostics.AddError("API Error Getting User ID by Emails", "User ID is not found for email: "+user.Email)
			return
		}

		users = append(users, UserByEmail{
			UserID: types.StringValue(user.UserID),
			Email:  types.StringValue(user.Email),
		})
	}
	data.Users = users

	data.Id = types.StringValue(common.ConstructID(common.DATA_SOURCE, common.USER_BY_EMAIL, ""))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
