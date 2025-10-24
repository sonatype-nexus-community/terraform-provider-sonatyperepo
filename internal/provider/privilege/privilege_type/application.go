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

type ApplicationPrivilegeType struct {
	BasePrivilegeType
}

func (pt *ApplicationPrivilegeType) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.PrivilegeApplicationModel)

	// Call API to Create
	return apiClient.SecurityManagementPrivilegesAPI.CreateApplicationPrivilege(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (pt *ApplicationPrivilegeType) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.PrivilegeApplicationModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.SecurityManagementPrivilegesAPI.GetPrivilege(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (pt *ApplicationPrivilegeType) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.PrivilegeApplicationModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.PrivilegeApplicationModel)

	// Call API to Create
	return apiClient.SecurityManagementPrivilegesAPI.UpdateApplicationPrivilege(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiCreateModel()).Execute()
}

func (pt *ApplicationPrivilegeType) GetPrivilegeTypeSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"domain": schema.StringAttribute{
			Description: "The domain (i.e. 'blobstores', 'capabilities' or even '*' for all) that this privilege is granting access to. Note that creating new privileges with a domain is only necessary when using plugins that define their own domain(s).",
			Required:    true,
			Optional:    false,
		},
		"actions": schema.SetAttribute{
			Description: "A set of actions to associate with the privilege, using BREAD syntax (browse,read,edit,add,delete,all) as well as 'run' for script privileges.",
			Required:    true,
			Optional:    false,
			ElementType: types.StringType,
			Validators: []validator.Set{
				setvalidator.ValueStringsAre([]validator.String{
					stringvalidator.OneOf(AllActionsExceptRun()...),
				}...),
			},
		},
	}
}

func (pt *ApplicationPrivilegeType) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.PrivilegeApplicationModel
	return planModel, plan.Get(ctx, &planModel)
}

func (pt *ApplicationPrivilegeType) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.PrivilegeApplicationModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (pt *ApplicationPrivilegeType) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.PrivilegeApplicationModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (pt *ApplicationPrivilegeType) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.PrivilegeApplicationModel)
	stateModel.FromApiModel((api).(sonatyperepo.ApiPrivilegeRequest))
	return stateModel
}
