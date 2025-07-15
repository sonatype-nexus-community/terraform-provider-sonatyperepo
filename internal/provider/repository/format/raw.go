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

type RawRepositoryFormat struct {
	BaseRepositoryFormat
}

type RawRepositoryFormatHosted struct {
	RawRepositoryFormat
}

type RawRepositoryFormatProxy struct {
	RawRepositoryFormat
}

type RawRepositoryFormatGroup struct {
	RawRepositoryFormat
}

// --------------------------------------------
// Generic Raw Format Functions
// --------------------------------------------
func (f *RawRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_RAW
}

func (f *RawRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted Raw Format Functions
// --------------------------------------------
func (f *RawRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRawHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRawHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RawRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRawHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRawHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *RawRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRawHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRawHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRawHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RawRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, getRawSchemaAttributes())
	return additionalAttributes
}

func (f *RawRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryRawHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RawRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryRawHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RawRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryRawHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RawRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryRawHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.RawHostedApiRepository))
	return stateModel
}

// --------------------------------------------
// PROXY Raw Format Functions
// --------------------------------------------
func (f *RawRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRawProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRawProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RawRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRawProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRawProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *RawRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRawProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRawProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRawProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RawRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonProxySchemaAttributes()
	maps.Copy(additionalAttributes, getRawSchemaAttributes())
	return additionalAttributes
}

func (f *RawRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryRawProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RawRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryRawProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RawRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryRawProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RawRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryRawProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.RawProxyApiRepository))
	return stateModel
}

// --------------------------------------------
// GORUP Raw Format Functions
// --------------------------------------------
func (f *RawRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRawGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateRawGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *RawRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRawGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetRawGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *RawRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryRawGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryRawGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateRawGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *RawRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonGroupSchemaAttributes(false)
	maps.Copy(additionalAttributes, getRawSchemaAttributes())
	return additionalAttributes
}

func (f *RawRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryRawGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RawRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryRawGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RawRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryRawGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RawRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryRawGroupModel)
	stateModel.FromApiModel((api).(sonatyperepo.RawGroupApiRepository))
	return stateModel
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getRawSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"raw": schema.SingleNestedAttribute{
			Description: "Raw specific configuration for this Repository",
			Required:    true,
			Attributes: map[string]schema.Attribute{
				"content_disposition": schema.StringAttribute{
					Description: "Content Disposition",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(common.CONTENT_DISPOSITION_ATTACHMENT, common.CONTENT_DISPOSITION_INLINE),
					},
				},
			},
		},
	}
}
