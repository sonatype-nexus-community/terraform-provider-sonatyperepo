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
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// --------------------------------------------
// Docker Repository GC Task
// --------------------------------------------
type RepositoryDockerGcTask struct {
	BaseTaskType
}

func NewRepositoryDockerGcTask() *RepositoryDockerGcTask {
	return &RepositoryDockerGcTask{
		BaseTaskType: BaseTaskType{
			publicName: "Docker - Delete unused manifests and images",
			taskType:   common.TASK_TYPE_REPOSITORY_DOCKER_GC,
		},
	}
}

// --------------------------------------------
// Docker Repository GC Format Functions
// --------------------------------------------
func (f *RepositoryDockerGcTask) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CreateTask201Response, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskRepositoryDockerGcModel)

	// Call API to Create
	return apiClient.TasksAPI.CreateTask(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *RepositoryDockerGcTask) DoUpdateRequest(plan any, state any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskRepositoryDockerGcModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.TaskRepositoryDockerGcModel)

	// Call API to Update
	return apiClient.TasksAPI.UpdateTask(ctx, stateModel.Id.ValueString()).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *RepositoryDockerGcTask) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.TaskRepositoryDockerGcModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RepositoryDockerGcTask) GetPropertiesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"deploy_offset": schema.Int32Attribute{
			MarkdownDescription: `Manifests and images deployed within this period before the task starts will not be deleted.`,
			Optional:            true,
			Computed:            true,
			Default:             int32default.StaticInt32(common.TASK_REPOSITORY_DOCKER_GC_DEFAULT_DEPLOY_OFFSET),
		},
		"repository_name": schema.StringAttribute{
			Description: "The Docker repository to clean up.",
			Required:    true,
			Optional:    false,
		},
	}
}

func (f *RepositoryDockerGcTask) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.TaskRepositoryDockerGcModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RepositoryDockerGcTask) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.TaskRepositoryDockerGcModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RepositoryDockerGcTask) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.TaskRepositoryDockerGcModel)
	apiModel := (api).(v3.CreateTask201Response)
	stateModel.Id = types.StringValue(apiModel.Id)
	return stateModel
}

func (f *RepositoryDockerGcTask) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.TaskRepositoryDockerGcModel)
	stateModel := (state).(model.TaskRepositoryDockerGcModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}

// --------------------------------------------
// Docker Repository Upload Purge
// --------------------------------------------
type RepositoryDockerUploadPurgeTask struct {
	BaseTaskType
}

func NewRepositoryDockerUploadPurgeTaskTask() *RepositoryDockerUploadPurgeTask {
	return &RepositoryDockerUploadPurgeTask{
		BaseTaskType: BaseTaskType{
			publicName: "Docker - Delete incomplete uploads",
			taskType:   common.TASK_TYPE_REPOSITORY_DOCKER_UPLOAD_PURGE,
		},
	}
}

// --------------------------------------------
// Docker Repository GC Format Functions
// --------------------------------------------
func (f *RepositoryDockerUploadPurgeTask) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CreateTask201Response, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskRepositoryDockerUploadPurgeModel)

	// Call API to Create
	return apiClient.TasksAPI.CreateTask(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *RepositoryDockerUploadPurgeTask) DoUpdateRequest(plan any, state any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskRepositoryDockerUploadPurgeModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.TaskRepositoryDockerUploadPurgeModel)

	// Call API to Update
	return apiClient.TasksAPI.UpdateTask(ctx, stateModel.Id.ValueString()).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *RepositoryDockerUploadPurgeTask) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.TaskRepositoryDockerUploadPurgeModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *RepositoryDockerUploadPurgeTask) GetPropertiesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"age": schema.Int32Attribute{
			MarkdownDescription: `Delete incomplete docker uploads that are older than the specified age in hours.`,
			Optional:            true,
			Computed:            true,
			Default:             int32default.StaticInt32(common.TASK_REPOSITORY_DOCKER_UPLOAD_PURGE_DEFAULT_AGE),
		},
	}
}

func (f *RepositoryDockerUploadPurgeTask) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.TaskRepositoryDockerUploadPurgeModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *RepositoryDockerUploadPurgeTask) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.TaskRepositoryDockerUploadPurgeModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *RepositoryDockerUploadPurgeTask) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.TaskRepositoryDockerUploadPurgeModel)
	apiModel := (api).(v3.CreateTask201Response)
	stateModel.Id = types.StringValue(apiModel.Id)
	return stateModel
}

func (f *RepositoryDockerUploadPurgeTask) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.TaskRepositoryDockerUploadPurgeModel)
	stateModel := (state).(model.TaskRepositoryDockerUploadPurgeModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
