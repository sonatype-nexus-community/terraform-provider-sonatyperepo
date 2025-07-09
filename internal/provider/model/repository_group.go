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

// RepositoryGroupModel
// --------------------------------------------------------
type RepositoryGroupModel struct {
	BasicRepositoryModel
	Storage repositoryStorageModel `tfsdk:"storage"`
	Group   repositoryGroupDetails `tfsdk:"group"`
}

// repositoryGroupDetails
// --------------------------------------------------------
type repositoryGroupDetails struct {
	MemberNames []types.String `tfsdk:"member_names"`
}

func (m *repositoryGroupDetails) MapFromApi(api *sonatyperepo.GroupAttributes) {
	m.MemberNames = make([]types.String, 0)
	for _, n := range api.GetMemberNames() {
		m.MemberNames = append(m.MemberNames, types.StringValue(n))
	}
}

func (m *repositoryGroupDetails) MapToApi(api *sonatyperepo.GroupAttributes) {
	api.MemberNames = make([]string, 0)
	for _, n := range m.MemberNames {
		api.MemberNames = append(api.MemberNames, n.ValueString())
	}
}
