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
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

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
	resp.Schema = schema.Schema{
		Description: "Use this data source to get a single Content Selector by name",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the Content Selector.",
				Required:    true,
				Optional:    false,
			},
			"description": schema.StringAttribute{
				Description: "The description of this Content Selector.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"expression": schema.StringAttribute{
				Description: "The Content Selector expression used to identify content.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
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

	ctx = d.GetAuthContext(ctx)

	contentSelectorsResponse, httpResponse, err := d.Client.ContentSelectorsAPI.GetContentSelector(ctx, data.Name.ValueString()).Execute()

	state := model.ContentSelectorModel{}
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning(
				"No Content Selector with supplied name",
				fmt.Sprintf("No Content Selector with supplied name: %s", httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error finding Content Selector",
				fmt.Sprintf("Error finding Content Selector with supplied name: %s", httpResponse.Status),
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
