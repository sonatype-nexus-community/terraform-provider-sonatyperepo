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

package model

import (
	"terraform-provider-sonatyperepo/internal/provider/common"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Task Blobstore Compact
// ----------------------------------------
type TaskLicenseExpirationNotificationModel struct {
	BaseTaskModel
}

func (m *TaskLicenseExpirationNotificationModel) ToApiCreateModel() *v3.TaskTemplateXO {
	api := m.toApiCreateModel()
	api.Type = common.TASK_TYPE_LICENSE_EXPIRATION_NOTIFICATION.String()
	return api
}

func (m *TaskLicenseExpirationNotificationModel) ToApiUpdateModel() *v3.UpdateTaskRequest {
	return m.toApiUpdateModel()
}
