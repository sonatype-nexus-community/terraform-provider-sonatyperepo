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
	"terraform-provider-sonatyperepo/internal/provider/model"
	"terraform-provider-sonatyperepo/internal/provider/validators"

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

// FirewallModeFromFlags derives the inline `firewall.mode` value (NXRM 3.94+) from the
// Terraform schema flags stored in `repository_firewall`.
func FirewallModeFromFlags(hasConfig, enabled, quarantine, pccsEnabled bool) common.FirewallMode {
	if !hasConfig || !enabled {
		return common.FirewallModeDisabled
	}
	if pccsEnabled {
		return common.FirewallModePccs
	}
	if quarantine {
		return common.FirewallModeQuarantine
	}
	return common.FirewallModeAudit
}

// FirewallFlagsFromMode derives the `repository_firewall` schema flags from the inline
// `firewall.mode` value (NXRM 3.94+).
func FirewallFlagsFromMode(mode common.FirewallMode) (enabled, quarantine, pccsEnabled bool) {
	switch mode {
	case common.FirewallModeAudit:
		return true, false, false
	case common.FirewallModeQuarantine:
		return true, true, false
	case common.FirewallModePccs:
		return true, false, true
	default:
		return false, false, false
	}
}

// ComputeFirewallMode derives the inline `firewall.mode` value (NXRM 3.94+) for a proxy
// repository's plan/state model. It only asks for quarantine/pccs flags when a
// `repository_firewall` block is actually present, since some formats' getters for those
// flags assume a non-nil block.
func ComputeFirewallMode(f RepositoryFormat, state any) common.FirewallMode {
	hasConfig := f.HasFirewallConfig(state)
	enabled := f.GetRepositoryFirewallEnabled(state)
	var quarantine, pccsEnabled bool
	if hasConfig {
		quarantine = f.GetRepositoryFirewallQuarantineEnabled(state)
		pccsEnabled = f.GetRepositoryFirewallPccsEnabled(state)
	}
	return FirewallModeFromFlags(hasConfig, enabled, quarantine, pccsEnabled)
}

// ResolveFirewallBlockFlags reconciles the `repository_firewall` state attribute against the
// inline `firewall.mode` (NXRM 3.94+) reported by DoReadRequest/DoImportRequest, given
// whether the incoming state/plan model already had a `repository_firewall` block
// configured.
//
// `repository_firewall` is Optional but not Computed, with no plan modifier (see
// commonProxyFirewallAuditQuarantineAttribute), so Terraform always builds its planned value
// directly from the practitioner's config. If the block was configured - even with
// `enabled = false` - the plan is a non-null object; unconditionally nulling the whole
// attribute out just because the server resolves the mode as disabled then makes apply fail
// with "Provider produced inconsistent result after apply" (see GH-469). A resolved Disabled
// mode should only clear the block when it was never configured to begin with - otherwise
// it's just another set of flags to report, the same as Audit/Quarantine/Pccs.
func ResolveFirewallBlockFlags(hadConfig bool, mode common.FirewallMode) (keep, enabled, quarantine, pccsEnabled bool) {
	enabled, quarantine, pccsEnabled = FirewallFlagsFromMode(mode)
	keep = hadConfig || mode != common.FirewallModeDisabled
	return keep, enabled, quarantine, pccsEnabled
}

// BackfillFirewallBlockFromPlan mirrors planFirewall into state (creating the block if
// necessary, with capability_id cleared) when the plan configured `repository_firewall` but
// stateFirewall is still nil. This covers the one case ResolveFirewallBlockFlags can't see on
// its own: Update() feeds UpdateStateFromApi the PRIOR state rather than the new plan (see
// repository_common.go), so adding the block for the first time with `enabled = false` in the
// same apply isn't visible there. UpdateStateFromPlanForNonApiFields has both the plan and the
// post-read state, so it's the right place to catch this up. See GH-469.
func BackfillFirewallBlockFromPlan(stateFirewall, planFirewall *model.FirewallAuditAndQuarantineModel) *model.FirewallAuditAndQuarantineModel {
	if stateFirewall != nil || planFirewall == nil {
		return stateFirewall
	}
	backfilled := *planFirewall
	backfilled.CapabilityId = types.StringNull()
	return &backfilled
}

