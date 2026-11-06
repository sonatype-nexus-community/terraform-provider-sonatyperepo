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

// samlServiceV397 implements SamlService against NXRM API client V397 (targets NXRM 3.97.0+).
type samlServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func NewSamlServiceV397(client *sonatyperepoV397.APIClient) SamlService {
	return &samlServiceV397{client: client}
}

func (s *samlServiceV397) GetSamlConfiguration(ctx context.Context) (*sonatyperepoV382.SamlConfigurationXO, *http.Response, error) {
	apiV397, httpResponse, err := s.client.SecurityManagementSAMLAPI.ListSecuritySaml(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SamlConfigurationXO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *samlServiceV397) PutSamlConfiguration(ctx context.Context, body sonatyperepoV382.SamlConfigurationXO) (*http.Response, error) {
	var v397Body sonatyperepoV397.SamlConfigurationXO
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.SecurityManagementSAMLAPI.UpdateSecuritySaml(ctx).SamlConfigurationXO(v397Body).Execute()
}

func (s *samlServiceV397) DeleteSamlConfiguration(ctx context.Context) (*http.Response, error) {
	return s.client.SecurityManagementSAMLAPI.DeleteSecuritySaml(ctx).Execute()
}
