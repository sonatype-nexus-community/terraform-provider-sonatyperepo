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

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// BaseTaskType that all task types build from
// --------------------------------------------
type BaseTaskType struct {
	publicName string
	taskType   common.TaskType
}

func (f *BaseTaskType) ApiCreateSuccessResponseCodes() []int {
	return []int{http.StatusCreated}
}

func (f *BaseTaskType) Key() string {
	return f.taskType.String()
}

func (f *BaseTaskType) MarkdownDescription() string {
	return fmt.Sprintf("Manage Task '%s' (%s)", f.PublicName(), f.Type().String())
}

func (f *BaseTaskType) PublicName() string {
	return f.publicName
}

func (f *BaseTaskType) ResourceName() string {
	return fmt.Sprintf("task_%s", strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(f.Key()), ".", "_"), "-", "_"))
}

func (f *BaseTaskType) Type() common.TaskType {
	return f.taskType
}

// TaskTypeI that all Repository Formats must implement
// --------------------------------------------
type TaskTypeI interface {
	DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CreateTask201Response, *http.Response, error)
	DoUpdateRequest(plan any, state any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error)
	ApiCreateSuccessResponseCodes() []int
	MarkdownDescription() string
	PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics)
	PropertiesSchema() map[string]tfschema.Attribute
	PublicName() string
	StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics)
	ResourceName() string
	Key() string
	Type() common.TaskType
	UpdatePlanForState(plan any) any
	UpdateStateFromApi(state any, api any) any
	UpdateStateFromPlanForUpdate(plan any, state any) any
}
