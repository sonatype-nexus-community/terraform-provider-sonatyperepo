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

// ContentSelectorService abstracts the Content Selectors API across NXRM API client generations.
// Every method's request/response shape is expressed in terms of the V382 generated types --
// the vocabulary the existing internal/provider/model package's mapping methods are already
// written against.
type ContentSelectorService interface {
	CreateContentSelector(ctx context.Context, body sonatyperepoV382.ContentSelectorApiCreateRequest) (*http.Response, error)
	GetContentSelector(ctx context.Context, name string) (*sonatyperepoV382.ContentSelectorApiResponse, *http.Response, error)
	UpdateContentSelector(ctx context.Context, name string, body sonatyperepoV382.ContentSelectorApiUpdateRequest) (*http.Response, error)
	DeleteContentSelector(ctx context.Context, name string) (*http.Response, error)
	GetContentSelectors(ctx context.Context) ([]sonatyperepoV382.ContentSelectorApiResponse, *http.Response, error)
}

// contentSelectorServiceV382 implements ContentSelectorService against NXRM API client V382 (targets NXRM < 3.94.0).
type contentSelectorServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *contentSelectorServiceV382) CreateContentSelector(ctx context.Context, body sonatyperepoV382.ContentSelectorApiCreateRequest) (*http.Response, error) {
	return s.client.ContentSelectorsAPI.CreateContentSelector(ctx).Body(body).Execute()
}

func (s *contentSelectorServiceV382) GetContentSelector(ctx context.Context, name string) (*sonatyperepoV382.ContentSelectorApiResponse, *http.Response, error) {
	return s.client.ContentSelectorsAPI.GetContentSelector(ctx, name).Execute()
}

func (s *contentSelectorServiceV382) UpdateContentSelector(ctx context.Context, name string, body sonatyperepoV382.ContentSelectorApiUpdateRequest) (*http.Response, error) {
	return s.client.ContentSelectorsAPI.UpdateContentSelector(ctx, name).Body(body).Execute()
}

func (s *contentSelectorServiceV382) DeleteContentSelector(ctx context.Context, name string) (*http.Response, error) {
	return s.client.ContentSelectorsAPI.DeleteContentSelector(ctx, name).Execute()
}

func (s *contentSelectorServiceV382) GetContentSelectors(ctx context.Context) ([]sonatyperepoV382.ContentSelectorApiResponse, *http.Response, error) {
	return s.client.ContentSelectorsAPI.GetContentSelectors(ctx).Execute()
}

// contentSelectorServiceV395 implements ContentSelectorService against NXRM API client V395 (targets NXRM 3.94.0+).
// Requests/responses are bridged to/from V382 shapes via JSON, since the internal/provider/model
// package's mapping methods are written in terms of V382 types.
//
// Field differences requiring explicit patching:
//   - ContentSelectorApiCreateRequest.Name and Expression changed from *string (V382) to string (V395).
//     jsonBridge handles this transparently for reads. For writes, we bridge and then manually
//     populate non-nil required fields from the V382 struct's optional fields.
//   - ContentSelectorApiUpdateRequest.Expression changed from *string (V382) to string (V395).
//     Handled similarly on writes.
type contentSelectorServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *contentSelectorServiceV395) CreateContentSelector(ctx context.Context, body sonatyperepoV382.ContentSelectorApiCreateRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.ContentSelectorApiCreateRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	// Patch required fields: Name and Expression changed from *string (V382) to string (V395)
	if body.Name != nil {
		v395Body.Name = *body.Name
	}
	if body.Expression != nil {
		v395Body.Expression = *body.Expression
	}

	return s.client.ContentSelectorsAPI.CreateSecurityContentSelectors(ctx).ContentSelectorApiCreateRequest(v395Body).Execute()
}

func (s *contentSelectorServiceV395) GetContentSelector(ctx context.Context, name string) (*sonatyperepoV382.ContentSelectorApiResponse, *http.Response, error) {
	apiV395, httpResponse, err := s.client.ContentSelectorsAPI.GetSecurityContentSelectors(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.ContentSelectorApiResponse
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *contentSelectorServiceV395) UpdateContentSelector(ctx context.Context, name string, body sonatyperepoV382.ContentSelectorApiUpdateRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.ContentSelectorApiUpdateRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	// Patch required field: Expression changed from *string (V382) to string (V395)
	if body.Expression != nil {
		v395Body.Expression = *body.Expression
	}

	return s.client.ContentSelectorsAPI.UpdateSecurityContentSelectors(ctx, name).ContentSelectorApiUpdateRequest(v395Body).Execute()
}

func (s *contentSelectorServiceV395) DeleteContentSelector(ctx context.Context, name string) (*http.Response, error) {
	return s.client.ContentSelectorsAPI.DeleteSecurityContentSelectors(ctx, name).Execute()
}

func (s *contentSelectorServiceV395) GetContentSelectors(ctx context.Context) ([]sonatyperepoV382.ContentSelectorApiResponse, *http.Response, error) {
	apiV395, httpResponse, err := s.client.ContentSelectorsAPI.ListSecurityContentSelectors(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.ContentSelectorApiResponse
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
