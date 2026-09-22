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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	"github.com/stretchr/testify/assert"
)

// TestHasRequiredSecret covers RepositoryHttpClientAuthenticationModel.hasRequiredSecret across
// every auth Type this provider supports, including the nil receiver (an unconfigured
// `authentication` block).
func TestHasRequiredSecret(t *testing.T) {
	tests := []struct {
		name string
		m    *RepositoryHttpClientAuthenticationModel
		want bool
	}{
		{"nil model", nil, false},
		{
			"username with password",
			&RepositoryHttpClientAuthenticationModel{
				Type:     types.StringValue(common.HTTP_AUTH_TYPE_USERNAME),
				Password: types.StringValue("pass"),
			},
			true,
		},
		{
			"username with null password - GH-491",
			&RepositoryHttpClientAuthenticationModel{
				Type:     types.StringValue(common.HTTP_AUTH_TYPE_USERNAME),
				Username: types.StringValue("user"),
				Password: types.StringNull(),
			},
			false,
		},
		{
			"username with unknown password",
			&RepositoryHttpClientAuthenticationModel{
				Type:     types.StringValue(common.HTTP_AUTH_TYPE_USERNAME),
				Password: types.StringUnknown(),
			},
			false,
		},
		{
			"ntlm with password",
			&RepositoryHttpClientAuthenticationModel{
				Type:     types.StringValue(common.HTTP_AUTH_TYPE_NTLM),
				Password: types.StringValue("pass"),
			},
			true,
		},
		{
			"ntlm with null password",
			&RepositoryHttpClientAuthenticationModel{
				Type: types.StringValue(common.HTTP_AUTH_TYPE_NTLM),
			},
			false,
		},
		{
			"bearerToken with token",
			&RepositoryHttpClientAuthenticationModel{
				Type:        types.StringValue(common.HTTP_AUTH_TYPE_BEARER_TOKEN),
				BearerToken: types.StringValue("token"),
			},
			true,
		},
		{
			"bearerToken with null token",
			&RepositoryHttpClientAuthenticationModel{
				Type: types.StringValue(common.HTTP_AUTH_TYPE_BEARER_TOKEN),
			},
			false,
		},
		{
			"null type",
			&RepositoryHttpClientAuthenticationModel{
				Type: types.StringNull(),
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.hasRequiredSecret())
		})
	}
}

// TestMapToApiHttpClientAttributesOmitsAuthenticationMissingSecret verifies that
// MapToApiHttpClientAttributes omits Authentication from the outbound request entirely when the
// model's authentication is missing the secret NXRM requires for its Type, rather than sending a
// half-populated object NXRM's own validation rejects (GH-491).
func TestMapToApiHttpClientAttributesOmitsAuthenticationMissingSecret(t *testing.T) {
	m := repositoryHttpClientModel{
		AutoBlock: types.BoolValue(true),
		Blocked:   types.BoolValue(false),
		Authentication: &RepositoryHttpClientAuthenticationModel{
			Type:     types.StringValue(common.HTTP_AUTH_TYPE_USERNAME),
			Username: types.StringValue("user"),
			Password: types.StringNull(),
		},
	}

	api := sonatyperepo.HttpClientAttributes{}
	m.MapToApiHttpClientAttributes(&api)

	assert.Nil(t, api.Authentication)
}

// TestMapToApiHttpClientAttributesSendsAuthenticationWithSecret is the inverse of the above:
// when the secret is present, Authentication must still be sent as normal.
func TestMapToApiHttpClientAttributesSendsAuthenticationWithSecret(t *testing.T) {
	m := repositoryHttpClientModel{
		AutoBlock: types.BoolValue(true),
		Blocked:   types.BoolValue(false),
		Authentication: &RepositoryHttpClientAuthenticationModel{
			Type:     types.StringValue(common.HTTP_AUTH_TYPE_USERNAME),
			Username: types.StringValue("user"),
			Password: types.StringValue("pass"),
		},
	}

	api := sonatyperepo.HttpClientAttributes{}
	m.MapToApiHttpClientAttributes(&api)

	if assert.NotNil(t, api.Authentication) {
		assert.Equal(t, "user", *api.Authentication.Username)
		assert.Equal(t, "pass", *api.Authentication.Password)
	}
}

