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

type CondaRepositoryFormat struct {
	BaseRepositoryFormat
}

type CondaRepositoryFormatProxy struct {
	CondaRepositoryFormat
}

// --------------------------------------------
// Generic Conda Format Functions
// --------------------------------------------
func (f *CondaRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_CONDA
}

func (f *CondaRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// PROXY Conda Format Functions
// --------------------------------------------
func (f *CondaRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryCondaProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateCondaProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *CondaRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryCondaProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetCondaProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *CondaRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryCondaProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryCondaProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateCondaProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *CondaRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return getCommonProxySchemaAttributes()
}

func (f *CondaRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryCondaProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *CondaRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryCondaProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *CondaRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryCondaProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *CondaRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryCondaProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryCondaProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Conda Proxy repositories
func (f *CondaRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetCondaProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}
