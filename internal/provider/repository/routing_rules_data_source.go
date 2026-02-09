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

package repository

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
	_ datasource.DataSource              = &routingRulesDataSource{}
	_ datasource.DataSourceWithConfigure = &routingRulesDataSource{}
)

// RoutingRulesDataSource is a helper function to simplify the provider implementation.
func RoutingRulesDataSource() datasource.DataSource {
	return &routingRulesDataSource{}
}

// routingRulesDataSource is the data source implementation.
type routingRulesDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *routingRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_rules"
}

// Schema defines the schema for the data source.
func (d *routingRulesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get all routing rules",
		Attributes: map[string]tfschema.Attribute{
			"routing_rules": schema.DataSourceComputedListNestedAttribute(
				"List of Routing Rules",
				tfschema.NestedAttributeObject{
					Attributes: map[string]tfschema.Attribute{
						"name":        schema.DataSourceComputedString("The name of the routing rule"),
						"description": schema.DataSourceComputedString("The description of the routing rule"),
						"mode":        schema.DataSourceComputedString("The mode of the routing rule (ALLOW or BLOCK)"),
						"matchers":    schema.DataSourceComputedStringSet("Regular expressions used to identify request paths"),
					},
				},
			),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *routingRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.RoutingRulesModel

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	routingRulesResponse, httpResponse, err := d.Client.RoutingRulesAPI.GetRoutingRules(ctx).Execute()
	if err != nil {
		errors.HandleAPIError(
			"Unable to list routing rules",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d routing rules", len(routingRulesResponse)))

	state.RoutingRules = make([]model.RoutingRuleModelDS, 0)
	for _, routingRule := range routingRulesResponse {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s routing rule", *routingRule.Name))
		newRr := model.RoutingRuleModelDS{}
		newRr.MapFromApi(&routingRule)

		state.RoutingRules = append(state.RoutingRules, newRr)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
