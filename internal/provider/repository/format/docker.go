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
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

const (
	lowercaseRepositoryNameRequiredError string = "Docker Repository Names must be lowercase for Sonatype Nexus Repository >= 3.89.0"
	pathEnabledSupportedError            string = "`path_enabled` is only supported for Sonatype Nexus Repository >= 3.83.0"
)

type DockerRepositoryFormat struct {
	BaseRepositoryFormat
}

type DockerRepositoryFormatHosted struct {
	DockerRepositoryFormat
}

type DockerRepositoryFormatProxy struct {
	DockerRepositoryFormat
}

type DockerRepositoryFormatGroup struct {
	DockerRepositoryFormat
}

// --------------------------------------------
// Generic Docker Format Functions
// --------------------------------------------
func (f *DockerRepositoryFormat) Key() string {
	return common.REPO_FORMAT_DOCKER
}

func (f *DockerRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// Hosted Docker Format Functions
// --------------------------------------------
func (f *DockerRepositoryFormatHosted) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerHostedModel)

	// Call API to Create
	return apiClient.CreateDockerHostedRepository(ctx, planModel.ToApiCreateModel())
}

func (f *DockerRepositoryFormatHosted) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetDockerHostedRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}

	// Temporary Workaround:
	// latest_policy not returned from READ API for Docker Hosted
	if stateModel.Storage.LatestPolicy.IsNull() {
		apiResponse.Storage.LatestPolicy = common.NewFalse()
	} else {
		apiResponse.Storage.LatestPolicy = stateModel.Storage.LatestPolicy.ValueBoolPointer()
	}

	return *apiResponse, httpResponse, err
}

func (f *DockerRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerHostedModel)

	// Call API to Create
	return apiClient.UpdateDockerHostedRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

// DoImportRequest implements the import functionality for Docker Hosted repositories
func (f *DockerRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetDockerHostedRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *DockerRepositoryFormatHosted) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, dockerSchemaAttributes())
	return additionalAttributes
}

func (f *DockerRepositoryFormatHosted) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryDockerHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *DockerRepositoryFormatHosted) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryDockerHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *DockerRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryDockerHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *DockerRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryDockerHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.DockerHostedApiRepository))
	return stateModel
}

func (f *DockerRepositoryFormatHosted) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	var planModel = (plan).(model.RepositoryDockerHostedModel)
	return validatePlanForDockerRespository(version, planModel.Docker.PathEnabled, planModel.Name.ValueString())
}

// --------------------------------------------
// PROXY Docker Format Functions
// --------------------------------------------
func (f *DockerRepositoryFormatProxy) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); ignored by pre-3.94 service implementations
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Create
	return apiClient.CreateDockerProxyRepository(ctx, planModel.ToApiCreateModel(), &firewallMode)
}

func (f *DockerRepositoryFormatProxy) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerProxyModel)

	// Call to API to Read
	apiResponse, firewallMode, httpResponse, err := apiClient.GetDockerProxyRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, err
}

func (f *DockerRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); ignored by pre-3.94 service implementations
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Create
	return apiClient.UpdateDockerProxyRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel(), &firewallMode)
}

// DoImportRequest implements the import functionality for Docker Proxy repositories
func (f *DockerRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, firewallMode, httpResponse, err := apiClient.GetDockerProxyRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, nil
}

func (f *DockerRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
	maps.Copy(additionalAttributes, dockerSchemaAttributes())
	maps.Copy(additionalAttributes, dockerProxySchemaAttributes())
	return additionalAttributes
}

func (f *DockerRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryDockerProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *DockerRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryDockerProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *DockerRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryDockerProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *DockerRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryDockerProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerProxyModel)
	}

	// NXRM 3.94+ returns the repository wrapped with its inline firewall mode; use that
	// directly instead of the Capability-based UpateStateWithCapability path.
	if wrapped, ok := api.(ProxyApiResponseWithFirewall); ok {
		stateModel.FromApiModel((wrapped.Repository).(sonatyperepo.DockerProxyApiRepository))
		if wrapped.FirewallMode != nil {
			if *wrapped.FirewallMode == common.FirewallModeDisabled {
				stateModel.FirewallAuditAndQuarantine = nil
			} else {
				if stateModel.FirewallAuditAndQuarantine == nil {
					stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineModelWithDefaults()
				}
				enabled, quarantine, _ := FirewallFlagsFromMode(*wrapped.FirewallMode)
				stateModel.FirewallAuditAndQuarantine.Enabled = types.BoolValue(enabled)
				stateModel.FirewallAuditAndQuarantine.Quarantine = types.BoolValue(quarantine)
			}
		}
		return stateModel
	}

	stateModel.FromApiModel((api).(sonatyperepo.DockerProxyApiRepository))
	return stateModel
}

func (f *DockerRepositoryFormatProxy) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryDockerProxyModel)
	var stateModel model.RepositoryDockerProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerProxyModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	return stateModel
}

func (f *DockerRepositoryFormatProxy) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	var planModel = (plan).(model.RepositoryDockerProxyModel)
	return validatePlanForDockerRespository(version, planModel.Docker.PathEnabled, planModel.Name.ValueString())
}

