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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

// --------------------------------------------
// Generic NPM Format Functions
// --------------------------------------------
func (f *NpmRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_NPM
}

func (f *NpmRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted NPM Format Functions
// --------------------------------------------
func (f *NpmRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateNpmHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *NpmRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetNpmHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *NpmRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateNpmHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *NpmRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonHostedSchemaAttributes()
	// maps.Copy(additionalAttributes, getMavenSchemaAttributes())
	return additionalAttributes
}

func (f *NpmRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNpmHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NpmRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNpmHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NpmRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNpmHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NpmRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryNpmHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// --------------------------------------------
// Hosted NPM Format Functions
// --------------------------------------------
func (f *NpmRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateNpmProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *NpmRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetNpmProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *NpmRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateNpmProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *NpmRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonProxySchemaAttributes()
	maps.Copy(additionalAttributes, getNpmSchemaAttributes())
	return additionalAttributes
}

func (f *NpmRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNpmProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NpmRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNpmProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NpmRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNpmProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NpmRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryNpmProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.NpmProxyApiRepository))
	return stateModel
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getNpmSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"npm": schema.SingleNestedAttribute{
			Description: "NPM specific configuration for this Repository",
			Required:    false,
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"remove_quarrantined": schema.BoolAttribute{
					Description: "Remove Quarantined Versions",
					Required:    true,
					Optional:    false,
				},
			},
		},
	}
}
