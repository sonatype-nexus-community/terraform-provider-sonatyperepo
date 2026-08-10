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

// SamlService abstracts SAML configuration API across NXRM API client generations.
type SamlService interface {
	GetSamlConfiguration(ctx context.Context) (*sonatyperepoV382.SamlConfigurationXO, *http.Response, error)
	PutSamlConfiguration(ctx context.Context, body sonatyperepoV382.SamlConfigurationXO) (*http.Response, error)
	DeleteSamlConfiguration(ctx context.Context) (*http.Response, error)
}

// samlServiceV382 implements SamlService against NXRM API client V382 (targets NXRM < 3.94.0).
type samlServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func NewSamlServiceV382(client *sonatyperepoV382.APIClient) SamlService {
	return &samlServiceV382{client: client}
}

func (s *samlServiceV382) GetSamlConfiguration(ctx context.Context) (*sonatyperepoV382.SamlConfigurationXO, *http.Response, error) {
	httpResponse, err := s.client.SecurityManagementSAMLAPI.GetSamlConfiguration(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, httpResponse, err
	}

	var result sonatyperepoV382.SamlConfigurationXO
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, httpResponse, err
	}

	return &result, httpResponse, nil
}

func (s *samlServiceV382) PutSamlConfiguration(ctx context.Context, body sonatyperepoV382.SamlConfigurationXO) (*http.Response, error) {
	return s.client.SecurityManagementSAMLAPI.PutSamlConfiguration(ctx).Body(body).Execute()
}

func (s *samlServiceV382) DeleteSamlConfiguration(ctx context.Context) (*http.Response, error) {
	return s.client.SecurityManagementSAMLAPI.DeleteSamlConfiguration(ctx).Execute()
}

// samlServiceV395 implements SamlService against NXRM API client V395 (targets NXRM 3.94.0+).
type samlServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func NewSamlServiceV395(client *sonatyperepoV395.APIClient) SamlService {
	return &samlServiceV395{client: client}
}

func (s *samlServiceV395) GetSamlConfiguration(ctx context.Context) (*sonatyperepoV382.SamlConfigurationXO, *http.Response, error) {
	apiV395, httpResponse, err := s.client.SecurityManagementSAMLAPI.ListSecuritySaml(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SamlConfigurationXO
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *samlServiceV395) PutSamlConfiguration(ctx context.Context, body sonatyperepoV382.SamlConfigurationXO) (*http.Response, error) {
	var v395Body sonatyperepoV395.SamlConfigurationXO
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.SecurityManagementSAMLAPI.UpdateSecuritySaml(ctx).SamlConfigurationXO(v395Body).Execute()
}

func (s *samlServiceV395) DeleteSamlConfiguration(ctx context.Context) (*http.Response, error) {
	return s.client.SecurityManagementSAMLAPI.DeleteSecuritySaml(ctx).Execute()
}
