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

type HuggingFaceRepositoryFormat struct {
	BaseRepositoryFormat
}

type HuggingFaceRepositoryFormatProxy struct {
	HuggingFaceRepositoryFormat
}

// --------------------------------------------
// Generic HuggingFace Format Functions
// --------------------------------------------
func (f *HuggingFaceRepositoryFormat) Key() string {
	return common.REPO_FORMAT_HUGGING_FACE
}

func (f *HuggingFaceRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// PROXY HuggingFace Format Functions
// --------------------------------------------
func (f *HuggingFaceRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryHuggingFaceProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateHuggingfaceProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *HuggingFaceRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryHuggingFaceProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetHuggingfaceProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *HuggingFaceRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryHuggingFaceProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryHuggingFaceProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateHuggingfaceProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *HuggingFaceRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonProxySchemaAttributes()
}

func (f *HuggingFaceRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryHuggingFaceProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *HuggingFaceRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryHuggingFaceProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *HuggingFaceRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryHuggingFaceProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *HuggingFaceRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryHuggingFaceProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryHuggingFaceProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for HuggingFace Proxy repositories
func (f *HuggingFaceRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetHuggingfaceProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}
