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

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &repositoriesDataSource{}
	_ datasource.DataSourceWithConfigure = &repositoriesDataSource{}
)

// RepositoriesDataSource is a helper function to simplify the provider implementation.
func RepositoriesDataSource() datasource.DataSource {
	return &repositoriesDataSource{}
}

// repositoriesDataSource is the data source implementation.
type repositoriesDataSource struct {
	baseDataSource
}

// Metadata returns the data source type name.
func (d *repositoriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repositories"
}

// Schema defines the schema for the data source.
func (d *repositoriesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get all Repositories",
		Attributes: map[string]schema.Attribute{
			"repositories": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Name of the Repository",
							Required:    true,
						},
						"format": schema.StringAttribute{
							Description: "Repository format",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "Repository type",
							Required:    true,
						},
						"url": schema.StringAttribute{
							Description: "URL to use this Repository",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *repositoriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.RepositoriesModel

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.auth,
	)

	repositories, httpResponse, err := d.client.RepositoryManagementAPI.GetAllRepositories(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Repositories",
			fmt.Sprintf("Unable to read Repositories: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Repositories", len(repositories)))

	for _, repository := range repositories {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s Repository", *repository.Name))

		state.Repositories = append(state.Repositories, model.RepositoryModel{
			Name:   types.StringValue(*repository.Name),
			Format: types.StringValue(*repository.Format),
			Type:   types.StringValue(*repository.Type),
			Url:    types.StringValue(*repository.Url),
		})
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