// BackfillFirewallBlockWithPccsFromPlan is BackfillFirewallBlockFromPlan for the PCCS-capable
// variant of the `repository_firewall` block (NPM, PyPI). See GH-469.
func BackfillFirewallBlockWithPccsFromPlan(stateFirewall, planFirewall *model.FirewallAuditAndQuarantineWithPccsModel) *model.FirewallAuditAndQuarantineWithPccsModel {
	if stateFirewall != nil || planFirewall == nil {
		return stateFirewall
	}
	backfilled := *planFirewall
	backfilled.CapabilityId = types.StringNull()
	return &backfilled
}

// ProxyApiResponseWithFirewall carries a proxy repository API response together with its
// inline `firewall.mode` value (NXRM 3.94+), threading the mode through DoReadRequest/
// DoImportRequest to UpdateStateFromApi without changing the RepositoryFormat interface.
type ProxyApiResponseWithFirewall struct {
	Repository   any
	FirewallMode *common.FirewallMode
}

func commonProxySchemaAttributes(supportsRepositoryFirewall, supportsPccs bool) map[string]tfschema.Attribute {
	thisAttr := map[string]tfschema.Attribute{
		"proxy": schema.ResourceRequiredSingleNestedAttribute(
			"Proxy specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"remote_url": schema.ResourceRequiredStringWithRegex(
					"Location of the remote repository being proxied",
					regexp.MustCompile(`^https?://`),
					"must be a valid HTTP URL (starting with http:// or https://)",
				),
				"content_max_age": schema.ResourceOptionalInt64WithDefault(
					"How long to cache artifacts before rechecking the remote repository (in minutes)",
					common.DEFAULT_PROXY_CONTENT_MAX_AGE,
				),
				"metadata_max_age": schema.ResourceOptionalInt64WithDefault(
					"How long to cache metadata before rechecking the remote repository (in minutes)",
					common.DEFAULT_PROXY_METADATA_MAX_AGE,
				),
			},
		),
		"negative_cache": schema.ResourceRequiredSingleNestedAttribute(
			"Negative Cache configuration for this Repository",
			map[string]tfschema.Attribute{
				"enabled": schema.ResourceOptionalBoolWithDefault(
					"Whether to cache responses for content not present in the proxied repository",
					common.DEFAULT_PROXY_NEGATIVE_CACHE_ENABLED,
				),
				"time_to_live": schema.ResourceOptionalInt64WithDefaultAndValidators(
					"How long to cache the fact that a file was not found in the repository (in minutes)",
					common.DEFAULT_PROXY_NEGATIVE_CACHE_TTL,
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

	if supportsRepositoryFirewall {
		thisAttr["repository_firewall"] = commonProxyFirewallAuditQuarantineAttribute(supportsPccs)
	}

	return thisAttr
}

func commonProxyFirewallAuditQuarantineAttribute(supportsPccs bool) tfschema.SingleNestedAttribute {
	thisAttr := schema.ResourceOptionalSingleNestedAttribute(
		`Sonatype Repository Firewall configuration for this Repository.
		
**Requires Sonatype Nexus Repository 3.84.0 or later.`,
		map[string]tfschema.Attribute{
			"capability_id": schema.ResourceComputedString("Internal ID of the Audit & Quarantine Capability created for this Repository"),
			"enabled":       schema.ResourceRequiredBool("Whether to enable Sonatype Repository Firewall for this Repository"),
			"quarantine":    schema.ResourceOptionalBoolWithDefault("Whether Quarantine functionallity is enabled (if false - just run in Audit mode) - see [documentation](https://help.sonatype.com/en/firewall-quarantine.html).", false),
		},
	)

	if supportsPccs {
		thisAttr.Attributes["pccs_enabled"] = schema.ResourceOptionalBoolWithDefault(
			`Whether Policy-Compliant Component Selection is enabled. See [documentatation](https://help.sonatype.com/en/policy-compliant-component-selection.html) for details.`,
			false,
		)

		// Add the validator to enforce the constraint
		thisAttr.Validators = []validator.Object{
			validators.PccsEnabledRequiresFirewallEnabled(),
		}
	}

	return thisAttr
}

func commonProxyConnectionAttribute() tfschema.SingleNestedAttribute {
	thisAttr := schema.ResourceOptionalSingleNestedAttribute(
		"HTTP Client Connection configuration for this Repository",
		map[string]tfschema.Attribute{
			"retries": schema.ResourceOptionalInt64WithDefaultAndValidators(
				"Total retries if the initial connection attempt suffers a timeout",
				common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES,
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
				common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT,
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
			"retries":                   types.Int64Value(common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES),
			"user_agent_suffix":         types.StringNull(),
			"timeout":                   types.Int64Value(common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT),
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
