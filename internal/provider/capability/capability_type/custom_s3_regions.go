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

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// --------------------------------------------
// Capabiltiy Type: Custom S3 Regions
// --------------------------------------------
type CustomS3RegionsCapability struct {
	BaseCapabilityType
}

func NewCustomS3RegionsCapability() *CustomS3RegionsCapability {
	return &CustomS3RegionsCapability{
		BaseCapabilityType: BaseCapabilityType{
			capabilityType: common.CAPABILITY_TYPE_CUSTOM_S3_REGIONS,
			publicName:     "Custom S3 Regions",
		},
	}
}

// --------------------------------------------
// Capabiltiy Type: Custom S3 Regions Functions
// --------------------------------------------
func (f *CustomS3RegionsCapability) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.CapabilityCustomS3RegionsModel)

	// Call API to Create
	return apiClient.CapabilitiesAPI.Create3(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *CustomS3RegionsCapability) DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.CapabilityCustomS3RegionsModel)

	// Call API to Update
	return apiClient.CapabilitiesAPI.Update3(ctx, capabilityId).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *CustomS3RegionsCapability) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.CapabilityCustomS3RegionsModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *CustomS3RegionsCapability) PropertiesSchema() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"regions": schema.ResourceRequiredStringSetWithValidator(
			"Custom S3 Regions.",
			setvalidator.SizeAtLeast(1),
		),
	}
}

func (f *CustomS3RegionsCapability) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.CapabilityCustomS3RegionsModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *CustomS3RegionsCapability) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.CapabilityCustomS3RegionsModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *CustomS3RegionsCapability) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.CapabilityCustomS3RegionsModel)
	apiModel := (api).(*v3.CapabilityDTO)
	stateModel.FromApiModel(apiModel)
	stateModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return stateModel
}

func (f *CustomS3RegionsCapability) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.CapabilityCustomS3RegionsModel)
	stateModel := (state).(model.CapabilityCustomS3RegionsModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
