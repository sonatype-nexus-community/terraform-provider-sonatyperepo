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

func TestNugetProxyUpdateStateFromApiUnwrapsFirewallMode(t *testing.T) {
	f := &NugetRepositoryFormatProxy{}
	mode := common.FirewallModeDisabled

	// An existing repository_firewall block in state must be cleared when the server
	// reports the mode as disabled.
	stateModel := f.UpdateStateFromApi(model.RepositoryNugetProxyModel{
		FirewallAuditAndQuarantine: model.NewFirewallAuditAndQuarantineModelWithDefaults(),
	}, ProxyApiResponseWithFirewall{
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

// TestRawProxyUpdateStateFromApiUnwrapsFirewallMode mirrors the tests above for Raw. Unlike
// the other formats, Raw's GetRawProxyRepository always returns a nil FirewallMode (its
// response type has no `firewall` field - see the comment on that function), so this
// exercises the defensive branch that would populate repository_firewall if a future client
// update ever started returning it, without regressing the format's actual fix
// (MapMissingApiFieldsFromPlan - see repository_raw_test.go) if that ever changes.
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
