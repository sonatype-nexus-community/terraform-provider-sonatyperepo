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

package format

import (
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// TestConanProxyUpdateStateFromApiUnwrapsFirewallMode, and its siblings below (including the
// Yum and CocoaPods variants further down - see GH-464), verify that UpdateStateFromApi
// correctly unwraps a ProxyApiResponseWithFirewall and populates repository_firewall in
// state. Prior to wiring these formats into the NXRM 3.94+ inline firewall path,
// DoReadRequest never wrapped its response, so UpdateStateFromApi silently dropped any
// firewall configuration - the same class of bug reported for Raw in
// https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/461.
func TestConanProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &ConanRepositoryFormatProxy{}
	mode := common.FirewallModeQuarantine

	stateModel := f.UpdateStateFromApi(model.RepositoryConanProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.ConanProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryConanProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsNull())
	}
}

func TestCondaProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &CondaRepositoryFormatProxy{}
	mode := common.FirewallModeAudit

	stateModel := f.UpdateStateFromApi(model.RepositoryCondaProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.SimpleApiProxyRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryCondaProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
	}
}

func TestMavenProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &MavenRepositoryFormatProxy{}
	mode := common.FirewallModeQuarantine

	stateModel := f.UpdateStateFromApi(model.RepositoryMavenProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.MavenProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryMavenProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
	}
}

// TestNugetProxyUpdateStateFromApiPreservesFirewallBlockWhenDisabledButConfigured verifies
// the GH-469 fix: repository_firewall is Optional-but-not-Computed, so once the practitioner
// has configured it (even with enabled = false), Terraform's plan is always a non-null
// object. Nulling the whole attribute out just because the server resolves the mode as
// disabled - the previous behavior - made apply fail with "Provider produced inconsistent
// result after apply". The block must be preserved, populated with the disabled flags,
// instead.
func TestNugetProxyUpdateStateFromApiPreservesFirewallBlockWhenDisabledButConfigured(t *testing.T) {
	f := &NugetRepositoryFormatProxy{}
	mode := common.FirewallModeDisabled

	stateModel := f.UpdateStateFromApi(model.RepositoryNugetProxyModel{
		FirewallAuditAndQuarantine: model.NewFirewallAuditAndQuarantineModelWithDefaults(),
	}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.NugetProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryNugetProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsNull())
	}
}

// TestNugetProxyUpdateStateFromApiClearsFirewallWhenDisabledAndUnconfigured is the case that
// legitimately clears repository_firewall: it was never configured to begin with, so the
// plan already expects null and there's nothing to preserve.
func TestNugetProxyUpdateStateFromApiClearsFirewallWhenDisabledAndUnconfigured(t *testing.T) {
	f := &NugetRepositoryFormatProxy{}
	mode := common.FirewallModeDisabled

	stateModel := f.UpdateStateFromApi(model.RepositoryNugetProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.NugetProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryNugetProxyModel)

	assert.Nil(t, stateModel.FirewallAuditAndQuarantine)
}

func TestRubyGemsProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &RubyGemsRepositoryFormatProxy{}
	mode := common.FirewallModeAudit

	stateModel := f.UpdateStateFromApi(model.RepositorRubyGemsProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.SimpleApiProxyRepository{},
		FirewallMode: &mode,
	}).(model.RepositorRubyGemsProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
	}
}

func TestYumProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &YumRepositoryFormatProxy{}
	mode := common.FirewallModeAudit

	stateModel := f.UpdateStateFromApi(model.RepositoryYumProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.YumProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryYumProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsNull())
	}
}

func TestCocoaPodsProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &CocoaPodsRepositoryFormatProxy{}
	mode := common.FirewallModeQuarantine

	stateModel := f.UpdateStateFromApi(model.RepositoryCocoaPodsProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.SimpleApiProxyRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryCocoaPodsProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
	}
}

// TestPyPiProxyUpdateStateFromApiUnwrapsFirewallMode verifies the same unwrapping for PyPI -
// see GH-466. Unlike its siblings above, PyPI's GetPypiProxyRepository previously always
// returned a nil FirewallMode because the vendored client's PyPiProxyApiRepository response
// type was missing the `firewall` field entirely (a client-generation staleness gap fixed
// upstream in nexus-repo-api-client-go v395.95.3, not a real NXRM API limitation - the live
// server's own OpenAPI schema always declared it). This also exercises PccsEnabled, since
// PyPI is the first of these formats under test that supports repository_firewall PCCS.
func TestPyPiProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &PyPiRepositoryFormatProxy{}
	mode := common.FirewallModePccs

	stateModel := f.UpdateStateFromApi(model.RepositoryPyPiProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.PyPiProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryPyPiProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.PccsEnabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsNull())
	}
}

