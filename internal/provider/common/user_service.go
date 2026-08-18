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

// UserService abstracts User management across NXRM API client generations.
// Every method's request/response shape is expressed in terms of the V382
// generated types -- the vocabulary the existing internal/provider/model
// package's FromApiModel/ToApiCreateModel/ToApiUpdateModel methods are already
// written against. The V395 adapter bridges to/from its own generated types
// internally so that callers never need to know which generation is behind
// the interface.
type UserService interface {
	CreateUser(ctx context.Context, body sonatyperepoV382.ApiCreateUser) (*sonatyperepoV382.ApiUser, *http.Response, error)
	GetUsers(ctx context.Context, userId string, source string) ([]sonatyperepoV382.ApiUser, *http.Response, error)
	UpdateUser(ctx context.Context, userId string, body sonatyperepoV382.ApiUser) (*http.Response, error)
	ChangePassword(ctx context.Context, userId string, password string) (*http.Response, error)
	DeleteUser(ctx context.Context, userId string) (*http.Response, error)
}

// userServiceV382 implements UserService against NXRM API client V382 (targets NXRM < 3.94.0).
type userServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *userServiceV382) CreateUser(ctx context.Context, body sonatyperepoV382.ApiCreateUser) (*sonatyperepoV382.ApiUser, *http.Response, error) {
	return s.client.SecurityManagementUsersAPI.CreateUser(ctx).Body(body).Execute()
}

func (s *userServiceV382) GetUsers(ctx context.Context, userId string, source string) ([]sonatyperepoV382.ApiUser, *http.Response, error) {
	req := s.client.SecurityManagementUsersAPI.GetUsers(ctx)
	// The generated builder methods always set the query param, even to an empty
	// string, which NXRM then treats as an exact-match filter for "" rather than
	// "unset" -- so an empty filter value must be omitted, not passed through.
	if userId != "" {
		req = req.UserId(userId)
	}
	if source != "" {
		req = req.Source(source)
	}
	return req.Execute()
}

func (s *userServiceV382) UpdateUser(ctx context.Context, userId string, body sonatyperepoV382.ApiUser) (*http.Response, error) {
	return s.client.SecurityManagementUsersAPI.UpdateUser(ctx, userId).Body(body).Execute()
}

func (s *userServiceV382) ChangePassword(ctx context.Context, userId string, password string) (*http.Response, error) {
	return s.client.SecurityManagementUsersAPI.ChangePassword(ctx, userId).Body(password).Execute()
}

func (s *userServiceV382) DeleteUser(ctx context.Context, userId string) (*http.Response, error) {
	return s.client.SecurityManagementUsersAPI.DeleteUser(ctx, userId).Execute()
}

// userServiceV395 implements UserService against NXRM API client V395 (targets NXRM 3.94.0+).
type userServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *userServiceV395) CreateUser(ctx context.Context, body sonatyperepoV382.ApiCreateUser) (*sonatyperepoV382.ApiUser, *http.Response, error) {
	// Bridge V382 ApiCreateUser -> V395 ApiCreateUser
	var v395Body sonatyperepoV395.ApiCreateUser
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, nil, err
	}

	// V395 ApiCreateUser has required fields that are optional in V382.
	// jsonBridge will drop fields that don't exist on the V395 side,
	// but V395 requires emailAddress, firstName, lastName, password, userId.
	// We need to explicitly map these from the V382 body (which may have nil pointers).
	if body.EmailAddress != nil {
		v395Body.EmailAddress = *body.EmailAddress
	}
	if body.FirstName != nil {
		v395Body.FirstName = *body.FirstName
	}
	if body.LastName != nil {
		v395Body.LastName = *body.LastName
	}
	if body.Password != nil {
		v395Body.Password = *body.Password
	}
	if body.UserId != nil {
		v395Body.UserId = *body.UserId
	}
	v395Body.Status = body.Status
	v395Body.Roles = body.Roles

	apiV395, httpResponse, err := s.client.SecurityManagementUsersAPI.CreateSecurityUsers(ctx).ApiCreateUser(v395Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 ApiUser -> V382 ApiUser
	var result sonatyperepoV382.ApiUser
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}

	return &result, httpResponse, nil
}

func (s *userServiceV395) GetUsers(ctx context.Context, userId string, source string) ([]sonatyperepoV382.ApiUser, *http.Response, error) {
	req := s.client.SecurityManagementUsersAPI.ListSecurityUsers(ctx)
	// See userServiceV382.GetUsers: an empty filter value must be omitted, not
	// passed through, since the builder always sets the query param otherwise.
	if userId != "" {
		req = req.UserId(userId)
	}
	if source != "" {
		req = req.Source(source)
	}
	apiV395, httpResponse, err := req.Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 []ApiUser -> V382 []ApiUser
	var result []sonatyperepoV382.ApiUser
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}

	return result, httpResponse, nil
}

func (s *userServiceV395) UpdateUser(ctx context.Context, userId string, body sonatyperepoV382.ApiUser) (*http.Response, error) {
	// Bridge V382 ApiUser -> V395 ApiUser
	var v395Body sonatyperepoV395.ApiUser
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}

	// V395 ApiUser has Source and UserId as required fields (not pointers).
	// jsonBridge will drop fields that don't exist on the V395 side,
	// but V395 requires source, status, userId.
	// We need to explicitly map these from the V382 body.
	if body.Source != nil {
		v395Body.Source = *body.Source
	}
	if body.UserId != nil {
		v395Body.UserId = *body.UserId
	}
	v395Body.Status = body.Status

	return s.client.SecurityManagementUsersAPI.UpdateSecurityUsers(ctx, userId).ApiUser(v395Body).Execute()
}

func (s *userServiceV395) ChangePassword(ctx context.Context, userId string, password string) (*http.Response, error) {
	return s.client.SecurityManagementUsersAPI.UpdateSecurityUsersChangePassword(ctx, userId).Body(password).Execute()
}

func (s *userServiceV395) DeleteUser(ctx context.Context, userId string) (*http.Response, error) {
	return s.client.SecurityManagementUsersAPI.DeleteSecurityUsers(ctx, userId).Execute()
}
