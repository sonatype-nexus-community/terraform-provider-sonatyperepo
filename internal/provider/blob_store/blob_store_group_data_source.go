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
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	tfschema "github.com/sonatype-nexus-community/terraform-provider-shared/schema"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &groupBlobStoreDataSource{}
	_ datasource.DataSourceWithConfigure = &groupBlobStoreDataSource{}
)

// BlobStoreGroupDataSource is a helper function to simplify the provider implementation.
func BlobStoreGroupDataSource() datasource.DataSource {
	return &groupBlobStoreDataSource{}
}

// groupBlobStoreDataSource is the data source implementation.
type groupBlobStoreDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *groupBlobStoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_group"
}

// Schema defines the schema for the data source.
func (d *groupBlobStoreDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Use this data source to get a specific File Blob Store by it's name",
		Attributes: map[string]dsschema.Attribute{
			"name":         tfschema.DataSourceRequiredString("Name of the Blob Store Group"),
			"fill_policy":  tfschema.DataSourceComputedString("Defines how writes are made to the member Blob Stores"),
			"last_updated": tfschema.DataSourceComputedString("The timestamp of when the resource was last updated"),
			"members": dsschema.SetAttribute{
				Description: "Set of the names of blob stores that are members of this group",
				ElementType: types.StringType,
				Computed:    true,
			},
			"soft_quota": tfschema.DataSourceComputedSingleNestedAttribute(
				"Soft Quota for this Blob Store",
				map[string]dsschema.Attribute{
					"type":  tfschema.DataSourceComputedString("Soft Quota type"),
					"limit": tfschema.DataSourceComputedInt64("Quota limit"),
				},
			),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *groupBlobStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.BlobStoreGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = d.GetAuthContext(ctx)

	if data.Name.IsNull() {
		resp.Diagnostics.AddError("Name must not be empty.", "Name must be provided.")
		return
	}

	apiResponse, httpResponse, err := d.Client.BlobStoreAPI.GetGroupBlobStoreConfiguration(ctx, data.Name.ValueString()).Execute()
	if err != nil {
		sharederr.HandleAPIError(
			"Unable to read group blob store",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	state := model.BlobStoreGroupModel{
		Name: types.StringValue(data.Name.ValueString()),
	}

	if apiResponse.SoftQuota != nil && apiResponse.SoftQuota.Type != nil {
		tflog.Debug(ctx, fmt.Sprintf("%v", apiResponse.SoftQuota))
		state.SoftQuota = &model.BlobStoreSoftQuota{
			Type:  types.StringValue(*apiResponse.SoftQuota.Type),
			Limit: types.Int64Value(*apiResponse.SoftQuota.Limit),
		}
	}
	if len(apiResponse.Members) > 0 {
		for _, m := range apiResponse.Members {
			state.Members = append(state.Members, types.StringValue(m))
		}
	}
	if apiResponse.FillPolicy != nil {
		state.FillPolicy = types.StringValue(*apiResponse.FillPolicy)
	}

	// Set state
	state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
