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

// TestRawProxyMapMissingApiFieldsPreservesFirewallFromPlan verifies that
// MapMissingApiFieldsFromPlan copies the repository_firewall block from the plan
// into state. On NXRM 3.94+ the GET response for a Raw proxy repository has no
// `firewall` field, so this is the only place the value can come from after a
// Create/Update; failing to do so leaves state's repository_firewall null while
// the plan has it set, which Terraform reports as "Provider produced inconsistent
// result after apply" (GitHub issue #461).
func TestRawProxyMapMissingApiFieldsPreservesFirewallFromPlan(t *testing.T) {
	// Simulate the state after UpdateStateFromApi: the GET API did not (and
	// cannot) return repository_firewall, so it is still nil.
	stateModel := RepositoryRawProxyModel{}

	planModel := RepositoryRawProxyModel{
		FirewallAuditAndQuarantine: &FirewallAuditAndQuarantineModel{
			Enabled:    types.BoolValue(true),
			Quarantine: types.BoolValue(true),
		},
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine, "repository_firewall should be non-nil after MapMissingApiFieldsFromPlan") {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsUnknown(), "capability_id must be known (null) - inline firewall mode has no capability")
	}
}

// TestRawProxyMapMissingApiFieldsUpdatesFirewallFromPlan verifies that an
// existing (non-nil) repository_firewall block in state is replaced with the
// plan's values, e.g. when a user flips quarantine on for an already-configured
// repository.
func TestRawProxyMapMissingApiFieldsUpdatesFirewallFromPlan(t *testing.T) {
	stateModel := RepositoryRawProxyModel{
		FirewallAuditAndQuarantine: &FirewallAuditAndQuarantineModel{
			Enabled:    types.BoolValue(true),
			Quarantine: types.BoolValue(false),
		},
	}

	planModel := RepositoryRawProxyModel{
		FirewallAuditAndQuarantine: &FirewallAuditAndQuarantineModel{
			Enabled:    types.BoolValue(true),
			Quarantine: types.BoolValue(true),
		},
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	assert.True(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
}

// TestRawProxyMapMissingApiFieldsNilFirewallInPlan verifies that a nil
// repository_firewall in the plan (user removed the block) is written to state
// without a panic.
func TestRawProxyMapMissingApiFieldsNilFirewallInPlan(t *testing.T) {
	stateModel := RepositoryRawProxyModel{
		FirewallAuditAndQuarantine: &FirewallAuditAndQuarantineModel{
			Enabled:    types.BoolValue(true),
			Quarantine: types.BoolValue(true),
		},
	}

	planModel := RepositoryRawProxyModel{} // FirewallAuditAndQuarantine is nil

	stateModel.MapMissingApiFieldsFromPlan(planModel)

	assert.Nil(t, stateModel.FirewallAuditAndQuarantine)
}
