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
	"net/http"
	"os"

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// LicenseService abstracts the Product Licensing API across NXRM API client generations.
type LicenseService interface {
	// GetLicenseStatus retrieves current license status.
	GetLicenseStatus(ctx context.Context) (*sonatyperepoV382.ApiLicenseDetailsXO, *http.Response, error)
	// RemoveLicense removes the current license.
	RemoveLicense(ctx context.Context) (*http.Response, error)
	// SetLicense uploads a new license file.
	SetLicense(ctx context.Context, licenseFile *os.File) (*sonatyperepoV382.ApiLicenseDetailsXO, *http.Response, error)
}

// licenseServiceV382 implements LicenseService against NXRM API client V382 (targets NXRM < 3.94.0).
type licenseServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *licenseServiceV382) GetLicenseStatus(ctx context.Context) (*sonatyperepoV382.ApiLicenseDetailsXO, *http.Response, error) {
	return s.client.ProductLicensingAPI.GetLicenseStatus(ctx).Execute()
}

func (s *licenseServiceV382) RemoveLicense(ctx context.Context) (*http.Response, error) {
	return s.client.ProductLicensingAPI.RemoveLicense(ctx).Execute()
}

func (s *licenseServiceV382) SetLicense(ctx context.Context, licenseFile *os.File) (*sonatyperepoV382.ApiLicenseDetailsXO, *http.Response, error) {
	return s.client.ProductLicensingAPI.SetLicense(ctx).Body(licenseFile).Execute()
}

// licenseServiceV395 implements LicenseService against NXRM API client V395 (targets NXRM 3.94.0+).
type licenseServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *licenseServiceV395) GetLicenseStatus(ctx context.Context) (*sonatyperepoV382.ApiLicenseDetailsXO, *http.Response, error) {
	v395Response, httpResponse, err := s.client.ProductLicensingAPI.ListSystemLicense(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 response to V382 shape
	var v382Response sonatyperepoV382.ApiLicenseDetailsXO
	if err := jsonBridge(v395Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

func (s *licenseServiceV395) RemoveLicense(ctx context.Context) (*http.Response, error) {
	return s.client.ProductLicensingAPI.DeleteSystemLicense(ctx).Execute()
}

// SetLicense cannot currently transmit the license file's bytes to a V395 server: the
// upstream OpenAPI spec (nexus-repo-api-client, POST /v1/system/license) declares the
// application/octet-stream request body as `schema: type: object`, so the Go generator
// emits CreateSystemLicense's body param as map[string]interface{} rather than a binary
// payload (*os.File/[]byte). The generated client's setBody() only serializes
// io.Reader/*os.File/[]byte/string or JSON-content-typed bodies, none of which match here,
// so any call to CreateSystemLicenseExecute fails with "invalid body type
// application/octet-stream" regardless of what is passed. Fixing this requires correcting
// the spec to `type: string, format: binary` and regenerating the V395 client in the
// sibling nexus-repo-api-client repo; licenseFile is accepted here to satisfy the
// LicenseService interface and is intentionally unused until that upstream fix lands.
func (s *licenseServiceV395) SetLicense(ctx context.Context, licenseFile *os.File) (*sonatyperepoV382.ApiLicenseDetailsXO, *http.Response, error) {
	_ = licenseFile
	v395Response, httpResponse, err := s.client.ProductLicensingAPI.CreateSystemLicense(ctx).Body(map[string]interface{}{}).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 response to V382 shape
	var v382Response sonatyperepoV382.ApiLicenseDetailsXO
	if err := jsonBridge(v395Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

// NewLicenseServiceV382 creates a V382 LicenseService adapter.
func NewLicenseServiceV382(client *sonatyperepoV382.APIClient) LicenseService {
	return &licenseServiceV382{client: client}
}

// NewLicenseServiceV395 creates a V395 LicenseService adapter.
func NewLicenseServiceV395(client *sonatyperepoV395.APIClient) LicenseService {
	return &licenseServiceV395{client: client}
}
