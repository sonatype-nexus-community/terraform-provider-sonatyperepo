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
	stateModel := (state).(model.RepositoryHelmHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
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
	stateModel := (state).(model.RepositoryHelmProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}
