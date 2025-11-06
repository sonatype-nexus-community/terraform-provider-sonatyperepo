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
	"fmt"
	"maps"
	"net/http"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

const (
	errRepositoryFormatNil      = "repository format is nil, expected '%s'"
	errRepositoryFormatMismatch = "repository format is '%s', expected '%s'"
	errRepositoryTypeNil        = "repository type is nil, expected '%s'"
	errRepositoryTypeMismatch   = "repository type is '%s', expected '%s'"
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
func (f *DockerRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_DOCKER
}

func (f *DockerRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
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

// ValidateRepositoryForImport validates that the imported repository is indeed a Docker Hosted repository
func (f *DockerRepositoryFormatHosted) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// Cast to Docker Hosted API Repository
	apiRepo, ok := repositoryData.(sonatyperepo.DockerHostedApiRepository)
	if !ok {
		return fmt.Errorf("repository data is not a Docker Hosted repository")
	}

	if apiRepo.Format == nil {
		return fmt.Errorf(errRepositoryFormatNil, expectedFormat)
	}
	// Convert both to lowercase for comparison
	actualFormat := strings.ToLower(*apiRepo.Format)
	expectedFormatLower := strings.ToLower(expectedFormat)
	if actualFormat != expectedFormatLower {
		return fmt.Errorf(errRepositoryFormatMismatch, *apiRepo.Format, expectedFormat)
	}

	// Validate type
	expectedTypeStr := expectedType.String()
	if apiRepo.Type == nil {
		return fmt.Errorf(errRepositoryTypeNil, expectedTypeStr)
	}
	if *apiRepo.Type != expectedTypeStr {
		return fmt.Errorf(errRepositoryTypeMismatch, *apiRepo.Type, expectedTypeStr)
	}

	return nil
}

func (f *DockerRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, getDockerSchemaAttributes())
	return additionalAttributes
}

func (f *DockerRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryDockerHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *DockerRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryDockerHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *DockerRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryDockerHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *DockerRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryDockerHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.DockerHostedApiRepository))
	return stateModel
}

func (f *DockerRepositoryFormatHosted) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	var planModel = (plan).(model.RepositoryDockerHostedModel)
	if !planModel.Docker.PathEnabled.IsNull() && version.OlderThan(3, 83, 0, 0) {
		return []string{
			"`path_enabled` is only supported for Sonatype Nexus Repository >= 3.83.0",
		}
	}
	return nil
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

// ValidateRepositoryForImport validates that the imported repository is indeed a Docker Proxy repository
func (f *DockerRepositoryFormatProxy) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// Cast to Docker Proxy API Repository
	apiRepo, ok := repositoryData.(sonatyperepo.DockerProxyApiRepository)
	if !ok {
		return fmt.Errorf("repository data is not a Docker Proxy repository")
	}

	// Validate format (case-insensitive)
	if apiRepo.Format == nil {
		return fmt.Errorf(errRepositoryFormatNil, expectedFormat)
	}
	// Convert both to lowercase for comparison
	actualFormat := strings.ToLower(*apiRepo.Format)
	expectedFormatLower := strings.ToLower(expectedFormat)
	if actualFormat != expectedFormatLower {
		return fmt.Errorf(errRepositoryFormatMismatch, *apiRepo.Format, expectedFormat)
	}

	// Validate type
	expectedTypeStr := expectedType.String()
	if apiRepo.Type == nil {
		return fmt.Errorf(errRepositoryTypeNil, expectedTypeStr)
	}
	if *apiRepo.Type != expectedTypeStr {
		return fmt.Errorf(errRepositoryTypeMismatch, *apiRepo.Type, expectedTypeStr)
	}

	return nil
}

func (f *DockerRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonProxySchemaAttributes()
	maps.Copy(additionalAttributes, getDockerSchemaAttributes())
	maps.Copy(additionalAttributes, getDockerProxySchemaAttributes())
	return additionalAttributes
}

