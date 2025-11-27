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

package user

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

// UsersDataSource is a helper function to simplify the provider implementation.
func UsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// fileBlobStoreDataSource is the data source implementation.
type usersDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *usersDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get all Users",
		Attributes: map[string]tfschema.Attribute{
			"users": schema.DataSourceComputedListNestedAttribute(
				"List of Users",
				tfschema.NestedAttributeObject{
					Attributes: map[string]tfschema.Attribute{
						"user_id":       schema.DataSourceRequiredString("The userid which is required for login. This value cannot be changed."),
						"first_name":    schema.DataSourceRequiredString("The first name of the user."),
						"last_name":     schema.DataSourceRequiredString("The last name of the user."),
						"email_address": schema.DataSourceRequiredString("The email address associated with the user."),
						"read_only":     schema.DataSourceRequiredBool("Indicates whether the user's properties could be modified by the Nexus Repository Manager. When false only roles are considered during update."),
						"source":        schema.DataSourceRequiredString("The user source which is the origin of this user. This value cannot be changed."),
						"status": schema.DataSourceRequiredStringWithValidators(
							"The user's status",
							stringvalidator.OneOf(common.AllUserStatusTypes()...),
						),
						"roles":          schema.DataSourceRequiredStringSet("The roles which the user has been assigned within Nexus."),
						"external_roles": schema.DataSourceRequiredStringSet("The roles which the user has been assigned in an external source, e.g. LDAP group. These cannot be changed within the Nexus Repository Manager."),
					},
				},
			),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.UsersModel

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	usersResponse, httpResponse, err := d.Client.SecurityManagementUsersAPI.GetUsers(ctx).Execute()
	if err != nil {
		errors.HandleAPIError(
			"Unable to list Users",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Users", len(usersResponse)))

	for _, u := range usersResponse {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s User", *u.UserId))
		newUser := model.UserModel{}
		newUser.MapFromApi(&u)
		state.Users = append(state.Users, newUser)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
