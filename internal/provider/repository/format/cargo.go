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

type CargoRepositoryFormat struct {
	BaseRepositoryFormat
}

type CargoRepositoryFormatHosted struct {
	CargoRepositoryFormat
}

type CargoRepositoryFormatProxy struct {
	CargoRepositoryFormat
}

type CargoRepositoryFormatGroup struct {
	CargoRepositoryFormat
}

// --------------------------------------------
// Generic Cargo Format Functions
// --------------------------------------------
func (f *CargoRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_CARGO
}

func (f *CargoRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted Cargo Format Functions
// --------------------------------------------
func (f *CargoRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorCargoHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateCargoHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *CargoRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorCargoHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetCargoHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *CargoRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorCargoHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorCargoHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateCargoHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *CargoRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return getCommonHostedSchemaAttributes()
}

func (f *CargoRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositorCargoHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *CargoRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositorCargoHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *CargoRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositorCargoHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *CargoRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositorCargoHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// --------------------------------------------
// PROXY Cargo Format Functions
// --------------------------------------------
func (f *CargoRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryCargoProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateCargoProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *CargoRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryCargoProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetCargoProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *CargoRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryCargoProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryCargoProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateCargoProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *CargoRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonProxySchemaAttributes()
	maps.Copy(additionalAttributes, getCargoSchemaAttributes())
	return additionalAttributes
}

func (f *CargoRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryCargoProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *CargoRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryCargoProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *CargoRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryCargoProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *CargoRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryCargoProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.CargoProxyApiRepository))
	return stateModel
}

// --------------------------------------------
// GORUP Cargo Format Functions
// --------------------------------------------
func (f *CargoRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryCargoGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateCargoGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *CargoRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryCargoGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetCargoGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *CargoRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryCargoGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryCargoGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateCargoGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *CargoRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttrs := getCommonGroupSchemaAttributes(false)
	maps.Copy(additionalAttrs, getCargoSchemaAttributes())
	return additionalAttrs
}

func (f *CargoRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryCargoGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *CargoRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryCargoGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *CargoRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryCargoGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *CargoRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryCargoGroupModel)
	stateModel.FromApiModel((api).(sonatyperepo.CargoGroupApiRepository))
	return stateModel
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getCargoSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"cargo": schema.SingleNestedAttribute{
			Description: "Cargo specific configuration for this Repository",
			Required:    true,
			Attributes: map[string]schema.Attribute{
				"require_authentication": schema.BoolAttribute{
					Description: "Indicates if this repository requires authentication overriding anonymous access.",
					Required:    true,
				},
			},
		},
	}
}
