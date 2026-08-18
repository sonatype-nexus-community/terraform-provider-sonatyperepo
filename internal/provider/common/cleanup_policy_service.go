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
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// CleanupPolicyService abstracts cleanup policy CRUD across NXRM API client generations.
// Every method's request/response shape is expressed in terms of the V382 generated types --
// the vocabulary the existing internal/provider/model package's methods are already written against.
// The V395 adapter bridges to/from its own generated types internally so that callers never need
// to know which generation is behind the interface.
type CleanupPolicyService interface {
	// Create creates a new cleanup policy.
	Create(ctx context.Context, body sonatyperepoV382.CleanupPolicyResourceXO) (*http.Response, error)
	// GetByName retrieves a cleanup policy by name.
	GetByName(ctx context.Context, name string) (*http.Response, error)
	// Update updates an existing cleanup policy.
	Update(ctx context.Context, policyName string, body sonatyperepoV382.CleanupPolicyResourceXO) (*http.Response, error)
	// DeleteByName deletes a cleanup policy by name.
	DeleteByName(ctx context.Context, name string) (*http.Response, error)
}

// cleanupPolicyServiceV382 implements CleanupPolicyService against NXRM API client V382 (targets NXRM < 3.94.0).
type cleanupPolicyServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func NewCleanupPolicyServiceV382(client *sonatyperepoV382.APIClient) CleanupPolicyService {
	return &cleanupPolicyServiceV382{client: client}
}

func (s *cleanupPolicyServiceV382) Create(ctx context.Context, body sonatyperepoV382.CleanupPolicyResourceXO) (*http.Response, error) {
	return s.client.CleanupPoliciesAPI.Create1(ctx).Body(body).Execute()
}

func (s *cleanupPolicyServiceV382) GetByName(ctx context.Context, name string) (*http.Response, error) {
	return s.client.CleanupPoliciesAPI.GetCleanupPolicyByName(ctx, name).Execute()
}

func (s *cleanupPolicyServiceV382) Update(ctx context.Context, policyName string, body sonatyperepoV382.CleanupPolicyResourceXO) (*http.Response, error) {
	return s.client.CleanupPoliciesAPI.Update2(ctx, policyName).Body(body).Execute()
}

func (s *cleanupPolicyServiceV382) DeleteByName(ctx context.Context, name string) (*http.Response, error) {
	return s.client.CleanupPoliciesAPI.DeletePolicyByName(ctx, name).Execute()
}

// cleanupPolicyServiceV395 implements CleanupPolicyService against NXRM API client V395 (targets NXRM 3.94.0+).
type cleanupPolicyServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func NewCleanupPolicyServiceV395(client *sonatyperepoV395.APIClient) CleanupPolicyService {
	return &cleanupPolicyServiceV395{client: client}
}

func (s *cleanupPolicyServiceV395) Create(ctx context.Context, body sonatyperepoV382.CleanupPolicyResourceXO) (*http.Response, error) {
	// Bridge the V382-shaped request into V395 shape
	var v395Body sonatyperepoV395.CleanupPolicyResourceXO
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}

	return s.client.CleanupPoliciesAPI.CreateCleanupPolicies(ctx).CleanupPolicyResourceXO(v395Body).Execute()
}

func (s *cleanupPolicyServiceV395) GetByName(ctx context.Context, name string) (*http.Response, error) {
	httpResponse, err := s.client.CleanupPoliciesAPI.GetCleanupPolicies(ctx, name).Execute()
	if err != nil {
		return httpResponse, err
	}

	// The API doesn't declare a typed response for this endpoint in either client generation,
	// so callers manually decode the body into the V382 CleanupPolicyResourceXO shape. V395's
	// response has an extra 'repositories' field V382's struct doesn't know about, so prune the
	// raw body down to V382's shape here before the caller ever sees it.
	body, readErr := io.ReadAll(httpResponse.Body)
	_ = httpResponse.Body.Close()
	if readErr != nil {
		return httpResponse, fmt.Errorf("could not read response body: %w", readErr)
	}

	pruned, pruneErr := pruneUnknownJSONFields(body, reflect.TypeOf(sonatyperepoV382.CleanupPolicyResourceXO{}))
	if pruneErr != nil {
		httpResponse.Body = io.NopCloser(bytes.NewReader(body))
		return httpResponse, pruneErr
	}

	httpResponse.Body = io.NopCloser(bytes.NewReader(pruned))
	return httpResponse, nil
}

func (s *cleanupPolicyServiceV395) Update(ctx context.Context, policyName string, body sonatyperepoV382.CleanupPolicyResourceXO) (*http.Response, error) {
	// Bridge the V382-shaped request into V395 shape
	var v395Body sonatyperepoV395.CleanupPolicyResourceXO
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}

	return s.client.CleanupPoliciesAPI.UpdateCleanupPolicies(ctx, policyName).CleanupPolicyResourceXO(v395Body).Execute()
}

func (s *cleanupPolicyServiceV395) DeleteByName(ctx context.Context, name string) (*http.Response, error) {
	return s.client.CleanupPoliciesAPI.DeleteCleanupPolicies(ctx, name).Execute()
}
