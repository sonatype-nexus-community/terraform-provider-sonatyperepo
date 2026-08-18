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
	"fmt"
	"net/http"

	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// errOAuth2Unsupported is returned by every OAuth2Service method when the connected NXRM
// server predates 3.94.0. Unlike every other Service in this package, the OAuth2/OIDC
// configuration API has no equivalent in the V382 client generation at all -- there is no
// bridging to do, only a clear error to surface.
var errOAuth2Unsupported = fmt.Errorf("OAuth2/OIDC configuration API requires Nexus Repository Manager 3.94.0 or later")

// OAuth2Service abstracts the OAuth2/OIDC configuration API. This API was only added to the
// V395 client generation (NXRM 3.94.0+); it does not exist in V382 at all.
type OAuth2Service interface {
	// GetOAuth2Configuration retrieves the current OAuth2/OIDC configuration.
	GetOAuth2Configuration(ctx context.Context) (*sonatyperepoV395.OAuth2OidcConfigurationXO, *http.Response, error)
	// PutOAuth2Configuration creates or updates the OAuth2/OIDC configuration.
	PutOAuth2Configuration(ctx context.Context, body sonatyperepoV395.OAuth2OidcConfigurationXO) (*http.Response, error)
	// DeleteOAuth2Configuration removes the OAuth2/OIDC configuration.
	DeleteOAuth2Configuration(ctx context.Context) (*http.Response, error)
}

// oauth2ServiceUnsupported implements OAuth2Service for NXRM servers older than 3.94.0, where
// the OAuth2/OIDC configuration REST API does not exist.
type oauth2ServiceUnsupported struct{}

// NewOAuth2ServiceUnsupported creates an OAuth2Service adapter for NXRM servers older than
// 3.94.0, which have no OAuth2/OIDC configuration API to call.
func NewOAuth2ServiceUnsupported() OAuth2Service {
	return &oauth2ServiceUnsupported{}
}

func (s *oauth2ServiceUnsupported) GetOAuth2Configuration(ctx context.Context) (*sonatyperepoV395.OAuth2OidcConfigurationXO, *http.Response, error) {
	return nil, nil, errOAuth2Unsupported
}

func (s *oauth2ServiceUnsupported) PutOAuth2Configuration(ctx context.Context, body sonatyperepoV395.OAuth2OidcConfigurationXO) (*http.Response, error) {
	return nil, errOAuth2Unsupported
}

func (s *oauth2ServiceUnsupported) DeleteOAuth2Configuration(ctx context.Context) (*http.Response, error) {
	return nil, errOAuth2Unsupported
}

// oauth2ServiceV395 implements OAuth2Service against NXRM API client V395 (targets NXRM 3.94.0+).
type oauth2ServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

// NewOAuth2ServiceV395 creates a V395 OAuth2Service adapter.
func NewOAuth2ServiceV395(client *sonatyperepoV395.APIClient) OAuth2Service {
	return &oauth2ServiceV395{client: client}
}

func (s *oauth2ServiceV395) GetOAuth2Configuration(ctx context.Context) (*sonatyperepoV395.OAuth2OidcConfigurationXO, *http.Response, error) {
	return s.client.SecurityOAuth2OIDCAPI.ListSecurityOauth2(ctx).Execute()
}

func (s *oauth2ServiceV395) PutOAuth2Configuration(ctx context.Context, body sonatyperepoV395.OAuth2OidcConfigurationXO) (*http.Response, error) {
	return s.client.SecurityOAuth2OIDCAPI.UpdateSecurityOauth2(ctx).OAuth2OidcConfigurationXO(body).Execute()
}

func (s *oauth2ServiceV395) DeleteOAuth2Configuration(ctx context.Context) (*http.Response, error) {
	return s.client.SecurityOAuth2OIDCAPI.DeleteSecurityOauth2(ctx).Execute()
}
