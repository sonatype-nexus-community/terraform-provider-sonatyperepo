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

type GoRepositoryFormat struct {
	BaseRepositoryFormat
}

type GoRepositoryFormatProxy struct {
	GoRepositoryFormat
}

type GoRepositoryFormatGroup struct {
	GoRepositoryFormat
}

// --------------------------------------------
// Generic NPM Format Functions
// --------------------------------------------
func (f *GoRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_GO
}

func (f *GoRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return resourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// PROXY Go Format Functions
// --------------------------------------------
func (f *GoRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryGoProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateGoProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *GoRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryGoProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetGoProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *GoRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryGoProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryGoProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateGoProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *GoRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonProxySchemaAttributes()
}

func (f *GoRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryGoProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *GoRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryGoProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *GoRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryGoProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *GoRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryGoProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryGoProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Go Proxy repositories
func (f *GoRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetGoProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// GORUP Go Format Functions
// --------------------------------------------
func (f *GoRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryGoGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateGoGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *GoRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryGoGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetGoGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *GoRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryGoGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryGoGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateGoGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *GoRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonGroupSchemaAttributes(false)
}

func (f *GoRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryGoGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *GoRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryGoGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *GoRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryGoGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *GoRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryGoGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryGoGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Go Group repositories
func (f *GoRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetGoGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}
