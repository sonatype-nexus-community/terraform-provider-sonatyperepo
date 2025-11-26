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

package tasktype

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

// --------------------------------------------
// Repair: Rebuild Repository Browse Nodes
// --------------------------------------------
type BaseRepairTask struct {
	BaseTaskType
}

func (f *BaseRepairTask) ResourceName() string {
	return fmt.Sprintf("task_repair_%s", strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(f.Key()), ".", "_"), "-", "_"))
}

// --------------------------------------------
// Repair: Rebuild Repository Browse Nodes
// --------------------------------------------
type RepairRebuildBrowseNodesTask struct {
	BaseRepairTask
}

func NewRepairRebuildBrowseNodesTask() *RepairRebuildBrowseNodesTask {
	return &RepairRebuildBrowseNodesTask{
		BaseRepairTask: BaseRepairTask{
			BaseTaskType: BaseTaskType{
				publicName: "Repair - Rebuild repository browse",
				taskType:   common.TASK_TYPE_CREATE_BROWSE_NODES,
			},
		},
	}
}

// --------------------------------------------
// Repair: Rebuild Repository Browse Nodes Functions
// --------------------------------------------
func (f *RepairRebuildBrowseNodesTask) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CreateTask201Response, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskRepairCreateBrowseNodesModel)

	// Call API to Create
	return apiClient.TasksAPI.CreateTask(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *RepairRebuildBrowseNodesTask) DoUpdateRequest(plan any, state any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskRepairCreateBrowseNodesModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.TaskRepairCreateBrowseNodesModel)

	// Call API to Update
	return apiClient.TasksAPI.UpdateTask(ctx, stateModel.Id.ValueString()).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *RepairRebuildBrowseNodesTask) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.TaskRepairCreateBrowseNodesModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RepairRebuildBrowseNodesTask) GetPropertiesSchema() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"repository_name": schema.ResourceRequiredString("The Repository to rebuild browse trees for."),
	}
}

func (f *RepairRebuildBrowseNodesTask) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.TaskRepairCreateBrowseNodesModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RepairRebuildBrowseNodesTask) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.TaskRepairCreateBrowseNodesModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RepairRebuildBrowseNodesTask) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.TaskRepairCreateBrowseNodesModel)
	apiModel := (api).(v3.CreateTask201Response)
	stateModel.Id = types.StringValue(apiModel.Id)
	return stateModel
}

func (f *RepairRebuildBrowseNodesTask) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.TaskRepairCreateBrowseNodesModel)
	stateModel := (state).(model.TaskRepairCreateBrowseNodesModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
