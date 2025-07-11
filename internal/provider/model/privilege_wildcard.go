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

// PrivilegeWildcardModel
// ------------------------------------
type PrivilegeWildcardModel struct {
	PrivilegeModelResource
	Pattern types.String `tfsdk:"pattern"`
}

func (p *PrivilegeWildcardModel) FromApiModel(api sonatyperepo.ApiPrivilegeRequest) {
	p.Name = types.StringValue(api.Name)
	p.Description = types.StringPointerValue(api.Description)
	p.ReadOnly = types.BoolPointerValue(api.ReadOnly)
	p.Pattern = types.StringPointerValue(api.Pattern)
}

func (p *PrivilegeWildcardModel) ToApiCreateModel() sonatyperepo.ApiPrivilegeWildcardRequest {
	apiModel := sonatyperepo.NewApiPrivilegeWildcardRequest()
	p.MapToApi(apiModel)
	return *apiModel
}

func (p *PrivilegeWildcardModel) MapToApi(api *sonatyperepo.ApiPrivilegeWildcardRequest) {
	api.Name = p.Name.ValueStringPointer()
	api.Description = p.Description.ValueStringPointer()
	api.Pattern = p.Pattern.ValueStringPointer()
}
