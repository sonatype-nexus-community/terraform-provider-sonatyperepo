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
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// SecurityServiceAccountTokenModel is the Terraform model of a Service Account Token.
//
// The Nexus Repository API has no get-by-id and no update, and returns the token value exactly
// once, when the token is created. Every user-settable attribute therefore forces replacement,
// and `Token` is only ever populated from the create response.
type SecurityServiceAccountTokenModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	RoleId         types.String `tfsdk:"role_id"`
	Description    types.String `tfsdk:"description"`
	ExpirationDays types.Int64  `tfsdk:"expiration_days"`
	Token          types.String `tfsdk:"token"`
	ExpiresAt      types.String `tfsdk:"expires_at"`
	CreatedAt      types.String `tfsdk:"created_at"`
	CreatedBy      types.String `tfsdk:"created_by"`
}

// MapToApi builds the create request from the plan.
func (m *SecurityServiceAccountTokenModel) MapToApi(api *sonatyperepoV397.ServiceAccountCreateXO) {
	api.Name = m.Name.ValueString()
	api.RoleId = m.RoleId.ValueString()

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		description := m.Description.ValueString()
		api.Description = &description
	}

	// Omitting expirationDays creates a token that never expires.
	if !m.ExpirationDays.IsNull() && !m.ExpirationDays.IsUnknown() {
		days := int32(m.ExpirationDays.ValueInt64())
		api.ExpirationDays = &days
	}
}

// MapFromCreate maps the create response. This is the only response that carries the token
// value, and it does not carry the role, description or creation details, which stay as
// planned (or are completed by MapFromApi). Attributes that this response does not provide are
// set to null so that no unknown value is left in the state.
func (m *SecurityServiceAccountTokenModel) MapFromCreate(api *sonatyperepoV397.ServiceAccountCreatedXO) {
	m.ID = types.StringPointerValue(api.Id)
	m.Token = types.StringPointerValue(api.Token)
	m.ExpiresAt = timeToString(api.ExpiresAt)
	m.CreatedAt = types.StringNull()
	m.CreatedBy = types.StringNull()
}

// MapFromApi maps a list entry. It never touches `Token` (secret material is never returned
// by the list) or `ExpirationDays` (the API reports the resulting expiry time, not the
// number of days it was requested with).
func (m *SecurityServiceAccountTokenModel) MapFromApi(api *sonatyperepoV397.ServiceAccountXO) {
	if api.Id != nil {
		m.ID = types.StringValue(*api.Id)
	}
	if api.Name != nil {
		m.Name = types.StringValue(*api.Name)
	}
	if api.RoleId != nil {
		m.RoleId = types.StringValue(*api.RoleId)
	}

	// The API omits an empty description. Keep an explicit empty string as configured so that
	// `description = ""` does not drift to null.
	if api.Description != nil && *api.Description != "" {
		m.Description = types.StringValue(*api.Description)
	} else if m.Description.IsNull() || m.Description.IsUnknown() || m.Description.ValueString() != "" {
		m.Description = types.StringNull()
	}

	m.ExpiresAt = timeToString(api.ExpiresAt)
	m.CreatedAt = timeToString(api.CreatedAt)
	m.CreatedBy = types.StringPointerValue(api.CreatedBy)
}

// timeToString renders a timestamp the same way wherever it comes from, so that the create
// response and later reads of the same token never differ in format.
func timeToString(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}
	return types.StringValue(t.UTC().Format(time.RFC3339Nano))
}
