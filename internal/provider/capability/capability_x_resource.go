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
	"github.com/hashicorp/terraform-plugin-framework/resource"

	capabilitytype "terraform-provider-sonatyperepo/internal/provider/capability/capability_type"
)

// NewCapabilityAuditResource is a helper function to simplify the provider implementation.
func NewCapabilityAuditResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewNewAuditCapability(),
	}
}

// NewCapabilityCoreBaseUrlResource is a helper function to simplify the provider implementation.
func NewCapabilityCoreBaseUrlResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewCoreBaseUrlCapability(),
	}
}

// NewCapabilityCoreStorageSettingsResource is a helper function to simplify the provider implementation.
func NewCapabilityCoreStorageSettingsResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewCoreStorageSettingsCapability(),
	}
}

// NewCapabilityCustomS3RegionsResource is a helper function to simplify the provider implementation.
func NewCapabilityCustomS3RegionsResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewCustomS3RegionsCapability(),
	}
}

// NewCapabilityDefaultRoleResource is a helper function to simplify the provider implementation.
func NewCapabilityDefaultRoleResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewDefaultRoleCapability(),
	}
}

// NewCapabilityFirewallAuditQuarantineResource is a helper function to simplify the provider implementation.
func NewCapabilityFirewallAuditQuarantineResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewFirewallAuditQuarantineCapability(),
	}
}

// NewCapabilityHealthcheckResource is a helper function to simplify the provider implementation.
func NewCapabilityHealthcheckResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewHealthcheckCapability(),
	}
}

// NewCapabilityOutreachResource is a helper function to simplify the provider implementation.
func NewCapabilityOutreachResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewOutreachCapability(),
	}
}

// NewCapabilityUiBrandingResource is a helper function to simplify the provider implementation.
func NewCapabilityUiBrandingResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewUiBrandingCapability(),
	}
}

// NewCapabilitySecurityRutAuthResource is a helper function to simplify the provider implementation.
func NewCapabilitySecurityRutAuthResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewRutAuthCapability(),
	}
}

// NewCapabilityUiSettingsResource is a helper function to simplify the provider implementation.
func NewCapabilityUiSettingsResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewUiSettingsCapability(),
	}
}

// NewCapabilityWebhookRepositoryResource is a helper function to simplify the provider implementation.
func NewCapabilityWebhookRepositoryResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewWebhookRepositoryCapability(),
	}
}

// NewCapabilityWebhookGlobalResource is a helper function to simplify the provider implementation.
func NewCapabilityWebhookGlobalResource() resource.Resource {
	return &capabilityResource{
		CapabilityType: capabilitytype.NewWebhookGlobalCapability(),
	}
}
