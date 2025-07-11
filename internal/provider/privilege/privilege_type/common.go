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
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// PrivilegeTypeType
// --------------------------------------------
type PrivilegeTypeType int8

const (
	TypeApplication PrivilegeTypeType = iota
	TypeRepositoryAdmin
	TypeRepositoryContentSelector
	TypeRepositoryView
	TypeScript
	TypeWildcard
)

func (pt PrivilegeTypeType) String() string {
	switch pt {
	case TypeApplication:
		return "application"
	case TypeRepositoryAdmin:
		return "repository-admin"
	case TypeRepositoryContentSelector:
		return "repository-content-selector"
	case TypeRepositoryView:
		return "repository-view"
	case TypeScript:
		return "script"
	case TypeWildcard:
		return "wildcard"
	}

	return "unknown"
}

func AllPrivilegeTypes() []string {
	return []string{
		TypeApplication.String(),
		TypeRepositoryAdmin.String(),
		TypeRepositoryContentSelector.String(),
		TypeRepositoryView.String(),
		TypeScript.String(),
		TypeWildcard.String(),
	}
}

// PrivilegeAction
// --------------------------------------------
type PrivilegeAction int8

const (
	ActionBrowse PrivilegeAction = iota
	ActionRead
	ActionEdit
	ActionAdd
	ActionDelete
	ActionAll
	ActionRun
)

func (pt PrivilegeAction) String() string {
	switch pt {
	case ActionAdd:
		return "ADD"
	case ActionAll:
		return "ALL"
	case ActionBrowse:
		return "BROWSE"
	case ActionDelete:
		return "DELETE"
	case ActionEdit:
		return "EDIT"
	case ActionRead:
		return "READ"
	case ActionRun:
		return "RUN"
	}

	return "unknown"
}

func AllActionsExceptRun() []string {
	return []string{
		ActionAdd.String(),
		ActionAll.String(),
		ActionBrowse.String(),
		ActionDelete.String(),
		ActionEdit.String(),
		ActionRead.String(),
	}
}

func BreadActions() []string {
	return []string{
		ActionAdd.String(),
		ActionBrowse.String(),
		ActionDelete.String(),
		ActionEdit.String(),
		ActionRead.String(),
	}
}

func BreadAndRunActions() []string {
	actions := BreadActions()
	actions = append(actions, ActionRun.String())
	return actions
}

// BasePrivilegeType that all Privilege Types build from
// --------------------------------------------
type BasePrivilegeType struct{}

func (f *BasePrivilegeType) DoDeleteRequest(privilegeName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Call API to Delete
	return apiClient.SecurityManagementPrivilegesAPI.DeletePrivilege(ctx, privilegeName).Execute()
}

func (f *BasePrivilegeType) GetApiCreateSuccessResponseCodes() []int {
	return []int{http.StatusCreated}
}

func (f *BasePrivilegeType) GetResourceName(privType PrivilegeTypeType) string {
	return fmt.Sprintf("privilege_%s", strings.ReplaceAll(privType.String(), "-", "_"))
}

func (f *BasePrivilegeType) IsDeprecated() bool {
	return false
}

// PrivilegeType that all Privilege Types must implement
// --------------------------------------------
type PrivilegeType interface {
	DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoDeleteRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error)
	IsDeprecated() bool
	GetApiCreateSuccessResponseCodes() []int
	GetPrivilegeTypeSchemaAttributes() map[string]schema.Attribute
	GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics)
	GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics)
	GetResourceName(privType PrivilegeTypeType) string
	UpdatePlanForState(plan any) any
	UpdateStateFromApi(state any, api any) any
}
