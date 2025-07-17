// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/aganisatria/terraform-provider-lark/internal/common"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &docsSpaceFolderResource{}

func NewDocsSpaceFolderResource() resource.Resource {
	return &docsSpaceFolderResource{}
}

// docsSpaceFolderResource defines the resource implementation.
type docsSpaceFolderResource struct {
	client *common.LarkClient
}

// docsSpaceFolderResourceModel describes the resource data model.
type docsSpaceFolderResourceModel struct {
	BaseResourceModel
	Name              types.String `tfsdk:"name"`
	Token             types.String `tfsdk:"token"`
	ParentFolderToken types.String `tfsdk:"parent_folder_token"`
}

func (r *docsSpaceFolderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_docs_space_folder"
}

func (r *docsSpaceFolderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	baseAttributes := BaseSchemaResourceAttributes()
	attributes := map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "The name of the folder.",
			MarkdownDescription: "The name of the folder.",
			Required:            true,
		},
		"token": schema.StringAttribute{
			Description:         "The token of the folder.",
			MarkdownDescription: "The token of the folder.",
			Computed:            true,
		},
		"parent_folder_token": schema.StringAttribute{
			Description:         "Parent folder token. If not provided, the folder will be created in the root folder.",
			MarkdownDescription: "Parent folder token. If not provided, the folder will be created in the root folder.",
			Optional:            true,
			Computed:            true,
		},
	}

	for k, v := range baseAttributes {
		attributes[k] = v
	}

	resp.Schema = schema.Schema{
		Description:         "Manages a folder in Lark Docs Space.",
		MarkdownDescription: "Manages a folder in Lark Docs Space.",
		Attributes:          attributes,
	}
}

