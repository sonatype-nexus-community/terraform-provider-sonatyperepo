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
// Capabiltiy Type: Webhook: Repository
// --------------------------------------------
type WebhookRepositoryCapability struct {
	BaseCapabilityType
}

func NewWebhookRepositoryCapability() *WebhookRepositoryCapability {
	return &WebhookRepositoryCapability{
		BaseCapabilityType: BaseCapabilityType{
			capabilityType: common.CAPABILITY_TYPE_WEBHOOK_REPOSITORY,
			publicName:     "Webhook Repository",
		},
	}
}

// --------------------------------------------
// Capabiltiy Type: Webhook: Repository
// --------------------------------------------
func (f *WebhookRepositoryCapability) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.WebhookRepositoryCapabilityModel)

	// Call API to Create
	return apiClient.CapabilitiesAPI.Create3(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *WebhookRepositoryCapability) DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.WebhookRepositoryCapabilityModel)

	// Call API to Update
	return apiClient.CapabilitiesAPI.Update3(ctx, capabilityId).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *WebhookRepositoryCapability) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.WebhookRepositoryCapabilityModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *WebhookRepositoryCapability) PropertiesSchema() map[string]tfschema.Attribute {
	return propertiesSchemaForWebhookCapability(common.AllRepositoryWebHookEventTypes(), true)
}

func (f *WebhookRepositoryCapability) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.WebhookRepositoryCapabilityModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *WebhookRepositoryCapability) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.WebhookRepositoryCapabilityModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *WebhookRepositoryCapability) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.WebhookRepositoryCapabilityModel)
	apiModel := (api).(*v3.CapabilityDTO)
	stateModel.FromApiModel(apiModel)
	stateModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return stateModel
}

func (ct *WebhookRepositoryCapability) MapFromPlanToState(plan any, state any) any {
	planModel := (plan).(model.WebhookRepositoryCapabilityModel)
	stateModel := (state).(model.WebhookRepositoryCapabilityModel)
	stateModel.Properties.Secret = types.StringValue(planModel.Properties.Secret.ValueString())
	return stateModel
}

func (f *WebhookRepositoryCapability) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.WebhookRepositoryCapabilityModel)
	stateModel := (state).(model.WebhookRepositoryCapabilityModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
