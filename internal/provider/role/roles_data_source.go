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

package role

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
	_ datasource.DataSource              = &rolesDataSource{}
	_ datasource.DataSourceWithConfigure = &rolesDataSource{}
)

// RolesDataSource is a helper function to simplify the provider implementation.
func RolesDataSource() datasource.DataSource {
	return &rolesDataSource{}
}

// fileBlobStoreDataSource is the data source implementation.
type rolesDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *rolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles"
}

// Schema defines the schema for the data source.
func (d *rolesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get all Roles",
		Attributes: map[string]tfschema.Attribute{
			"roles": schema.DataSourceComputedListNestedAttribute(
				"List of Roles",
				tfschema.NestedAttributeObject{
					Attributes: map[string]tfschema.Attribute{
						"id":          schema.DataSourceRequiredString("The id of the role."),
						"name":        schema.DataSourceRequiredString("The name of the role."),
						"description": schema.DataSourceRequiredString("The description of this role."),
						"read_only":   schema.DataSourceComputedBool("Indicates whether the role can be changed. The system will ignore any supplied external values."),
						"source":      schema.DataSourceRequiredString("The user source which is the origin of this role."),
						"privileges":  schema.DataSourceComputedStringSet("The set of privileges assigned to this role."),
						"roles":       schema.DataSourceComputedStringSet("The set of roles assigned to this role."),
					},
				},
			),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *rolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.RolesModel

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	rolesResponse, httpResponse, err := d.Client.SecurityManagementRolesAPI.GetRoles(ctx).Execute()
	if err != nil {
		errors.HandleAPIError(
			"Unable to list roles",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Roles", len(rolesResponse)))

	for _, role := range rolesResponse {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s Role", *role.Name))
		newRole := model.RoleModelIncludingReadOnly{}
		newRole.MapFromApi(&role)
		state.Roles = append(state.Roles, newRole)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
