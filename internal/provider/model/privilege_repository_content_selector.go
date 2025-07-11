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

// PrivilegeRepositoryAdminModel
// ------------------------------------
type PrivilegeRepositoryContentSelectorModel struct {
	PrivilegeRepositoryAdminModel
	ContentSelector types.String `tfsdk:"content_selector"`
}

func (p *PrivilegeRepositoryContentSelectorModel) FromApiModel(api sonatyperepo.ApiPrivilegeRequest) {
	p.Name = types.StringValue(api.Name)
	p.Description = types.StringPointerValue(api.Description)
	p.ReadOnly = types.BoolPointerValue(api.ReadOnly)
	p.Actions = make([]types.String, 0)
	for _, a := range api.Actions {
		p.Actions = append(p.Actions, types.StringValue(a))
	}
	p.Format = types.StringPointerValue(api.Format)
	p.Repository = types.StringPointerValue(api.Repository)
	p.ContentSelector = types.StringPointerValue(api.ContentSelector)
}

func (p *PrivilegeRepositoryContentSelectorModel) ToApiCreateModel() sonatyperepo.ApiPrivilegeRepositoryContentSelectorRequest {
	apiModel := sonatyperepo.NewApiPrivilegeRepositoryContentSelectorRequest()
	p.MapToApi(apiModel)
	return *apiModel
}

func (p *PrivilegeRepositoryContentSelectorModel) MapToApi(api *sonatyperepo.ApiPrivilegeRepositoryContentSelectorRequest) {
	api.Name = p.Name.ValueStringPointer()
	api.Description = p.Description.ValueStringPointer()
	api.Actions = make([]string, 0)
	for _, a := range p.Actions {
		api.Actions = append(api.Actions, a.ValueString())
	}
	api.Format = p.Format.ValueStringPointer()
	api.Repository = p.Repository.ValueStringPointer()
	api.ContentSelector = p.ContentSelector.ValueStringPointer()
}
