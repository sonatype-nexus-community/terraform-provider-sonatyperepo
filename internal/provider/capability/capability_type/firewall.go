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
// Capabiltiy Type: Firewall Audit & Quarantine
// --------------------------------------------
type FirewallAuditQuarantineCapability struct {
	BaseCapabilityType
}

func NewFirewallAuditQuarantineCapability() *FirewallAuditQuarantineCapability {
	return &FirewallAuditQuarantineCapability{
		BaseCapabilityType: BaseCapabilityType{
			capabilityType: common.CAPABILITY_TYPE_FIREWALL_AUDIT_QUARANTINE,
			publicName:     "Firewall Audit and Quarantine",
		},
	}
}

// --------------------------------------------
// Capabiltiy Type: Base URL Functions
// --------------------------------------------
func (f *FirewallAuditQuarantineCapability) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.CapabilityFirewallAuditQuarantineModel)

	// Call API to Create
	return apiClient.CapabilitiesAPI.Create3(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *FirewallAuditQuarantineCapability) DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.CapabilityFirewallAuditQuarantineModel)

	// Call API to Update
	return apiClient.CapabilitiesAPI.Update3(ctx, capabilityId).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *FirewallAuditQuarantineCapability) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.CapabilityFirewallAuditQuarantineModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *FirewallAuditQuarantineCapability) GetPropertiesSchema() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"repository": schema.ResourceRequiredString("The repository to be evaluated."),
		"quarantine": schema.ResourceRequiredBool(`Whether enable Quarantine for this repository. 
			
**Note:** If enabled and later disabled, all quarantined components will be made available in the repository; those components cannot be re-quarantined.`),
	}
}

func (f *FirewallAuditQuarantineCapability) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.CapabilityFirewallAuditQuarantineModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *FirewallAuditQuarantineCapability) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.CapabilityFirewallAuditQuarantineModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *FirewallAuditQuarantineCapability) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.CapabilityFirewallAuditQuarantineModel)
	apiModel := (api).(*v3.CapabilityDTO)
	stateModel.FromApiModel(apiModel)
	stateModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return stateModel
}

func (f *FirewallAuditQuarantineCapability) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.CapabilityFirewallAuditQuarantineModel)
	stateModel := (state).(model.CapabilityFirewallAuditQuarantineModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
