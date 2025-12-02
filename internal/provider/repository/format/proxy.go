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
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

func commonProxySchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"proxy": schema.ResourceRequiredSingleNestedAttribute(
			"Proxy specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"remote_url": schema.ResourceRequiredStringWithRegex(
					"Location of the remote repository being proxied",
					regexp.MustCompile(`^https?://`),
					"must be a valid HTTP URL (starting with http:// or https://)",
				),
				"content_max_age":  schema.ResourceRequiredInt64("How long to cache artifacts before rechecking the remote repository (in minutes)"),
				"metadata_max_age": schema.ResourceRequiredInt64("How long to cache metadata before rechecking the remote repository (in minutes)"),
			},
		),
		"negative_cache": schema.ResourceRequiredSingleNestedAttribute(
			"Negative Cache configuration for this Repository",
			map[string]tfschema.Attribute{
				"enabled": schema.ResourceRequiredBool("Whether to cache responses for content not present in the proxied repository"),
				"time_to_live": schema.ResourceRequiredInt64WithValidators(
					"How long to cache the fact that a file was not found in the repository (in minutes)",
					[]validator.Int64{
						int64validator.AtLeast(0),
					}...,
				),
			},
		),
		"http_client": schema.ResourceRequiredSingleNestedAttribute(
			"HTTP Client configuration for this Repository",
			map[string]tfschema.Attribute{
				"blocked":        schema.ResourceRequiredBool("Whether to block outbound connections on the repository"),
				"auto_block":     schema.ResourceRequiredBool("Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive"),
				"connection":     commonProxyConnectionAttribute(),
				"authentication": commonProxyAuthenticationAttribute(),
			},
		),
		"routing_rule": schema.ResourceOptionalString("Routing Rule"),
		"replication":  commonProxyReplicationAttribute(),
	}
}

func commonProxyConnectionAttribute() tfschema.SingleNestedAttribute {
	thisAttr := schema.ResourceOptionalSingleNestedAttribute(
		"HTTP Client Connection configuration for this Repository",
		map[string]tfschema.Attribute{
			"retries": schema.ResourceOptionalInt64WithDefaultAndValidators(
				"Total retries if the initial connection attempt suffers a timeout",
				common.DEFAULT_HTTP_CONNECTION_RETRIES,
				[]validator.Int64{
					int64validator.Between(
						common.REPOSITORY_HTTP_CLIENT_CONNECTION_RETRIES_MIN,
						common.REPOSITORY_HTTP_CLIENT_CONNECTION_RETRIES_MAX,
					),
				}...,
			),
			"user_agent_suffix": schema.ResourceOptionalString("Custom fragment to append to User-Agent header in HTTP requests"),
			"timeout": schema.ResourceOptionalInt64WithDefaultAndValidators(
				"Seconds to wait for activity before stopping and retrying the connection",
				common.DEFAULT_HTTP_CONNECTION_TIMEOUT,
				[]validator.Int64{
					int64validator.Between(
						common.REPOSITORY_HTTP_CLIENT_CONNECTION_TIMEOUT_MIN,
						common.REPOSITORY_HTTP_CLIENT_CONNECTION_TIMEOUT_MAX,
					),
				}...,
			),
			"enable_circular_redirects": schema.ResourceComputedOptionalBoolWithDefault(
				"Whether to enable redirects to the same location (may be required by some servers)",
				false,
			),
			"enable_cookies": schema.ResourceComputedOptionalBoolWithDefault(
				"Whether to allow cookies to be stored and used",
				false,
			),
			"use_trust_store": schema.ResourceComputedOptionalBoolWithDefault(
				"Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
				false,
			),
		},
	)
	thisAttr.Computed = true
	thisAttr.Default = objectdefault.StaticValue(types.ObjectValueMust(
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
	))
	thisAttr.PlanModifiers = []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	}

	return thisAttr
}

func commonProxyAuthenticationAttribute() tfschema.SingleNestedAttribute {
	return schema.ResourceOptionalSingleNestedAttribute(
		"Authentication to upstream Repository",
		map[string]tfschema.Attribute{
			"type": schema.ResourceOptionalStringEnum(
				"Authentication type",
				common.HTTP_AUTH_TYPE_USERNAME,
				common.HTTP_AUTH_TYPE_NTLM,
				common.HTTP_AUTH_TYPE_BEARER_TOKEN,
			),
			"username": schema.ResourceOptionalString("Username"),
			"password": schema.ResourceSensitiveOptionalStringWithPlanModifier(
				"Password",
				stringplanmodifier.UseStateForUnknown(),
			),
			"ntlm_host":   schema.ResourceOptionalString("NTLM Host"),
			"ntlm_domain": schema.ResourceOptionalString("NTLM Domain"),
			"preemptive":  schema.ResourceOptionalBool("Whether to use pre-emptive authentication. Use with caution. Defaults to false."),
			"bearer_token": schema.ResourceSensitiveOptionalStringWithPlanModifier(
				"Bearer Token used when Authentication Type == bearerToken",
				stringplanmodifier.UseStateForUnknown(),
			),
		},
	)
}

func commonProxyReplicationAttribute() tfschema.SingleNestedAttribute {
	thisAttr := schema.ResourceOptionalSingleNestedAttribute(
		"Replication configuration for this Repository",
		map[string]tfschema.Attribute{
			"preemptive_pull_enabled": schema.ResourceRequiredBool("Whether pre-emptive pull is enabled"),
			"asset_path_regex":        schema.ResourceOptionalString("Regular Expression of Asset Paths to pull pre-emptively pull"),
		},
	)
	thisAttr.Computed = true
	thisAttr.Default = objectdefault.StaticValue(
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
	)
	return thisAttr
}