func (r *docsSpaceFolderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *docsSpaceFolderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data docsSpaceFolderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentToken := data.ParentFolderToken.ValueString()

	if data.ParentFolderToken.IsUnknown() || data.ParentFolderToken.IsNull() {
		rootFolderResp, err := common.RootFolderMetaGetAPI(ctx, r.client)
		if err != nil {
			resp.Diagnostics.AddError("API Error Getting Root Folder", err.Error())
			return
		}
		parentToken = rootFolderResp.Data.Token
	}

	createRequest := common.FolderCreateRequest{
		Name:        data.Name.ValueString(),
		FolderToken: parentToken,
	}

	createResponse, err := common.FolderCreateAPI(ctx, r.client, createRequest)
	if err != nil {
		resp.Diagnostics.AddError("API Error Creating Docs Space Folder", err.Error())
		return
	}

	data.Id = types.StringValue(createResponse.Data.Token)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))

	data.Token = types.StringValue(createResponse.Data.Token)
	data.ParentFolderToken = types.StringValue(parentToken)
	data.Name = types.StringValue(createRequest.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *docsSpaceFolderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data docsSpaceFolderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folderToken := data.Token.ValueString()
	metaResponse, err := common.FolderMetaGetAPI(ctx, r.client, folderToken)
	if err != nil {
		resp.Diagnostics.AddWarning("API Error Reading Docs Space Folder", fmt.Sprintf("Unable to get folder metadata: %s. The resource may have been deleted.", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}

	data.Name = types.StringValue(metaResponse.Data.Name)
	data.ParentFolderToken = types.StringValue(metaResponse.Data.ParentID)
	data.Token = types.StringValue(metaResponse.Data.Token)

	tflog.Debug(ctx, "Reading docs_space_folder state2", map[string]interface{}{
		"state_details": map[string]interface{}{
			"token":        data.Token.ValueString(),
			"name":         data.Name.ValueString(),
			"parent":       data.ParentFolderToken.ValueString(),
			"id":           data.Id.ValueString(),
			"last_updated": data.LastUpdated.ValueString(),
		},
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *docsSpaceFolderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state docsSpaceFolderResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oldFolderToken := state.Token.ValueString()

	hasNameChanged := !plan.Name.Equal(state.Name)
	hasParentChanged := !plan.ParentFolderToken.Equal(state.ParentFolderToken)

	// Since we don't have such update API. We need to do workaround
	// case 1: folder name change.
	// step 1: create new folder based on plan
	// step 2: get list of its contents on old folder
	// step 3: move every contents to new folder
	// step 4: delete old folder
	if hasNameChanged {
		//step 1
		destinationParentToken := plan.ParentFolderToken.ValueString()
		if plan.ParentFolderToken.IsUnknown() || plan.ParentFolderToken.IsNull() {
			rootFolderResp, err := common.RootFolderMetaGetAPI(ctx, r.client)
			if err != nil {
				resp.Diagnostics.AddError("API Error Getting Root Folder for New Folder", err.Error())
				return
			}
			destinationParentToken = rootFolderResp.Data.Token
		}

		createRequest := common.FolderCreateRequest{
			Name:        plan.Name.ValueString(),
			FolderToken: destinationParentToken,
		}
		newFolder, err := common.FolderCreateAPI(ctx, r.client, createRequest)
		if err != nil {
			resp.Diagnostics.AddError("API Error Creating New Folder for Rename", err.Error())
			return
		}
		newFolderToken := newFolder.Data.Token

		// step 2
		oldFolderChildren, err := common.FolderChildrenListAPI(ctx, r.client, oldFolderToken)
		if err != nil {
			resp.Diagnostics.AddError("API Error Listing Old Folder Children", fmt.Sprintf("Failed to list items in old folder %s: %s. New folder %s was created but cannot be populated.", oldFolderToken, err.Error(), newFolderToken))
			return
		}

		// step 3
		if len(oldFolderChildren.Data.Files) > 0 {
			for _, child := range oldFolderChildren.Data.Files {
				moveReq := common.FileMoveRequest{
					Type:        child.Type,
					FolderToken: newFolderToken,
				}
				_, err := common.FileMoveAPI(ctx, r.client, child.Token, moveReq)
				if err != nil {
					resp.Diagnostics.AddWarning("Failed to Move Item", fmt.Sprintf("Could not move item %s (%s) from old folder to new folder. You may need to move it manually.", child.Name, child.Token))
				}
			}
		}

		// step 4
		_, err = common.FileDeleteAPI(ctx, r.client, oldFolderToken, "folder")
		if err != nil {
			resp.Diagnostics.AddWarning("Failed to Delete Old Folder", fmt.Sprintf("The old folder %s could not be deleted after rename. You may need to delete it manually.", oldFolderToken))
		}

		plan.Id = types.StringValue(newFolderToken)
		plan.Token = types.StringValue(newFolderToken)
		plan.ParentFolderToken = types.StringValue(destinationParentToken)

		// case 2: only parent folder change
	} else if hasParentChanged {
		destinationParentToken := plan.ParentFolderToken.ValueString()
		if plan.ParentFolderToken.IsUnknown() || plan.ParentFolderToken.IsNull() {
			rootFolderResp, err := common.RootFolderMetaGetAPI(ctx, r.client)
			if err != nil {
				resp.Diagnostics.AddError("API Error Getting Root Folder for Move", err.Error())
				return
			}
			destinationParentToken = rootFolderResp.Data.Token
		}

		moveReq := common.FileMoveRequest{
			Type:        "folder",
			FolderToken: destinationParentToken,
		}
		_, err := common.FileMoveAPI(ctx, r.client, oldFolderToken, moveReq)
		if err != nil {
			resp.Diagnostics.AddError("API Error Moving Docs Space Folder", err.Error())
			return
		}

		plan.Id = state.Id
		plan.Token = state.Token
		plan.ParentFolderToken = types.StringValue(destinationParentToken)
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC3339))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *docsSpaceFolderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data docsSpaceFolderResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	folderToken := data.Token.ValueString()
	_, err := common.FileDeleteAPI(ctx, r.client, folderToken, "folder")
	if err != nil {
		resp.Diagnostics.AddError("API Error Deleting Docs Space Folder", err.Error())
		return
	}
}
