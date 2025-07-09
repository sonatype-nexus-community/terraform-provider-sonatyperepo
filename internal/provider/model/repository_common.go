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
	Name   types.String `tfsdk:"name"`
	Format types.String `tfsdk:"format"`
	Type   types.String `tfsdk:"type"`
	Url    types.String `tfsdk:"url"`
	// LastUpdated types.String `tfsdk:"last_updated"`
}

type RepositoriesModel struct {
	Repositories []RepositoryModel `tfsdk:"repositories"`
}

type BasicRepositoryModel struct {
	Name        types.String            `tfsdk:"name"`
	Online      types.Bool              `tfsdk:"online"`
	Url         types.String            `tfsdk:"url"`
	Cleanup     *RepositoryCleanupModel `tfsdk:"cleanup"`
	LastUpdated types.String            `tfsdk:"last_updated"`
}

type RepositoryCleanupModel struct {
	PolicyNames []types.String `tfsdk:"policy_names"`
}

func NewRepositoryCleanupModel() *RepositoryCleanupModel {
	return &RepositoryCleanupModel{
		PolicyNames: make([]types.String, 0),
	}
}

func mapCleanupFromApi(api *sonatyperepo.CleanupPolicyAttributes, m *RepositoryCleanupModel) {
	for _, p := range api.GetPolicyNames() {
		m.PolicyNames = append(m.PolicyNames, types.StringValue(p))
	}
}

func mapCleanupToApi(m *RepositoryCleanupModel, api *sonatyperepo.CleanupPolicyAttributes) {
	if m != nil {
		for _, p := range m.PolicyNames {
			api.PolicyNames = append(api.PolicyNames, p.ValueString())
		}
	}
}

// repositoryStorageModel
// ----------------------------------------
type repositoryStorageModel struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
	// WritePolicy                 types.String `tfsdk:"write_policy"`
}

func (m *repositoryStorageModel) MapFromApi(api *sonatyperepo.StorageAttributes) {
	m.BlobStoreName = types.StringValue(api.BlobStoreName)
	m.StrictContentTypeValidation = types.BoolValue(api.StrictContentTypeValidation)
	// m.WritePolicy = types.StringPointerValue(api.WritePolicy)

}

func (m *repositoryStorageModel) MapToApi(api *sonatyperepo.StorageAttributes) {
	api.BlobStoreName = m.BlobStoreName.ValueString()
	api.StrictContentTypeValidation = m.StrictContentTypeValidation.ValueBool()
	// api.WritePolicy = m.WritePolicy.ValueStringPointer()
}

// func mapStorageNonGroupFromApi(api *sonatyperepo.StorageAttributes, m *repositoryStorageModelNonGroup) {
// 	m.BlobStoreName = types.StringValue(api.BlobStoreName)
// 	m.StrictContentTypeValidation = types.BoolValue(api.StrictContentTypeValidation)
// 	m.WritePolicy = types.StringPointerValue(api.WritePolicy)
// }

// func mapStorageNonGroupToApi(m *repositoryStorageModelNonGroup, api *sonatyperepo.StorageAttributes) {
// 	api.BlobStoreName = m.BlobStoreName.ValueString()
// 	api.StrictContentTypeValidation = m.StrictContentTypeValidation.ValueBool()
// 	api.WritePolicy = m.WritePolicy.ValueStringPointer()
// }

// func mapHostedStorageAttributesFromApi(api *sonatyperepo.HostedStorageAttributes, m *repositoryStorageModelNonGroup) {
// 	m.BlobStoreName = types.StringValue(api.BlobStoreName)
// 	m.StrictContentTypeValidation = types.BoolValue(api.StrictContentTypeValidation)
// 	m.WritePolicy = types.StringValue(api.WritePolicy)
// }

// func mapHostedStorageAttributesToApi(m repositoryStorageModelNonGroup, api *sonatyperepo.HostedStorageAttributes) {
// 	api.BlobStoreName = m.BlobStoreName.ValueString()
// 	api.StrictContentTypeValidation = m.StrictContentTypeValidation.ValueBool()
// 	api.WritePolicy = m.WritePolicy.ValueString()
// }
