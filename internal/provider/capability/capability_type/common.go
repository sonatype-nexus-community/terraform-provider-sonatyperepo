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

package capabilitytype

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// --------------------------------------------
// BaseCapabilityType that all capability types build from
// --------------------------------------------
type BaseCapabilityType struct {
	capabilityType common.CapabilityType
	publicName     string
}

func (ct *BaseCapabilityType) ApiCreateSuccessResponseCodes() []int {
	return []int{http.StatusOK}
}

func (ct *BaseCapabilityType) Key() string {
	return ct.capabilityType.String()
}

func (ct *BaseCapabilityType) GetMarkdownDescription() string {
	return fmt.Sprintf("Manage Capability: %s", ct.publicName)
}

func (ct *BaseCapabilityType) GetPublicName() string {
	return ct.publicName
}

func (ct *BaseCapabilityType) ResourceName() string {
	return fmt.Sprintf("capability_%s", common.SanitiseStringForResourceName(ct.GetPublicName()))
}

func (ct *BaseCapabilityType) GetType() common.CapabilityType {
	return ct.capabilityType
}

// --------------------------------------------
// CapabilityTypeI that all Capability Types must implement
// --------------------------------------------
type CapabilityTypeI interface {
	DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error)
	DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error)
	ApiCreateSuccessResponseCodes() []int
	GetMarkdownDescription() string
	PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics)
	PropertiesSchema() map[string]tfschema.Attribute
	StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics)
	ResourceName() string
	Key() string
	GetPublicName() string
	GetType() common.CapabilityType
	UpdatePlanForState(plan any) any
	UpdateStateFromApi(state any, api any) any
	UpdateStateFromPlanForUpdate(plan any, state any) any
}

// --------------------------------------------
// Helper method to generate schema for Webhook Capabilities
// --------------------------------------------
func propertiesSchemaForWebhookCapability(permissibleEventTypes []string, includeRepository bool) map[string]tfschema.Attribute {
	defaultProps := map[string]tfschema.Attribute{
		"names": schema.ResourceRequiredStringSetWithValidator(
			"Event types which trigger this Webhook.",
			setvalidator.All(
				setvalidator.SizeBetween(1, 2),
				setvalidator.ValueStringsAre(
					stringvalidator.OneOf(permissibleEventTypes...),
				),
			),
		),
		"secret": schema.ResourceSensitiveString("Key to use for HMAC payload digest."),
		"url": schema.ResourceRequiredStringWithRegex(
			"Send a HTTP POST request to this URL.",
			regexp.MustCompile(`^https?://[^\s]+$`),
			"Must be a valid http:// or https:// URL",
		),
	}

	if includeRepository {
		defaultProps["repository"] = schema.ResourceRequiredStringWithRegex(
			"Repository to discriminate events from.",
			regexp.MustCompile(`^[a-zA-Z0-9\-]{1}[a-zA-Z0-9_\-\.]*$`),
			"Must be a valid repository name",
		)
	}

	return defaultProps
}
