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

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type NugetRepositoryFormat struct {
	BaseRepositoryFormat
}

type NugetRepositoryFormatHosted struct {
	NugetRepositoryFormat
}

type NugetRepositoryFormatProxy struct {
	NugetRepositoryFormat
}

type NugetRepositoryFormatGroup struct {
	NugetRepositoryFormat
}

// --------------------------------------------
// Generic Nuget Format Functions
// --------------------------------------------
func (f *NugetRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_NUGET
}

func (f *NugetRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted Nuget Format Functions
// --------------------------------------------
func (f *NugetRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNugetHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateNugetHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *NugetRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNugetHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetNugetHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *NugetRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNugetHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNugetHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateNugetHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *NugetRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonHostedSchemaAttributes()
	return additionalAttributes
}

func (f *NugetRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNugetHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NugetRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNugetHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NugetRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNugetHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NugetRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryNugetHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// --------------------------------------------
// PROXY Nuget Format Functions
// --------------------------------------------
func (f *NugetRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNugetProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateNugetProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *NugetRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNugetProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetNugetProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *NugetRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNugetProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNugetProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateNugetProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *NugetRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonProxySchemaAttributes()
	maps.Copy(additionalAttributes, getNugetProxySchemaAttributes())
	return additionalAttributes
}

func (f *NugetRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNugetProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NugetRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNugetProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NugetRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNugetProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NugetRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryNugetProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.NugetProxyApiRepository))
	return stateModel
}

// --------------------------------------------
// GORUP Nuget Format Functions
// --------------------------------------------
func (f *NugetRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNugetGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateNugetGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *NugetRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNugetGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetNugetGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *NugetRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNugetGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNugetGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateNugetGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *NugetRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return getCommonGroupSchemaAttributes(false)
}

func (f *NugetRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNugetGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NugetRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNugetGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NugetRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNugetGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NugetRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryNugetGroupModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupRepository))
	return stateModel
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getNugetProxySchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"nuget_proxy": schema.SingleNestedAttribute{
			Description: "Nuget specific configuration for this Repository",
			Required:    true,
			Attributes: map[string]schema.Attribute{
				"nuget_version": schema.StringAttribute{
					Description: "Nuget Protocol Versions",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(common.NUGET_PROTOCOL_V2, common.NUGET_PROTOCOL_V3),
					},
				},
				"query_cache_item_max_age": schema.Int32Attribute{
					Description: "How long to cache query results from the proxied repository (in seconds)",
					Optional:    true,
				},
			},
		},
	}
}
