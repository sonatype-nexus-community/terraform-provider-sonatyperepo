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

package common_test

import (
	"fmt"
	"testing"

	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/stretchr/testify/assert"

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// TestNewServicesClientGenerationByVersion guards the NXRM version thresholds that select
// which generated API client generation serves a connected server. A wrong threshold would
// silently route a server to a client whose strict JSON decoding rejects its responses.
func TestNewServicesClientGenerationByVersion(t *testing.T) {
	tests := []struct {
		header   string
		expected string
	}{
		{"Nexus/3.82.0-08 (PRO)", "V382"},
		{"Nexus/3.93.0-06 (PRO)", "V382"},
		{"Nexus/3.93.2-99 (OSS)", "V382"},
		{"Nexus/3.94.0-01 (PRO)", "V395"},
		{"Nexus/3.94.1-06 (PRO)", "V395"},
		{"Nexus/3.95.1-01 (PRO)", "V395"},
		{"Nexus/3.96.0-09 (PRO)", "V395"},
		{"Nexus/3.96.3-99 (OSS)", "V395"},
		{"Nexus/3.97.0-13 (PRO)", "V397"},
		{"Nexus/3.97.0-13 (OSS)", "V397"},
		{"Nexus/3.98.0-01 (PRO)", "V397"},
	}

	for _, tc := range tests {
		t.Run(tc.header, func(t *testing.T) {
			version := common.ParseServerHeaderToVersion(tc.header)
			svc := common.NewServices(
				version,
				sonatyperepoV382.NewAPIClient(sonatyperepoV382.NewConfiguration()),
				sonatyperepoV395.NewAPIClient(sonatyperepoV395.NewConfiguration()),
				sonatyperepoV397.NewAPIClient(sonatyperepoV397.NewConfiguration()),
			)

			// Every service in the selected set must come from the same client generation.
			for name, impl := range map[string]any{
				"Repository": svc.Repository,
				"Status":     svc.Status,
				"Task":       svc.Task,
				"BlobStore":  svc.BlobStore,
				"User":       svc.User,
			} {
				got := fmt.Sprintf("%T", impl)
				assert.Contains(t, got, tc.expected, "%s service for %s", name, tc.header)
			}

			// The OAuth2/OIDC API does not exist before 3.94.0, so older servers get an
			// "unsupported" adapter rather than a V382 one.
			if tc.expected == "V382" {
				assert.Contains(t, fmt.Sprintf("%T", svc.OAuth2), "Unsupported", "OAuth2 service for %s", tc.header)
			} else {
				assert.Contains(t, fmt.Sprintf("%T", svc.OAuth2), tc.expected, "OAuth2 service for %s", tc.header)
			}
		})
	}
}
