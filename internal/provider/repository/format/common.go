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

// Error message constants for repository validation during import
const (
	errRepositoryFormatNil      = "repository format is nil, expected '%s'"
	errRepositoryFormatMismatch = "repository format is '%s', expected '%s'"
	errRepositoryTypeNil        = "repository type is nil, expected '%s'"
	errRepositoryTypeMismatch   = "repository type is '%s', expected '%s'"
)

// BaseRepositoryFormat that all formats build from
// --------------------------------------------
type BaseRepositoryFormat struct{}

func (f *BaseRepositoryFormat) DoDeleteRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Call API to Delete
	return apiClient.RepositoryManagementAPI.DeleteRepository(ctx, repositoryName).Execute()
}

func (f *BaseRepositoryFormat) GetApiCreateSuccessResponseCodes() []int {
	return []int{http.StatusCreated}
}

func (f *BaseRepositoryFormat) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	return nil
}

// DoImportRequest provides a base implementation for repository import
// This can be overridden by specific formats if needed
func (f *BaseRepositoryFormat) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// For base implementation, we can't determine the specific repository type
	// This should be overridden by each format implementation
	return nil, nil, fmt.Errorf("import not implemented for this repository format")
}

// ValidateRepositoryForImport validates that the repository matches the expected format and type
func (f *BaseRepositoryFormat) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// This is a base implementation that can be overridden by specific formats
	// Each format should implement its own validation logic
	return nil
}

// RepositoryFormat that all Repository Formats must implement
// --------------------------------------------
type RepositoryFormat interface {
	DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoDeleteRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error)
	DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error)
	DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error)
	ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error
	GetApiCreateSuccessResponseCodes() []int
	GetFormatSchemaAttributes() map[string]schema.Attribute
	GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics)
	GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics)
	GetResourceName(repoType RepositoryType) string
	GetKey() string
	UpdatePlanForState(plan any) any
	UpdateStateFromApi(state any, api any) any
	ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string
}

func getResourceName(format string, repoType RepositoryType) string {
	return fmt.Sprintf("repository_%s_%s", strings.ToLower(format), repoType.String())
}