// TestPyPiProxyUpdateStateFromApiPreservesFirewallBlockWhenDisabledButConfigured mirrors
// TestNugetProxyUpdateStateFromApiPreservesFirewallBlockWhenDisabledButConfigured above for
// PyPI - see GH-469.
func TestPyPiProxyUpdateStateFromApiPreservesFirewallBlockWhenDisabledButConfigured(t *testing.T) {
	f := &PyPiRepositoryFormatProxy{}
	mode := common.FirewallModeDisabled

	stateModel := f.UpdateStateFromApi(model.RepositoryPyPiProxyModel{
		FirewallAuditAndQuarantine: model.NewFirewallAuditAndQuarantineWithPccsModelWithDefaults(),
	}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.PyPiProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryPyPiProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.PccsEnabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsNull())
	}
}

// TestPyPiProxyUpdateStateFromApiClearsFirewallWhenDisabledAndUnconfigured mirrors
// TestNugetProxyUpdateStateFromApiClearsFirewallWhenDisabledAndUnconfigured above for PyPI.
func TestPyPiProxyUpdateStateFromApiClearsFirewallWhenDisabledAndUnconfigured(t *testing.T) {
	f := &PyPiRepositoryFormatProxy{}
	mode := common.FirewallModeDisabled

	stateModel := f.UpdateStateFromApi(model.RepositoryPyPiProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.PyPiProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryPyPiProxyModel)

	assert.Nil(t, stateModel.FirewallAuditAndQuarantine)
}

// TestRawProxyUpdateStateFromApiUnwrapsFirewallMode mirrors the tests above for Raw. Unlike
// the other formats, Raw's GetRawProxyRepository always returns a nil FirewallMode (its
// response type has no `firewall` field - see the comment on that function), so this
// exercises the defensive branch that would populate repository_firewall if a future client
// update ever started returning it, without regressing the format's actual fix
// (MapMissingApiFieldsFromPlan - see repository_raw_test.go) if that ever changes.
// TestOciProxyUpdateStateFromApiUnwrapsFirewallMode verifies the same unwrapping for OCI.
func TestOciProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &OciRepositoryFormatProxy{}
	mode := common.FirewallModeQuarantine

	stateModel := f.UpdateStateFromApi(model.RepositoryOciProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   common.OciProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryOciProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsNull())
	}
}

// TestOciProxyUpdateStateFromApiPreservesFirewallBlockWhenDisabledButConfigured mirrors
// TestNugetProxyUpdateStateFromApiPreservesFirewallBlockWhenDisabledButConfigured above for
// OCI - see GH-469.
func TestOciProxyUpdateStateFromApiPreservesFirewallBlockWhenDisabledButConfigured(t *testing.T) {
	f := &OciRepositoryFormatProxy{}
	mode := common.FirewallModeDisabled

	stateModel := f.UpdateStateFromApi(model.RepositoryOciProxyModel{
		FirewallAuditAndQuarantine: model.NewFirewallAuditAndQuarantineModelWithDefaults(),
	}, ProxyApiResponseWithFirewall{
		Repository:   common.OciProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryOciProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsNull())
	}
}

// TestOciProxyUpdateStateFromApiClearsFirewallWhenDisabledAndUnconfigured is the case that
// legitimately clears repository_firewall: it was never configured to begin with, so the
// plan already expects null and there's nothing to preserve.
func TestOciProxyUpdateStateFromApiClearsFirewallWhenDisabledAndUnconfigured(t *testing.T) {
	f := &OciRepositoryFormatProxy{}
	mode := common.FirewallModeDisabled

	stateModel := f.UpdateStateFromApi(model.RepositoryOciProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   common.OciProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryOciProxyModel)

	assert.Nil(t, stateModel.FirewallAuditAndQuarantine)
}

func TestRawProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &RawRepositoryFormatProxy{}
	mode := common.FirewallModeQuarantine

	stateModel := f.UpdateStateFromApi(model.RepositoryRawProxyModel{}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.RawProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryRawProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
	}
}

// TestRawProxyUpdateStateFromApiPreservesExistingFirewallWhenModeUnknown covers Raw's actual
// runtime path: GetRawProxyRepository always wraps with a nil FirewallMode, so
// UpdateStateFromApi must leave repository_firewall exactly as it came in on state, rather
// than wiping it - it's MapMissingApiFieldsFromPlan's job (called right after, in both
// Create and Update) to correct it from the plan.
func TestRawProxyUpdateStateFromApiPreservesExistingFirewallWhenModeUnknown(t *testing.T) {
	f := &RawRepositoryFormatProxy{}

	stateModel := f.UpdateStateFromApi(model.RepositoryRawProxyModel{
		FirewallAuditAndQuarantine: &model.FirewallAuditAndQuarantineModel{
			Enabled:    types.BoolValue(true),
			Quarantine: types.BoolValue(true),
		},
	}, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.RawProxyApiRepository{},
		FirewallMode: nil,
	}).(model.RepositoryRawProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
	}
}

