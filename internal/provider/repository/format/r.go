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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type RRepositoryFormat struct {
	BaseRepositoryFormat
}

type RRepositoryFormatHosted struct {
	RRepositoryFormat
}

type RRepositoryFormatProxy struct {
	RRepositoryFormat
}

type RRepositoryFormatGroup struct {
	RRepositoryFormat
}

// --------------------------------------------
// Generic R(CRAN) Format Functions
// --------------------------------------------
func (f *RRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_R
}

func (f *RRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return resourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted R(CRAN) Format Functions
// --------------------------------------------
func (f *RRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *RRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return commonHostedSchemaAttributes()
}

func (f *RRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryRHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryRHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryRHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryRHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// --------------------------------------------
// PROXY R(CRAN) Format Functions
// --------------------------------------------
func (f *RRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorRProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorRProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *RRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorRProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorRProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return commonProxySchemaAttributes()
}

func (f *RRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositorRProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositorRProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositorRProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositorRProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

// --------------------------------------------
// GORUP R(CRAN) Format Functions
// --------------------------------------------
func (f *RRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *RRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return commonGroupSchemaAttributes(false)
}

func (f *RRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryRGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryRGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryRGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryRGroupModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupRepository))
	return stateModel
}