func (f *DockerRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryDockerProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *DockerRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
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
	if !planModel.Docker.PathEnabled.IsNull() && version.OlderThan(3, 83, 0, 0) {
		return []string{
			"`path_enabled` is only supported for Sonatype Nexus Repository >= 3.83.0",
		}
	}
	return nil
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

// ValidateRepositoryForImport validates that the imported repository is indeed a Docker Group repository
func (f *DockerRepositoryFormatGroup) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// Cast to Docker Group API Repository
	apiRepo, ok := repositoryData.(sonatyperepo.DockerGroupApiRepository)
	if !ok {
		return fmt.Errorf("repository data is not a Docker Group repository")
	}

	if apiRepo.Format == nil {
		return fmt.Errorf(errRepositoryFormatNil, expectedFormat)
	}
	// Convert both to lowercase for comparison
	actualFormat := strings.ToLower(*apiRepo.Format)
	expectedFormatLower := strings.ToLower(expectedFormat)
	if actualFormat != expectedFormatLower {
		return fmt.Errorf(errRepositoryFormatMismatch, *apiRepo.Format, expectedFormat)
	}

	// Validate type
	expectedTypeStr := expectedType.String()
	if apiRepo.Type == nil {
		return fmt.Errorf(errRepositoryTypeNil, expectedTypeStr)
	}
	if *apiRepo.Type != expectedTypeStr {
		return fmt.Errorf(errRepositoryTypeMismatch, *apiRepo.Type, expectedTypeStr)
	}

	return nil
}

func (f *DockerRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonGroupSchemaAttributes(true)
	maps.Copy(additionalAttributes, getDockerSchemaAttributes())
	return additionalAttributes
}

func (f *DockerRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryDockerroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *DockerRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryDockerroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *DockerRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryDockerroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *DockerRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryDockerroupModel)
	stateModel.FromApiModel((api).(sonatyperepo.DockerGroupApiRepository))
	return stateModel
}

func (f *DockerRepositoryFormatGroup) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	var planModel = (plan).(model.RepositoryDockerroupModel)
	if !planModel.Docker.PathEnabled.IsNull() && version.OlderThan(3, 83, 0, 0) {
		return []string{
			"`path_enabled` is only supported for Sonatype Nexus Repository >= 3.83.0",
		}
	}
	return nil
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getDockerSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"docker": schema.SingleNestedAttribute{
			Description: "Docker specific configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes: map[string]schema.Attribute{
				"force_basic_auth": schema.BoolAttribute{
					Description: "Whether to force authentication (Docker Bearer Token Realm required if false)",
					Required:    true,
				},
				"http_port": schema.Int32Attribute{
					Description: "Create an HTTP connector at specified port",
					Optional:    true,
				},
				"https_port": schema.Int32Attribute{
					Description: "Create an HTTPS connector at specified port",
					Optional:    true,
				},
				"path_enabled": schema.BoolAttribute{
					Description: "Allows to use repository name in Docker image paths (only supply for Sonatype Nexus Repository Manager >= 3.83.0)",
					Optional:    true,
					PlanModifiers: []planmodifier.Bool{
						boolplanmodifier.UseStateForUnknown(),
					},
				},
				"subdomain": schema.StringAttribute{
					Description: "Allows to use subdomain",
					Optional:    true,
				},
				"v1_enabled": schema.BoolAttribute{
					Description: "Whether to allow clients to use the V1 API to interact with this repository",
					Required:    true,
				},
			},
		},
	}
}

func getDockerProxySchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"docker_proxy": schema.SingleNestedAttribute{
			Description: "Docker Proxy specific configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes: map[string]schema.Attribute{
				"cache_foreign_layers": schema.BoolAttribute{
					Description: "Allow Nexus Repository Manager to download and cache foreign layers",
					Optional:    true,
					Computed:    true,
					Default:     booldefault.StaticBool(false),
				},
				"foreign_layer_url_whitelist": schema.SetAttribute{
					Description: "Foreign Layer URL Whitelist",
					Optional:    true,
					Computed:    true,
					ElementType: types.StringType,
					Default:     setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
				},
				"index_type": schema.StringAttribute{
					Description: "Type of Docker Index",
					Optional:    true,
					Computed:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							common.DOCKER_PROXY_INDEX_TYPE_HUB,
							common.DOCKER_PROXY_INDEX_TYPE_REGISTRY,
							common.DOCKER_PROXY_INDEX_TYPE_CUSTOM,
						),
					},
					Default: stringdefault.StaticString(common.DOCKER_PROXY_INDEX_TYPE_REGISTRY),
				},
				"index_url": schema.StringAttribute{
					Description: "Url of Docker Index to use",
					Optional:    true,
				},
			},
		},
	}
}