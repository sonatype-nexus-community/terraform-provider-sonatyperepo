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

	"github.com/hashicorp/terraform-plugin-framework/types"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Properties for repository.docker.gc
// ----------------------------------------
type TaskPropertiesRepositoryDockerGc struct {
	DeployOffset   types.Int32  `tfsdk:"deploy_offset" nxrm:"deployOffset"`
	RepositoryName types.String `tfsdk:"repository_name" nxrm:"repositoryName"`
}

func (p *TaskPropertiesRepositoryDockerGc) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Task Repositor Docker GC
// ----------------------------------------
type TaskRepositoryDockerGcModel struct {
	BaseTaskModel
	Properties TaskPropertiesRepositoryDockerGc `tfsdk:"properties"`
}

func (m *TaskRepositoryDockerGcModel) ToApiCreateModel(version common.SystemVersion) *v3.TaskTemplateXO {
	api := m.toApiCreateModel()
	api.Type = common.TASK_TYPE_REPOSITORY_DOCKER_GC.String()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *TaskRepositoryDockerGcModel) ToApiUpdateModel(version common.SystemVersion) *v3.UpdateTaskRequest {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

// Properties for repository.docker.upload-purge
// ----------------------------------------
type TaskPropertiesRepositoryDockerUploadPurge struct {
	Age types.Int32 `tfsdk:"age" nxrm:"age"`
}

func (p *TaskPropertiesRepositoryDockerUploadPurge) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	return StructToMap(p)
}

// Task Repositor Docker Upload Purge
// ----------------------------------------
type TaskRepositoryDockerUploadPurgeModel struct {
	BaseTaskModel
	Properties TaskPropertiesRepositoryDockerUploadPurge `tfsdk:"properties"`
}

func (m *TaskRepositoryDockerUploadPurgeModel) ToApiCreateModel(version common.SystemVersion) *v3.TaskTemplateXO {
	api := m.toApiCreateModel()
	api.Type = common.TASK_TYPE_REPOSITORY_DOCKER_UPLOAD_PURGE.String()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *TaskRepositoryDockerUploadPurgeModel) ToApiUpdateModel(version common.SystemVersion) *v3.UpdateTaskRequest {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

// Properties for repository.maven.remove-snapshots
// ----------------------------------------
type TaskPropertiesRepositoryMavenRemoveSnapshots struct {
	RepositoryName        types.String `tfsdk:"repository_name" nxrm:"repositoryName"`
	MinimumRetained       types.Int32  `tfsdk:"minimum_retained" nxrm:"minimumRetained"`
	SnapshotRetentionDays types.Int32  `tfsdk:"snapshot_retention_days" nxrm:"snapshotRetentionDays"`
	RemoveIfReleased      types.Bool   `tfsdk:"remove_if_released" nxrm:"removeIfReleased"`
	GracePeriodInDays     types.Int32  `tfsdk:"grace_period_in_days" nxrm:"gracePeriodInDays"`
}

func (p *TaskPropertiesRepositoryMavenRemoveSnapshots) GetFilteredPropertiesAsMap(version common.SystemVersion) *map[string]string {
	properties := StructToMap(p)

	if p.GracePeriodInDays.ValueInt32() == 0 {
		delete(*properties, "gracePeriodInDays")
	}

	return properties
}

// Task Repository Maven Remove Snapshots
// ----------------------------------------
type TaskRepositoryMavenRemoveSnapshotsModel struct {
	BaseTaskModel
	Properties TaskPropertiesRepositoryMavenRemoveSnapshots `tfsdk:"properties"`
}

func (m *TaskRepositoryMavenRemoveSnapshotsModel) ToApiCreateModel(version common.SystemVersion) *v3.TaskTemplateXO {
	api := m.toApiCreateModel()
	api.Type = common.TASK_TYPE_REPOSITORY_MAVEN_REMOVE_SNAPSHOTS.String()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}

func (m *TaskRepositoryMavenRemoveSnapshotsModel) ToApiUpdateModel(version common.SystemVersion) *v3.UpdateTaskRequest {
	api := m.toApiUpdateModel()
	api.Properties = m.Properties.GetFilteredPropertiesAsMap(version)
	return api
}
