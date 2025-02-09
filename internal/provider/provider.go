// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/common"
)

// Ensure LarkProvider satisfies various provider interfaces.
var _ provider.Provider = &LarkProvider{}

// LarkProvider defines the provider implementation.
type LarkProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// LarkProviderModel describes the provider data model.
type LarkProviderModel struct {
	AppId     types.String `tfsdk:"app_id"`
	AppSecret types.String `tfsdk:"app_secret"`
}

func (p *LarkProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "lark"
	resp.Version = p.version
}

func (p *LarkProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				Description:         "The app ID for authenticating with Lark API",
				MarkdownDescription: "The App ID for authenticating with Lark API",
			},
			"app_secret": schema.StringAttribute{
				Required:            true,
				Sensitive:           true,
				Description:         "The app Secret for authenticating with Lark API",
				MarkdownDescription: "The App Secret for authenticating with Lark API",
			},
		},
	}
}

func (p *LarkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LarkProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.AppId.IsNull() || data.AppSecret.IsNull() || data.AppId.ValueString() == "" || data.AppSecret.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Lark API credentials",
			"The Lark API credentials (app_id and app_secret) are missing or invalid. Please check your provider configuration.",
		)
		return
	}

	tenantAccessToken, appAccessToken, err := common.GetAccessTokenAPI(data.AppId.ValueString(), data.AppSecret.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication"),
			"Failed to Authenticate",
			fmt.Sprintf("Unable to retrieve access token from Lark API: %s", err.Error()),
		)
		return
	}

	client := common.NewLarkClient(tenantAccessToken, appAccessToken, common.BASE_DELAY)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LarkProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *LarkProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *LarkProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *LarkProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LarkProvider{
			version: version,
		}
	}
}
