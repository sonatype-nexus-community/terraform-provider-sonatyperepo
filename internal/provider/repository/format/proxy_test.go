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
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

type firewallModeFromFlagsTestCase struct {
	description    string
	hasConfig      bool
	enabled        bool
	quarantine     bool
	pccsEnabled    bool
	expectedResult common.FirewallMode
}

func TestFirewallModeFromFlags(t *testing.T) {
	testCases := []firewallModeFromFlagsTestCase{
		{
			"No repository_firewall block",
			false, false, false, false,
			common.FirewallModeDisabled,
		},
		{
			"repository_firewall present but disabled",
			true, false, false, false,
			common.FirewallModeDisabled,
		},
		{
			"Enabled, no quarantine, no pccs -> Audit",
			true, true, false, false,
			common.FirewallModeAudit,
		},
		{
			"Enabled with quarantine -> Quarantine",
			true, true, true, false,
			common.FirewallModeQuarantine,
		},
		{
			"Enabled with pccs -> Pccs",
			true, true, false, true,
			common.FirewallModePccs,
		},
		{
			"Enabled with both quarantine and pccs -> Pccs takes precedence",
			true, true, true, true,
			common.FirewallModePccs,
		},
	}

	for i, testCase := range testCases {
		assert.Equal(
			t,
			testCase.expectedResult,
			FirewallModeFromFlags(testCase.hasConfig, testCase.enabled, testCase.quarantine, testCase.pccsEnabled),
			fmt.Sprintf("%d: %s - FirewallMode not as expected", i, testCase.description),
		)
	}
}

type firewallFlagsFromModeTestCase struct {
	mode               common.FirewallMode
	expectedEnabled    bool
	expectedQuarantine bool
	expectedPccs       bool
}

func TestFirewallFlagsFromMode(t *testing.T) {
	testCases := []firewallFlagsFromModeTestCase{
		{common.FirewallModeDisabled, false, false, false},
		{common.FirewallModeAudit, true, false, false},
		{common.FirewallModeQuarantine, true, true, false},
		{common.FirewallModePccs, true, false, true},
	}

	for i, testCase := range testCases {
		enabled, quarantine, pccsEnabled := FirewallFlagsFromMode(testCase.mode)
		assert.Equal(t, testCase.expectedEnabled, enabled, fmt.Sprintf("%d: %s - enabled not as expected", i, testCase.mode))
		assert.Equal(t, testCase.expectedQuarantine, quarantine, fmt.Sprintf("%d: %s - quarantine not as expected", i, testCase.mode))
		assert.Equal(t, testCase.expectedPccs, pccsEnabled, fmt.Sprintf("%d: %s - pccsEnabled not as expected", i, testCase.mode))
	}
}

type computeFirewallModeTestCase struct {
	description    string
	firewall       *model.FirewallAuditAndQuarantineWithPccsModel
	expectedResult common.FirewallMode
}

func TestComputeFirewallMode(t *testing.T) {
	testCases := []computeFirewallModeTestCase{
		{
			"No repository_firewall block",
			nil,
			common.FirewallModeDisabled,
		},
		{
			"repository_firewall present but disabled",
			&model.FirewallAuditAndQuarantineWithPccsModel{
				FirewallAuditAndQuarantineModel: model.FirewallAuditAndQuarantineModel{
					Enabled:    types.BoolValue(false),
					Quarantine: types.BoolValue(false),
				},
				PccsEnabled: types.BoolValue(false),
			},
			common.FirewallModeDisabled,
		},
		{
			"Enabled, no quarantine, no pccs -> Audit",
			&model.FirewallAuditAndQuarantineWithPccsModel{
				FirewallAuditAndQuarantineModel: model.FirewallAuditAndQuarantineModel{
					Enabled:    types.BoolValue(true),
					Quarantine: types.BoolValue(false),
				},
				PccsEnabled: types.BoolValue(false),
			},
			common.FirewallModeAudit,
		},
		{
			"Enabled with quarantine -> Quarantine",
			&model.FirewallAuditAndQuarantineWithPccsModel{
				FirewallAuditAndQuarantineModel: model.FirewallAuditAndQuarantineModel{
					Enabled:    types.BoolValue(true),
					Quarantine: types.BoolValue(true),
				},
				PccsEnabled: types.BoolValue(false),
			},
			common.FirewallModeQuarantine,
		},
		{
			"Enabled with pccs -> Pccs",
			&model.FirewallAuditAndQuarantineWithPccsModel{
				FirewallAuditAndQuarantineModel: model.FirewallAuditAndQuarantineModel{
					Enabled:    types.BoolValue(true),
					Quarantine: types.BoolValue(false),
				},
				PccsEnabled: types.BoolValue(true),
			},
			common.FirewallModePccs,
		},
	}

	f := &NpmRepositoryFormatProxy{}

	for i, testCase := range testCases {
		stateModel := model.RepositoryNpmProxyModel{
			FirewallAuditAndQuarantine: testCase.firewall,
		}
		assert.Equal(
			t,
			testCase.expectedResult,
			ComputeFirewallMode(f, stateModel),
			fmt.Sprintf("%d: %s - FirewallMode not as expected", i, testCase.description),
		)
	}
}
