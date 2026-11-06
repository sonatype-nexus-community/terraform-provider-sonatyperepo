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

	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// oauth2ServiceV397 implements OAuth2Service against NXRM API client V397 (targets NXRM 3.97.0+).
// The OAuth2Service interface is typed against the V395 OAuth2OidcConfigurationXO, which has
// identical fields in V397, so values are bridged across rather than changing the interface.
type oauth2ServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

// NewOAuth2ServiceV397 creates a V397 OAuth2Service adapter.
func NewOAuth2ServiceV397(client *sonatyperepoV397.APIClient) OAuth2Service {
	return &oauth2ServiceV397{client: client}
}

func (s *oauth2ServiceV397) GetOAuth2Configuration(ctx context.Context) (*sonatyperepoV395.OAuth2OidcConfigurationXO, *http.Response, error) {
	apiV397, httpResponse, err := s.client.SecurityOAuth2OIDCAPI.ListSecurityOauth2(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV395.OAuth2OidcConfigurationXO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *oauth2ServiceV397) PutOAuth2Configuration(ctx context.Context, body sonatyperepoV395.OAuth2OidcConfigurationXO) (*http.Response, error) {
	var bodyV397 sonatyperepoV397.OAuth2OidcConfigurationXO
	if err := jsonBridge(body, &bodyV397); err != nil {
		return nil, err
	}
	return s.client.SecurityOAuth2OIDCAPI.UpdateSecurityOauth2(ctx).OAuth2OidcConfigurationXO(bodyV397).Execute()
}

func (s *oauth2ServiceV397) DeleteOAuth2Configuration(ctx context.Context) (*http.Response, error) {
	return s.client.SecurityOAuth2OIDCAPI.DeleteSecurityOauth2(ctx).Execute()
}
