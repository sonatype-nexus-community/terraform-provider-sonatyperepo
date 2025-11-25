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

type CocoaPodsRepositoryFormat struct {
	BaseRepositoryFormat
}

type CocoaPodsRepositoryFormatProxy struct {
	CocoaPodsRepositoryFormat
}

// --------------------------------------------
// Generic CocoaPods Format Functions
// --------------------------------------------
func (f *CocoaPodsRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_COCOAPODS
}

func (f *CocoaPodsRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// PROXY CocoaPods Format Functions
// --------------------------------------------
func (f *CocoaPodsRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryCocoaPodsProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateCocoapodsProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *CocoaPodsRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryCocoaPodsProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetCocoapodsProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *CocoaPodsRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryCocoaPodsProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryCocoaPodsProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateCocoapodsProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *CocoaPodsRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]tfschema.Attribute {
	return getCommonProxySchemaAttributes()
}

func (f *CocoaPodsRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryCocoaPodsProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *CocoaPodsRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryCocoaPodsProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *CocoaPodsRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryCocoaPodsProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *CocoaPodsRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryCocoaPodsProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}
