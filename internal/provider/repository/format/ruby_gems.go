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

type RubyGemsRepositoryFormat struct {
	BaseRepositoryFormat
}

type RubyGemsRepositoryFormatHosted struct {
	RubyGemsRepositoryFormat
}

type RubyGemsRepositoryFormatProxy struct {
	RubyGemsRepositoryFormat
}

type RubyGemsRepositoryFormatGroup struct {
	RubyGemsRepositoryFormat
}

// --------------------------------------------
// Generic RubyGems Format Functions
// --------------------------------------------
func (f *RubyGemsRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_RUBY_GEMS
}

func (f *RubyGemsRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	// Override to maintain backward compatibility with resource name containing underscore
	return resourceName("ruby_gems", repoType)
}

// --------------------------------------------
// Hosted RubyGems Format Functions
// --------------------------------------------
func (f *RubyGemsRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRubyGemsHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRubygemsHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RubyGemsRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRubyGemsHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRubygemsHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *RubyGemsRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRubyGemsHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRubyGemsHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRubygemsHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RubyGemsRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonHostedSchemaAttributes()
}

func (f *RubyGemsRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryRubyGemsHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RubyGemsRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryRubyGemsHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RubyGemsRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryRubyGemsHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RubyGemsRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryRubyGemsHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryRubyGemsHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for RubyGems Hosted repositories
func (f *RubyGemsRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRubygemsHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// PROXY RubyGems Format Functions
// --------------------------------------------
func (f *RubyGemsRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorRubyGemsProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRubygemsProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RubyGemsRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorRubyGemsProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRubygemsProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *RubyGemsRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorRubyGemsProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorRubyGemsProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRubygemsProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RubyGemsRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonProxySchemaAttributes()
}

func (f *RubyGemsRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositorRubyGemsProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RubyGemsRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositorRubyGemsProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RubyGemsRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositorRubyGemsProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RubyGemsRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositorRubyGemsProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositorRubyGemsProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for RubyGems Proxy repositories
func (f *RubyGemsRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRubygemsProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// GORUP RubyGems Format Functions
// --------------------------------------------
func (f *RubyGemsRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRubyGemsGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRubygemsGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RubyGemsRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRubyGemsGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRubygemsGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *RubyGemsRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRubyGemsGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRubyGemsGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRubygemsGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RubyGemsRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonGroupSchemaAttributes(false)
}

func (f *RubyGemsRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryRubyGemsGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RubyGemsRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryRubyGemsGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RubyGemsRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryRubyGemsGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RubyGemsRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryRubyGemsGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryRubyGemsGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for RubyGems Group repositories
func (f *RubyGemsRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRubygemsGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}
