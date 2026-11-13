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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// serviceAccountTokenServiceV397 implements ServiceAccountTokenService against NXRM API client
// V397 (targets NXRM 3.97.0+).
type serviceAccountTokenServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

// NewServiceAccountTokenServiceV397 creates a V397 ServiceAccountTokenService adapter.
func NewServiceAccountTokenServiceV397(client *sonatyperepoV397.APIClient) ServiceAccountTokenService {
	return &serviceAccountTokenServiceV397{client: client}
}

func (s *serviceAccountTokenServiceV397) CreateServiceAccountToken(ctx context.Context, body sonatyperepoV397.ServiceAccountCreateXO) (*sonatyperepoV397.ServiceAccountCreatedXO, *http.Response, error) {
	return s.client.SecurityManagementServiceAccountTokensAPI.CreateSecurityServiceAccountTokens(ctx).ServiceAccountCreateXO(body).Execute()
}

// ListServiceAccountTokens lists all service account tokens.
//
// The generated client declares the list response as a single ServiceAccountXO, but NXRM
// returns a JSON array (verified against NXRM 3.98.1), so the typed result cannot be used: the
// generated Execute() fails to decode an array into a struct. The generated client leaves the
// response body readable after a decode failure, so the body is decoded here instead.
func (s *serviceAccountTokenServiceV397) ListServiceAccountTokens(ctx context.Context) ([]sonatyperepoV397.ServiceAccountXO, *http.Response, error) {
	_, httpResponse, err := s.client.SecurityManagementServiceAccountTokensAPI.ListSecurityServiceAccountTokens(ctx).Execute()

	// A genuine failure (no response, or a non-2xx status) is returned as-is. A 2xx response
	// with an error can only be the expected decode failure described above.
	if httpResponse == nil || httpResponse.StatusCode < 200 || httpResponse.StatusCode >= 300 {
		return nil, httpResponse, err
	}

	body, readErr := io.ReadAll(httpResponse.Body)
	if readErr != nil {
		return nil, httpResponse, readErr
	}
	httpResponse.Body = io.NopCloser(bytes.NewReader(body))

	tokens, decodeErr := decodeServiceAccountTokenList(body)
	if decodeErr != nil {
		return nil, httpResponse, decodeErr
	}
	return tokens, httpResponse, nil
}

func (s *serviceAccountTokenServiceV397) DeleteServiceAccountToken(ctx context.Context, id string) (*http.Response, error) {
	return s.client.SecurityManagementServiceAccountTokensAPI.DeleteSecurityServiceAccountTokens(ctx, id).Execute()
}

// decodeServiceAccountTokenList decodes the list response body. NXRM returns a JSON array; a
// single JSON object (which is what the published API specification describes) is accepted
// as a one-element list so that the provider keeps working if the specification is ever
// followed, and an empty body is an empty list.
func decodeServiceAccountTokenList(body []byte) ([]sonatyperepoV397.ServiceAccountXO, error) {
	trimmed := bytes.TrimSpace(body)
	if len(trimmed) == 0 {
		return []sonatyperepoV397.ServiceAccountXO{}, nil
	}

	switch trimmed[0] {
	case '[':
		var tokens []sonatyperepoV397.ServiceAccountXO
		if err := json.Unmarshal(trimmed, &tokens); err != nil {
			return nil, fmt.Errorf("unable to decode service account token list: %w", err)
		}
		return tokens, nil
	case '{':
		var token sonatyperepoV397.ServiceAccountXO
		if err := json.Unmarshal(trimmed, &token); err != nil {
			return nil, fmt.Errorf("unable to decode service account token: %w", err)
		}
		return []sonatyperepoV397.ServiceAccountXO{token}, nil
	default:
		return nil, fmt.Errorf("unexpected service account token list response: %.80s", string(trimmed))
	}
}
