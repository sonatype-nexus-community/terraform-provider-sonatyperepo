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

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

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
	resp.Schema = schema.Schema{
		Description: "Use this data source to get all Roles",
		Attributes: map[string]schema.Attribute{
			"roles": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The id of the role.",
							Required:    true,
							Optional:    false,
						},
						"name": schema.StringAttribute{
							Description: "The name of the role.",
							Required:    true,
							Optional:    false,
						},
						"description": schema.StringAttribute{
							Description: "The description of this role.",
							Required:    true,
							Optional:    false,
						},
						"read_only": schema.BoolAttribute{
							Description: "Indicates whether the role can be changed. The system will ignore any supplied external values.",
							Required:    true,
							Optional:    false,
						},
						"source": schema.StringAttribute{
							Description: "The user source which is the origin of this role.",
							Required:    true,
							Optional:    false,
						},
						"privileges": schema.ListAttribute{
							Description: "The list of privileges assigned to this role.",
							Required:    true,
							Optional:    false,
							ElementType: types.StringType,
							Validators: []validator.List{
								listvalidator.UniqueValues(),
							},
						},
						"roles": schema.ListAttribute{
							Description: "The list of roles assigned to this role.",
							Required:    true,
							Optional:    false,
							ElementType: types.StringType,
							Validators: []validator.List{
								listvalidator.UniqueValues(),
							},
						},
					},
				},
			},
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
		resp.Diagnostics.AddError(
			"Unable list Roles",
			fmt.Sprintf("Unable to read Roles: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Roles", len(rolesResponse)))

	for _, role := range rolesResponse {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s Role", *role.Name))

		newRole := model.RoleModel{
			Id:          types.StringValue(*role.Id),
			Name:        types.StringValue(*role.Name),
			Description: types.StringValue(*role.Description),
			ReadOnly:    types.BoolValue(*role.ReadOnly),
			Source:      types.StringValue(*role.Source),
			Privileges:  make([]types.String, 0),
			Roles:       make([]types.String, 0),
		}

		for _, privilege := range role.Privileges {
			newRole.Privileges = append(newRole.Privileges, types.StringValue(privilege))
		}

		for _, r := range role.Roles {
			newRole.Roles = append(newRole.Roles, types.StringValue(r))
		}

		state.Roles = append(state.Roles, newRole)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
