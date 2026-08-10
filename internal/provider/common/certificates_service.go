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
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// CertificatesService abstracts the Security Certificates API across NXRM API client generations.
type CertificatesService interface {
	// AddCertificate adds a certificate to the truststore.
	AddCertificate(ctx context.Context, pemBody string) (*sonatyperepoV382.ApiCertificate, *http.Response, error)
	// GetTrustStoreCertificates retrieves all certificates from the truststore.
	GetTrustStoreCertificates(ctx context.Context) ([]sonatyperepoV382.ApiCertificate, *http.Response, error)
	// RemoveCertificate removes a certificate from the truststore by ID.
	RemoveCertificate(ctx context.Context, id string) (*http.Response, error)
}

// certificatesServiceV382 implements CertificatesService against NXRM API client V382 (targets NXRM < 3.94.0).
type certificatesServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *certificatesServiceV382) AddCertificate(ctx context.Context, pemBody string) (*sonatyperepoV382.ApiCertificate, *http.Response, error) {
	return s.client.SecurityCertificatesAPI.AddCertificate(ctx).Body(pemBody).Execute()
}

func (s *certificatesServiceV382) GetTrustStoreCertificates(ctx context.Context) ([]sonatyperepoV382.ApiCertificate, *http.Response, error) {
	return s.client.SecurityCertificatesAPI.GetTrustStoreCertificates(ctx).Execute()
}

func (s *certificatesServiceV382) RemoveCertificate(ctx context.Context, id string) (*http.Response, error) {
	return s.client.SecurityCertificatesAPI.RemoveCertificate(ctx, id).Execute()
}

// certificatesServiceV395 implements CertificatesService against NXRM API client V395 (targets NXRM 3.94.0+).
type certificatesServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *certificatesServiceV395) AddCertificate(ctx context.Context, pemBody string) (*sonatyperepoV382.ApiCertificate, *http.Response, error) {
	v395Response, httpResponse, err := s.client.SecurityCertificatesAPI.CreateSecuritySslTruststore(ctx).Body(pemBody).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 response to V382 shape
	var v382Response sonatyperepoV382.ApiCertificate
	if err := jsonBridge(v395Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

func (s *certificatesServiceV395) GetTrustStoreCertificates(ctx context.Context) ([]sonatyperepoV382.ApiCertificate, *http.Response, error) {
	v395Response, httpResponse, err := s.client.SecurityCertificatesAPI.ListSecuritySslTruststore(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 response to V382 shape
	v382Response := make([]sonatyperepoV382.ApiCertificate, len(v395Response))
	for i, v395Cert := range v395Response {
		if err := jsonBridge(v395Cert, &v382Response[i]); err != nil {
			return nil, httpResponse, err
		}
	}

	return v382Response, httpResponse, nil
}

func (s *certificatesServiceV395) RemoveCertificate(ctx context.Context, id string) (*http.Response, error) {
	return s.client.SecurityCertificatesAPI.DeleteSecuritySslTruststore(ctx, id).Execute()
}

// NewCertificatesServiceV382 creates a V382 CertificatesService adapter.
func NewCertificatesServiceV382(client *sonatyperepoV382.APIClient) CertificatesService {
	return &certificatesServiceV382{client: client}
}

// NewCertificatesServiceV395 creates a V395 CertificatesService adapter.
func NewCertificatesServiceV395(client *sonatyperepoV395.APIClient) CertificatesService {
	return &certificatesServiceV395{client: client}
}
