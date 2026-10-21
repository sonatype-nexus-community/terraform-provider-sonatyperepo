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

// contentSelectorServiceV397 implements ContentSelectorService against NXRM API client V397 (targets NXRM 3.97.0+).
//   - ContentSelectorApiCreateRequest.Name and Expression changed from *string (V382) to string (V397).
//   - ContentSelectorApiUpdateRequest.Expression changed from *string (V382) to string (V397).
type contentSelectorServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *contentSelectorServiceV397) CreateContentSelector(ctx context.Context, body sonatyperepoV382.ContentSelectorApiCreateRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.ContentSelectorApiCreateRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	// Patch required fields: Name and Expression changed from *string (V382) to string (V397)
	if body.Name != nil {
		v397Body.Name = *body.Name
	}
	if body.Expression != nil {
		v397Body.Expression = *body.Expression
	}

	return s.client.ContentSelectorsAPI.CreateSecurityContentSelectors(ctx).ContentSelectorApiCreateRequest(v397Body).Execute()
}

func (s *contentSelectorServiceV397) GetContentSelector(ctx context.Context, name string) (*sonatyperepoV382.ContentSelectorApiResponse, *http.Response, error) {
	apiV397, httpResponse, err := s.client.ContentSelectorsAPI.GetSecurityContentSelectors(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.ContentSelectorApiResponse
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *contentSelectorServiceV397) UpdateContentSelector(ctx context.Context, name string, body sonatyperepoV382.ContentSelectorApiUpdateRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.ContentSelectorApiUpdateRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	// Patch required field: Expression changed from *string (V382) to string (V397)
	if body.Expression != nil {
		v397Body.Expression = *body.Expression
	}

	return s.client.ContentSelectorsAPI.UpdateSecurityContentSelectors(ctx, name).ContentSelectorApiUpdateRequest(v397Body).Execute()
}

func (s *contentSelectorServiceV397) DeleteContentSelector(ctx context.Context, name string) (*http.Response, error) {
	return s.client.ContentSelectorsAPI.DeleteSecurityContentSelectors(ctx, name).Execute()
}

func (s *contentSelectorServiceV397) GetContentSelectors(ctx context.Context) ([]sonatyperepoV382.ContentSelectorApiResponse, *http.Response, error) {
	apiV397, httpResponse, err := s.client.ContentSelectorsAPI.ListSecurityContentSelectors(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.ContentSelectorApiResponse
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
