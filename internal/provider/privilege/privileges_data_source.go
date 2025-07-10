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

package privilege

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"terraform-provider-sonatyperepo/internal/provider/privilege/privilege_type"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &privilegesDataSource{}
	_ datasource.DataSourceWithConfigure = &privilegesDataSource{}
)

// PrivilegesDataSource is a helper function to simplify the provider implementation.
func PrivilegesDataSource() datasource.DataSource {
	return &privilegesDataSource{}
}

// privilegesDataSource is the data source implementation.
type privilegesDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *privilegesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_privileges"
}

// Schema defines the schema for the data source.
func (d *privilegesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get all Privileges",
		Attributes: map[string]schema.Attribute{
			"privileges": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the privilege. This value cannot be changed.",
							Required:    true,
							Optional:    false,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^[a-zA-Z0-9\-]{1}[a-zA-Z0-9_\-\.]*$`),
									`Please provide a name that complies with the Regular Expression: '^[a-zA-Z0-9\-]{1}[a-zA-Z0-9_\-\.]*$'`,
								),
							},
						},
						"description": schema.StringAttribute{
							Description: "Friendly description of this Privilege",
							Required:    true,
							Optional:    false,
						},
						"read_only": schema.BoolAttribute{
							Description: "Indicates whether the privilege can be changed. External values supplied to this will be ignored by the system.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "The email address associated with the user.",
							Required:    true,
							Optional:    false,
							Validators: []validator.String{
								stringvalidator.OneOf(privilege_type.AllPrivilegeTypes()...),
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *privilegesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.PrivilegesModel

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	apiResponse, httpResponse, err := d.Client.SecurityManagementPrivilegesAPI.GetAllPrivileges(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable list Privileges",
			fmt.Sprintf("Unable to read Privileges: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Privileges", len(apiResponse)))

	state.Privileges = make([]model.BasePrivilegeModel, 0)
	for _, p := range apiResponse {
		privilege := model.BasePrivilegeModel{}
		privilege.MapFromApi(&p)
		state.Privileges = append(state.Privileges, privilege)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
