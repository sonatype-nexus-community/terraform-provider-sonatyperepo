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

package common

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

const licenseTestContent = "license-bytes-GH-545"

// Issue 545: SetLicense on the V395 and V397 clients must transmit the license file's bytes as
// application/octet-stream (previously it failed client-side with "invalid body type").
func TestLicenseServiceSetLicenseSendsFileBytes(t *testing.T) {
	var gotMethod, gotPath, gotContentType, gotBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		gotMethod, gotPath, gotContentType, gotBody = r.Method, r.URL.Path, r.Header.Get("Content-Type"), string(body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"contactCompany":"Acme","licenseType":"PRO"}`))
	}))
	defer server.Close()

	configurationV395 := sonatyperepoV395.NewConfiguration()
	configurationV395.Servers = sonatyperepoV395.ServerConfigurations{{URL: server.URL}}
	configurationV397 := sonatyperepoV397.NewConfiguration()
	configurationV397.Servers = sonatyperepoV397.ServerConfigurations{{URL: server.URL}}

	services := map[string]LicenseService{
		"V395": NewLicenseServiceV395(sonatyperepoV395.NewAPIClient(configurationV395)),
		"V397": NewLicenseServiceV397(sonatyperepoV397.NewAPIClient(configurationV397)),
	}

	for name, service := range services {
		t.Run(name, func(t *testing.T) {
			gotMethod, gotPath, gotContentType, gotBody = "", "", "", ""

			licenseFile, err := os.CreateTemp(t.TempDir(), "license")
			require.NoError(t, err)
			_, err = licenseFile.WriteString(licenseTestContent)
			require.NoError(t, err)
			_, err = licenseFile.Seek(0, io.SeekStart)
			require.NoError(t, err)
			defer func() { _ = licenseFile.Close() }()

			details, httpResponse, err := service.SetLicense(context.Background(), licenseFile)

			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
			assert.Equal(t, http.MethodPost, gotMethod)
			assert.Equal(t, "/v1/system/license", gotPath)
			assert.Equal(t, "application/octet-stream", gotContentType)
			assert.Equal(t, licenseTestContent, gotBody)
			require.NotNil(t, details)
			assert.Equal(t, "Acme", details.GetContactCompany())
		})
	}
}