// TestMapToApiHttpClientAttributesWithPreemptiveAuthOmitsAuthenticationMissingSecret mirrors
// TestMapToApiHttpClientAttributesOmitsAuthenticationMissingSecret for Maven's
// HttpClientAttributesWithPreemptiveAuth variant.
func TestMapToApiHttpClientAttributesWithPreemptiveAuthOmitsAuthenticationMissingSecret(t *testing.T) {
	m := repositoryHttpClientModel{
		Authentication: &RepositoryHttpClientAuthenticationModel{
			Type:        types.StringValue(common.HTTP_AUTH_TYPE_BEARER_TOKEN),
			BearerToken: types.StringNull(),
		},
	}

	api := sonatyperepo.HttpClientAttributesWithPreemptiveAuth{}
	m.MapToApiHttpClientAttributesWithPreemptiveAuth(&api)

	assert.Nil(t, api.Authentication)
}

// TestHttpClientMapMissingApiFieldsFromPlanRestoresWholeAuthenticationWhenSecretMissing
// reproduces the second half of GH-491's fix: when the plan's authentication is missing its
// required secret (so MapToApiHttpClientAttributes omitted it from the request NXRM saw - see
// above), a subsequent Read() sees NXRM's response with Authentication nil (NXRM's PUT is a full
// replace, so omitting the field genuinely clears it server-side). Since `authentication` is not
// a Computed schema attribute, state must still mirror the plan's authentication object exactly,
// or Terraform's post-apply consistency check fails with "produced an unexpected new value".
func TestHttpClientMapMissingApiFieldsFromPlanRestoresWholeAuthenticationWhenSecretMissing(t *testing.T) {
	planModel := repositoryHttpClientModel{
		Authentication: &RepositoryHttpClientAuthenticationModel{
			Type:     types.StringValue(common.HTTP_AUTH_TYPE_USERNAME),
			Username: types.StringValue("manually-configured-user"),
			Password: types.StringNull(),
		},
	}

	// Simulates UpdateStateFromApi after a Read() that reflects NXRM's now-nil authentication.
	stateModel := repositoryHttpClientModel{Authentication: nil}

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	if assert.NotNil(t, stateModel.Authentication) {
		assert.Equal(t, common.HTTP_AUTH_TYPE_USERNAME, stateModel.Authentication.Type.ValueString())
		assert.Equal(t, "manually-configured-user", stateModel.Authentication.Username.ValueString())
		assert.True(t, stateModel.Authentication.Password.IsNull())
	}

	// The restored value must be a copy, not an alias of the plan's own struct - mutating one
	// must never mutate the other (see GH-489, which was exactly this failure mode elsewhere in
	// this same model).
	stateModel.Authentication.Username = types.StringValue("mutated")
	assert.Equal(t, "manually-configured-user", planModel.Authentication.Username.ValueString())
}

// TestHttpClientMapMissingApiFieldsFromPlanRestoresSecretsWhenPresent covers the pre-existing,
// unchanged behavior: when the plan's authentication has its required secret, only the
// never-returned-by-the-API fields (password/bearerToken/preemptive) are restored onto state's
// own (API-derived) authentication - state's Type/Username, freshly read from NXRM, are left
// alone.
func TestHttpClientMapMissingApiFieldsFromPlanRestoresSecretsWhenPresent(t *testing.T) {
	planModel := repositoryHttpClientModel{
		Authentication: &RepositoryHttpClientAuthenticationModel{
			Type:     types.StringValue(common.HTTP_AUTH_TYPE_USERNAME),
			Username: types.StringValue("user"),
			Password: types.StringValue("pass"),
		},
	}

	stateModel := repositoryHttpClientModel{
		Authentication: &RepositoryHttpClientAuthenticationModel{
			Type:     types.StringValue(common.HTTP_AUTH_TYPE_USERNAME),
			Username: types.StringValue("user"),
			// Password/BearerToken/Preemptive as left by MapFromApiHttpClientAttributes -
			// NXRM never returns them.
		},
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	if assert.NotNil(t, stateModel.Authentication) {
		assert.Equal(t, "pass", stateModel.Authentication.Password.ValueString())
	}
}

// TestHttpClientMapMissingApiFieldsFromPlanNilPlanAuthentication verifies the no-op case: a nil
// plan authentication (user config has no `authentication` block) leaves state untouched.
func TestHttpClientMapMissingApiFieldsFromPlanNilPlanAuthentication(t *testing.T) {
	planModel := repositoryHttpClientModel{Authentication: nil}
	stateModel := repositoryHttpClientModel{Authentication: nil}

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	assert.Nil(t, stateModel.Authentication)
}
