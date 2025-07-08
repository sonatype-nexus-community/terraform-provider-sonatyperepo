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

package format

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

type RepositoryType int64

const (
	REPO_TYPE_HOSTED RepositoryType = iota
	REPO_TYPE_GROUP
	REPO_TYPE_PROXY
)

func (rt RepositoryType) String() string {
	switch rt {
	case REPO_TYPE_HOSTED:
		return "hosted"
	case REPO_TYPE_GROUP:
		return "group"
	case REPO_TYPE_PROXY:
		return "proxy"
	}
	return "unknown"
}

// BaseRepositoryFormat that all formats build from
// --------------------------------------------
type BaseRepositoryFormat struct{}

func (f *BaseRepositoryFormat) DoDeleteRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Call API to Delete
	return apiClient.RepositoryManagementAPI.DeleteRepository(ctx, repositoryName).Execute()
}

func (f *BaseRepositoryFormat) GetApiCreateSuccessResposneCodes() []int {
	return []int{http.StatusCreated}
}

// RepositoryFormat that all Repository Formats must implement
// --------------------------------------------
type RepositoryFormat interface {
	DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoDeleteRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error)
	GetApiCreateSuccessResposneCodes() []int
	GetFormatSchemaAttributes() map[string]schema.Attribute
	GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics)
	GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics)
	GetResourceName(repoType RepositoryType) string
	GetKey() string
	UpdatePlanForState(plan any) any
	UpdateStateFromApi(state any, api any) any
}

// var RepositoryFormats map[string]RepositoryFormat = map[string]RepositoryFormat{
// 	common.REPO_FORMAT_NPM: &NpmRepositoryFormat{},
// }

func getResourceName(format string, repoType RepositoryType) string {
	return fmt.Sprintf("repository_%s_%s", strings.ToLower(format), repoType.String())
}
