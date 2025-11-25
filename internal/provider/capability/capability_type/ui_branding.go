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
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// --------------------------------------------
// Capabiltiy Type: UI: Branding
// --------------------------------------------
type UiBrandingCapability struct {
	BaseCapabilityType
}

func NewUiBrandingCapability() *UiBrandingCapability {
	return &UiBrandingCapability{
		BaseCapabilityType: BaseCapabilityType{
			capabilityType: common.CAPABILITY_TYPE_UI_BRANDING,
			publicName:     "UI Branding",
		},
	}
}

// --------------------------------------------
// Capabiltiy Type: UI: Branding Functions
// --------------------------------------------
func (f *UiBrandingCapability) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.UiBrandingCapabilityModel)

	// Call API to Create
	return apiClient.CapabilitiesAPI.Create3(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *UiBrandingCapability) DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.UiBrandingCapabilityModel)

	// Call API to Update
	return apiClient.CapabilitiesAPI.Update3(ctx, capabilityId).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *UiBrandingCapability) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.UiBrandingCapabilityModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *UiBrandingCapability) GetPropertiesSchema() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"footer_enabled": schema.ResourceComputedOptionalBoolWithDefault(
			"Enable branding header HTML snippet.",
			common.CAPABILITY_UI_BRANDING_DEFAULT_FOOTER_ENABLED,
		),
		"footer_html": schema.ResourceOptionalStringWithDefaultAndPlanModifier(
			"An HTML snippet to be included in branding header. Use '$baseUrl' to insert the base URL of the server (e.g. to reference an image).",
			common.CAPABILITY_UI_BRANDING_DEFAULT_FOOTER_HTML,
		),
		"header_enabled": schema.ResourceComputedOptionalBoolWithDefault(
			"Enable branding header HTML snippet.",
			common.CAPABILITY_UI_BRANDING_DEFAULT_HEADER_ENABLED,
		),
		"header_html": schema.ResourceOptionalStringWithDefaultAndPlanModifier(
			"An HTML snippet to be included in branding header. Use '$baseUrl' to insert the base URL of the server (e.g. to reference an image).",
			"",
		),
	}
}

func (f *UiBrandingCapability) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.UiBrandingCapabilityModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *UiBrandingCapability) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.UiBrandingCapabilityModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *UiBrandingCapability) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.UiBrandingCapabilityModel)
	apiModel := (api).(*v3.CapabilityDTO)
	stateModel.FromApiModel(apiModel)
	stateModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return stateModel
}

func (f *UiBrandingCapability) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.UiBrandingCapabilityModel)
	stateModel := (state).(model.UiBrandingCapabilityModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
