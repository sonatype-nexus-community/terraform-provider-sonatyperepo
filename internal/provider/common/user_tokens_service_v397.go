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

// userTokensServiceV397 implements UserTokensService against NXRM API client V397 (targets NXRM 3.97.0+).
type userTokensServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *userTokensServiceV397) ServiceStatus(ctx context.Context) (*sonatyperepoV382.UserTokensApiModel, *http.Response, error) {
	v397Response, httpResponse, err := s.client.SecurityManagementUserTokensAPI.ListSecurityUserTokens(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V397 response to V382 shape
	var v382Response sonatyperepoV382.UserTokensApiModel
	if err := jsonBridge(v397Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

func (s *userTokensServiceV397) SetServiceStatus(ctx context.Context, body sonatyperepoV382.UserTokensApiModel) (*sonatyperepoV382.UserTokensApiModel, *http.Response, error) {
	// Bridge V382 request to V397 shape
	var v397Body sonatyperepoV397.UserTokensApiModel
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, nil, err
	}

	v397Response, httpResponse, err := s.client.SecurityManagementUserTokensAPI.UpdateSecurityUserTokens(ctx).UserTokensApiModel(v397Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V397 response to V382 shape
	var v382Response sonatyperepoV382.UserTokensApiModel
	if err := jsonBridge(v397Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

// NewUserTokensServiceV397 creates a V397 UserTokensService adapter.
func NewUserTokensServiceV397(client *sonatyperepoV397.APIClient) UserTokensService {
	return &userTokensServiceV397{client: client}
}
