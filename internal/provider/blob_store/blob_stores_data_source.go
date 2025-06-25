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
	_ datasource.DataSource              = &blobStoresDataSource{}
	_ datasource.DataSourceWithConfigure = &blobStoresDataSource{}
)

// BlobStoresDataSource is a helper function to simplify the provider implementation.
func BlobStoresDataSource() datasource.DataSource {
	return &blobStoresDataSource{}
}

// blobStoresDataSource is the data source implementation.
type blobStoresDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *blobStoresDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_stores"
}

// Schema defines the schema for the data source.
func (d *blobStoresDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get all Blob Stores",
		Attributes: map[string]schema.Attribute{
			"blob_stores": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the Blob Store",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "Blob Store type",
							Required:    true,
						},
						"unavailable": schema.BoolAttribute{
							Description: "Whether the Blob Store is unavailable for use",
							Required:    true,
						},
						"blob_count": schema.Int64Attribute{
							Description: "Number of blobs in the Blob Store",
							Computed:    true,
						},
						"total_size_in_bytes": schema.Int64Attribute{
							Description: "Total size in bytes of the Blob Store",
							Computed:    true,
						},
						"available_space_in_bytes": schema.Int64Attribute{
							Description: "Available space in bytes for the Blob Store",
							Computed:    true,
						},
						"soft_quota": schema.SingleNestedAttribute{
							Description: "Soft Quota for this Blob Store",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"type": schema.StringAttribute{
									Description: "Soft Quota type",
									Required:    true,
								},
								"limit": schema.Int64Attribute{
									Description: "Quota limit",
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *blobStoresDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.BlobStoresModel

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	blobStores, _, err := d.Client.BlobStoreAPI.ListBlobStores(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Blob Stores",
			err.Error(),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Blob Stores", len(blobStores)))

	for _, blobStore := range blobStores {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s Blob Store", *blobStore.Name))

		blobStoreState := model.BlobStoreModel{
			Name:                  types.StringValue(*blobStore.Name),
			Type:                  types.StringValue(*blobStore.Type),
			Unavailable:           types.BoolValue(*blobStore.Unavailable),
			BlobCount:             types.Int64Value(*blobStore.BlobCount),
			TotalSizeInBytes:      types.Int64Value(*blobStore.TotalSizeInBytes),
			AvailableSpaceInBytes: types.Int64Value(*blobStore.AvailableSpaceInBytes),
		}

		if blobStore.SoftQuota != nil && blobStore.SoftQuota.Type != nil {
			tflog.Debug(ctx, fmt.Sprintf("%v", blobStore.SoftQuota))
			blobStoreState.SoftQuota = model.BlobStoreSoftQuota{
				Type:  types.StringValue(*blobStore.SoftQuota.Type),
				Limit: types.Int64Value(*blobStore.SoftQuota.Limit),
			}
		}

		state.BlobStores = append(state.BlobStores, blobStoreState)

		tflog.Debug(ctx, fmt.Sprintf("   Appended: %p", state.BlobStores))
	}

	// For test framework
	// state.ID = types.StringValue("placeholder")

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
