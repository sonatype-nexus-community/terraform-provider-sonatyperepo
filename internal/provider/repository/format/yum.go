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
	"maps"
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

type YumRepositoryFormat struct {
	BaseRepositoryFormat
}

type YumRepositoryFormatHosted struct {
	YumRepositoryFormat
}

type YumRepositoryFormatProxy struct {
	YumRepositoryFormat
}

type YumRepositoryFormatGroup struct {
	YumRepositoryFormat
}

// --------------------------------------------
// Generic YUM Format Functions
// --------------------------------------------
func (f *YumRepositoryFormat) Key() string {
	return common.REPO_FORMAT_YUM
}

func (f *YumRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// Hosted YUM Format Functions
// --------------------------------------------
func (f *YumRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryYumHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateYumHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *YumRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryYumHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetYumHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *YumRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryYumHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryYumHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateYumHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *YumRepositoryFormatHosted) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, yumSchemaAttributes(true))
	return additionalAttributes
}

func (f *YumRepositoryFormatHosted) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryYumHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *YumRepositoryFormatHosted) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryYumHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *YumRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryYumHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *YumRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryYumHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryYumHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.YumHostedApiRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for YUM Hosted repositories
func (f *YumRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetYumHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// PROXY YUM Format Functions
// --------------------------------------------
func (f *YumRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryYumProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateYumProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *YumRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryYumProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetYumProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *YumRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryYumProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryYumProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateYumProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *YumRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonProxySchemaAttributes()
	maps.Copy(additionalAttributes, yumSchemaAttributes(false))
	return additionalAttributes
}

func (f *YumRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryYumProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *YumRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryYumProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *YumRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryYumProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *YumRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryYumProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryYumProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for YUM Proxy repositories
func (f *YumRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetYumProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// GORUP YUM Format Functions
// --------------------------------------------
func (f *YumRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryYumGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateYumGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *YumRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryYumGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetYumGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *YumRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryYumGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryYumGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateYumGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *YumRepositoryFormatGroup) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonGroupSchemaAttributes(false)
	maps.Copy(additionalAttributes, yumSchemaAttributes(false))
	return additionalAttributes
}

func (f *YumRepositoryFormatGroup) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryYumGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *YumRepositoryFormatGroup) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryYumGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *YumRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryYumGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *YumRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryYumGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryYumGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for YUM Group repositories
func (f *YumRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetYumGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func yumSchemaAttributes(isHosted bool) map[string]tfschema.Attribute {
	if isHosted {
		return map[string]tfschema.Attribute{
			"yum": schema.ResourceRequiredSingleNestedAttribute(
				"YUM specific configuration for this Repository",
				map[string]tfschema.Attribute{
					"deploy_policy": schema.ResourceOptionalStringEnum(
						"Validate that all paths are RPMs or yum metadata.",
						common.DEPLOY_POLICY_PERMISSIVE,
						common.DEPLOY_POLICY_STRICT,
					),
					"repo_data_depth": schema.ResourceRequiredInt32WithRange(
						"Specifies the repository depth where repodata folder(s) are created",
						0,
						5,
					),
				},
			),
		}
	} else {
		return map[string]tfschema.Attribute{
			"yum": schema.ResourceOptionalSingleNestedAttribute(
				"YUM specific configuration for this Repository",
				map[string]tfschema.Attribute{
					"key_pair": schema.ResourceOptionalStringWithPlanModifier(
						"PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)",
						stringplanmodifier.UseStateForUnknown(),
					),
					"passphrase": schema.ResourceSensitiveOptionalStringWithPlanModifier(
						"Passphrase to access PGP signing key",
						stringplanmodifier.UseStateForUnknown(),
					),
				},
			),
		}
	}
}
