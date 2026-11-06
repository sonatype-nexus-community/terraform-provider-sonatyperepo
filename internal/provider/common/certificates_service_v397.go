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

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// certificatesServiceV397 implements CertificatesService against NXRM API client V397 (targets NXRM 3.97.0+).
type certificatesServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *certificatesServiceV397) AddCertificate(ctx context.Context, pemBody string) (*sonatyperepoV382.ApiCertificate, *http.Response, error) {
	v397Response, httpResponse, err := s.client.SecurityCertificatesAPI.CreateSecuritySslTruststore(ctx).Body(pemBody).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V397 response to V382 shape
	var v382Response sonatyperepoV382.ApiCertificate
	if err := jsonBridge(v397Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

func (s *certificatesServiceV397) GetTrustStoreCertificates(ctx context.Context) ([]sonatyperepoV382.ApiCertificate, *http.Response, error) {
	v397Response, httpResponse, err := s.client.SecurityCertificatesAPI.ListSecuritySslTruststore(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V397 response to V382 shape
	v382Response := make([]sonatyperepoV382.ApiCertificate, len(v397Response))
	for i, v397Cert := range v397Response {
		if err := jsonBridge(v397Cert, &v382Response[i]); err != nil {
			return nil, httpResponse, err
		}
	}

	return v382Response, httpResponse, nil
}

func (s *certificatesServiceV397) RemoveCertificate(ctx context.Context, id string) (*http.Response, error) {
	return s.client.SecurityCertificatesAPI.DeleteSecuritySslTruststore(ctx, id).Execute()
}

// NewCertificatesServiceV397 creates a V397 CertificatesService adapter.
func NewCertificatesServiceV397(client *sonatyperepoV397.APIClient) CertificatesService {
	return &certificatesServiceV397{client: client}
}
