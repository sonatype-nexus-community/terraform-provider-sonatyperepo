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

// PrivilegeService abstracts the Security Management Privileges API across NXRM API client generations.
type PrivilegeService interface {
	// GetAllPrivileges retrieves a list of all privileges.
	GetAllPrivileges(ctx context.Context) ([]sonatyperepoV382.ApiPrivilegeRequest, *http.Response, error)
}

// privilegeServiceV382 implements PrivilegeService against NXRM API client V382 (targets NXRM < 3.94.0).
type privilegeServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *privilegeServiceV382) GetAllPrivileges(ctx context.Context) ([]sonatyperepoV382.ApiPrivilegeRequest, *http.Response, error) {
	return s.client.SecurityManagementPrivilegesAPI.GetAllPrivileges(ctx).Execute()
}

// privilegeServiceV395 implements PrivilegeService against NXRM API client V395 (targets NXRM 3.94.0+).
type privilegeServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *privilegeServiceV395) GetAllPrivileges(ctx context.Context) ([]sonatyperepoV382.ApiPrivilegeRequest, *http.Response, error) {
	apiV395, httpResponse, err := s.client.SecurityManagementPrivilegesAPI.GetAllPrivileges(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.ApiPrivilegeRequest
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
