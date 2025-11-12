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

package capability

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &capabilitiesDataSource{}
	_ datasource.DataSourceWithConfigure = &capabilitiesDataSource{}
)

// CapabilitiesDataSource is a helper function to simplify the provider implementation.
func CapabilitiesDataSource() datasource.DataSource {
	return &capabilitiesDataSource{}
}

// capabilitiesDataSource is the data source implementation.
type capabilitiesDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *capabilitiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_capabilities"
}

// Schema defines the schema for the data source.
func (d *capabilitiesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Use this data source to get all Capabilities.
		
**NOTE:** Requires Sonatype Nexus Repostiory 3.84.0 or later.`,
		Attributes: map[string]schema.Attribute{
			"capabilities": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Internal ID of the Capability.",
							Required:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the Capability.",
							Required:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: "Whether the Capability is enabled.",
							Required:    true,
						},
						"notes": schema.StringAttribute{
							Description: "Notes about the configured Capability.",
							Optional:    true,
						},
						"properties": schema.MapAttribute{
							Description: "Properties of the Capability.",
							Required:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *capabilitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.CapabilitiesListModel

	ctx = d.GetAuthContext(ctx)

	apiResponse, httpResponse, err := d.Client.CapabilitiesAPI.List(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to list Capabilities",
			fmt.Sprintf("Unable to read Capabilities: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Capabilities", len(apiResponse)))

	for _, capability := range apiResponse {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s Capability", *capability.Id))

		state.Capabilities = append(state.Capabilities, model.CapabilityModel{
			CapabilitCommonModel: model.CapabilitCommonModel{
				Id:      types.StringPointerValue(capability.Id),
				Notes:   types.StringPointerValue(capability.Notes),
				Enabled: types.BoolPointerValue(capability.Enabled),
			},
			Type:       types.StringPointerValue(capability.Type),
			Properties: capability.GetProperties(),
		})
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
