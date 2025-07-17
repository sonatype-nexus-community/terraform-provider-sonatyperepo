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

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
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
func (f *YumRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_YUM
}

func (f *YumRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
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

func (f *YumRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, getYumSchemaAttributes(true))
	return additionalAttributes
}

func (f *YumRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryYumHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *YumRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryYumHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *YumRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryYumHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *YumRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryYumHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.YumHostedApiRepository))
	return stateModel
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

func (f *YumRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonProxySchemaAttributes()
	maps.Copy(additionalAttributes, getYumSchemaAttributes(false))
	return additionalAttributes
}

func (f *YumRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryYumProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *YumRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryYumProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *YumRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryYumProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *YumRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryYumProxyModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiProxyRepository))
	return stateModel
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

func (f *YumRepositoryFormatGroup) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonGroupSchemaAttributes(false)
	maps.Copy(additionalAttributes, getYumSchemaAttributes(false))
	return additionalAttributes
}

func (f *YumRepositoryFormatGroup) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryYumGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *YumRepositoryFormatGroup) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryYumGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *YumRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryYumGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *YumRepositoryFormatGroup) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryYumGroupModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiGroupRepository))
	return stateModel
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getYumSchemaAttributes(isHosted bool) map[string]schema.Attribute {
	var attrs = make(map[string]schema.Attribute, 0)
	if isHosted {
		attrs["deploy_policy"] = schema.StringAttribute{
			Description: "Validate that all paths are RPMs or yum metadata.",
			Optional:    true,
			Validators: []validator.String{
				stringvalidator.OneOf(common.DEPLOY_POLICY_PERMISSIVE, common.DEPLOY_POLICY_STRICT),
			},
		}
		attrs["repo_data_depth"] = schema.Int32Attribute{
			Description: "Specifies the repository depth where repodata folder(s) are created",
			Required:    true,
			Validators: []validator.Int32{
				int32validator.AtLeast(0),
				int32validator.AtMost(5),
			},
		}
	} else {
		attrs["key_pair"] = schema.StringAttribute{
			Description: "PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)",
			Optional:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		attrs["passphrase"] = schema.StringAttribute{
			Description: "Passphrase to access PGP signing key",
			Optional:    true,
			Sensitive:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
	}
	return map[string]schema.Attribute{
		"yum": schema.SingleNestedAttribute{
			Description: "YUM specific configuration for this Repository",
			Required:    isHosted,
			Optional:    !isHosted,
			Attributes:  attrs,
		},
	}
}
