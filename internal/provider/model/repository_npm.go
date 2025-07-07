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

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type RepositoryNpmHostedModel struct {
	RepositoryHostedModel
}

func (m *RepositoryNpmHostedModel) FromApiModel(api sonatyperepo.SimpleApiHostedRepository) {
	m.mapSimpleApiHostedRepository(api)
}

func (m *RepositoryNpmHostedModel) ToApiCreateModel() sonatyperepo.NpmHostedRepositoryApiRequest {
	apiModel := sonatyperepo.NpmHostedRepositoryApiRequest{
		Name:   m.Name.ValueString(),
		Online: m.Online.ValueBool(),
		Storage: sonatyperepo.HostedStorageAttributes{
			BlobStoreName:               m.Storage.BlobStoreName.ValueString(),
			StrictContentTypeValidation: m.Storage.StrictContentTypeValidation.ValueBool(),
			WritePolicy:                 m.Storage.WritePolicy.ValueString(),
		},
		Component: &sonatyperepo.ComponentAttributes{
			ProprietaryComponents: common.NewFalse(),
		},
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: make([]string, 0),
		},
	}

	if m.Component != nil {
		apiModel.Component.ProprietaryComponents = m.Component.ProprietaryComponents.ValueBoolPointer()
	}

	if m.Cleanup != nil {
		for _, p := range m.Cleanup.PolicyNames {
			apiModel.Cleanup.PolicyNames = append(apiModel.Cleanup.PolicyNames, p.ValueString())
		}
	}

	return apiModel
}

func (m *RepositoryNpmHostedModel) ToApiUpdateModel() sonatyperepo.NpmHostedRepositoryApiRequest {
	return m.ToApiCreateModel()
}
