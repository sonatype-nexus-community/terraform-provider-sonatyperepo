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

// UserTokensService abstracts the User Tokens API across NXRM API client generations.
type UserTokensService interface {
	// ServiceStatus retrieves current user token configuration.
	ServiceStatus(ctx context.Context) (*sonatyperepoV382.UserTokensApiModel, *http.Response, error)
	// SetServiceStatus updates user token configuration.
	SetServiceStatus(ctx context.Context, body sonatyperepoV382.UserTokensApiModel) (*sonatyperepoV382.UserTokensApiModel, *http.Response, error)
}

// userTokensServiceV382 implements UserTokensService against NXRM API client V382 (targets NXRM < 3.94.0).
type userTokensServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *userTokensServiceV382) ServiceStatus(ctx context.Context) (*sonatyperepoV382.UserTokensApiModel, *http.Response, error) {
	return s.client.SecurityManagementUserTokensAPI.ServiceStatus(ctx).Execute()
}

func (s *userTokensServiceV382) SetServiceStatus(ctx context.Context, body sonatyperepoV382.UserTokensApiModel) (*sonatyperepoV382.UserTokensApiModel, *http.Response, error) {
	return s.client.SecurityManagementUserTokensAPI.SetServiceStatus(ctx).Body(body).Execute()
}

// userTokensServiceV395 implements UserTokensService against NXRM API client V395 (targets NXRM 3.94.0+).
type userTokensServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *userTokensServiceV395) ServiceStatus(ctx context.Context) (*sonatyperepoV382.UserTokensApiModel, *http.Response, error) {
	v395Response, httpResponse, err := s.client.SecurityManagementUserTokensAPI.ListSecurityUserTokens(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 response to V382 shape
	var v382Response sonatyperepoV382.UserTokensApiModel
	if err := jsonBridge(v395Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

func (s *userTokensServiceV395) SetServiceStatus(ctx context.Context, body sonatyperepoV382.UserTokensApiModel) (*sonatyperepoV382.UserTokensApiModel, *http.Response, error) {
	// Bridge V382 request to V395 shape
	var v395Body sonatyperepoV395.UserTokensApiModel
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, nil, err
	}

	v395Response, httpResponse, err := s.client.SecurityManagementUserTokensAPI.UpdateSecurityUserTokens(ctx).UserTokensApiModel(v395Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 response to V382 shape
	var v382Response sonatyperepoV382.UserTokensApiModel
	if err := jsonBridge(v395Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

// NewUserTokensServiceV382 creates a V382 UserTokensService adapter.
func NewUserTokensServiceV382(client *sonatyperepoV382.APIClient) UserTokensService {
	return &userTokensServiceV382{client: client}
}

// NewUserTokensServiceV395 creates a V395 UserTokensService adapter.
func NewUserTokensServiceV395(client *sonatyperepoV395.APIClient) UserTokensService {
	return &userTokensServiceV395{client: client}
}
