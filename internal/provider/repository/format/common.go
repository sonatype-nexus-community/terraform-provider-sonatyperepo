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
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

type RepositoryFormat interface {
	DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoDeleteRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error)
	GetApiCreateSuccessResposneCodes() []int
	GetFormatSchemaAttributes() map[string]schema.Attribute
	GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics)
	GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics)
	GetResourceName(req resource.MetadataRequest) string
	GetKey() string
	UpdatePlanForState(plan any) any
	UpdateStateFromApi(state any, api any) any
}

var RepositoryFormats map[string]RepositoryFormat = map[string]RepositoryFormat{
	common.REPO_FORMAT_NPM: &NpmRepositoryFormat{},
}

func getResourceName(req resource.MetadataRequest, format string) string {
	return fmt.Sprintf("%s_repository_%s_hosted", req.ProviderTypeName, strings.ToLower(format))
}
