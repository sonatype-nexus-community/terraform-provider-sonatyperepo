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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type GitLfsRepositoryFormat struct {
	BaseRepositoryFormat
}

type GitLfsRepositoryFormatHosted struct {
	GitLfsRepositoryFormat
}

// --------------------------------------------
// Generic Git LFS Format Functions
// --------------------------------------------
func (f *GitLfsRepositoryFormat) GetKey() string {
	return common.REPO_FORMAT_GIT_LFS
}

func (f *GitLfsRepositoryFormat) GetResourceName(repoType RepositoryType) string {
	return getResourceName(f.GetKey(), repoType)
}

// --------------------------------------------
// Hosted Git LFS Format Functions
// --------------------------------------------
func (f *GitLfsRepositoryFormatHosted) DoCreateRequest(plan any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryGitLfsHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.CreateGitlfsHostedRepository(ctx).Body(planModel.ToApiCreateModel()).Execute()
}

func (f *GitLfsRepositoryFormatHosted) DoReadRequest(state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryGitLfsHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.RepositoryManagementAPI.GetGitlfsHostedRepository(ctx, stateModel.Name.ValueString()).Execute()
	return *apiResponse, httpResponse, err
}

func (f *GitLfsRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient *sonatyperepo.APIClient, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryGitLfsHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryGitLfsHostedModel)

	// Call API to Create
	return apiClient.RepositoryManagementAPI.UpdateGitlfsHostedRepository(ctx, stateModel.Name.ValueString()).Body(planModel.ToApiUpdateModel()).Execute()
}

func (f *GitLfsRepositoryFormatHosted) GetFormatSchemaAttributes() map[string]schema.Attribute {
	return getCommonHostedSchemaAttributes()
}

func (f *GitLfsRepositoryFormatHosted) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryGitLfsHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *GitLfsRepositoryFormatHosted) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryGitLfsHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *GitLfsRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryGitLfsHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *GitLfsRepositoryFormatHosted) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.RepositoryGitLfsHostedModel)
	stateModel.FromApiModel((api).(sonatyperepo.SimpleApiHostedRepository))
	return stateModel
}
