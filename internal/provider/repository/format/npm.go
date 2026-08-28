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

type NpmRepositoryFormat struct {
	BaseRepositoryFormat
}

type NpmRepositoryFormatHosted struct {
	NpmRepositoryFormat
}

type NpmRepositoryFormatProxy struct {
	NpmRepositoryFormat
}

type NpmRepositoryFormatGroup struct {
	NpmRepositoryFormat
}

// --------------------------------------------
// Generic NPM Format Functions
// --------------------------------------------
func (f *NpmRepositoryFormat) Key() string {
	return common.REPO_FORMAT_NPM
}

func (f *NpmRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// Hosted NPM Format Functions
// --------------------------------------------
func (f *NpmRepositoryFormatHosted) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmHostedModel)

	// Call API to Create
	return apiClient.CreateNpmHostedRepository(ctx, planModel.ToApiCreateModel())
}

func (f *NpmRepositoryFormatHosted) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetNpmHostedRepository(ctx, stateModel.Name.ValueString())
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *NpmRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmHostedModel)

	// Call API to Create
	return apiClient.UpdateNpmHostedRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

func (f *NpmRepositoryFormatHosted) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonHostedSchemaAttributes()
	return additionalAttributes
}

func (f *NpmRepositoryFormatHosted) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNpmHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NpmRepositoryFormatHosted) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNpmHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NpmRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNpmHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NpmRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryNpmHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for NPM Hosted repositories
func (f *NpmRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetNpmHostedRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// PROXY NPM Format Functions
// --------------------------------------------
func (f *NpmRepositoryFormatProxy) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); ignored by pre-3.94 service implementations
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Create
	return apiClient.CreateNpmProxyRepository(ctx, planModel.ToApiCreateModel(), &firewallMode)
}

func (f *NpmRepositoryFormatProxy) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmProxyModel)

	// Get the repository name
	repoName := stateModel.Name.ValueString()

	// Call to API to Read
	apiResponse, firewallMode, httpResponse, err := apiClient.GetNpmProxyRepository(ctx, repoName)
	if err != nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, err
}

func (f *NpmRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); ignored by pre-3.94 service implementations
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Create
	return apiClient.UpdateNpmProxyRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel(), &firewallMode)
}

// DoImportRequest implements the import functionality for NPM Proxy repositories
func (f *NpmRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, firewallMode, httpResponse, err := apiClient.GetNpmProxyRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, nil
}

func (f *NpmRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
}

func (f *NpmRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNpmProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NpmRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNpmProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NpmRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNpmProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NpmRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryNpmProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmProxyModel)
	}

	// NXRM 3.94+ returns the repository wrapped with its inline firewall mode; use that
	// directly instead of the Capability-based UpateStateWithCapability path.
	if wrapped, ok := api.(ProxyApiResponseWithFirewall); ok {
		stateModel.FromApiModel((wrapped.Repository).(sonatyperepo.NpmProxyApiRepository))
		if wrapped.FirewallMode != nil {
			keep, enabled, quarantine, pccsEnabled := ResolveFirewallBlockFlags(stateModel.FirewallAuditAndQuarantine != nil, *wrapped.FirewallMode)
			if !keep {
				stateModel.FirewallAuditAndQuarantine = nil
			} else {
				if stateModel.FirewallAuditAndQuarantine == nil {
					stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineWithPccsModelWithDefaults()
				}
				stateModel.FirewallAuditAndQuarantine.CapabilityId = types.StringNull()
				stateModel.FirewallAuditAndQuarantine.Enabled = types.BoolValue(enabled)
				stateModel.FirewallAuditAndQuarantine.Quarantine = types.BoolValue(quarantine)
				stateModel.FirewallAuditAndQuarantine.PccsEnabled = types.BoolValue(pccsEnabled)
			}
		}
		return stateModel
	}

	stateModel.FromApiModel((api).(sonatyperepo.NpmProxyApiRepository))
	return stateModel
}

func (f *NpmRepositoryFormatProxy) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryNpmProxyModel)
	var stateModel model.RepositoryNpmProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmProxyModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	stateModel.FirewallAuditAndQuarantine = ReconcileFirewallBlockWithPccsPlan(stateModel.FirewallAuditAndQuarantine, planModel.FirewallAuditAndQuarantine)
	return stateModel
}

// NPM Proxy Repositories support Repository Firewall PCCS
func (f *NpmRepositoryFormatProxy) SupportsRepositoryFirewallPccs() bool {
	return true
}

func (f *NpmRepositoryFormatProxy) GetRepositoryId(state any) string {
	var stateModel model.RepositoryNpmProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmProxyModel)
	}
	return stateModel.Name.ValueString()
}

func (f *NpmRepositoryFormatProxy) UpateStateWithCapability(state any, capability *sonatyperepo.CapabilityDTO) any {
	var stateModel = (state).(model.RepositoryNpmProxyModel)
	if capability != nil {
		if stateModel.FirewallAuditAndQuarantine == nil {
			stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineWithPccsModelWithDefaults()
		}
		stateModel.FirewallAuditAndQuarantine.MapFromCapabilityDTO(capability)
	}
	return stateModel
}

// Returns true only if `repository_firewall` block is supplied
func (f *NpmRepositoryFormatProxy) HasFirewallConfig(state any) bool {
	var stateModel model.RepositoryNpmProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine != nil {
		return true
	}
	return false
}

func (f *NpmRepositoryFormatProxy) GetRepositoryFirewallEnabled(state any) bool {
	var stateModel model.RepositoryNpmProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine == nil {
		return false
	}
	return stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool()
}

func (f *NpmRepositoryFormatProxy) GetRepositoryFirewallQuarantineEnabled(state any) bool {
	var stateModel model.RepositoryNpmProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmProxyModel)
	}
	return stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool()
}

func (f *NpmRepositoryFormatProxy) GetRepositoryFirewallPccsEnabled(state any) bool {
	var stateModel model.RepositoryNpmProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine != nil {
		return stateModel.FirewallAuditAndQuarantine.PccsEnabled.ValueBool()
	} else {
		return false
	}
}

// --------------------------------------------
// GORUP NPM Format Functions
// --------------------------------------------
func (f *NpmRepositoryFormatGroup) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmGroupModel)

	// Call API to Create
	return apiClient.CreateNpmGroupRepository(ctx, planModel.ToApiCreateModel())
}

func (f *NpmRepositoryFormatGroup) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetNpmGroupRepository(ctx, stateModel.Name.ValueString())
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *NpmRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmGroupModel)

	// Call API to Create
	return apiClient.UpdateNpmGroupRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

func (f *NpmRepositoryFormatGroup) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonGroupSchemaAttributes(true)
}

func (f *NpmRepositoryFormatGroup) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNpmGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NpmRepositoryFormatGroup) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNpmGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NpmRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNpmGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NpmRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryNpmGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryNpmGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupDeployRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for NPM Group repositories
func (f *NpmRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetNpmGroupRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}
