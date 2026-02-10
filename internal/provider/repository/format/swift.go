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

type SwiftRepositoryFormat struct {
	BaseRepositoryFormat
}

type SwiftRepositoryFormatProxy struct {
	SwiftRepositoryFormat
}

// --------------------------------------------
// Generic Swift Format Functions
// --------------------------------------------
func (f *SwiftRepositoryFormat) Key() string {
	return common.REPO_FORMAT_SWIFT
}

func (f *SwiftRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

// --------------------------------------------
// PROXY Swift Format Functions
// --------------------------------------------
func (f *SwiftRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorySwiftProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateSwiftProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *SwiftRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorySwiftProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetSwiftProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *SwiftRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositorySwiftProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositorySwiftProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateSwiftProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

// DoImportRequest implements the import functionality for HuggingFace Proxy repositories
func (f *SwiftRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetSwiftProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *SwiftRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
	maps.Copy(additionalAttributes, swiftProxySchemaAttributes())
	return additionalAttributes
}

func (f *SwiftRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositorySwiftProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *SwiftRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositorySwiftProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *SwiftRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositorySwiftProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *SwiftRepositoryFormatProxy) UpdateStateFromApi(state, api any) any {
	var stateModel model.RepositorySwiftProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositorySwiftProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.SwiftProxyApiRepository))
	return stateModel
}

func (f *SwiftRepositoryFormatProxy) SupportsRepositoryFirewall() bool {
	return false
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func swiftProxySchemaAttributes() map[string]tfschema.Attribute {
	swiftAttrs := map[string]tfschema.Attribute{
		"swift": schema.ResourceRequiredSingleNestedAttribute(
			"Swift specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"require_authentication": schema.ResourceOptionalBoolWithDefault(
					"Indicates if this repository requires authentication overriding anonymous access.",
					false,
				),
			},
		),
	}

	return swiftAttrs
}
