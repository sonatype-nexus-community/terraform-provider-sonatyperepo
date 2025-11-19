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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &contentSelectorsDataSource{}
	_ datasource.DataSourceWithConfigure = &contentSelectorsDataSource{}
)

// ContentSelectorsDataSource is a helper function to simplify the provider implementation.
func ContentSelectorsDataSource() datasource.DataSource {
	return &contentSelectorsDataSource{}
}

// contentSelectorsDataSource is the data source implementation.
type contentSelectorsDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *contentSelectorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_content_selectors"
}

// Schema defines the schema for the data source.
func (d *contentSelectorsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dsschema.Schema{
		Description: "Use this data source to get all Content Selectors",
		Attributes: map[string]dsschema.Attribute{
			"content_selectors": dsschema.ListNestedAttribute{
				Computed: true,
				NestedObject: dsschema.NestedAttributeObject{
					Attributes: map[string]dsschema.Attribute{
						"name":        dsschema.StringAttribute{Description: "The name of the Content Selector.", Computed: true},
						"description": dsschema.StringAttribute{Description: "The description of this Content Selector.", Computed: true},
						"expression":  dsschema.StringAttribute{Description: "The Content Selector expression used to identify content.", Computed: true},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *contentSelectorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.ContentSelectorsModel

	ctx = d.GetAuthContext(ctx)

	contentSelectorsResponse, httpResponse, err := d.Client.ContentSelectorsAPI.GetContentSelectors(ctx).Execute()
	if err != nil {
		sharederr.HandleAPIError(
			"Unable list Content Selectors",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Content Selectors", len(contentSelectorsResponse)))

	state.ContentSelectors = make([]model.ContentSelectorModel, 0)
	for _, contentSelector := range contentSelectorsResponse {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s Content Selector", *contentSelector.Name))
		newCs := model.ContentSelectorModel{}
		newCs.MapFromApi(&contentSelector)

		state.ContentSelectors = append(state.ContentSelectors, newCs)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
