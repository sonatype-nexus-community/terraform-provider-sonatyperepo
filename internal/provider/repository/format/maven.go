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
	"maps"
	"net/http"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type MavenRepositoryFormat struct {
	BaseRepositoryFormat
}

type MavenRepositoryFormatHosted struct {
	MavenRepositoryFormat
}

type MavenRepositoryFormatProxy struct {
	MavenRepositoryFormat
}

type MavenRepositoryFormatGroup struct {
	MavenRepositoryFormat
}

// --------------------------------------------
// Generic Maven Format Functions
// --------------------------------------------
func (f *MavenRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_MAVEN
}

func (f *MavenRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// HOSTED Maven Format Functions
// --------------------------------------------
func (f *MavenRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryMavenHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateMavenHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *MavenRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryMavenHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetMavenHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *MavenRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryMavenHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryMavenHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateMavenHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *MavenRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, getMavenSchemaAttributes())
	return additionalAttributes
}

func (f *MavenRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryMavenHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *MavenRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryMavenHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *MavenRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryMavenHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *MavenRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryMavenHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryMavenHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.MavenHostedApiRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Maven Hosted repositories
func (f *MavenRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetMavenHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// ValidateRepositoryForImport validates that the imported repository is indeed a Maven Hosted repository
func (f *MavenRepositoryFormatHosted) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// Cast to Maven Hosted API Repository
	apiRepo, ok := repositoryData.(sonatyperepo.MavenHostedApiRepository)
	if !ok {
		return fmt.Errorf("repository data is not a Maven Hosted repository")
	}

	if apiRepo.Format == nil {
		return fmt.Errorf(errRepositoryFormatNil, expectedFormat)
	}
	// Convert both to lowercase for comparison
	// Note: Maven repositories may return "maven2" from the API
	actualFormat := strings.ToLower(*apiRepo.Format)
	expectedFormatLower := strings.ToLower(expectedFormat)
	// Accept both "maven" and "maven2" as valid Maven formats
	if actualFormat != expectedFormatLower && actualFormat != "maven2" && expectedFormatLower != "maven2" {
		return fmt.Errorf(errRepositoryFormatMismatch, *apiRepo.Format, expectedFormat)
	}

	// Validate type
	expectedTypeStr := expectedType.String()
	if apiRepo.Type == nil {
		return fmt.Errorf(errRepositoryTypeNil, expectedTypeStr)
	}
	if *apiRepo.Type != expectedTypeStr {
		return fmt.Errorf(errRepositoryTypeMismatch, *apiRepo.Type, expectedTypeStr)
	}

	return nil
}

// --------------------------------------------
// PROXY Maven Format Functions
// --------------------------------------------
func (f *MavenRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryMavenProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateMavenProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *MavenRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryMavenProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetMavenProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *MavenRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryMavenProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryMavenProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateMavenProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *MavenRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonProxySchemaAttributes()
	maps.Copy(additionalAttributes, getMavenSchemaAttributes())
	return additionalAttributes
}

func (f *MavenRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryMavenProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *MavenRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryMavenProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *MavenRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryMavenProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *MavenRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryMavenProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryMavenProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.MavenProxyApiRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Maven Proxy repositories
func (f *MavenRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetMavenProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// ValidateRepositoryForImport validates that the imported repository is indeed a Maven Proxy repository
func (f *MavenRepositoryFormatProxy) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// Cast to Maven Proxy API Repository
	apiRepo, ok := repositoryData.(sonatyperepo.MavenProxyApiRepository)
	if !ok {
		return fmt.Errorf("repository data is not a Maven Proxy repository")
	}

	// Validate format (case-insensitive)
	if apiRepo.Format == nil {
		return fmt.Errorf(errRepositoryFormatNil, expectedFormat)
	}
	// Convert both to lowercase for comparison
	// Note: Maven repositories may return "maven2" from the API
	actualFormat := strings.ToLower(*apiRepo.Format)
	expectedFormatLower := strings.ToLower(expectedFormat)
	// Accept both "maven" and "maven2" as valid Maven formats
	if actualFormat != expectedFormatLower && actualFormat != "maven2" && expectedFormatLower != "maven2" {
		return fmt.Errorf(errRepositoryFormatMismatch, *apiRepo.Format, expectedFormat)
	}

	// Validate type
	expectedTypeStr := expectedType.String()
	if apiRepo.Type == nil {
		return fmt.Errorf(errRepositoryTypeNil, expectedTypeStr)
	}
	if *apiRepo.Type != expectedTypeStr {
		return fmt.Errorf(errRepositoryTypeMismatch, *apiRepo.Type, expectedTypeStr)
	}

	return nil
}

// --------------------------------------------
// GROUP Maven Format Functions
// --------------------------------------------
func (f *MavenRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryMavenGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateMavenGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *MavenRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryMavenGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetMavenGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *MavenRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryMavenGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryMavenGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateMavenGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *MavenRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return getCommonGroupSchemaAttributes(false)
}

func (f *MavenRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryMavenGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *MavenRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryMavenGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *MavenRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryMavenGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *MavenRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryMavenGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryMavenGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Maven Group repositories
func (f *MavenRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetMavenGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// ValidateRepositoryForImport validates that the imported repository is indeed a Maven Group repository
func (f *MavenRepositoryFormatGroup) ValidateRepositoryForImport(repositoryData any, expectedFormat string, expectedType RepositoryType) error {
	// Cast to Maven Group API Repository
	apiRepo, ok := repositoryData.(sonatyperepo.SimpleApiGroupRepository)
	if !ok {
		return fmt.Errorf("repository data is not a Maven Group repository")
	}

	if apiRepo.Format == nil {
		return fmt.Errorf(errRepositoryFormatNil, expectedFormat)
	}
	// Convert both to lowercase for comparison
	// Note: Maven repositories may return "maven2" from the API
	actualFormat := strings.ToLower(*apiRepo.Format)
	expectedFormatLower := strings.ToLower(expectedFormat)
	// Accept both "maven" and "maven2" as valid Maven formats
	if actualFormat != expectedFormatLower && actualFormat != "maven2" && expectedFormatLower != "maven2" {
		return fmt.Errorf(errRepositoryFormatMismatch, *apiRepo.Format, expectedFormat)
	}

	// Validate type
	expectedTypeStr := expectedType.String()
	if apiRepo.Type == nil {
		return fmt.Errorf(errRepositoryTypeNil, expectedTypeStr)
	}
	if *apiRepo.Type != expectedTypeStr {
		return fmt.Errorf(errRepositoryTypeMismatch, *apiRepo.Type, expectedTypeStr)
	}

	return nil
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getMavenSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"maven": schema.SingleNestedAttribute{
			Description: "Maven specific configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes: map[string]schema.Attribute{
				"version_policy": schema.StringAttribute{
					Description: "What type of artifacts does this repository store?",
					Required:    false,
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							common.MAVEN_VERSION_POLICY_RELEASE,
							common.MAVEN_VERSION_POLICY_SNAPSHOT,
							common.MAVEN_VERSION_POLICY_MIXED,
						),
					},
				},
				"layout_policy": schema.StringAttribute{
					Description: "Validate that all paths are maven artifact or metadata paths",
					Required:    false,
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							common.MAVEN_LAYOUT_STRICT, common.MAVEN_LAYOUT_PERMISSIVE,
						),
					},
				},
				"content_disposition": schema.StringAttribute{
					Description: "Add Content-Disposition header as 'ATTACHMENT' to disable some content from being inline in a browser.",
					Required:    false,
					Optional:    true,
					Validators: []validator.String{
						stringvalidator.OneOf(
							common.MAVEN_CONTENT_DISPOSITION_INLINE,
							common.MAVEN_CONTENT_DISPOSITION_ATTACHMENT,
						),
					},
				},
			},
		},
	}
}