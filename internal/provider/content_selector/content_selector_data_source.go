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

package content_selector

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
	_ datasource.DataSource              = &contentSelectorDataSource{}
	_ datasource.DataSourceWithConfigure = &contentSelectorDataSource{}
)

// ContentSelectorDataSource is a helper function to simplify the provider implementation.
func ContentSelectorDataSource() datasource.DataSource {
	return &contentSelectorDataSource{}
}

// contentSelectorDataSource is the data source implementation.
type contentSelectorDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *contentSelectorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_content_selector"
}

// Schema defines the schema for the data source.
func (d *contentSelectorDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get a single Content Selector by name",
		Attributes: map[string]tfschema.Attribute{
			"name":        schema.DataSourceRequiredString("The name of the Content Selector."),
			"description": schema.DataSourceComputedString("The description of this Content Selector."),
			"expression":  schema.DataSourceComputedString("The Content Selector expression used to identify content."),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *contentSelectorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.ContentSelectorModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	if data.Name.IsNull() {
		resp.Diagnostics.AddError("Name must not be empty.", "Name must be provided.")
		return
	}

	ctx = d.AuthContext(ctx)

	contentSelectorsResponse, httpResponse, err := d.Client.ContentSelectorsAPI.GetContentSelector(ctx, data.Name.ValueString()).Execute()

	state := model.ContentSelectorModel{}
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			errors.HandleAPIWarning(
				"No Content Selector with supplied name",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				"Error finding Content Selector",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
			return
		}
	} else if httpResponse.StatusCode == http.StatusOK {
		state.MapFromApi(contentSelectorsResponse)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
