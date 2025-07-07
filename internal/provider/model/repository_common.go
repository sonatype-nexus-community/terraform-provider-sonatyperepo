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
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type RepositoryModel struct {
	Name        types.String `tfsdk:"name"`
	Format      types.String `tfsdk:"format"`
	Type        types.String `tfsdk:"type"`
	Url         types.String `tfsdk:"url"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

type RepositoriesModel struct {
	Repositories []RepositoryModel `tfsdk:"repositories"`
}

type RepositoryHostedModel struct {
	Name        types.String                   `tfsdk:"name"`
	Online      types.Bool                     `tfsdk:"online"`
	Storage     repositoryStorageModelNonGroup `tfsdk:"storage"`
	Url         types.String                   `tfsdk:"url"`
	Cleanup     *RepositoryCleanupModel        `tfsdk:"cleanup"`
	Component   *RepositoryComponentModel      `tfsdk:"component"`
	LastUpdated types.String                   `tfsdk:"last_updated"`
}

func (m *RepositoryHostedModel) mapSimpleApiHostedRepository(api sonatyperepo.SimpleApiHostedRepository) {
	m.Name = types.StringPointerValue(api.Name)
	// m.Format = types.StringValue(*api.Format)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)

	// Cleanup
	if api.Cleanup != nil && len(api.Cleanup.PolicyNames) > 0 {
		m.Cleanup = &RepositoryCleanupModel{
			PolicyNames: make([]types.String, 0),
		}
		for _, p := range api.Cleanup.GetPolicyNames() {
			m.Cleanup.PolicyNames = append(m.Cleanup.PolicyNames, types.StringValue(p))
		}
	}

	// Component
	if api.Component != nil {
		m.Component = &RepositoryComponentModel{
			ProprietaryComponents: types.BoolValue(*api.Component.ProprietaryComponents),
		}
	}

	// Storage
	m.Storage = repositoryStorageModelNonGroup{
		BlobStoreName:               types.StringValue(api.Storage.BlobStoreName),
		StrictContentTypeValidation: types.BoolValue(api.Storage.StrictContentTypeValidation),
		WritePolicy:                 types.StringValue(api.Storage.WritePolicy),
	}
}
