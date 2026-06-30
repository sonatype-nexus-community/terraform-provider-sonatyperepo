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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// TestYumProxyMapMissingApiFieldsPreservesSigningFromPlan verifies that
// MapMissingApiFieldsFromPlan copies Yum signing fields (keypair and passphrase)
// from the plan into state.  These fields are never returned by the Nexus GET
// API, so they must be carried over from the plan after every create/update;
// failing to do so causes Terraform to report "inconsistent values for
// sensitive attribute" (GitHub issue #436).
func TestYumProxyMapMissingApiFieldsPreservesSigningFromPlan(t *testing.T) {
	// Simulate the state after UpdateStateFromApi: the GET API did not return
	// YumSigning, so m.Yum is nil (first apply) or has stale values.
	stateModel := RepositoryYumProxyModel{}

	planModel := RepositoryYumProxyModel{
		Yum: &yumSigningModel{
			KeyPair:    types.StringValue("my-keypair"),
			Passphrase: types.StringValue("super-secret"),
		},
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	assert.NotNil(t, stateModel.Yum, "Yum should be non-nil after MapMissingApiFieldsFromPlan")
	assert.Equal(t, "my-keypair", stateModel.Yum.KeyPair.ValueString())
	assert.Equal(t, "super-secret", stateModel.Yum.Passphrase.ValueString())
}

// TestYumProxyMapMissingApiFieldsUpdatesSigningFromPlan verifies that an
// existing (non-nil) Yum block in state is replaced with the plan's values.
// This covers the case where a user changes keypair/passphrase in their config.
func TestYumProxyMapMissingApiFieldsUpdatesSigningFromPlan(t *testing.T) {
	stateModel := RepositoryYumProxyModel{
		Yum: &yumSigningModel{
			KeyPair:    types.StringValue("old-keypair"),
			Passphrase: types.StringValue("old-secret"),
		},
	}

	planModel := RepositoryYumProxyModel{
		Yum: &yumSigningModel{
			KeyPair:    types.StringValue("new-keypair"),
			Passphrase: types.StringValue("new-secret"),
		},
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	assert.Equal(t, "new-keypair", stateModel.Yum.KeyPair.ValueString())
	assert.Equal(t, "new-secret", stateModel.Yum.Passphrase.ValueString())
}

// TestYumProxyMapMissingApiFieldsNilYumInPlan verifies that a nil Yum in the
// plan (user removed the yum block) is written to state without a panic.
func TestYumProxyMapMissingApiFieldsNilYumInPlan(t *testing.T) {
	stateModel := RepositoryYumProxyModel{
		Yum: &yumSigningModel{
			KeyPair:    types.StringValue("old-keypair"),
			Passphrase: types.StringValue("old-secret"),
		},
	}

	planModel := RepositoryYumProxyModel{} // Yum is nil

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	assert.Nil(t, stateModel.Yum)
}
