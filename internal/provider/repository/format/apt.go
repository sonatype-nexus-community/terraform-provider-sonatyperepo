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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type AptRepositoryFormat struct {
	BaseRepositoryFormat
}

type AptRepositoryFormatHosted struct {
	AptRepositoryFormat
}

type AptRepositoryFormatProxy struct {
	AptRepositoryFormat
}

// --------------------------------------------
// Generic APT Format Functions
// --------------------------------------------
func (f *AptRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_APT
}

func (f *AptRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// HOSTED APT Format Functions
// --------------------------------------------
func (f *AptRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAptHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateAptHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *AptRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAptHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAptHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *AptRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAptHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAptHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateAptHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

// DoImportRequest implements the import functionality for APT Hosted repositories
func (f *AptRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAptHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *AptRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, getAptSchemaAttributes(false))
	return additionalAttributes
}

func (f *AptRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAptHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *AptRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAptHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *AptRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAptHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *AptRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryAptHostedModel
	var preserveAptSigning bool
	var existingKeyPair types.String
	var existingPassphrase types.String
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAptHostedModel)
		// Check if apt_signing was in the plan/state (either field is not null)
		preserveAptSigning = stateModel.AptSigning != nil && (!stateModel.AptSigning.KeyPair.IsNull() || !stateModel.AptSigning.Passphrase.IsNull())
		if preserveAptSigning {
			// Preserve apt_signing values (API doesn't return sensitive passphrase)
			existingKeyPair = stateModel.AptSigning.KeyPair
			existingPassphrase = stateModel.AptSigning.Passphrase
		}
	}
	stateModel.FromApiModel((api).(sonatyperepo.AptHostedApiRepository))
	// Restore apt_signing from plan/state if it was provided (API doesn't return sensitive passphrase)
	if preserveAptSigning {
		stateModel.AptSigning.KeyPair = existingKeyPair
		stateModel.AptSigning.Passphrase = existingPassphrase
	}
	return stateModel
}

// --------------------------------------------
// PROXY Maven Format Functions
// --------------------------------------------
func (f *AptRepositoryFormatProxy) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAptProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateAptProxyRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *AptRepositoryFormatProxy) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAptProxyModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAptProxyRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *AptRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryAptProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryAptProxyModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateAptProxyRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

// DoImportRequest implements the import functionality for APT Proxy repositories
func (f *AptRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetAptProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *AptRepositoryFormatProxy) GetFormatSchemaAttributes() map[string]schema.Attribute {
	additionalAttributes := getCommonProxySchemaAttributes()
	maps.Copy(additionalAttributes, getAptSchemaAttributes(true))
	return additionalAttributes
}

func (f *AptRepositoryFormatProxy) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryAptProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *AptRepositoryFormatProxy) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryAptProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *AptRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryAptProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *AptRepositoryFormatProxy) UpdateStateFromApi(state any, api any) any {
	var stateModel model.RepositoryAptProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryAptProxyModel)
	}
	stateModel.FromApiModel((api).(sonatyperepo.AptProxyApiRepository))
	return stateModel
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func getAptSchemaAttributes(isProxy bool) map[string]schema.Attribute {
	aptAttrs := map[string]schema.Attribute{
		"distribution": schema.StringAttribute{
			Description: "Distribution to fetch",
			Required:    true,
		},
	}
	if isProxy {
		aptAttrs["flat"] = schema.BoolAttribute{
			Description: "Whether this repository is flat",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		}
	}

	attrs := map[string]schema.Attribute{
		"apt": schema.SingleNestedAttribute{
			Description: "APT specific configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes:  aptAttrs,
		},
	}

	if !isProxy {
		attrs["apt_signing"] = schema.SingleNestedAttribute{
			Description: "APT signing configuration for this Repository",
			Optional:    true,
			Attributes: map[string]schema.Attribute{
				"key_pair": schema.StringAttribute{
					Description: "PGP signing key pair (armored private key e.g. gpg --export-secret-key --armor)",
					Required:    true,
				},
				"passphrase": schema.StringAttribute{
					Description: "Passphrase to access PGP signing key",
					Required:    true,
					Sensitive:   true,
				},
			},
		}
	}

	return attrs
}
