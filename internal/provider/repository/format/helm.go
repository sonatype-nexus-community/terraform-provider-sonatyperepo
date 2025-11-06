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
	"net/http"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type HelmRepositoryFormat struct {
	BaseRepositoryFormat
}

type HelmRepositoryFormatHosted struct {
	HelmRepositoryFormat
}

type HelmRepositoryFormatProxy struct {
	HelmRepositoryFormat
}

// --------------------------------------------
// Generic Helm Format Functions
// --------------------------------------------
func (f *HelmRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_HELM
}

func (f *HelmRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted Helm Format Functions
// --------------------------------------------
func (f *HelmRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryHelmHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateHelmHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *HelmRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryHelmHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetHelmHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *HelmRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryHelmHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryHelmHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateHelmHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *HelmRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return getCommonHostedSchemaAttributes()
}

func (f *HelmRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryHelmHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *HelmRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryHelmHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *HelmRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryHelmHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *HelmRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryHelmHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryHelmHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Helm Hosted repositories
func (f *HelmRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetHelmHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// ValidateRepositoryForImport validates that the imported repository is indeed a Helm Hosted repository
func (f *HelmRepositoryFormatHosted) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// Cast to Helm Hosted API Repository
	apiRepo, ok := repositoryData.(sonatyperepo.SimpleApiHostedRepository)
	if !ok {
		return fmt.Errorf("repository data is not a Helm Hosted repository")
	}

	if apiRepo.Format == nil {
		return fmt.Errorf(errRepositoryFormatNil, expectedFormat)
	}
	// Case-insensitive format comparison
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

// --------------------------------------------
// PROXY Helm Format Functions
// --------------------------------------------
func (f *HelmRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryHelmProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateHelmProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *HelmRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryHelmProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetHelmProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *HelmRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryHelmProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryHelmProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateHelmProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *HelmRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return getCommonProxySchemaAttributes()
}

func (f *HelmRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryHelmProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *HelmRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryHelmProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *HelmRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryHelmProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *HelmRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryHelmProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryHelmProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Helm Proxy repositories
func (f *HelmRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetHelmProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// ValidateRepositoryForImport validates that the imported repository is indeed a Helm Proxy repository
func (f *HelmRepositoryFormatProxy) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// Cast to Helm Proxy API Repository
	apiRepo, ok := repositoryData.(sonatyperepo.SimpleApiProxyRepository)
	if !ok {
		return fmt.Errorf("repository data is not a Helm Proxy repository")
	}

	if apiRepo.Format == nil {
		return fmt.Errorf(errRepositoryFormatNil, expectedFormat)
	}
	// Case-insensitive format comparison
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
