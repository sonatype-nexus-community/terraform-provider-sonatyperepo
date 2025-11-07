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

package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// RoutingRulesModel represents a list of routing rules for data source
type RoutingRulesModel struct {
	RoutingRules []RoutingRuleModel `tfsdk:"routing_rules"`
}

// RoutingRuleModel represents the Terraform model for a routing rule (used by data source)
type RoutingRuleModel struct {
	Name        types.String   `tfsdk:"name"`
	Description types.String   `tfsdk:"description"`
	Mode        types.String   `tfsdk:"mode"`
	Matchers    []types.String `tfsdk:"matchers"`
	LastUpdated types.String   `tfsdk:"last_updated"`
}

// MapFromApi maps API response to model
func (m *RoutingRuleModel) MapFromApi(api *sonatyperepo.RoutingRuleXO) {
	m.Name = types.StringPointerValue(api.Name)
	m.Description = types.StringPointerValue(api.Description)
	m.Mode = types.StringPointerValue(api.Mode)

	if api.Matchers != nil {
		m.Matchers = make([]types.String, len(api.Matchers))
		for i, matcher := range api.Matchers {
			m.Matchers[i] = types.StringValue(matcher)
		}
	}
}
