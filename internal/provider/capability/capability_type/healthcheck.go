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
// Capabiltiy Type: Healthcheck
// --------------------------------------------
type HealthcheckCapability struct {
	BaseCapabilityType
}

func NewHealthcheckCapability() *HealthcheckCapability {
	return &HealthcheckCapability{
		BaseCapabilityType: BaseCapabilityType{
			capabilityType: common.CAPABILITY_TYPE_HEALTHCHECK,
			publicName:     "Healthcheck",
		},
	}
}

// --------------------------------------------
// Capabiltiy Type: Healthcheck Functions
// --------------------------------------------
func (f *HealthcheckCapability) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.CapabilityHealthcheckModel)

	// Call API to Create
	return apiClient.CapabilitiesAPI.Create3(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *HealthcheckCapability) DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.CapabilityHealthcheckModel)

	// Call API to Update
	return apiClient.CapabilitiesAPI.Update3(ctx, capabilityId).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *HealthcheckCapability) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.CapabilityHealthcheckModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *HealthcheckCapability) PropertiesSchema() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"configured_for_all_proxies": schema.ResourceOptionalBoolWithDefault(
			`Configure all supported proxy repositories to regularly check with [Sonatype Repository Healthcheck](https://help.sonatype.com/en/repository-health-check.html) for updates by default. Newly added repositories will automatically be configured as well, for as long as this is selected.`,
			common.CAPABILITY_HEALTHCHECK_DEFAULT_CONFIGURED_FOR_ALL,
		),
		"use_nexus_truststore": schema.ResourceOptionalBoolWithDefault(
			`Whether to use Nexus Truststore when communicating with Sonatype Repository Healthcheck.
			
  The RHC service works by performing calls to the following Sonatype data services depending on the Nexus Repository license agreement in use. 
  Network administrators need to allow these URLs through their network firewall to receive updates. 
  - https://rhc-pro.sonatype.com
  - https://rhc.sonatype.com`,
			common.CAPABILITY_HEALTHCHECK_DEFAULT_USE_NEXUS_TRUSTSTORE,
		),
	}
}

func (f *HealthcheckCapability) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.CapabilityHealthcheckModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *HealthcheckCapability) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.CapabilityHealthcheckModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *HealthcheckCapability) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.CapabilityHealthcheckModel)
	apiModel := (api).(*v3.CapabilityDTO)
	stateModel.FromApiModel(apiModel)
	stateModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return stateModel
}

func (f *HealthcheckCapability) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.CapabilityHealthcheckModel)
	stateModel := (state).(model.CapabilityHealthcheckModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
