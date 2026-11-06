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

package repository

import (
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"terraform-provider-sonatyperepo/internal/provider/repository/format"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func pccsFirewall(enabled, quarantine, pccs bool) *model.FirewallAuditAndQuarantineWithPccsModel {
	fw := model.NewFirewallAuditAndQuarantineWithPccsModelWithDefaults()
	fw.Enabled = types.BoolValue(enabled)
	fw.Quarantine = types.BoolValue(quarantine)
	fw.PccsEnabled = types.BoolValue(pccs)
	return fw
}

// TestValidateFirewallQuarantineAndPccs covers GH-511: quarantine and pccs_enabled are mutually
// exclusive Firewall modes on NXRM 3.94+, but remain independent (and therefore valid together)
// on the Capability-based firewall of earlier versions.
func TestValidateFirewallQuarantineAndPccs(t *testing.T) {
	v393 := common.SystemVersion{Major: 3, Minor: 93, Patch: 99}
	v394 := common.SystemVersion{Major: 3, Minor: 94, Patch: 1}
	v397 := common.SystemVersion{Major: 3, Minor: 97, Patch: 0}

	npm := func(fw *model.FirewallAuditAndQuarantineWithPccsModel) model.RepositoryNpmProxyModel {
		return model.RepositoryNpmProxyModel{FirewallAuditAndQuarantine: fw}
	}
	pypi := func(fw *model.FirewallAuditAndQuarantineWithPccsModel) model.RepositoryPyPiProxyModel {
		return model.RepositoryPyPiProxyModel{FirewallAuditAndQuarantine: fw}
	}

	tests := []struct {
		name    string
		repoFmt format.RepositoryFormat
		repoTyp format.RepositoryType
		version common.SystemVersion
		plan    any
		wantErr bool
	}{
		{"npm both true on 3.94 rejected", &format.NpmRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v394, npm(pccsFirewall(true, true, true)), true},
		{"npm both true on 3.97 rejected", &format.NpmRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v397, npm(pccsFirewall(true, true, true)), true},
		{"pypi both true on 3.97 rejected", &format.PyPiRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v397, pypi(pccsFirewall(true, true, true)), true},
		{"npm both true on 3.93 allowed (capability-based, independent flags)", &format.NpmRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v393, npm(pccsFirewall(true, true, true)), false},
		{"npm pccs only on 3.97 allowed", &format.NpmRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v397, npm(pccsFirewall(true, false, true)), false},
		{"npm quarantine only on 3.97 allowed", &format.NpmRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v397, npm(pccsFirewall(true, true, false)), false},
		{"npm audit only on 3.97 allowed", &format.NpmRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v397, npm(pccsFirewall(true, false, false)), false},
		{"npm firewall disabled ignores flags on 3.97", &format.NpmRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v397, npm(pccsFirewall(false, true, true)), false},
		{"npm no firewall block on 3.97 allowed", &format.NpmRepositoryFormatProxy{}, format.REPO_TYPE_PROXY, v397, npm(nil), false},
		{"npm hosted never validated", &format.NpmRepositoryFormatHosted{}, format.REPO_TYPE_HOSTED, v397, npm(pccsFirewall(true, true, true)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repositoryResource{RepositoryFormat: tt.repoFmt, RepositoryType: tt.repoTyp}
			r.NxrmVersion = tt.version
			got := r.validateFirewallQuarantineAndPccs(tt.plan)
			if tt.wantErr && len(got) == 0 {
				t.Fatalf("expected a validation message, got none")
			}
			if !tt.wantErr && len(got) != 0 {
				t.Fatalf("expected no validation message, got %v", got)
			}
		})
	}
}
