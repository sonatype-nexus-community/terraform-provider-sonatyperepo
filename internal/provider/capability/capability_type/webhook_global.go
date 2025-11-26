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
)

// --------------------------------------------
// Capabiltiy Type: Webhook: Global
// --------------------------------------------
type WebhookGlobalCapability struct {
	BaseCapabilityType
}

func NewWebhookGlobalCapability() *WebhookGlobalCapability {
	return &WebhookGlobalCapability{
		BaseCapabilityType: BaseCapabilityType{
			capabilityType: common.CAPABILITY_TYPE_WEBHOOK_GLOBAL,
			publicName:     "Webhook Global",
		},
	}
}

// --------------------------------------------
// Capabiltiy Type: Webhook: Repository
// --------------------------------------------
func (f *WebhookGlobalCapability) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.WebhookGlobalCapabilityModel)

	// Call API to Create
	return apiClient.CapabilitiesAPI.Create3(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *WebhookGlobalCapability) DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.WebhookGlobalCapabilityModel)

	// Call API to Update
	return apiClient.CapabilitiesAPI.Update3(ctx, capabilityId).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *WebhookGlobalCapability) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.WebhookGlobalCapabilityModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *WebhookGlobalCapability) PropertiesSchema() map[string]tfschema.Attribute {
	return propertiesSchemaForWebhookCapability(common.AllGlobalWebHookEventTypes(), false)
}

func (f *WebhookGlobalCapability) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.WebhookGlobalCapabilityModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *WebhookGlobalCapability) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.WebhookGlobalCapabilityModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *WebhookGlobalCapability) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.WebhookGlobalCapabilityModel)
	apiModel := (api).(*v3.CapabilityDTO)
	stateModel.FromApiModel(apiModel)
	stateModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return stateModel
}

func (f *WebhookGlobalCapability) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.WebhookGlobalCapabilityModel)
	stateModel := (state).(model.WebhookGlobalCapabilityModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
