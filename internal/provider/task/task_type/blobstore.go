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
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type BlobstoreCompactTask struct {
	BaseTaskType
}

func NewBlobstoreCompactTask() *BlobstoreCompactTask {
	return &BlobstoreCompactTask{
		BaseTaskType: BaseTaskType{
			taskType: TASK_TYPE_BLOBSTORE_COMPACT,
		},
	}
}

// --------------------------------------------
// Blobstore Compact Format Functions
// --------------------------------------------
func (f *BlobstoreCompactTask) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskBlobstoreCompactModel)

	// Call API to Create
	return apiClient.TasksAPI.CreateTask(ctx).Body(*planModel.ToApiCreateModel()).Execute()
}

func (f *BlobstoreCompactTask) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	// 	planModel := (plan).(model.RepositoryAptHostedModel)

	// 	// Cast to correct State Model Type
	// 	stateModel := (state).(model.RepositoryAptHostedModel)

	// // Call API to Create
	// return apiClient.RepositoryManagementAPI.UpdateAptHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
	return nil, nil
}

func (f *BlobstoreCompactTask) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAptHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *BlobstoreCompactTask) GetPropertiesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"blob_store_name": schema.StringAttribute{
			Description: "The blob store to compact",
			Required:    true,
			Optional:    false,
		},
		"blobs_older_than": schema.Int32Attribute{
			Description: "The number of days a blob should kept before permanent deletion (default 0)",
			Optional:    true,
			Computed:    true,
			Default:     int32default.StaticInt32(0),
		},
	}
}

func (f *BlobstoreCompactTask) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAptHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *BlobstoreCompactTask) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAptHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *BlobstoreCompactTask) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryAptHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.AptHostedApiRepository))
	return stateModel
}
