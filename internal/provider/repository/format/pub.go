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

package format

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

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type PubRepositoryFormat struct {
	BaseRepositoryFormat
}

type PubRepositoryFormatProxy struct {
	PubRepositoryFormat
}

type PubRepositoryFormatGroup struct {
	PubRepositoryFormat
}

type PubRepositoryFormatHosted struct {
	PubRepositoryFormat
}

// --------------------------------------------
// Generic Pub Format Functions
// --------------------------------------------
func (f *PubRepositoryFormat) Key() string {
	return common.REPO_FORMAT_PUB
}

func (f *PubRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// PROXY Pub Format Functions
// --------------------------------------------
func (f *PubRepositoryFormatProxy) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPubProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); ignored by pre-3.94 service implementations
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Create
	return apiClient.CreatePubProxyRepository(ctx, planModel.ToApiCreateModel(), &firewallMode)
}

func (f *PubRepositoryFormatProxy) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPubProxyModel)

	// Call to API to Read
	apiResponse, firewallMode, httpResponse, err := apiClient.GetPubProxyRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, err
}

func (f *PubRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPubProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPubProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); ignored by pre-3.94 service implementations
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Create
	return apiClient.UpdatePubProxyRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel(), &firewallMode)
}

// DoImportRequest implements the import functionality for Pub Proxy repositories
func (f *PubRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, firewallMode, httpResponse, err := apiClient.GetPubProxyRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, nil
}

func (f *PubRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
}

func (f *PubRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryPubProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *PubRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryPubProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *PubRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryPubProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *PubRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryPubProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPubProxyModel)
	}

	// NXRM 3.94+ returns the repository wrapped with its inline firewall mode; use that
	// directly instead of the Capability-based UpateStateWithCapability path.
	if wrapped, ok := api.(ProxyApiResponseWithFirewall); ok {
		stateModel.FromApiModel((wrapped.Repository).(sonatyperepo.SimpleApiProxyRepository))
		if wrapped.FirewallMode != nil {
			if *wrapped.FirewallMode == common.FirewallModeDisabled {
				stateModel.FirewallAuditAndQuarantine = nil
			} else {
				if stateModel.FirewallAuditAndQuarantine == nil {
					stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineModelWithDefaults()
				}
				enabled, quarantine, _ := FirewallFlagsFromMode(*wrapped.FirewallMode)
				stateModel.FirewallAuditAndQuarantine.CapabilityId = types.StringNull()
				stateModel.FirewallAuditAndQuarantine.Enabled = types.BoolValue(enabled)
				stateModel.FirewallAuditAndQuarantine.Quarantine = types.BoolValue(quarantine)
			}
		}
		return stateModel
	}

	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

func (f *PubRepositoryFormatProxy) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryPubProxyModel)
	var stateModel model.RepositoryPubProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPubProxyModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	return stateModel
}

func (f *PubRepositoryFormatProxy) GetRepositoryId(state any) string {
	var stateModel model.RepositoryPubProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPubProxyModel)
	}
	return stateModel.Name.ValueString()
}

func (f *PubRepositoryFormatProxy) UpateStateWithCapability(state any, capability *sonatyperepo.CapabilityDTO) any {
	var stateModel = (state).(model.RepositoryPubProxyModel)
	if capability != nil {
		if stateModel.FirewallAuditAndQuarantine == nil {
			stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineModelWithDefaults()
		}
		stateModel.FirewallAuditAndQuarantine.MapFromCapabilityDTO(capability)
	}
	return stateModel
}

// Returns true only if `repository_firewall` block is supplied
func (f *PubRepositoryFormatProxy) HasFirewallConfig(state any) bool {
	var stateModel model.RepositoryPubProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPubProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine != nil {
		return true
	}
	return false
}

func (f *PubRepositoryFormatProxy) GetRepositoryFirewallEnabled(state any) bool {
	var stateModel model.RepositoryPubProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPubProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine == nil {
		return false
	}
	return stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool()
}

func (f *PubRepositoryFormatProxy) GetRepositoryFirewallQuarantineEnabled(state any) bool {
	var stateModel model.RepositoryPubProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPubProxyModel)
	}
	return stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool()
}

// --------------------------------------------
// GROUP Pub Format Functions
// --------------------------------------------
func (f *PubRepositoryFormatGroup) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPubGroupModel)

	// Call API to Create
	return apiClient.CreatePubGroupRepository(ctx, planModel.ToApiCreateModel())
}

func (f *PubRepositoryFormatGroup) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPubGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetPubGroupRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *PubRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPubGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPubGroupModel)

	// Call API to Create
	return apiClient.UpdatePubGroupRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

func (f *PubRepositoryFormatGroup) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonGroupSchemaAttributes(false)
}

func (f *PubRepositoryFormatGroup) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryPubGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *PubRepositoryFormatGroup) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryPubGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *PubRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryPubGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *PubRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryPubGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPubGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Pub Group repositories
func (f *PubRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetPubGroupRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// HOSTED Pub Format Functions
// --------------------------------------------
func (f *PubRepositoryFormatHosted) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPubHostedModel)

	// Call API to Create
	return apiClient.CreatePubHostedRepository(ctx, planModel.ToApiCreateModel())
}

func (f *PubRepositoryFormatHosted) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPubHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetPubHostedRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *PubRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPubHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPubHostedModel)

	// Call API to Update
	return apiClient.UpdatePubHostedRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

func (f *PubRepositoryFormatHosted) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonHostedSchemaAttributes()
}

func (f *PubRepositoryFormatHosted) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryPubHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *PubRepositoryFormatHosted) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryPubHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *PubRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryPubHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *PubRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryPubHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPubHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Pub Hosted repositories
func (f *PubRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetPubHostedRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}
