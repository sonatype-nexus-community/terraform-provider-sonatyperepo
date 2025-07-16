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

package model

import (
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Git LFS Hosted
// ----------------------------------------
type RepositoryGitLfsHostedModel struct {
	RepositoryHostedModel
}

func (m *RepositoryGitLfsHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)
	m.Storage.MapFromApi(&api.Storage)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = &RepositoryCleanupModel{}
		mapCleanupFromApi(api.Cleanup, m.Cleanup)
	}

	// Component
	if api.Component != nil {
		m.Component = &RepositoryComponentModel{}
		m.Component.MapFromApi(api.Component)
	}
}

func (m *RepositoryGitLfsHostedModel) ToApiCreateModel() sonatyperepo.GitLfsHostedRepositoryApiRequest {
	apiModel := sonatyperepo.GitLfsHostedRepositoryApiRequest{
		Name:    m.Name.ValueString(),
		Online:  m.Online.ValueBool(),
		Storage: sonatyperepo.HostedStorageAttributes{},
		Component: &sonatyperepo.ComponentAttributes{
			ProprietaryComponents: common.NewFalse(),
		},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
	}
	m.Storage.MapToApi(&apiModel.Storage)
	mapCleanupToApi(m.Cleanup, apiModel.Cleanup)
	m.Component.MapToApi(apiModel.Component)

	return apiModel
}

func (m *RepositoryGitLfsHostedModel) ToApiUpdateModel() sonatyperepo.GitLfsHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}
