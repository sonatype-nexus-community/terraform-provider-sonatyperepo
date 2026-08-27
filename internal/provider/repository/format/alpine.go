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
	"maps"
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

type AlpineRepositoryFormat struct {
	BaseRepositoryFormat
}

type AlpineRepositoryFormatHosted struct {
	AlpineRepositoryFormat
}

type AlpineRepositoryFormatProxy struct {
	AlpineRepositoryFormat
}

type AlpineRepositoryFormatGroup struct {
	AlpineRepositoryFormat
}

// --------------------------------------------
// Generic Alpine Format Functions
// --------------------------------------------
func (f *AlpineRepositoryFormat) Key() string {
	return common.REPO_FORMAT_ALPINE
}

func (f *AlpineRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// Hosted Alpine Format Functions
// --------------------------------------------
func (f *AlpineRepositoryFormatHosted) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAlpineHostedModel)

	// Call API to Create
	return apiClient.CreateAlpineHostedRepository(ctx, planModel.ToApiCreateModel())
}

func (f *AlpineRepositoryFormatHosted) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAlpineHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetAlpineHostedRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *AlpineRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAlpineHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAlpineHostedModel)

	// Call API to Update
	return apiClient.UpdateAlpineHostedRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

func (f *AlpineRepositoryFormatHosted) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return alpineSchemaAttributes()
}

func (f *AlpineRepositoryFormatHosted) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAlpineHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *AlpineRepositoryFormatHosted) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAlpineHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *AlpineRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAlpineHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *AlpineRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryAlpineHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.AlpineHostedApiRepository))
	return stateModel
}

func (f *AlpineRepositoryFormatHosted) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryAlpineHostedModel)
	var stateModel model.RepositoryAlpineHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineHostedModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	return stateModel
}

// DoImportRequest implements the import functionality for Alpine Hosted repositories
func (f *AlpineRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetAlpineHostedRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// PROXY Alpine Format Functions
// --------------------------------------------
func (f *AlpineRepositoryFormatProxy) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAlpineProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); ignored by pre-3.94 service implementations
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Create
	return apiClient.CreateAlpineProxyRepository(ctx, planModel.ToApiCreateModel(), &firewallMode)
}

func (f *AlpineRepositoryFormatProxy) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAlpineProxyModel)

	// Call to API to Read
	apiResponse, firewallMode, httpResponse, err := apiClient.GetAlpineProxyRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, err
}

func (f *AlpineRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAlpineProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAlpineProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); ignored by pre-3.94 service implementations
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Update
	return apiClient.UpdateAlpineProxyRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel(), &firewallMode)
}

// DoImportRequest implements the import functionality for Alpine Proxy repositories
func (f *AlpineRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, firewallMode, httpResponse, err := apiClient.GetAlpineProxyRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, nil
}

func (f *AlpineRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
	maps.Copy(additionalAttributes, alpineSchemaAttributes())
	return additionalAttributes
}

func (f *AlpineRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAlpineProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *AlpineRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAlpineProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *AlpineRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAlpineProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *AlpineRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryAlpineProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineProxyModel)
	}

	// NXRM 3.94+ returns the repository wrapped with its inline firewall mode; use that
	// directly instead of the Capability-based UpateStateWithCapability path.
	if wrapped, ok := api.(ProxyApiResponseWithFirewall); ok {
		stateModel.FromApiModel((wrapped.Repository).(sonatyperepo.AlpineProxyApiRepository))
		if wrapped.FirewallMode != nil {
			keep, enabled, quarantine, _ := ResolveFirewallBlockFlags(stateModel.FirewallAuditAndQuarantine != nil, *wrapped.FirewallMode)
			if !keep {
				stateModel.FirewallAuditAndQuarantine = nil
			} else {
				if stateModel.FirewallAuditAndQuarantine == nil {
					stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineModelWithDefaults()
				}
				stateModel.FirewallAuditAndQuarantine.CapabilityId = types.StringNull()
				stateModel.FirewallAuditAndQuarantine.Enabled = types.BoolValue(enabled)
				stateModel.FirewallAuditAndQuarantine.Quarantine = types.BoolValue(quarantine)
			}
		}
		return stateModel
	}

	stateModel.FromApiModel((api).(sonatyperepo.AlpineProxyApiRepository))
	return stateModel
}

