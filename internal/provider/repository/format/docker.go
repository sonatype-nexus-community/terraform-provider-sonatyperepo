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

type DockerRepositoryFormat struct {
	BaseRepositoryFormat
}

type DockerRepositoryFormatHosted struct {
	DockerRepositoryFormat
}

type DockerRepositoryFormatProxy struct {
	DockerRepositoryFormat
}

type DockerRepositoryFormatGroup struct {
	DockerRepositoryFormat
}

// --------------------------------------------
// Generic Docker Format Functions
// --------------------------------------------
func (f *DockerRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_DOCKER
}

func (f *DockerRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted Docker Format Functions
// --------------------------------------------
func (f *DockerRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateDockerHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *DockerRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetDockerHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *DockerRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryDockerHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryDockerHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateDockerHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *DockerRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, getDockerSchemaAttributes())
	return additionalAttributes
}

func (f *DockerRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryDockerHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *DockerRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryDockerHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *DockerRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryDockerHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *DockerRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryDockerHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.DockerHostedApiRepository))
	return stateModel
}

// --------------------------------------------
// PROXY Docker Format Functions
// --------------------------------------------
// func (f *NpmRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
// 	// Cast to correct Plan Model Type
// 	planModel := (plan).(model.RepositoryNpmProxyModel)

// 	// Call API to Create
// 	return apiClient.RepositoryManagementAPI.CreateNpmProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
// }

// func (f *NpmRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
// 	// Cast to correct State Model Type
// 	stateModel := (state).(model.RepositoryNpmProxyModel)

// 	// Call to API to Read
// 	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetNpmProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
// 	return *apiResponse, httpResponse, err
// }

// func (f *NpmRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
// 	// Cast to correct Plan Model Type
// 	planModel := (plan).(model.RepositoryNpmProxyModel)

// 	// Cast to correct State Model Type
// 	stateModel := (state).(model.RepositoryNpmProxyModel)

// 	// Call API to Create
// 	return apiClient.RepositoryManagementAPI.UpdateNpmProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
// }

// func (f *NpmRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
// 	additionalAttributes := getCommonProxySchemaAttributes()
// 	maps.Copy(additionalAttributes, getNpmSchemaAttributes())
// 	return additionalAttributes
// }

// func (f *NpmRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
// 	var planModel model.RepositoryNpmProxyModel
// 	return planModel, plan.Get(ctx, &planModel)
// }

// func (f *NpmRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
// 	var stateModel model.RepositoryNpmProxyModel
// 	return stateModel, state.Get(ctx, &stateModel)
// }

// func (f *NpmRepositoryFormatProxy) UpdatePlanForState(plan any) any {
// 	var planModel = (plan).(model.RepositoryNpmProxyModel)
// 	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
// 	return planModel
// }

// func (f *NpmRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
// 	stateModel := (state).(model.RepositoryNpmProxyModel)
// 	stateModel.FromApiModel((api).(sonatyperepo.NpmProxyApiRepository))
// 	return stateModel
// }

// // --------------------------------------------
// // GORUP Docker Format Functions
// // --------------------------------------------
// func (f *NpmRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
// 	// Cast to correct Plan Model Type
// 	planModel := (plan).(model.RepositoryNpmGroupModel)

// 	// Call API to Create
// 	return apiClient.RepositoryManagementAPI.CreateNpmGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
// }

// func (f *NpmRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
// 	// Cast to correct State Model Type
// 	stateModel := (state).(model.RepositoryNpmGroupModel)

// 	// Call to API to Read
// 	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetNpmGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
// 	return *apiResponse, httpResponse, err
// }

// func (f *NpmRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
// 	// Cast to correct Plan Model Type
// 	planModel := (plan).(model.RepositoryNpmGroupModel)

// 	// Cast to correct State Model Type
// 	stateModel := (state).(model.RepositoryNpmGroupModel)

// 	// Call API to Create
// 	return apiClient.RepositoryManagementAPI.UpdateNpmGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
// }

// func (f *NpmRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]schema.Attribute {
// 	return getCommonGroupSchemaAttributes(true)
// }

// func (f *NpmRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
// 	var planModel model.RepositoryNpmGroupModel
// 	return planModel, plan.Get(ctx, &planModel)
// }

// func (f *NpmRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
// 	var stateModel model.RepositoryNpmGroupModel
// 	return stateModel, state.Get(ctx, &stateModel)
// }

// func (f *NpmRepositoryFormatGroup) UpdatePlanForState(plan any) any {
// 	var planModel = (plan).(model.RepositoryNpmGroupModel)
// 	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
// 	return planModel
// }

// func (f *NpmRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
// 	stateModel := (state).(model.RepositoryNpmGroupModel)
// 	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupDeployRepository))
// 	return stateModel
// }

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getDockerSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"docker": schema.SingleNestedAttribute{
			Description: "Docker specific configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes: map[string]schema.Attribute{
				"force_basic_auth": schema.BoolAttribute{
					Description: "Whether to force authentication (Docker Bearer Token Realm required if false)",
					Required:    true,
				},
				"http_port": schema.Int32Attribute{
					Description: "Create an HTTP connector at specified port",
					Optional:    true,
				},
				"https_port": schema.Int32Attribute{
					Description: "Create an HTTPS connector at specified port",
					Optional:    true,
				},
				"subdomain": schema.StringAttribute{
					Description: "Allows to use subdomain",
					Optional:    true,
				},
				"v1_enabled": schema.BoolAttribute{
					Description: "Whether to allow clients to use the V1 API to interact with this repository",
					Required:    true,
				},
			},
		},
	}
}
