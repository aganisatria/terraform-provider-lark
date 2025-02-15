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

var _ datasource.DataSource = &UserBasedOnUnionIDDataSource{}

func NewUserBasedOnUnionIDDataSource() datasource.DataSource {
	return &UserBasedOnUnionIDDataSource{}
}

// UserBasedOnUnionIDDataSource defines the data source implementation.
type UserBasedOnUnionIDDataSource struct {
	client *common.LarkClient
}

type UserBasedOnUnionID struct {
	UnionID types.String `tfsdk:"union_id"`
	OpenID  types.String `tfsdk:"open_id"`
}

// UserDataSourceModel describes the data source data model.
type UserBasedOnUnionIDDataSourceModel struct {
	BaseResourceModel
	Users []UserBasedOnUnionID `tfsdk:"users"`
}

func (d *UserBasedOnUnionIDDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_based_on_union_id"
}

func (d *UserBasedOnUnionIDDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
					"union_id": schema.StringAttribute{
						Description:         "Unique identifier of the user in the tenant, namely the user's union_id.",
						MarkdownDescription: "Unique identifier of the user in the tenant, namely the user's union_id.",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.RegexMatches(regexp.MustCompile("^on_[a-zA-Z0-9]+$"), "Union ID must be in the format of 'on_<alphanumeric characters>'"),
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

func (d *UserBasedOnUnionIDDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserBasedOnUnionIDDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data UserBasedOnUnionIDDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	UnionIDs := []string{}
	for _, user := range data.Users {
		UnionIDs = append(UnionIDs, user.UnionID.ValueString())
	}

	response, err := common.GetUsersByUnionIDAPI(ctx, d.client, UnionIDs)
	if err != nil {
		resp.Diagnostics.AddError("API Error Getting User ID by User ID", err.Error())
		return
	}

	users := make([]UserBasedOnUnionID, 0, len(response.Data.Items))
	for _, user := range response.Data.Items {
		if user.UnionID == "" {
			resp.Diagnostics.AddError("API Error Getting User ID by User ID", "User ID is not found for user ID: "+user.UnionID)
			return
		}

		users = append(users, UserBasedOnUnionID{
			UnionID: types.StringValue(user.UnionID),
			OpenID:  types.StringValue(user.OpenID),
		})
	}
	data.Users = users

	data.Id = types.StringValue("user-data-source-union-id" + "_" + time.Now().Format("20060102150405"))
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