func (f *AlpineRepositoryFormatProxy) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryAlpineProxyModel)
	var stateModel model.RepositoryAlpineProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineProxyModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	stateModel.FirewallAuditAndQuarantine = BackfillFirewallBlockFromPlan(stateModel.FirewallAuditAndQuarantine, planModel.FirewallAuditAndQuarantine)
	return stateModel
}

func (f *AlpineRepositoryFormatProxy) GetRepositoryId(state any) string {
	var stateModel model.RepositoryAlpineProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineProxyModel)
	}
	return stateModel.Name.ValueString()
}

func (f *AlpineRepositoryFormatProxy) UpateStateWithCapability(state any, capability *sonatyperepo.CapabilityDTO) any {
	var stateModel = (state).(model.RepositoryAlpineProxyModel)
	if capability != nil {
		if stateModel.FirewallAuditAndQuarantine == nil {
			stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineModelWithDefaults()
		}
		stateModel.FirewallAuditAndQuarantine.MapFromCapabilityDTO(capability)
	}
	return stateModel
}

// Returns true only if `repository_firewall` block is supplied
func (f *AlpineRepositoryFormatProxy) HasFirewallConfig(state any) bool {
	var stateModel model.RepositoryAlpineProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine != nil {
		return true
	}
	return false
}

func (f *AlpineRepositoryFormatProxy) GetRepositoryFirewallEnabled(state any) bool {
	var stateModel model.RepositoryAlpineProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine == nil {
		return false
	}
	return stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool()
}

func (f *AlpineRepositoryFormatProxy) GetRepositoryFirewallQuarantineEnabled(state any) bool {
	var stateModel model.RepositoryAlpineProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineProxyModel)
	}
	return stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool()
}

// --------------------------------------------
// GROUP Alpine Format Functions
// --------------------------------------------
func (f *AlpineRepositoryFormatGroup) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAlpineGroupModel)

	// Call API to Create
	return apiClient.CreateAlpineGroupRepository(ctx, planModel.ToApiCreateModel())
}

func (f *AlpineRepositoryFormatGroup) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAlpineGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetAlpineGroupRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *AlpineRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAlpineGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAlpineGroupModel)

	// Call API to Update
	return apiClient.UpdateAlpineGroupRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

func (f *AlpineRepositoryFormatGroup) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonGroupSchemaAttributes(false)
	maps.Copy(additionalAttributes, alpineSchemaAttributes())
	return additionalAttributes
}

func (f *AlpineRepositoryFormatGroup) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAlpineGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *AlpineRepositoryFormatGroup) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAlpineGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *AlpineRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAlpineGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *AlpineRepositoryFormatGroup) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryAlpineGroupModel)
	var stateModel model.RepositoryAlpineGroupModel
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineGroupModel)
	}
	stateModel.MapMissingApiFieldsFromPlan(planModel)
	return stateModel
}

func (f *AlpineRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryAlpineGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAlpineGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.AlpineGroupApiRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Alpine Group repositories
func (f *AlpineRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetAlpineGroupRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func alpineSchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"alpine": schema.ResourceOptionalSingleNestedAttribute(
			"Alpine specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"key_pair": schema.ResourceRequiredString(
					"RSA private key in PEM format used to sign the repository index",
				),
				"passphrase": schema.ResourceSensitiveOptionalStringWithPlanModifier(
					"Passphrase to access the RSA signing key",
					stringplanmodifier.UseStateForUnknown(),
				),
			},
		),
	}
}
