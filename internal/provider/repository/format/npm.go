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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type NpmRepositoryFormat struct{}

func (f *NpmRepositoryFormat) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateNpmHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *NpmRepositoryFormat) DoDeleteRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.DeleteRepository(ctx, stateModel.Name.ValueString()).Execute()
}

func (f *NpmRepositoryFormat) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetNpmHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *NpmRepositoryFormat) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryNpmHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryNpmHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateNpmHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *NpmRepositoryFormat) GetApiCreateSuccessResposneCodes() []int {
	return []int{http.StatusCreated}
}

func (f *NpmRepositoryFormat) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		// "format": schema.StringAttribute{
		// 	Description: fmt.Sprintf("Format of this Repository - will always be '%s'", f.GetKey()),
		// 	Optional:    true,
		// 	Computed:    true,
		// 	Default:     stringdefault.StaticString(f.GetKey()),
		// 	PlanModifiers: []planmodifier.String{
		// 		stringplanmodifier.UseStateForUnknown(),
		// 	},
		// },
	}
}

func (f *NpmRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_NPM
}

func (f *NpmRepositoryFormat) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryNpmHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *NpmRepositoryFormat) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryNpmHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *NpmRepositoryFormat) GetResourceName(req resource.MetadataRequest) string {
	return getResourceName(req, f.GetKey())
}

func (f *NpmRepositoryFormat) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryNpmHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *NpmRepositoryFormat) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryNpmHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}
