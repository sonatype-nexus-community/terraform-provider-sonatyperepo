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
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// The V397 adapter bridges to/from its own generated types internally so that callers never need
// cleanupPolicyServiceV397 implements CleanupPolicyService against NXRM API client V397 (targets NXRM 3.97.0+).
type cleanupPolicyServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func NewCleanupPolicyServiceV397(client *sonatyperepoV397.APIClient) CleanupPolicyService {
	return &cleanupPolicyServiceV397{client: client}
}

func (s *cleanupPolicyServiceV397) Create(ctx context.Context, body sonatyperepoV382.CleanupPolicyResourceXO) (*http.Response, error) {
	// Bridge the V382-shaped request into V397 shape
	var v397Body sonatyperepoV397.CleanupPolicyResourceXO
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}

	return s.client.CleanupPoliciesAPI.CreateCleanupPolicies(ctx).CleanupPolicyResourceXO(v397Body).Execute()
}

func (s *cleanupPolicyServiceV397) GetByName(ctx context.Context, name string) (*http.Response, error) {
	_, httpResponse, err := s.client.CleanupPoliciesAPI.GetCleanupPolicies(ctx, name).Execute()
	if err != nil {
		return httpResponse, err
	}

	// V397 now declares a typed response for this endpoint (as of the client's v397.95.2 release;
	// v397.95.0/.1 left it untyped), but callers still manually decode the body into the V382
	// CleanupPolicyResourceXO shape for cross-generation compatibility. V397's response has an
	// extra 'repositories' field V382's struct doesn't know about, so prune the raw body down to
	// V382's shape here before the caller ever sees it.
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

func (s *cleanupPolicyServiceV397) Update(ctx context.Context, policyName string, body sonatyperepoV382.CleanupPolicyResourceXO) (*http.Response, error) {
	// Bridge the V382-shaped request into V397 shape
	var v397Body sonatyperepoV397.CleanupPolicyResourceXO
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}

	return s.client.CleanupPoliciesAPI.UpdateCleanupPolicies(ctx, policyName).CleanupPolicyResourceXO(v397Body).Execute()
}

func (s *cleanupPolicyServiceV397) DeleteByName(ctx context.Context, name string) (*http.Response, error) {
	return s.client.CleanupPoliciesAPI.DeleteCleanupPolicies(ctx, name).Execute()
}
