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

package system

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &securityUserTokenDataSource{}
	_ datasource.DataSourceWithConfigure = &securityUserTokenDataSource{}
)

// SecurityUserTokenDataSource is a helper function to simplify the provider implementation.
func SecurityUserTokenDataSource() datasource.DataSource {
	return &securityUserTokenDataSource{}
}

// securityUserTokenDataSource is the data source implementation.
type securityUserTokenDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *securityUserTokenDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_security_user_tokens"
}

// Schema defines the schema for the data source.
func (d *securityUserTokenDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get the current User Token configuration",
		Attributes: map[string]tfschema.Attribute{
			"enabled":            schema.DataSourceComputedBool("Whether or not User Tokens feature is enabled"),
			"expiration_days":    schema.DataSourceComputedInt32("User token expiration days (1-999)"),
			"expiration_enabled": schema.DataSourceComputedBool("Whether user tokens expiration is enabled"),
			"protect_content":    schema.DataSourceComputedBool("Whether user tokens are required for repository authentication"),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *securityUserTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	// Read API Call
	apiResponse, httpResponse, err := d.Client.SecurityManagementUserTokensAPI.ServiceStatus(ctx).Execute()

	if err != nil {
		errors.HandleAPIError(
			"Unable to read User Token settings",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Map API response to state
	var state model.SecurityUserTokenModelDS
	state.MapFromApi(apiResponse)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Setting state has errors: %v", resp.Diagnostics.Errors()))
		return
	}
}
