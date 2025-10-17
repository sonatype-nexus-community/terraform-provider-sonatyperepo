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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type TaskType string

const (
	TASK_TYPE_BLOBSTORE_COMPACT                                 TaskType = "blobstore.compact"
	TASK_TYPE_BLOBSTORE_DELETE_TEMP_FILES                       TaskType = "blobstore.delete-temp-files"
	TASK_TYPE_BLOBSTORE_EXECUTERECONCILIATIONPLAN               TaskType = "blobstore.executeReconciliationPlan"
	TASK_TYPE_BLOBSTORE_GCS_SOFT_DELETED_BLOBS_MIGRATION        TaskType = "blobstore.gcs.soft-deleted-blobs.migration"
	TASK_TYPE_BLOBSTORE_GROUP_MEMBERREMOVAL                     TaskType = "blobstore.group.memberRemoval"
	TASK_TYPE_BLOBSTORE_METRICS_RECONCILE                       TaskType = "blobstore.metrics.reconcile"
	TASK_TYPE_BLOBSTORE_PLANRECONCILIATION                      TaskType = "blobstore.planReconciliation"
	TASK_TYPE_CREATE_BROWSE_NODES                               TaskType = "create.browse.nodes"
	TASK_TYPE_H2_BACKUP_TASK                                    TaskType = "h2.backup.task"
	TASK_TYPE_MALWARE_REMEDIATOR                                TaskType = "malware.remediator"
	TASK_TYPE_REPOSITORY_APT_REBUILD_METADATA                   TaskType = "repository.apt.rebuild.metadata"
	TASK_TYPE_REPOSITORY_DOCKER_GC                              TaskType = "repository.docker.gc"
	TASK_TYPE_REPOSITORY_DOCKER_UPLOAD_PURGE                    TaskType = "repository.docker.upload-purge"
	TASK_TYPE_REPOSITORY_EXPORT                                 TaskType = "repository.export"
	TASK_TYPE_REPOSITORY_HELM_REBUILD_METADATA                  TaskType = "repository.helm.rebuild.metadata"
	TASK_TYPE_REPOSITORY_IMPORT                                 TaskType = "repository.import"
	TASK_TYPE_REPOSITORY_MAVEN_PUBLISH_DOTINDEX                 TaskType = "repository.maven.publish-dotindex"
	TASK_TYPE_REPOSITORY_MAVEN_PURGE_UNUSED_SNAPSHOTS           TaskType = "repository.maven.purge-unused-snapshots"
	TASK_TYPE_REPOSITORY_MAVEN_REBUILD_METADATA                 TaskType = "repository.maven.rebuild-metadata"
	TASK_TYPE_REPOSITORY_MAVEN_REMOVE_SNAPSHOTS                 TaskType = "repository.maven.remove-snapshots"
	TASK_TYPE_REPOSITORY_MAVEN_REPAIR_BASE_VERSION              TaskType = "repository.maven.repair-base-version"
	TASK_TYPE_REPOSITORY_MAVEN_UNPUBLISH_DOTINDEX               TaskType = "repository.maven.unpublish-dotindex"
	TASK_TYPE_REPOSITORY_MOVE                                   TaskType = "repository.move"
	TASK_TYPE_REPOSITORY_NPM_REBUILD_METADATA                   TaskType = "repository.npm.rebuild-metadata"
	TASK_TYPE_REPOSITORY_NPM_REINDEX                            TaskType = "repository.npm.reindex"
	TASK_TYPE_REPOSITORY_PURGE_UNUSED                           TaskType = "repository.purge-unused"
	TASK_TYPE_REPOSITORY_PYPI_GENERATE_MISSING_SHA256_CHECKSUMS TaskType = "repository.pypi.generate-missing-sha256-checksums"
	TASK_TYPE_REPOSITORY_PYPI_REBUILD_METADATA                  TaskType = "repository.pypi.rebuild-metadata"
	TASK_TYPE_REPOSITORY_REBUILD_INDEX                          TaskType = "repository.rebuild-index"
	TASK_TYPE_REPOSITORY_RUBY_REBUILD_VERSIONS                  TaskType = "repository.ruby.rebuild.versions"
	TASK_TYPE_REPOSITORY_YUM_REBUILD_METADATA                   TaskType = "repository.yum.rebuild.metadata"
	TASK_TYPE_S3_COMPACT_TASK_SCHEDULING_MIGRATION              TaskType = "s3.compact.task.scheduling.migration"
	TASK_TYPE_SECURITY_PURGE_API_KEYS                           TaskType = "security.purge-api-keys"
	TASK_TYPE_TAGS_CLEANUP                                      TaskType = "tags.cleanup"
)

func (tt TaskType) String() string {
	return string(tt)
}

// BaseTaskType that all task types build from
// --------------------------------------------
type BaseTaskType struct {
	taskType TaskType
}

func (f *BaseTaskType) GetApiCreateSuccessResponseCodes() []int {
	return []int{http.StatusCreated}
}

func (f *BaseTaskType) GetKey() string {
	return f.taskType.String()
}

func (f *BaseTaskType) GetResourceName() string {
	return fmt.Sprintf("task_%s", strings.ReplaceAll(strings.ToLower(f.GetKey()), ".", "_"))
}

func (f *BaseTaskType) GetType() TaskType {
	return f.taskType
}

// TaskTypeI that all Repository Formats must implement
// --------------------------------------------
type TaskTypeI interface {
	DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoUpdateRequest(state any, plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	GetApiCreateSuccessResponseCodes() []int
	GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics)
	GetPropertiesSchema() map[string]schema.Attribute
	GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics)
	GetResourceName() string
	GetKey() string
	GetType() TaskType
	UpdatePlanForState(plan any) any
	UpdateStateFromApi(state any, api any) any
	// ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string
}
