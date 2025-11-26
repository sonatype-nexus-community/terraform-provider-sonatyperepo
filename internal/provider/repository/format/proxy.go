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
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func commonProxySchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"proxy": schema.SingleNestedAttribute{
			Description: "Proxy specific configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes: map[string]schema.Attribute{
				"remote_url": schema.StringAttribute{
					Description: "Location of the remote repository being proxied",
					Required:    true,
					Optional:    false,
				},
				"content_max_age": schema.Int64Attribute{
					Description: "How long to cache artifacts before rechecking the remote repository (in minutes)",
					Required:    true,
					Optional:    false,
				},
				"metadata_max_age": schema.Int64Attribute{
					Description: "How long to cache metadata before rechecking the remote repository (in minutes)",
					Required:    true,
					Optional:    false,
				},
			},
		},
		"negative_cache": schema.SingleNestedAttribute{
			Description: "Negative Cache configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes: map[string]schema.Attribute{
				"enabled": schema.BoolAttribute{
					Description: "Whether to cache responses for content not present in the proxied repository",
					Required:    true,
					Optional:    false,
				},
				"time_to_live": schema.Int64Attribute{
					Description: "How long to cache the fact that a file was not found in the repository (in minutes)",
					Required:    true,
					Optional:    false,
				},
			},
		},
		"http_client": schema.SingleNestedAttribute{
			Description: "HTTP Client configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes: map[string]schema.Attribute{
				"blocked": schema.BoolAttribute{
					Description: "Whether to block outbound connections on the repository",
					Required:    true,
					Optional:    false,
				},
				"auto_block": schema.BoolAttribute{
					Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
					Required:    true,
					Optional:    false,
				},
				"connection": schema.SingleNestedAttribute{
					Description: "HTTP Client Connection configuration for this Repository",
					Required:    false,
					Optional:    true,
					Computed:    true,
					Attributes: map[string]schema.Attribute{
						"retries": schema.Int64Attribute{
							Description: "Total retries if the initial connection attempt suffers a timeout",
							Required:    false,
							Optional:    true,
							Computed:    true,
							Default:     int64default.StaticInt64(common.DEFAULT_HTTP_CONNECTION_RETRIES),
							Validators: []validator.Int64{
								int64validator.Between(0, 10),
							},
						},
						"user_agent_suffix": schema.StringAttribute{
							Description: "Custom fragment to append to User-Agent header in HTTP requests",
							Required:    false,
							Optional:    true,
						},
						"timeout": schema.Int64Attribute{
							Description: "Seconds to wait for activity before stopping and retrying the connection",
							Required:    false,
							Optional:    true,
							Computed:    true,
							Default:     int64default.StaticInt64(common.DEFAULT_HTTP_CONNECTION_TIMEOUT),
							Validators: []validator.Int64{
								int64validator.Between(1, 3600),
							},
						},
						"enable_circular_redirects": schema.BoolAttribute{
							Description: "Whether to enable redirects to the same location (may be required by some servers)",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"enable_cookies": schema.BoolAttribute{
							Description: "Whether to allow cookies to be stored and used",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
						"use_trust_store": schema.BoolAttribute{
							Description: "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
							Optional:    true,
							Computed:    true,
							Default:     booldefault.StaticBool(false),
						},
					},
					Default: objectdefault.StaticValue(types.ObjectValueMust(
						map[string]attr.Type{
							"retries":                   types.Int64Type,
							"user_agent_suffix":         types.StringType,
							"timeout":                   types.Int64Type,
							"enable_circular_redirects": types.BoolType,
							"enable_cookies":            types.BoolType,
							"use_trust_store":           types.BoolType,
						},
						map[string]attr.Value{
							"retries":                   types.Int64Value(common.DEFAULT_HTTP_CONNECTION_RETRIES),
							"user_agent_suffix":         types.StringNull(),
							"timeout":                   types.Int64Value(common.DEFAULT_HTTP_CONNECTION_TIMEOUT),
							"enable_circular_redirects": types.BoolValue(false),
							"enable_cookies":            types.BoolValue(false),
							"use_trust_store":           types.BoolValue(false),
						},
					)),
					PlanModifiers: []planmodifier.Object{
						objectplanmodifier.UseStateForUnknown(),
					},
				},
				"authentication": schema.SingleNestedAttribute{
					Description: "Authentication to upstream Repository",
					Required:    false,
					Optional:    true,
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Authentication type",
							Required:    false,
							Optional:    true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									common.HTTP_AUTH_TYPE_USERNAME,
									common.HTTP_AUTH_TYPE_NTLM,
									common.HTTP_AUTH_TYPE_BEARER_TOKEN,
								),
							},
						},
						"username": schema.StringAttribute{
							Description: "Username",
							Required:    false,
							Optional:    true,
						},
						"password": schema.StringAttribute{
							Description: "Password",
							Required:    false,
							Optional:    true,
							Sensitive:   true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"ntlm_host": schema.StringAttribute{
							Description: "NTLM Host",
							Required:    false,
							Optional:    true,
						},
						"ntlm_domain": schema.StringAttribute{
							Description: "NTLM Domain",
							Required:    false,
							Optional:    true,
						},
						"preemptive": schema.BoolAttribute{
							Description: "Whether to use pre-emptive authentication. Use with caution. Defaults to false.",
							Required:    false,
							Optional:    true,
							// Computed:    true,
							// Default:     booldefault.StaticBool(false),
						},
						"bearer_token": schema.StringAttribute{
							Description: "Bearer Token used when Authentication Type == bearerToken",
							Required:    false,
							Optional:    true,
							Sensitive:   true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
		"routing_rule": schema.StringAttribute{
			Description: "Routing Rule",
			Required:    false,
			Optional:    true,
		},
		"replication": schema.SingleNestedAttribute{
			Description: "Replication configuration for this Repository",
			Required:    false,
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"preemptive_pull_enabled": schema.BoolAttribute{
					Description: "Whether pre-emptive pull is enabled",
					Required:    true,
					Optional:    false,
				},
				"asset_path_regex": schema.StringAttribute{
					Description: "Regular Expression of Asset Paths to pull pre-emptively pull",
					Required:    false,
					Optional:    true,
				},
			},
			Computed: true,
			Default: objectdefault.StaticValue(
				types.ObjectValueMust(
					map[string]attr.Type{
						"preemptive_pull_enabled": types.BoolType,
						"asset_path_regex":        types.StringType,
					},
					map[string]attr.Value{
						"preemptive_pull_enabled": types.BoolValue(false),
						"asset_path_regex":        types.StringNull(),
					},
				),
			),
		},
	}
}
