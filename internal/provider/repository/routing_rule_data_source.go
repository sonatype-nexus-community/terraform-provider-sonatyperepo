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
	"net/http"

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
	_ datasource.DataSource              = &routingRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &routingRuleDataSource{}
)

// RoutingRuleDataSource is a helper function to simplify the provider implementation.
func RoutingRuleDataSource() datasource.DataSource {
	return &routingRuleDataSource{}
}

// routingRuleDataSource is the data source implementation.
type routingRuleDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *routingRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_rule"
}

// Schema defines the schema for the data source.
func (d *routingRuleDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get a single routing rule by name",
		Attributes: map[string]tfschema.Attribute{
			"name":        schema.DataSourceRequiredStringWithLengthAtLeast("The name of the routing rule", 1),
			"description": schema.DataSourceComputedString("The description of the routing rule"),
			"mode":        schema.DataSourceComputedString("The mode of the routing rule (ALLOW or BLOCK)"),
			"matchers":    schema.DataSourceComputedStringSet("Regular expressions used to identify request paths"),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *routingRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.RoutingRuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	routingRuleResponse, httpResponse, err := d.Client.RoutingRulesAPI.GetRoutingRule(ctx, data.Name.ValueString()).Execute()

	state := model.RoutingRuleModel{}
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			errors.HandleAPIWarning(
				"No routing rule with supplied name",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				"Error finding routing rule",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
			return
		}
	} else if httpResponse.StatusCode == http.StatusOK {
		state.MapFromApi(routingRuleResponse)
	} else {
		errors.HandleAPIError(
			"Unexpected response when reading routing rule",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
