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
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type AnsibleGalaxyRepositoryFormat struct {
	BaseRepositoryFormat
}

type AnsibleGalaxyRepositoryFormatHosted struct {
	AnsibleGalaxyRepositoryFormat
}

type AnsibleGalaxyRepositoryFormatProxy struct {
	AnsibleGalaxyRepositoryFormat
}

type AnsibleGalaxyRepositoryFormatGroup struct {
	AnsibleGalaxyRepositoryFormat
}

// --------------------------------------------
// Generic Ansible Galaxy Format Functions
// --------------------------------------------
func (f *AnsibleGalaxyRepositoryFormat) Key() string {
	return common.REPO_FORMAT_ANSIBLE_GALAXY
}

func (f *AnsibleGalaxyRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// Hosted Ansible Galaxy Format Functions
// --------------------------------------------
func (f *AnsibleGalaxyRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAnsibleGalaxyHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateAnsiblegalaxyHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *AnsibleGalaxyRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAnsibleGalaxyHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAnsiblegalaxyHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *AnsibleGalaxyRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAnsibleGalaxyHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAnsibleGalaxyHostedModel)

	// Call API to Update
	return apiClient.RepositoryManagementAPI.UpdateAnsiblegalaxyHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *AnsibleGalaxyRepositoryFormatHosted) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{}
}

func (f *AnsibleGalaxyRepositoryFormatHosted) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAnsibleGalaxyHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *AnsibleGalaxyRepositoryFormatHosted) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAnsibleGalaxyHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *AnsibleGalaxyRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAnsibleGalaxyHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *AnsibleGalaxyRepositoryFormatHosted) UpdateStateFromApi(state, api any) any {
	var stateModel model.RepositoryAnsibleGalaxyHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAnsibleGalaxyHostedModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.AnsibleGalaxyHostedApiRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Ansible Galaxy Hosted repositories
func (f *AnsibleGalaxyRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAnsiblegalaxyHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

// --------------------------------------------
// Proxy Ansible Galaxy Format Functions
// --------------------------------------------
func (f *AnsibleGalaxyRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAnsibleGalaxyProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateAnsiblegalaxyProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *AnsibleGalaxyRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAnsibleGalaxyProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAnsiblegalaxyProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *AnsibleGalaxyRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAnsibleGalaxyProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAnsibleGalaxyProxyModel)

	// Call API to Update
	return apiClient.RepositoryManagementAPI.UpdateAnsiblegalaxyProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

// DoImportRequest implements the import functionality for Ansible Galaxy Proxy repositories
func (f *AnsibleGalaxyRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAnsiblegalaxyProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *AnsibleGalaxyRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
}

func (f *AnsibleGalaxyRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAnsibleGalaxyProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *AnsibleGalaxyRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAnsibleGalaxyProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *AnsibleGalaxyRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAnsibleGalaxyProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *AnsibleGalaxyRepositoryFormatProxy) UpdateStateFromApi(state, api any) any {
	var stateModel model.RepositoryAnsibleGalaxyProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAnsibleGalaxyProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.AnsibleGalaxyProxyApiRepository))
	return stateModel
}

func (f *AnsibleGalaxyRepositoryFormatProxy) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryAnsibleGalaxyProxyModel)
	var stateModel model.RepositoryAnsibleGalaxyProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAnsibleGalaxyProxyModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	return stateModel
}

// Ansible Galaxy Proxy repos are not supported by Repository Firewall
func (f *AnsibleGalaxyRepositoryFormatProxy) SupportsRepositoryFirewall() bool {
	return false
}

// --------------------------------------------
// Group Ansible Galaxy Format Functions
// --------------------------------------------
func (f *AnsibleGalaxyRepositoryFormatGroup) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAnsibleGalaxyGroupModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateAnsiblegalaxyGroupRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *AnsibleGalaxyRepositoryFormatGroup) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAnsibleGalaxyGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAnsiblegalaxyGroupRepository(ctx, stateModel.Name.ValueString()).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *AnsibleGalaxyRepositoryFormatGroup) DoUpdateRequest(plan, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAnsibleGalaxyGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAnsibleGalaxyGroupModel)

	// Call API to Update
	return apiClient.RepositoryManagementAPI.UpdateAnsiblegalaxyGroupRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *AnsibleGalaxyRepositoryFormatGroup) FormatSchemaAttributes() map[string]tfschema.Attribute {
	return commonGroupSchemaAttributes(false)
}

func (f *AnsibleGalaxyRepositoryFormatGroup) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAnsibleGalaxyGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *AnsibleGalaxyRepositoryFormatGroup) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAnsibleGalaxyGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *AnsibleGalaxyRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAnsibleGalaxyGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *AnsibleGalaxyRepositoryFormatGroup) UpdateStateFromApi(state, api any) any {
	var stateModel model.RepositoryAnsibleGalaxyGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAnsibleGalaxyGroupModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.AnsibleGalaxyGroupApiRepository))
	return stateModel
}

// DoImportRequest implements the import functionality for Ansible Galaxy Group repositories
func (f *AnsibleGalaxyRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAnsiblegalaxyGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}
