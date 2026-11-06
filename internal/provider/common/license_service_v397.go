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
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// licenseServiceV397 implements LicenseService against NXRM API client V397 (targets NXRM 3.97.0+).
type licenseServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *licenseServiceV397) GetLicenseStatus(ctx context.Context) (*sonatyperepoV382.ApiLicenseDetailsXO, *http.Response, error) {
	v397Response, httpResponse, err := s.client.ProductLicensingAPI.ListSystemLicense(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V397 response to V382 shape
	var v382Response sonatyperepoV382.ApiLicenseDetailsXO
	if err := jsonBridge(v397Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

func (s *licenseServiceV397) RemoveLicense(ctx context.Context) (*http.Response, error) {
	return s.client.ProductLicensingAPI.DeleteSystemLicense(ctx).Execute()
}

// SetLicense cannot currently transmit the license file's bytes to a V397 server: the
// the spec to `type: string, format: binary` and regenerating the V397 client in the
func (s *licenseServiceV397) SetLicense(ctx context.Context, licenseFile *os.File) (*sonatyperepoV382.ApiLicenseDetailsXO, *http.Response, error) {
	_ = licenseFile
	v397Response, httpResponse, err := s.client.ProductLicensingAPI.CreateSystemLicense(ctx).Body(map[string]interface{}{}).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V397 response to V382 shape
	var v382Response sonatyperepoV382.ApiLicenseDetailsXO
	if err := jsonBridge(v397Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

// NewLicenseServiceV397 creates a V397 LicenseService adapter.
func NewLicenseServiceV397(client *sonatyperepoV397.APIClient) LicenseService {
	return &licenseServiceV397{client: client}
}
