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

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type ScriptPrivilegeType struct {
	BasePrivilegeType
}

func (pt *ScriptPrivilegeType) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.PrivilegeScriptModel)

	// Call API to Create
	return apiClient.SecurityManagementPrivilegesAPI.CreateScriptPrivilege(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (pt *ScriptPrivilegeType) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.PrivilegeScriptModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.SecurityManagementPrivilegesAPI.GetPrivilege(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (pt *ScriptPrivilegeType) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.PrivilegeScriptModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.PrivilegeScriptModel)

	// Call API to Create
	return apiClient.SecurityManagementPrivilegesAPI.UpdateScriptPrivilege(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiCreateModel()).Execute()
}

func (pt *ScriptPrivilegeType) GetPrivilegeTypeSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"actions": schema.SetAttribute{
			Description: "A set of actions to associate with the privilege, using BREAD syntax (browse,read,edit,add,delete,all) as well as 'run' for script privileges.",
			Required:    true,
			Optional:    false,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre([]validator.String{
					stringvalidator.OneOf(BreadAndRunActions()...),
				}...),
			},
		},
		"script_name": schema.StringAttribute{
			Description: "The name of a script to give access to.",
			Required:    true,
			Optional:    false,
		},
	}
}

func (pt *ScriptPrivilegeType) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.PrivilegeScriptModel
	return planModel, plan.Get(ctx, &planModel)
}

func (pt *ScriptPrivilegeType) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.PrivilegeScriptModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *ScriptPrivilegeType) IsDeprecated() bool {
	return true
}

func (pt *ScriptPrivilegeType) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.PrivilegeScriptModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (pt *ScriptPrivilegeType) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.PrivilegeScriptModel)
	stateModel.FromApiModel((api).(sonatyperepo.ApiPrivilegeRequest))
	return stateModel
}
