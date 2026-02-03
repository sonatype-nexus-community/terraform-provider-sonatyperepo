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

type ComposerRepositoryFormat struct {
	BaseRepositoryFormat
}

type ComposerRepositoryFormatProxy struct {
	ComposerRepositoryFormat
}

// --------------------------------------------
// Generic Composer Format Functions
// --------------------------------------------
func (f *ComposerRepositoryFormat) Key() string {
	return common.REPO_FORMAT_COMPOSER
}

func (f *ComposerRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// PROXY Composer Format Functions
// --------------------------------------------
func (f *ComposerRepositoryFormat) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryComposerProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateComposerProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *ComposerRepositoryFormat) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryComposerProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetComposerProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *ComposerRepositoryFormat) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryComposerProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryComposerProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateComposerProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *ComposerRepositoryFormat) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
}

func (f *ComposerRepositoryFormat) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryComposerProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *ComposerRepositoryFormat) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryComposerProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *ComposerRepositoryFormat) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryComposerProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *ComposerRepositoryFormat) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryComposerProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}
