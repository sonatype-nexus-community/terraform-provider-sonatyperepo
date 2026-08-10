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

// AnonymousAccessService abstracts anonymous access settings API across NXRM API client generations.
type AnonymousAccessService interface {
	GetAnonymousAccessSettings(ctx context.Context) (*sonatyperepoV382.AnonymousAccessSettingsXO, *http.Response, error)
	UpdateAnonymousAccessSettings(ctx context.Context, body sonatyperepoV382.AnonymousAccessSettingsXO) (*sonatyperepoV382.AnonymousAccessSettingsXO, *http.Response, error)
}

// anonymousAccessServiceV382 implements AnonymousAccessService against NXRM API client V382 (targets NXRM < 3.94.0).
type anonymousAccessServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func NewAnonymousAccessServiceV382(client *sonatyperepoV382.APIClient) AnonymousAccessService {
	return &anonymousAccessServiceV382{client: client}
}

func (s *anonymousAccessServiceV382) GetAnonymousAccessSettings(ctx context.Context) (*sonatyperepoV382.AnonymousAccessSettingsXO, *http.Response, error) {
	return s.client.SecurityManagementAnonymousAccessAPI.Read1(ctx).Execute()
}

func (s *anonymousAccessServiceV382) UpdateAnonymousAccessSettings(ctx context.Context, body sonatyperepoV382.AnonymousAccessSettingsXO) (*sonatyperepoV382.AnonymousAccessSettingsXO, *http.Response, error) {
	return s.client.SecurityManagementAnonymousAccessAPI.Update1(ctx).Body(body).Execute()
}

// anonymousAccessServiceV395 implements AnonymousAccessService against NXRM API client V395 (targets NXRM 3.94.0+).
type anonymousAccessServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func NewAnonymousAccessServiceV395(client *sonatyperepoV395.APIClient) AnonymousAccessService {
	return &anonymousAccessServiceV395{client: client}
}

func (s *anonymousAccessServiceV395) GetAnonymousAccessSettings(ctx context.Context) (*sonatyperepoV382.AnonymousAccessSettingsXO, *http.Response, error) {
	apiV395, httpResponse, err := s.client.SecurityManagementAnonymousAccessAPI.ListSecurityAnonymous(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnonymousAccessSettingsXO
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *anonymousAccessServiceV395) UpdateAnonymousAccessSettings(ctx context.Context, body sonatyperepoV382.AnonymousAccessSettingsXO) (*sonatyperepoV382.AnonymousAccessSettingsXO, *http.Response, error) {
	var v395Body sonatyperepoV395.AnonymousAccessSettingsXO
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, nil, err
	}
	apiV395, httpResponse, err := s.client.SecurityManagementAnonymousAccessAPI.UpdateSecurityAnonymous(ctx).AnonymousAccessSettingsXO(v395Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnonymousAccessSettingsXO
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}
