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

type PrivilegeType int8

const (
	TypeApplication PrivilegeType = iota
	TypeRepositoryAdmin
	TypeRepositoryContentSelector
	TypeRepositoryView
	TypeScript
	TypeWildcard
)

func (pt PrivilegeType) String() string {
	switch pt {
	case TypeApplication:
		return "application"
	case TypeRepositoryAdmin:
		return "repository-admin"
	case TypeRepositoryContentSelector:
		return "repository-content-selector"
	case TypeRepositoryView:
		return "repository-view"
	case TypeScript:
		return "script"
	case TypeWildcard:
		return "wildcard"
	}

	return "unknown"
}

func AllPrivilegeTypes() []string {
	return []string{
		TypeApplication.String(),
		TypeRepositoryAdmin.String(),
		TypeRepositoryContentSelector.String(),
		TypeRepositoryView.String(),
		TypeScript.String(),
		TypeWildcard.String(),
	}
}

// BasePrivilegeModel
// ------------------------------------
type BasePrivilegeModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ReadOnly    types.Bool   `tfsdk:"read_only"`
	Type        types.String `tfsdk:"type"`
}

func (m *BasePrivilegeModel) MapFromApi(api *sonatyperepo.ApiPrivilege) {
	m.Name = types.StringPointerValue(api.Name)
	m.Description = types.StringPointerValue(api.Description)
	m.ReadOnly = types.BoolPointerValue(api.ReadOnly)
	m.Type = types.StringPointerValue(api.Type)
}

// PrivilegesModel (used by DataSource)
// ------------------------------------
type PrivilegesModel struct {
	Privileges []BasePrivilegeModel `tfsdk:"privileges"`
}