// TestResolveFirewallBlockFlags exercises the decision table behind the GH-469 fix directly:
// repository_firewall must only be cleared (keep=false) when the server resolves the mode as
// Disabled AND it was never configured to begin with. Every other combination must keep the
// block, since either the practitioner configured it (so Terraform's plan already expects a
// non-null object) or the resolved mode is one of Audit/Quarantine/Pccs (so there's a real
// firewall configuration to report, configured or not - e.g. on import/drift).
func TestResolveFirewallBlockFlags(t *testing.T) {
	tests := []struct {
		name         string
		hadConfig    bool
		mode         common.FirewallMode
		expectedKeep bool
	}{
		{"disabled and unconfigured clears the block", false, common.FirewallModeDisabled, false},
		{"disabled but configured keeps the block", true, common.FirewallModeDisabled, true},
		{"audit always keeps the block", false, common.FirewallModeAudit, true},
		{"quarantine always keeps the block", false, common.FirewallModeQuarantine, true},
		{"pccs always keeps the block", false, common.FirewallModePccs, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keep, enabled, quarantine, pccsEnabled := ResolveFirewallBlockFlags(tt.hadConfig, tt.mode)
			assert.Equal(t, tt.expectedKeep, keep)

			expectedEnabled, expectedQuarantine, expectedPccsEnabled := FirewallFlagsFromMode(tt.mode)
			assert.Equal(t, expectedEnabled, enabled)
			assert.Equal(t, expectedQuarantine, quarantine)
			assert.Equal(t, expectedPccsEnabled, pccsEnabled)
		})
	}
}

// TestReconcileFirewallBlockWithPlan covers the second half of the GH-469 fix in both
// directions. Update() feeds UpdateStateFromApi the PRIOR state rather than the new plan (see
// repository_common.go), so ResolveFirewallBlockFlags there can be wrong in either direction
// once the real plan is known:
//   - adding repository_firewall for the first time with enabled = false in the same apply
//     isn't visible there - state ends up nil even though the plan expects an object.
//   - removing repository_firewall entirely (rather than setting enabled = false) leaves
//     ResolveFirewallBlockFlags looking at the prior state's now-stale non-nil block and
//     concluding it's still configured - state ends up non-nil even though the plan expects
//     null.
//
// UpdateStateFromPlanForNonApiFields (which has both the plan and the post-read state) is
// what catches both of these up.
func TestReconcileFirewallBlockWithPlan(t *testing.T) {
	t.Run("backfills from plan when state is nil but plan configured the block", func(t *testing.T) {
		planFirewall := &model.FirewallAuditAndQuarantineModel{
			CapabilityId: types.StringValue("should-be-cleared"),
			Enabled:      types.BoolValue(false),
			Quarantine:   types.BoolValue(false),
		}

		result := ReconcileFirewallBlockWithPlan(nil, planFirewall)

		if assert.NotNil(t, result) {
			assert.False(t, result.Enabled.ValueBool())
			assert.False(t, result.Quarantine.ValueBool())
			assert.True(t, result.CapabilityId.IsNull())
		}
	})

	t.Run("leaves state nil when plan never configured the block either", func(t *testing.T) {
		assert.Nil(t, ReconcileFirewallBlockWithPlan(nil, nil))
	})

	t.Run("never overwrites a non-nil state value the plan still wants", func(t *testing.T) {
		stateFirewall := &model.FirewallAuditAndQuarantineModel{Enabled: types.BoolValue(true)}
		planFirewall := &model.FirewallAuditAndQuarantineModel{Enabled: types.BoolValue(false)}

		result := ReconcileFirewallBlockWithPlan(stateFirewall, planFirewall)

		assert.Same(t, stateFirewall, result)
	})

	t.Run("clears a stale non-nil state value when the plan removed the block", func(t *testing.T) {
		stateFirewall := &model.FirewallAuditAndQuarantineModel{Enabled: types.BoolValue(true)}

		result := ReconcileFirewallBlockWithPlan(stateFirewall, nil)

		assert.Nil(t, result)
	})
}

