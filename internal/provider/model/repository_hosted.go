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

type RepositoryHostedModel struct {
	BasicRepositoryModel
	Component *RepositoryComponentModel `tfsdk:"component"`
}

// RepositoryComponentModel
// ----------------------------------------
type RepositoryComponentModel struct {
	ProprietaryComponents types.Bool `tfsdk:"proprietary_components"`
}

func (m *RepositoryComponentModel) MapFromApi(api *sonatyperepo.ComponentAttributes) {
	if api != nil {
		m.ProprietaryComponents = types.BoolValue(*api.ProprietaryComponents)
	}
}

func (m *RepositoryComponentModel) MapToApi(api *sonatyperepo.ComponentAttributes) {
	if m != nil {
		api.ProprietaryComponents = m.ProprietaryComponents.ValueBoolPointer()
	}
}

func (m *RepositoryHostedModel) mapSimpleApiHostedRepository(api sonatyperepo.SimpleApiHostedRepository) {
	m.Name = types.StringPointerValue(api.Name)
	m.Online = types.BoolValue(api.Online)
	m.Url = types.StringPointerValue(api.Url)
	m.Storage = repositoryStorageModelNonGroup{}
	mapHostedStorageAttributesFromApi(&api.Storage, &m.Storage)

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
