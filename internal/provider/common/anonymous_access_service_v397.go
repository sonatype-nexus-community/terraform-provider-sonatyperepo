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

// anonymousAccessServiceV397 implements AnonymousAccessService against NXRM API client V397 (targets NXRM 3.97.0+).
type anonymousAccessServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func NewAnonymousAccessServiceV397(client *sonatyperepoV397.APIClient) AnonymousAccessService {
	return &anonymousAccessServiceV397{client: client}
}

func (s *anonymousAccessServiceV397) GetAnonymousAccessSettings(ctx context.Context) (*sonatyperepoV382.AnonymousAccessSettingsXO, *http.Response, error) {
	apiV397, httpResponse, err := s.client.SecurityManagementAnonymousAccessAPI.ListSecurityAnonymous(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnonymousAccessSettingsXO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *anonymousAccessServiceV397) UpdateAnonymousAccessSettings(ctx context.Context, body sonatyperepoV382.AnonymousAccessSettingsXO) (*sonatyperepoV382.AnonymousAccessSettingsXO, *http.Response, error) {
	var v397Body sonatyperepoV397.AnonymousAccessSettingsXO
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, nil, err
	}
	apiV397, httpResponse, err := s.client.SecurityManagementAnonymousAccessAPI.UpdateSecurityAnonymous(ctx).AnonymousAccessSettingsXO(v397Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnonymousAccessSettingsXO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}
