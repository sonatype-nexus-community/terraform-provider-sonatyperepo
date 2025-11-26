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

package privilege_type

import (
	"context"
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type RepositoryViewPrivilegeType struct {
	BasePrivilegeType
}

func (pt *RepositoryViewPrivilegeType) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.PrivilegeRepositoryViewModel)

	// Call API to Create
	return apiClient.SecurityManagementPrivilegesAPI.CreateRepositoryViewPrivilege(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (pt *RepositoryViewPrivilegeType) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.PrivilegeRepositoryViewModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.SecurityManagementPrivilegesAPI.GetPrivilege(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (pt *RepositoryViewPrivilegeType) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.PrivilegeRepositoryViewModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.PrivilegeRepositoryViewModel)

	// Call API to Create
	return apiClient.SecurityManagementPrivilegesAPI.UpdateRepositoryViewPrivilege(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiCreateModel()).Execute()
}

func (pt *RepositoryViewPrivilegeType) PrivilegeTypeSchemaAttributes() map[string]tfschema.Attribute {
	return schemaAttributesActionFormatRepository()
}

func (pt *RepositoryViewPrivilegeType) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.PrivilegeRepositoryViewModel
	return planModel, plan.Get(ctx, &planModel)
}

func (pt *RepositoryViewPrivilegeType) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.PrivilegeRepositoryViewModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (pt *RepositoryViewPrivilegeType) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.PrivilegeRepositoryViewModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (pt *RepositoryViewPrivilegeType) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.PrivilegeRepositoryViewModel)
	stateModel.FromApiModel((api).(sonatyperepo.ApiPrivilegeRequest))
	return stateModel
}