// TestReconcileFirewallBlockWithPccsPlan mirrors TestReconcileFirewallBlockWithPlan for the
// PCCS-capable variant of the block (NPM, PyPI).
func TestReconcileFirewallBlockWithPccsPlan(t *testing.T) {
	t.Run("backfills from plan when state is nil but plan configured the block", func(t *testing.T) {
		planFirewall := &model.FirewallAuditAndQuarantineWithPccsModel{
			FirewallAuditAndQuarantineModel: model.FirewallAuditAndQuarantineModel{
				CapabilityId: types.StringValue("should-be-cleared"),
				Enabled:      types.BoolValue(false),
				Quarantine:   types.BoolValue(false),
			},
			PccsEnabled: types.BoolValue(false),
		}

		result := ReconcileFirewallBlockWithPccsPlan(nil, planFirewall)

		if assert.NotNil(t, result) {
			assert.False(t, result.Enabled.ValueBool())
			assert.False(t, result.Quarantine.ValueBool())
			assert.False(t, result.PccsEnabled.ValueBool())
			assert.True(t, result.CapabilityId.IsNull())
		}
	})

	t.Run("leaves state nil when plan never configured the block either", func(t *testing.T) {
		assert.Nil(t, ReconcileFirewallBlockWithPccsPlan(nil, nil))
	})

	t.Run("clears a stale non-nil state value when the plan removed the block", func(t *testing.T) {
		stateFirewall := &model.FirewallAuditAndQuarantineWithPccsModel{
			FirewallAuditAndQuarantineModel: model.FirewallAuditAndQuarantineModel{Enabled: types.BoolValue(true)},
		}

		result := ReconcileFirewallBlockWithPccsPlan(stateFirewall, nil)

		assert.Nil(t, result)
	})
}

// TestNpmProxyUpdateStateFromPlanForNonApiFieldsBackfillsNewlyConfiguredFirewall reproduces
// the exact GH-469 scenario end-to-end for the format it was originally filed against:
// adding `repository_firewall = { enabled = false }` for the first time in the same apply
// that changes other attributes. Update()'s call into UpdateStateFromApi is fed the prior
// state (which has no repository_firewall block), so the state produced there is nil even
// though the new plan configured it - UpdateStateFromPlanForNonApiFields must reconcile that
// before it reaches Terraform's post-apply consistency check.
func TestNpmProxyUpdateStateFromPlanForNonApiFieldsBackfillsNewlyConfiguredFirewall(t *testing.T) {
	f := &NpmRepositoryFormatProxy{}

	planModel := model.RepositoryNpmProxyModel{
		FirewallAuditAndQuarantine: &model.FirewallAuditAndQuarantineWithPccsModel{
			FirewallAuditAndQuarantineModel: model.FirewallAuditAndQuarantineModel{
				Enabled:    types.BoolValue(false),
				Quarantine: types.BoolValue(false),
			},
			PccsEnabled: types.BoolValue(false),
		},
	}
	// Simulates the state UpdateStateFromApi would have produced from prior (unconfigured) state.
	stateAfterApiUpdate := model.RepositoryNpmProxyModel{}

	stateModel := f.UpdateStateFromPlanForNonApiFields(planModel, stateAfterApiUpdate).(model.RepositoryNpmProxyModel)

	if assert.NotNil(t, stateModel.FirewallAuditAndQuarantine) {
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool())
		assert.False(t, stateModel.FirewallAuditAndQuarantine.PccsEnabled.ValueBool())
		assert.True(t, stateModel.FirewallAuditAndQuarantine.CapabilityId.IsNull())
	}
}

// TestNugetProxyUpdateFlowClearsFirewallWhenBlockRemovedAfterBeingEnabled reproduces, end to
// end through the exact two calls repository_common.go's Update() makes
// (UpdateStateFromApi then UpdateStateFromPlanForNonApiFields), the opposite-direction
// regression this fix's first version introduced: disabling the firewall by removing the
// `repository_firewall` block entirely (rather than setting enabled = false) after it was
// previously enabled. Caught by reviewing this fix against
// TestAccRepositoryGenericProxyFirewallToggle's live step 2 -> step 3 transition before that
// acceptance test could be run against a real IQ Server.
func TestNugetProxyUpdateFlowClearsFirewallWhenBlockRemovedAfterBeingEnabled(t *testing.T) {
	f := &NugetRepositoryFormatProxy{}

	// Prior state: repository_firewall was enabled+quarantine from a previous apply.
	priorState := model.RepositoryNugetProxyModel{
		FirewallAuditAndQuarantine: &model.FirewallAuditAndQuarantineModel{
			Enabled:    types.BoolValue(true),
			Quarantine: types.BoolValue(true),
		},
	}
	// New plan: block removed entirely (disable by deleting the config block, not enabled = false).
	planModel := model.RepositoryNugetProxyModel{}

	mode := common.FirewallModeDisabled
	stateAfterApi := f.UpdateStateFromApi(priorState, ProxyApiResponseWithFirewall{
		Repository:   sonatyperepo.NugetProxyApiRepository{},
		FirewallMode: &mode,
	}).(model.RepositoryNugetProxyModel)

	finalState := f.UpdateStateFromPlanForNonApiFields(planModel, stateAfterApi).(model.RepositoryNugetProxyModel)

	assert.Nil(t, finalState.FirewallAuditAndQuarantine)
}
