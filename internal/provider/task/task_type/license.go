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
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type LicenseExpirationNotificationTask struct {
	BaseTaskType
}

func NewLicenseExpirationNotificationTask() *LicenseExpirationNotificationTask {
	return &LicenseExpirationNotificationTask{
		BaseTaskType: BaseTaskType{
			publicName: "License - Check license expiration and send notifications",
			taskType:   common.TASK_TYPE_LICENSE_EXPIRATION_NOTIFICATION,
		},
	}
}

// --------------------------------------------
// Blobstore Compact Format Functions
// --------------------------------------------
func (f *LicenseExpirationNotificationTask) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.TaskXO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskLicenseExpirationNotificationModel)

	// Call API to Create
	return apiClient.TasksAPI.CreateTask(ctx).Body(*planModel.ToApiCreateModel()).Execute()
}

func (f *LicenseExpirationNotificationTask) DoUpdateRequest(plan any, state any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.TaskLicenseExpirationNotificationModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.TaskLicenseExpirationNotificationModel)

	// Call API to Update
	return apiClient.TasksAPI.UpdateTask(ctx, stateModel.Id.ValueString()).Body(*planModel.ToApiUpdateModel()).Execute()
}

func (f *LicenseExpirationNotificationTask) MarkdownDescription() string {
	return fmt.Sprintf(
		`Manage Task '%s' (%s)
	
This task requires the Sonatype Nexus Repository >= 3.86.0 - see [official documentation](https://help.sonatype.com/en/sonatype-nexus-repository-3-86-0-release-notes.html#license-expiry-notification-and-status-check).`,
		f.PublicName(), f.Type().String(),
	)
}

func (f *LicenseExpirationNotificationTask) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.TaskLicenseExpirationNotificationModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *LicenseExpirationNotificationTask) PropertiesSchema() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{}
}

func (f *LicenseExpirationNotificationTask) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.TaskLicenseExpirationNotificationModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *LicenseExpirationNotificationTask) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.TaskLicenseExpirationNotificationModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *LicenseExpirationNotificationTask) UpdateStateFromApi(state, api any) any {
	stateModel := (state).(model.TaskLicenseExpirationNotificationModel)
	apiModel := (api).(v3.TaskXO)
	stateModel.Id = types.StringPointerValue(apiModel.Id)
	return stateModel
}

func (f *LicenseExpirationNotificationTask) UpdateStateFromPlanForUpdate(plan, state any) any {
	planModel := (plan).(model.TaskLicenseExpirationNotificationModel)
	stateModel := (state).(model.TaskLicenseExpirationNotificationModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
