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

package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// pccsEnabledRequiresFirewallEnabled validates that pccs_enabled cannot be true when enabled is false
type pccsEnabledRequiresFirewallEnabled struct{}

// Description returns a plain text description of the validator's behavior.
func (v pccsEnabledRequiresFirewallEnabled) Description(ctx context.Context) string {
	return "pccs_enabled cannot be true when enabled is false"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior.
func (v pccsEnabledRequiresFirewallEnabled) MarkdownDescription(ctx context.Context) string {
	return "pccs_enabled cannot be true when enabled is false"
}

// ValidateObject performs the validation.
func (v pccsEnabledRequiresFirewallEnabled) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	// If the entire object is null or unknown, skip validation
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Get the object's attributes
	data := req.ConfigValue.Attributes()

	// Get the enabled field
	enabledAttr, hasEnabled := data["enabled"]
	if !hasEnabled {
		return
	}

	if enabledAttr.IsNull() || enabledAttr.IsUnknown() {
		return
	}

	// Convert to Bool type
	var enabled types.Bool
	if err := tfsdk.ValueAs(ctx, enabledAttr, &enabled); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	// Get the pccs_enabled field
	pccsEnabledAttr, hasPccsEnabled := data["pccs_enabled"]
	if !hasPccsEnabled {
		return
	}

	if pccsEnabledAttr.IsNull() || pccsEnabledAttr.IsUnknown() {
		return
	}

	// Convert to Bool type
	var pccsEnabled types.Bool
	if err := tfsdk.ValueAs(ctx, pccsEnabledAttr, &pccsEnabled); err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	// Validation: pccs_enabled cannot be true if enabled is false
	if pccsEnabled.ValueBool() && !enabled.ValueBool() {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			path.Root("repository_firewall").AtName("pccs_enabled"),
			"Invalid Attribute Combination",
			"pccs_enabled cannot be true when enabled is false. Policy-Compliant Component Selection requires the Repository Firewall to be enabled.",
		))
	}
}

// PccsEnabledRequiresFirewallEnabled returns a validator that ensures pccs_enabled is only true when enabled is true
func PccsEnabledRequiresFirewallEnabled() validator.Object {
	return pccsEnabledRequiresFirewallEnabled{}
}
