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

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Capability: Audit
// ----------------------------------------
type CapabilityAuditModel struct {
	BaseCapabilityModel
}

func (m *CapabilityAuditModel) FromApiModel(api *v3.CapabilityDTO) {
	m.Id = types.StringPointerValue(api.Id)
	m.Notes = types.StringValue(*api.Notes)
	m.Enabled = types.BoolPointerValue(api.Enabled)
}

func (m *CapabilityAuditModel) ToApiCreateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiCreateModel()
	api.Type = common.CAPABILITY_TYPE_AUDIT.StringPointer()
	return api
}

func (m *CapabilityAuditModel) ToApiUpdateModel(version common.SystemVersion) *v3.CapabilityDTO {
	api := m.toApiUpdateModel()
	return api
}
