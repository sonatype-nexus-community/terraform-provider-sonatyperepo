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
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

type TerraformRepositoryFormat struct {
	BaseRepositoryFormat
}

type TerraformRepositoryFormatProxy struct {
	TerraformRepositoryFormat
}

// --------------------------------------------
// Generic Terraform Format Functions
// --------------------------------------------
func (f *TerraformRepositoryFormatProxy) Key() string {
	return common.REPO_FORMAT_TERRAFORM
}

func (f *TerraformRepositoryFormatProxy) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// PROXY Terraform Format Functions
// --------------------------------------------
func (f *TerraformRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryTerraformProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateTerraformProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *TerraformRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryTerraformProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetTerraformProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *TerraformRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryTerraformProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryTerraformProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateTerraformProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

// DoImportRequest implements the import functionality for HuggingFace Proxy repositories
func (f *TerraformRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetTerraformProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *TerraformRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
	maps.Copy(additionalAttributes, terraformSchemaAttributes())
	return additionalAttributes
}

func (f *TerraformRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryTerraformProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *TerraformRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryTerraformProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *TerraformRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryTerraformProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *TerraformRepositoryFormatProxy) UpdateStateFromApi(state, api any) any {
	var stateModel model.RepositoryTerraformProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryTerraformProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.TerraformProxyApiRepository))
	return stateModel
}

func (f *TerraformRepositoryFormatProxy) SupportsRepositoryFirewall() bool {
	return false
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func terraformSchemaAttributes() map[string]tfschema.Attribute {
	terraformAttrs := map[string]tfschema.Attribute{
		"terraform": schema.ResourceRequiredSingleNestedAttribute(
			"Terraform specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"require_authentication": schema.ResourceOptionalBoolWithDefault(
					"Indicates if this repository requires authentication overriding anonymous access.",
					false,
				),
			},
		),
	}

	return terraformAttrs
}
