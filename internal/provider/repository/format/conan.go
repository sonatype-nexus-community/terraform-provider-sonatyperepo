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
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

type ConanRepositoryFormat struct {
	BaseRepositoryFormat
}

type ConanRepositoryFormatHosted struct {
	ConanRepositoryFormat
}

type ConanRepositoryFormatProxy struct {
	ConanRepositoryFormat
}

type ConanRepositoryFormatGroup struct {
	ConanRepositoryFormat
}

// --------------------------------------------
// Generic Conan Format Functions
// --------------------------------------------
func (f *ConanRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_CONAN
}

func (f *ConanRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return resourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted Conan Format Functions
// --------------------------------------------
func (f *ConanRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorConanHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateConanHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *ConanRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorConanHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetConanHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *ConanRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorConanHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorConanHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateConanHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *ConanRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonHostedSchemaAttributes()
}

func (f *ConanRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositorConanHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *ConanRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositorConanHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *ConanRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositorConanHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *ConanRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositorConanHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}

// --------------------------------------------
// PROXY Conan Format Functions
// --------------------------------------------
func (f *ConanRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryConanProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateConanProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *ConanRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryConanProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetConanProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *ConanRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryConanProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryConanProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateConanProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *ConanRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonProxySchemaAttributes()
	maps.Copy(additionalAttributes, conanProxySchemaAttributes())
	return additionalAttributes
}

func (f *ConanRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryConanProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *ConanRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryConanProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *ConanRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryConanProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *ConanRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryConanProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

// --------------------------------------------
// GORUP Conan Format Functions
// --------------------------------------------
func (f *ConanRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryConanGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateConanGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *ConanRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryConanGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetConanGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *ConanRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryConanGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryConanGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateConanGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *ConanRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonGroupSchemaAttributes(true)
}

func (f *ConanRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryConanGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *ConanRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryConanGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *ConanRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryConanGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *ConanRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryConanGroupModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupDeployRepository))
	return stateModel
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func conanProxySchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"conan": schema.ResourceRequiredSingleNestedAttribute(
			"Conan Proxy specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"conan_version": func() tfschema.StringAttribute {
					thisAttr := schema.ResourceRequiredStringWithValidators(
						"Conan protocol version. Cannot be changed once repository is created.",
						stringvalidator.OneOf(common.CONAN_PROTOCOL_V1, common.CONAN_PROTOCOL_V2),
					)
					thisAttr.PlanModifiers = []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					}
					return thisAttr
				}(),
			},
		),
	}
}
