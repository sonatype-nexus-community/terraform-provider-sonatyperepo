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

// written against. The V397 adapter bridges to/from its own generated types
// userServiceV397 implements UserService against NXRM API client V397 (targets NXRM 3.97.0+).
type userServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *userServiceV397) CreateUser(ctx context.Context, body sonatyperepoV382.ApiCreateUser) (*sonatyperepoV382.ApiUser, *http.Response, error) {
	// Bridge V382 ApiCreateUser -> V397 ApiCreateUser
	var v397Body sonatyperepoV397.ApiCreateUser
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, nil, err
	}

	// V397 ApiCreateUser has required fields that are optional in V382.
	// jsonBridge will drop fields that don't exist on the V397 side,
	// but V397 requires emailAddress, firstName, lastName, password, userId.
	// We need to explicitly map these from the V382 body (which may have nil pointers).
	if body.EmailAddress != nil {
		v397Body.EmailAddress = *body.EmailAddress
	}
	if body.FirstName != nil {
		v397Body.FirstName = *body.FirstName
	}
	if body.LastName != nil {
		v397Body.LastName = *body.LastName
	}
	if body.Password != nil {
		v397Body.Password = *body.Password
	}
	if body.UserId != nil {
		v397Body.UserId = *body.UserId
	}
	v397Body.Status = body.Status
	v397Body.Roles = body.Roles

	apiV397, httpResponse, err := s.client.SecurityManagementUsersAPI.CreateSecurityUsers(ctx).ApiCreateUser(v397Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V397 ApiUser -> V382 ApiUser
	var result sonatyperepoV382.ApiUser
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}

	return &result, httpResponse, nil
}

func (s *userServiceV397) GetUsers(ctx context.Context, userId string, source string) ([]sonatyperepoV382.ApiUser, *http.Response, error) {
	req := s.client.SecurityManagementUsersAPI.ListSecurityUsers(ctx)
	// See userServiceV382.GetUsers: an empty filter value must be omitted, not
	// passed through, since the builder always sets the query param otherwise.
	if userId != "" {
		req = req.UserId(userId)
	}
	if source != "" {
		req = req.Source(source)
	}
	apiV397, httpResponse, err := req.Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V397 []ApiUser -> V382 []ApiUser
	var result []sonatyperepoV382.ApiUser
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}

	return result, httpResponse, nil
}

func (s *userServiceV397) UpdateUser(ctx context.Context, userId string, body sonatyperepoV382.ApiUser) (*http.Response, error) {
	// Bridge V382 ApiUser -> V397 ApiUser
	var v397Body sonatyperepoV397.ApiUser
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}

	// V397 ApiUser has Source and UserId as required fields (not pointers).
	// jsonBridge will drop fields that don't exist on the V397 side,
	// but V397 requires source, status, userId.
	// We need to explicitly map these from the V382 body.
	if body.Source != nil {
		v397Body.Source = *body.Source
	}
	if body.UserId != nil {
		v397Body.UserId = *body.UserId
	}
	v397Body.Status = body.Status

	return s.client.SecurityManagementUsersAPI.UpdateSecurityUsers(ctx, userId).ApiUser(v397Body).Execute()
}

func (s *userServiceV397) ChangePassword(ctx context.Context, userId string, password string) (*http.Response, error) {
	return s.client.SecurityManagementUsersAPI.UpdateSecurityUsersChangePassword(ctx, userId).Body(password).Execute()
}

func (s *userServiceV397) DeleteUser(ctx context.Context, userId string) (*http.Response, error) {
	return s.client.SecurityManagementUsersAPI.DeleteSecurityUsers(ctx, userId).Execute()
}
