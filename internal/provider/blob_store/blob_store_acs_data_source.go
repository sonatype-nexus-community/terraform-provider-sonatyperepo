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
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &s3BlobStoreDataSource{}
	_ datasource.DataSourceWithConfigure = &s3BlobStoreDataSource{}
)

// BlobStoreAcsDataSource is a helper function to simplify the provider implementation.
func BlobStoreAcsDataSource() datasource.DataSource {
	return &acsBlobStoreDataSource{}
}

// acsBlobStoreDataSource is the data source implementation.
type acsBlobStoreDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *acsBlobStoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blob_store_acs"
}

// Schema defines the schema for the data source.
func (d *acsBlobStoreDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get a specific S3 Blob Store by it's name",
		Attributes: map[string]tfschema.Attribute{
			"name": schema.DataSourceRequiredString("Name of the Blob Store"),
			// "type": schema.DataSourceOptionalString(fmt.Sprintf("Type of this Blob Store - will always be '%s'", common.BLOB_STORE_TYPE_S3)),
			"soft_quota": schema.DataSourceComputedOptionalSingleNestedAttribute("Soft Quota for this Blob Store", map[string]tfschema.Attribute{
				"type":  schema.DataSourceOptionalString("Soft Quota type"),
				"limit": schema.DataSourceOptionalInt64("Quota limit"),
			}),
			"bucket_configuration": schema.DataSourceComputedSingleNestedAttribute("Bucket Configuration for this Blob Store", map[string]tfschema.Attribute{
				"account_name":   schema.DataSourceComputedString("Account name found under Access keys for the storage account."),
				"container_name": schema.DataSourceComputedString("The name of an existing container to be used for storage."),
				"authentication": schema.DataSourceComputedSingleNestedAttribute("Authentication to Azure Cloud", map[string]tfschema.Attribute{
					"authentication_method": schema.DataSourceComputedString("The type of Azure authentication to use."),
					"account_key":           schema.DataSourceComputedString("The account key"),
				}),
			}),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *acsBlobStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.BlobStoreAcsModelDS

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API
	apiResponse, httpResponse, err := d.Client.BlobStoreAPI.GetBlobStore1(d.AuthContext(ctx), data.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Azure Cloud Storage Blob Store does not exist",
				fmt.Sprintf("Azure Cloud Storage Blob Store does not exist: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
			errors.HandleAPIError(
				fmt.Sprintf("No Azure Cloud Storage Blob Store with name: %s", data.Name.ValueString()),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
			return
		} else {
			errors.HandleAPIError(
				fmt.Sprintf("Unexpected error reading Azure Cloud Storage Blob Store with name: %s", data.Name.ValueString()),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
			return
		}
	} else {
		// Update State
		data.MapFromApi(apiResponse)
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}
