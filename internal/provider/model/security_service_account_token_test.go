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

package model_test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

func strPtr(s string) *string { return &s }

func TestServiceAccountTokenMapToApi(t *testing.T) {
	t.Run("required arguments only: no description and no expiry are sent", func(t *testing.T) {
		m := model.SecurityServiceAccountTokenModel{
			Name:           types.StringValue("ci"),
			RoleId:         types.StringValue("nx-anonymous"),
			Description:    types.StringNull(),
			ExpirationDays: types.Int64Null(),
		}
		var api sonatyperepoV397.ServiceAccountCreateXO

		m.MapToApi(&api)

		assert.Equal(t, "ci", api.Name)
		assert.Equal(t, "nx-anonymous", api.RoleId)
		assert.Nil(t, api.Description)
		assert.Nil(t, api.ExpirationDays, "omitting expirationDays creates a token that never expires")
	})

	t.Run("all arguments", func(t *testing.T) {
		m := model.SecurityServiceAccountTokenModel{
			Name:           types.StringValue("ci"),
			RoleId:         types.StringValue("nx-anonymous"),
			Description:    types.StringValue("for CI"),
			ExpirationDays: types.Int64Value(30),
		}
		var api sonatyperepoV397.ServiceAccountCreateXO

		m.MapToApi(&api)

		assert.Equal(t, "for CI", *api.Description)
		assert.Equal(t, int32(30), *api.ExpirationDays)
	})
}

func TestServiceAccountTokenMapFromCreate(t *testing.T) {
	expires := time.Date(2026, 12, 8, 0, 15, 20, 838723000, time.UTC)
	m := model.SecurityServiceAccountTokenModel{}

	m.MapFromCreate(&sonatyperepoV397.ServiceAccountCreatedXO{
		Id:        strPtr("id-1"),
		Name:      strPtr("ci"),
		Token:     strPtr("sat.secret"),
		ExpiresAt: &expires,
	})

	assert.Equal(t, "id-1", m.ID.ValueString())
	assert.Equal(t, "sat.secret", m.Token.ValueString())
	assert.Equal(t, "2026-12-08T00:15:20.838723Z", m.ExpiresAt.ValueString())
	// Not in the create response, so must be null (known) rather than left unknown in state.
	assert.True(t, m.CreatedAt.IsNull())
	assert.True(t, m.CreatedBy.IsNull())
}

func TestServiceAccountTokenMapFromCreateNeverExpires(t *testing.T) {
	m := model.SecurityServiceAccountTokenModel{}

	m.MapFromCreate(&sonatyperepoV397.ServiceAccountCreatedXO{Id: strPtr("id-1"), Token: strPtr("sat.secret")})

	assert.True(t, m.ExpiresAt.IsNull())
}

func TestServiceAccountTokenMapFromApi(t *testing.T) {
	created := time.Date(2026, 11, 8, 0, 15, 20, 794601000, time.UTC)

	t.Run("fills in the read-only details and leaves the token alone", func(t *testing.T) {
		m := model.SecurityServiceAccountTokenModel{
			Token:          types.StringValue("sat.secret"),
			ExpirationDays: types.Int64Value(30),
		}

		m.MapFromApi(&sonatyperepoV397.ServiceAccountXO{
			Id:          strPtr("id-1"),
			Name:        strPtr("ci"),
			RoleId:      strPtr("nx-anonymous"),
			Description: strPtr("for CI"),
			CreatedBy:   strPtr("admin"),
			CreatedAt:   &created,
		})

		assert.Equal(t, "id-1", m.ID.ValueString())
		assert.Equal(t, "ci", m.Name.ValueString())
		assert.Equal(t, "nx-anonymous", m.RoleId.ValueString())
		assert.Equal(t, "for CI", m.Description.ValueString())
		assert.Equal(t, "admin", m.CreatedBy.ValueString())
		assert.Equal(t, "2026-11-08T00:15:20.794601Z", m.CreatedAt.ValueString())
		assert.True(t, m.ExpiresAt.IsNull())
		assert.Equal(t, "sat.secret", m.Token.ValueString(), "the list never returns the token, so it must be kept")
		assert.Equal(t, int64(30), m.ExpirationDays.ValueInt64(), "the API does not report the requested number of days")
	})

	t.Run("absent description stays null", func(t *testing.T) {
		m := model.SecurityServiceAccountTokenModel{Description: types.StringNull()}
		m.MapFromApi(&sonatyperepoV397.ServiceAccountXO{Id: strPtr("id-1")})
		assert.True(t, m.Description.IsNull())
	})

	t.Run("an explicitly empty description does not drift to null", func(t *testing.T) {
		m := model.SecurityServiceAccountTokenModel{Description: types.StringValue("")}
		m.MapFromApi(&sonatyperepoV397.ServiceAccountXO{Id: strPtr("id-1")})
		assert.False(t, m.Description.IsNull())
		assert.Equal(t, "", m.Description.ValueString())
	})

	t.Run("a description removed on the server is reported as null", func(t *testing.T) {
		m := model.SecurityServiceAccountTokenModel{Description: types.StringValue("was set")}
		m.MapFromApi(&sonatyperepoV397.ServiceAccountXO{Id: strPtr("id-1")})
		assert.True(t, m.Description.IsNull())
	})
}
