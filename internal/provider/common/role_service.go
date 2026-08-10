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

// RoleService abstracts the Security Management Roles API across NXRM API client generations.
// Every method's request/response shape is expressed in terms of the V382 generated types --
// the vocabulary the existing internal/provider/model package's mapping methods are already
// written against.
type RoleService interface {
	CreateRole(ctx context.Context, body sonatyperepoV382.RoleXORequest) (*sonatyperepoV382.RoleXOResponse, *http.Response, error)
	GetRole(ctx context.Context, id string) (*sonatyperepoV382.RoleXOResponse, *http.Response, error)
	UpdateRole(ctx context.Context, id string, body sonatyperepoV382.RoleXORequest) (*http.Response, error)
	DeleteRole(ctx context.Context, id string) (*http.Response, error)
	GetRoles(ctx context.Context) ([]sonatyperepoV382.RoleXOResponse, *http.Response, error)
}

// roleServiceV382 implements RoleService against NXRM API client V382 (targets NXRM < 3.94.0).
type roleServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *roleServiceV382) CreateRole(ctx context.Context, body sonatyperepoV382.RoleXORequest) (*sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	return s.client.SecurityManagementRolesAPI.Create(ctx).Body(body).Execute()
}

func (s *roleServiceV382) GetRole(ctx context.Context, id string) (*sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	return s.client.SecurityManagementRolesAPI.GetRole(ctx, id).Execute()
}

func (s *roleServiceV382) UpdateRole(ctx context.Context, id string, body sonatyperepoV382.RoleXORequest) (*http.Response, error) {
	return s.client.SecurityManagementRolesAPI.Update(ctx, id).Body(body).Execute()
}

func (s *roleServiceV382) DeleteRole(ctx context.Context, id string) (*http.Response, error) {
	return s.client.SecurityManagementRolesAPI.Delete(ctx, id).Execute()
}

func (s *roleServiceV382) GetRoles(ctx context.Context) ([]sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	return s.client.SecurityManagementRolesAPI.GetRoles(ctx).Execute()
}

// roleServiceV395 implements RoleService against NXRM API client V395 (targets NXRM 3.94.0+).
// Requests/responses are bridged to/from V382 shapes via JSON, since the internal/provider/model
// package's mapping methods are written in terms of V382 types.
//
// Field differences requiring explicit patching:
//   - RoleXORequest.Id and Name changed from *string (V382) to string (V395). jsonBridge handles
//     this transparently for reads (V395->V382). For writes (V382->V395), we bridge and then
//     manually populate non-nil required fields from the V382 struct's optional fields.
type roleServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *roleServiceV395) CreateRole(ctx context.Context, body sonatyperepoV382.RoleXORequest) (*sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	var v395Body sonatyperepoV395.RoleXORequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, nil, err
	}
	// Patch required fields: Id and Name changed from *string (V382) to string (V395)
	if body.Id != nil {
		v395Body.Id = *body.Id
	}
	if body.Name != nil {
		v395Body.Name = *body.Name
	}

	apiV395, httpResponse, err := s.client.SecurityManagementRolesAPI.CreateSecurityRoles(ctx).RoleXORequest(v395Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RoleXOResponse
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *roleServiceV395) GetRole(ctx context.Context, id string) (*sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	apiV395, httpResponse, err := s.client.SecurityManagementRolesAPI.GetSecurityRoles(ctx, id).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RoleXOResponse
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *roleServiceV395) UpdateRole(ctx context.Context, id string, body sonatyperepoV382.RoleXORequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RoleXORequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	// Patch required fields: Id and Name changed from *string (V382) to string (V395)
	if body.Id != nil {
		v395Body.Id = *body.Id
	}
	if body.Name != nil {
		v395Body.Name = *body.Name
	}

	return s.client.SecurityManagementRolesAPI.UpdateSecurityRoles(ctx, id).RoleXORequest(v395Body).Execute()
}

func (s *roleServiceV395) DeleteRole(ctx context.Context, id string) (*http.Response, error) {
	return s.client.SecurityManagementRolesAPI.DeleteSecurityRoles(ctx, id).Execute()
}

func (s *roleServiceV395) GetRoles(ctx context.Context) ([]sonatyperepoV382.RoleXOResponse, *http.Response, error) {
	apiV395, httpResponse, err := s.client.SecurityManagementRolesAPI.ListSecurityRoles(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.RoleXOResponse
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
