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
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getHostedStandardSchema(format string) schema.Schema {
	return schema.Schema{
		Description: fmt.Sprintf("Manage Hosted %s Repositories", format),
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the Repository",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL to access the Repository",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"online": schema.BoolAttribute{
				Description: "Whether this Repository is online and accepting incoming requests",
				Required:    true,
			},
			"storage": schema.SingleNestedAttribute{
				Description: "Storage configuration for this Repository",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"blob_store_name": schema.StringAttribute{
						Description: "Name of the Blob Store to use",
						Required:    true,
						Optional:    false,
					},
					"strict_content_type_validation": schema.BoolAttribute{
						Description: "Whether this Repository validates that all content uploaded to this repository is of a MIME type appropriate for the repository format",
						Required:    true,
					},
					"write_policy": schema.StringAttribute{
						Description: "Controls if deployments of and updates to assets are allowed",
						Required:    true,
						Optional:    false,
						Validators: []validator.String{
							stringvalidator.OneOf("ALLOW", "ALLOW_ONCE", "DENY"),
						},
					},
				},
			},
			"cleanup": schema.SingleNestedAttribute{
				Description: "Repository Cleanup configuration",
				Required:    false,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"policy_names": schema.ListAttribute{
						Description: "Components that match any of the applied policies will be deleted",
						ElementType: types.StringType,
						Required:    false,
						Optional:    true,
					},
				},
			},
			"component": schema.SingleNestedAttribute{
				Description: "Component configuration for this Repository",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"proprietary_components": schema.BoolAttribute{
						Description: "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Firewall)",
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
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
