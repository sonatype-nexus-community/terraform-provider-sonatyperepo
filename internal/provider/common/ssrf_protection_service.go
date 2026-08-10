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
	"encoding/json"
	"io"
	"net/http"

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// SsrfProtectionService abstracts SSRF protection configuration API across NXRM API client generations.
type SsrfProtectionService interface {
	GetSsrfProtectionConfiguration(ctx context.Context) (*sonatyperepoV382.SsrfProtectionConfigurationXO, *http.Response, error)
	UpdateSsrfProtectionConfiguration(ctx context.Context, body sonatyperepoV382.SsrfProtectionConfigurationXO) (*http.Response, error)
}

// ssrfProtectionServiceV382 implements SsrfProtectionService against NXRM API client V382 (targets NXRM < 3.94.0).
type ssrfProtectionServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func NewSsrfProtectionServiceV382(client *sonatyperepoV382.APIClient) SsrfProtectionService {
	return &ssrfProtectionServiceV382{client: client}
}

func (s *ssrfProtectionServiceV382) GetSsrfProtectionConfiguration(ctx context.Context) (*sonatyperepoV382.SsrfProtectionConfigurationXO, *http.Response, error) {
	httpResponse, err := s.client.SecurityManagementSSRFProtectionAPI.GetConfiguration(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, httpResponse, err
	}

	var result sonatyperepoV382.SsrfProtectionConfigurationXO
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, httpResponse, err
	}

	return &result, httpResponse, nil
}

func (s *ssrfProtectionServiceV382) UpdateSsrfProtectionConfiguration(ctx context.Context, body sonatyperepoV382.SsrfProtectionConfigurationXO) (*http.Response, error) {
	return s.client.SecurityManagementSSRFProtectionAPI.UpdateConfiguration(ctx).Body(body).Execute()
}

// ssrfProtectionServiceV395 implements SsrfProtectionService against NXRM API client V395 (targets NXRM 3.94.0+).
type ssrfProtectionServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func NewSsrfProtectionServiceV395(client *sonatyperepoV395.APIClient) SsrfProtectionService {
	return &ssrfProtectionServiceV395{client: client}
}

func (s *ssrfProtectionServiceV395) GetSsrfProtectionConfiguration(ctx context.Context) (*sonatyperepoV382.SsrfProtectionConfigurationXO, *http.Response, error) {
	apiV395, httpResponse, err := s.client.SecurityManagementSSRFProtectionAPI.ListSecuritySsrfProtection(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SsrfProtectionConfigurationXO
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *ssrfProtectionServiceV395) UpdateSsrfProtectionConfiguration(ctx context.Context, body sonatyperepoV382.SsrfProtectionConfigurationXO) (*http.Response, error) {
	var v395Body sonatyperepoV395.SsrfProtectionConfigurationXO
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	_, httpResponse, err := s.client.SecurityManagementSSRFProtectionAPI.UpdateSecuritySsrfProtection(ctx).SsrfProtectionConfigurationXO(v395Body).Execute()
	return httpResponse, err
}
