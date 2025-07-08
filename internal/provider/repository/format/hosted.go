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

package format

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getCommonHostedSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"component": schema.SingleNestedAttribute{
			Description: "Component configuration for this Repository",
			Required:    false,
			Optional:    true,
			Computed:    true,
			Attributes: map[string]schema.Attribute{
				"proprietary_components": schema.BoolAttribute{
					Description: "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Repository Firewall)",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
			},
			Default: objectdefault.StaticValue(types.ObjectValueMust(
				map[string]attr.Type{
					"proprietary_components": types.BoolType,
				},
				map[string]attr.Value{
					"proprietary_components": types.BoolValue(false),
				},
			)),
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
		},
	}
}
