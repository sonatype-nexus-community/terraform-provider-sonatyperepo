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

// roleServiceV397 implements RoleService against NXRM API client V397 (targets NXRM 3.97.0+).
//   - RoleXORequest.Id and Name changed from *string (V382) to string (V397). jsonBridge handles
//     this transparently for reads (V397->V382). For writes (V382->V397), we bridge and then
type roleServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *roleServiceV397) CreateRole(ctx context.Context, body sonatyperepoV382.RoleXORequest) (*sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	var v397Body sonatyperepoV397.RoleXORequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, nil, err
	}
	// Patch required fields: Id and Name changed from *string (V382) to string (V397)
	if body.Id != nil {
		v397Body.Id = *body.Id
	}
	if body.Name != nil {
		v397Body.Name = *body.Name
	}

	apiV397, httpResponse, err := s.client.SecurityManagementRolesAPI.CreateSecurityRoles(ctx).RoleXORequest(v397Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RoleXOResponse
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *roleServiceV397) GetRole(ctx context.Context, id string) (*sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	apiV397, httpResponse, err := s.client.SecurityManagementRolesAPI.GetSecurityRoles(ctx, id).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RoleXOResponse
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *roleServiceV397) UpdateRole(ctx context.Context, id string, body sonatyperepoV382.RoleXORequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RoleXORequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	// Patch required fields: Id and Name changed from *string (V382) to string (V397)
	if body.Id != nil {
		v397Body.Id = *body.Id
	}
	if body.Name != nil {
		v397Body.Name = *body.Name
	}

	return s.client.SecurityManagementRolesAPI.UpdateSecurityRoles(ctx, id).RoleXORequest(v397Body).Execute()
}

func (s *roleServiceV397) DeleteRole(ctx context.Context, id string) (*http.Response, error) {
	return s.client.SecurityManagementRolesAPI.DeleteSecurityRoles(ctx, id).Execute()
}

func (s *roleServiceV397) GetRoles(ctx context.Context) ([]sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	apiV397, httpResponse, err := s.client.SecurityManagementRolesAPI.ListSecurityRoles(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.RoleXOResponse
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
