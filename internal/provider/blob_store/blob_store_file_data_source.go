/*
 * Copyright (c) 2019-present Sonatype, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package blob_store

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &fileBlobStoreDataSource{}
	_ datasource.DataSourceWithConfigure = &fileBlobStoreDataSource{}
)

// BlobStoreFileDataSource is a helper function to simplify the provider implementation.
func BlobStoreFileDataSource() datasource.DataSource {
	return &fileBlobStoreDataSource{}
}

// fileBlobStoreDataSource is the data source implementation.
type fileBlobStoreDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *fileBlobStoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_file"
}

// Schema defines the schema for the data source.
func (d *fileBlobStoreDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get a specific File Blob Store by it's name",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the Blob Store",
				Required:    true,
			},
			"path": schema.StringAttribute{
				Description: "The Path on disk of this File Blob Store",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"soft_quota": schema.SingleNestedAttribute{
				Description: "Soft Quota for this Blob Store",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						Description: "Soft Quota type",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
					"limit": schema.Int64Attribute{
						Description: "Quota limit",
						Required:    false,
						Optional:    false,
						Computed:    true,
					},
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *fileBlobStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.BlobStoreFileModel
	// var state model.BlobStoreFileModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	if data.Name.IsNull() {
		resp.Diagnostics.AddError("Name must not be empty.", "Name must be provided.")
		return
	}

	blobStore, httpResponse, err := d.Client.BlobStoreAPI.GetFileBlobStoreConfiguration(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		common.HandleApiError(
			common.ERROR_UNABLE_TO_READ_BLOB_STORE_FILE,
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	state := model.BlobStoreFileModel{
		Name: types.StringValue(data.Name.ValueString()),
		Path: types.StringValue(*blobStore.Path),
	}

	if blobStore.SoftQuota != nil && blobStore.SoftQuota.Type != nil {
		tflog.Debug(ctx, fmt.Sprintf("%v", blobStore.SoftQuota))
		state.SoftQuota = &model.BlobStoreSoftQuota{
			Type:  types.StringValue(*blobStore.SoftQuota.Type),
			Limit: types.Int64Value(*blobStore.SoftQuota.Limit),
		}
	}

	// Set state
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