func (f *DockerRepositoryFormatProxy) GetRepositoryId(state any) string {
	var stateModel model.RepositoryDockerProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerProxyModel)
	}
	return stateModel.Name.ValueString()
}

func (f *DockerRepositoryFormatProxy) UpateStateWithCapability(state any, capability *sonatyperepo.CapabilityDTO) any {
	var stateModel = (state).(model.RepositoryDockerProxyModel)
	if capability != nil {
		if stateModel.FirewallAuditAndQuarantine == nil {
			stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineModelWithDefaults()
		}
		stateModel.FirewallAuditAndQuarantine.MapFromCapabilityDTO(capability)
	} else {
		stateModel.FirewallAuditAndQuarantine = nil
	}
	return stateModel
}

// Returns true only if `repository_firewall` block is supplied
func (f *DockerRepositoryFormatProxy) HasFirewallConfig(state any) bool {
	var stateModel model.RepositoryDockerProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine != nil {
		return true
	}
	return false
}

func (f *DockerRepositoryFormatProxy) GetRepositoryFirewallEnabled(state any) bool {
	var stateModel model.RepositoryDockerProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine == nil {
		return false
	}
	return stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool()
}

func (f *DockerRepositoryFormatProxy) GetRepositoryFirewallQuarantineEnabled(state any) bool {
	var stateModel model.RepositoryDockerProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine != nil {
		return stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool()
	}
	return false
}

// --------------------------------------------
// GROUP Docker Format Functions
// --------------------------------------------
func (f *DockerRepositoryFormatGroup) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerGroupModel)

	// Call API to Create
	return apiClient.CreateDockerGroupRepository(ctx, planModel.ToApiCreateModel())
}

func (f *DockerRepositoryFormatGroup) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetDockerGroupRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *DockerRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerGroupModel)

	// Call API to Create
	return apiClient.UpdateDockerGroupRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

// DoImportRequest implements the import functionality for Docker Group repositories
func (f *DockerRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetDockerGroupRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *DockerRepositoryFormatGroup) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonGroupSchemaAttributes(true)
	maps.Copy(additionalAttributes, dockerSchemaAttributes())
	return additionalAttributes
}

func (f *DockerRepositoryFormatGroup) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryDockerGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *DockerRepositoryFormatGroup) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryDockerGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *DockerRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryDockerGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *DockerRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryDockerGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.DockerGroupApiRepository))
	return stateModel
}

func (f *DockerRepositoryFormatGroup) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	var planModel = (plan).(model.RepositoryDockerGroupModel)
	return validatePlanForDockerRespository(version, planModel.Docker.PathEnabled, planModel.Name.ValueString())
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func dockerSchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"docker": schema.ResourceRequiredSingleNestedAttribute(
			"Docker specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"force_basic_auth": schema.ResourceRequiredBool("Whether to force authentication (Docker Bearer Token Realm required if false)"),
				"http_port":        schema.ResourceOptionalInt32("Create an HTTP connector at specified port"),
				"https_port":       schema.ResourceOptionalInt32("Create an HTTPS connector at specified port"),
				"path_enabled": schema.ResourceOptionalBoolWithPlanModifier(
					"Allows to use repository name in Docker image paths (only supply for Sonatype Nexus Repository Manager >= 3.83.0)",
					boolplanmodifier.UseStateForUnknown(),
				),
				"subdomain":  schema.ResourceOptionalString("Allows to use subdomain"),
				"v1_enabled": schema.ResourceRequiredBool("Whether to allow clients to use the V1 API to interact with this repository"),
			},
		),
	}
}

func dockerProxySchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"docker_proxy": schema.ResourceRequiredSingleNestedAttribute(
			"Docker Proxy specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"cache_foreign_layers": schema.ResourceComputedOptionalBoolWithDefault(
					"Allow Nexus Repository Manager to download and cache foreign layers",
					false,
				),
				"foreign_layer_url_whitelist": func() tfschema.SetAttribute {
					thisAttr := schema.ResourceOptionalStringSet("Foreign Layer URL Whitelist")
					thisAttr.Computed = true
					thisAttr.Default = setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{}))
					return thisAttr
				}(),
				"index_type": schema.ResourceStringEnumWithDefault(
					"Type of Docker Index",
					common.DOCKER_PROXY_INDEX_TYPE_REGISTRY,
					common.DOCKER_PROXY_INDEX_TYPE_HUB,
					common.DOCKER_PROXY_INDEX_TYPE_REGISTRY,
					common.DOCKER_PROXY_INDEX_TYPE_CUSTOM,
				),
				"index_url": schema.ResourceOptionalString("Url of Docker Index to use"),
			},
		),
	}
}

func validatePlanForDockerRespository(version common.SystemVersion, pathEnabled types.Bool, repositoryName string) []string {
	if version.RequiresLowerCaseRepostioryNameDocker() && strings.IndexFunc(repositoryName, unicode.IsUpper) != -1 {
		return []string{lowercaseRepositoryNameRequiredError}
	}

	if !pathEnabled.IsNull() && version.OlderThan(3, 83, 0, 0) {
		return []string{pathEnabledSupportedError}
	}

	return nil
}
