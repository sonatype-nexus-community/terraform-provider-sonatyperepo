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
	tfschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
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
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get all Blob Stores",
		Attributes: map[string]tfschema.Attribute{
			"blob_stores": schema.DataSourceComputedListNestedAttribute(
				"List of Blob Stores",
				tfschema.NestedAttributeObject{
					Attributes: map[string]tfschema.Attribute{
						"name":                     schema.DataSourceRequiredString("Name of the Blob Store"),
						"type":                     schema.DataSourceRequiredString("Blob Store type"),
						"unavailable":              schema.DataSourceRequiredBool("Whether the Blob Store is unavailable for use"),
						"blob_count":               schema.DataSourceComputedInt64("Number of blobs in the Blob Store"),
						"total_size_in_bytes":      schema.DataSourceComputedInt64("Total size in bytes of the Blob Store"),
						"available_space_in_bytes": schema.DataSourceComputedInt64("Available space in bytes for the Blob Store"),
						"soft_quota": schema.DataSourceOptionalSingleNestedAttribute(
							"Soft Quota for this Blob Store",
							map[string]tfschema.Attribute{
								"type":  schema.DataSourceRequiredString("Soft Quota type"),
								"limit": schema.DataSourceComputedInt64("Quota limit"),
							},
						),
					},
				},
			),
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

	blobStores, httpResponse, err := d.Client.BlobStoreAPI.ListBlobStores(ctx).Execute()
	if err != nil {
		sharederr.HandleAPIError(
			"Unable to Read Blob Stores",
			&err,
			httpResponse,
			&resp.Diagnostics,
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

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
