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

// RealmsService abstracts security realms API across NXRM API client generations.
type RealmsService interface {
	GetActiveRealms(ctx context.Context) ([]string, *http.Response, error)
	SetActiveRealms(ctx context.Context, activeRealms []string) (*http.Response, error)
}

// realmsServiceV382 implements RealmsService against NXRM API client V382 (targets NXRM < 3.94.0).
type realmsServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func NewRealmsServiceV382(client *sonatyperepoV382.APIClient) RealmsService {
	return &realmsServiceV382{client: client}
}

func (s *realmsServiceV382) GetActiveRealms(ctx context.Context) ([]string, *http.Response, error) {
	return s.client.SecurityManagementRealmsAPI.GetActiveRealms(ctx).Execute()
}

func (s *realmsServiceV382) SetActiveRealms(ctx context.Context, activeRealms []string) (*http.Response, error) {
	return s.client.SecurityManagementRealmsAPI.SetActiveRealms(ctx).Body(activeRealms).Execute()
}

// realmsServiceV395 implements RealmsService against NXRM API client V395 (targets NXRM 3.94.0+).
type realmsServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func NewRealmsServiceV395(client *sonatyperepoV395.APIClient) RealmsService {
	return &realmsServiceV395{client: client}
}

func (s *realmsServiceV395) GetActiveRealms(ctx context.Context) ([]string, *http.Response, error) {
	return s.client.SecurityManagementRealmsAPI.ListSecurityRealmsActive(ctx).Execute()
}

func (s *realmsServiceV395) SetActiveRealms(ctx context.Context, activeRealms []string) (*http.Response, error) {
	return s.client.SecurityManagementRealmsAPI.UpdateSecurityRealmsActive(ctx).RequestBody(activeRealms).Execute()
}
