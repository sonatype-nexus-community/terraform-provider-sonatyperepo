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
// Capabiltiy Type: RUT Auth
// --------------------------------------------
type RutAuthCapability struct {
	BaseCapabilityType
}

func NewRutAuthCapability() *RutAuthCapability {
	return &RutAuthCapability{
		BaseCapabilityType: BaseCapabilityType{
			capabilityType: common.CAPABILITY_TYPE_RUT_AUTH,
			publicName:     "RUT Auth",
		},
	}
}

// --------------------------------------------
// Capabiltiy Type: RUT Auth Functions
// --------------------------------------------
func (f *RutAuthCapability) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.CapabilityRutAuthModel)

	// Call API to Create
	return apiClient.CapabilitiesAPI.Create3(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *RutAuthCapability) DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.CapabilityRutAuthModel)
	planModel.Id = types.StringValue(capabilityId)

	// Call API to Update
	return apiClient.CapabilitiesAPI.Update3(ctx, capabilityId).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *RutAuthCapability) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.CapabilityRutAuthModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RutAuthCapability) PropertiesSchema() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"http_header": schema.ResourceRequiredStringWithLengthAtLeast(
			"Handled HTTP Header should contain the name of the header that is used to source the principal of already authenticated user.",
			1,
		),
	}
}

func (f *RutAuthCapability) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.CapabilityRutAuthModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RutAuthCapability) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.CapabilityRutAuthModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RutAuthCapability) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.CapabilityRutAuthModel)
	apiModel := (api).(*v3.CapabilityDTO)
	stateModel.FromApiModel(apiModel)
	stateModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return stateModel
}

func (f *RutAuthCapability) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.CapabilityRutAuthModel)
	stateModel := (state).(model.CapabilityRutAuthModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
