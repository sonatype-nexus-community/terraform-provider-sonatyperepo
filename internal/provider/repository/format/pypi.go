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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

type PyPiRepositoryFormat struct {
	BaseRepositoryFormat
}

type PyPiRepositoryFormatHosted struct {
	PyPiRepositoryFormat
}

type PyPiRepositoryFormatProxy struct {
	PyPiRepositoryFormat
}

type PyPiRepositoryFormatGroup struct {
	PyPiRepositoryFormat
}

// --------------------------------------------
// Generic PyPi Format Functions
// --------------------------------------------
func (f *PyPiRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_PYPI
}

func (f *PyPiRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return resourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted PyPi Format Functions
// --------------------------------------------
func (f *PyPiRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPyPiHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreatePypiHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *PyPiRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPyPiHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetPypiHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *PyPiRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPyPiHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPyPiHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdatePypiHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *PyPiRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonHostedSchemaAttributes()
}

func (f *PyPiRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryPyPiHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *PyPiRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryPyPiHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *PyPiRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryPyPiHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *PyPiRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryPyPiHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPyPiHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for PyPI Hosted repositories
func (f *PyPiRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetPypiHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// PROXY PyPi Format Functions
// --------------------------------------------
func (f *PyPiRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPyPiProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreatePypiProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *PyPiRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPyPiProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetPypiProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *PyPiRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPyPiProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPyPiProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdatePypiProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *PyPiRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonProxySchemaAttributes()
	maps.Copy(additionalAttributes, pyPiSchemaAttributes())
	return additionalAttributes
}

func (f *PyPiRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryPyPiProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *PyPiRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryPyPiProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *PyPiRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryPyPiProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *PyPiRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryPyPiProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPyPiProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.PyPiProxyApiRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for PyPI Proxy repositories
func (f *PyPiRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetPypiProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// GORUP PyPi Format Functions
// --------------------------------------------
func (f *PyPiRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPyPiGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreatePypiGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *PyPiRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPyPiGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetPypiGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *PyPiRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryPyPiGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryPyPiGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdatePypiGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *PyPiRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonGroupSchemaAttributes(true)
}

func (f *PyPiRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryPyPiGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *PyPiRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryPyPiGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *PyPiRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryPyPiGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *PyPiRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryPyPiGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryPyPiGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupDeployRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for PyPI Group repositories
func (f *PyPiRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetPypiGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func pyPiSchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"pypi": schema.ResourceRequiredSingleNestedAttribute(
			"PyPi specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"remove_quarrantined": schema.ResourceRequiredBool("Remove Quarantined Versions"),
			},
		),
	}
}
