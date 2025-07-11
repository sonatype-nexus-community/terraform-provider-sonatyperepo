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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type WildcardPrivilegeType struct {
	BasePrivilegeType
}

func (pt *WildcardPrivilegeType) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.PrivilegeWildcardModel)

	// Call API to Create
	return apiClient.SecurityManagementPrivilegesAPI.CreateWildcardPrivilege(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (pt *WildcardPrivilegeType) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.PrivilegeWildcardModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.SecurityManagementPrivilegesAPI.GetPrivilege(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (pt *WildcardPrivilegeType) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.PrivilegeWildcardModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.PrivilegeWildcardModel)

	// Call API to Create
	return apiClient.SecurityManagementPrivilegesAPI.UpdateWildcardPrivilege(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiCreateModel()).Execute()
}

func (pt *WildcardPrivilegeType) GetPrivilegeTypeSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"pattern": schema.StringAttribute{
			Description: "A colon separated list of parts that create a permission string.",
			Required:    true,
			Optional:    false,
		},
	}
}

func (pt *WildcardPrivilegeType) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.PrivilegeWildcardModel
	return planModel, plan.Get(ctx, &planModel)
}

func (pt *WildcardPrivilegeType) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.PrivilegeWildcardModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (pt *WildcardPrivilegeType) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.PrivilegeWildcardModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (pt *WildcardPrivilegeType) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.PrivilegeWildcardModel)
	stateModel.FromApiModel((api).(sonatyperepo.ApiPrivilegeRequest))
	return stateModel
}
