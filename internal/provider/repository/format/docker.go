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
func (f *DockerRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateDockerHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *DockerRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetDockerHostedRepository(ctx, stateModel.Name.ValueString()).Execute()

	// Temporary Workaround:
	// latest_policy not returned from READ API for Docker Hosted
	if stateModel.Storage.LatestPolicy.IsNull() {
		apiResponse.Storage.LatestPolicy = common.NewFalse()
	} else {
		apiResponse.Storage.LatestPolicy = stateModel.Storage.LatestPolicy.ValueBoolPointer()
	}

	return *apiResponse, httpResponse, err
}

func (f *DockerRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateDockerHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

// DoImportRequest implements the import functionality for Docker Hosted repositories
func (f *DockerRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetDockerHostedRepository(ctx, repositoryName).Execute()
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
func (f *DockerRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateDockerProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *DockerRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetDockerProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *DockerRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateDockerProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

// DoImportRequest implements the import functionality for Docker Proxy repositories
func (f *DockerRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetDockerProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
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
	stateModel.FromApiModel((api).(sonatyperepo.DockerProxyApiRepository))
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
	stateModel.FirewallAuditAndQuarantine.MapFromCapabilityDTO(capability)
	return stateModel
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
	return stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool()
}

// --------------------------------------------
// GROUP Docker Format Functions
// --------------------------------------------
func (f *DockerRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateDockerGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *DockerRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetDockerGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *DockerRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateDockerGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

// DoImportRequest implements the import functionality for Docker Group repositories
func (f *DockerRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetDockerGroupRepository(ctx, repositoryName).Execute()
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
	var planModel model.RepositoryDockerroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *DockerRepositoryFormatGroup) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryDockerroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *DockerRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryDockerroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *DockerRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryDockerroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryDockerroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.DockerGroupApiRepository))
	return stateModel
}

func (f *DockerRepositoryFormatGroup) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	var planModel = (plan).(model.RepositoryDockerHostedModel)
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
