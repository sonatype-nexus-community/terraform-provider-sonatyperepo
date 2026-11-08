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
	"fmt"
	"net/http"

	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// errServiceAccountTokensUnsupported is returned by every ServiceAccountTokenService method
// when the connected NXRM server predates 3.97.0. Service account tokens have no equivalent in
// the V382 or V395 client generations at all -- there is no bridging to do, only a clear error
// to surface.
var errServiceAccountTokensUnsupported = fmt.Errorf("service account tokens require Nexus Repository Manager 3.97.0 or later")

// ServiceAccountTokenService abstracts the service account token API
// (/v1/security/service-account-tokens). This API only exists in the V397 client generation
// (NXRM 3.97.0+), and only has create, list and delete operations: there is no get-by-id and
// no update.
type ServiceAccountTokenService interface {
	// CreateServiceAccountToken creates a token. The token value is only ever returned by
	// this call; it cannot be retrieved later.
	CreateServiceAccountToken(ctx context.Context, body sonatyperepoV397.ServiceAccountCreateXO) (*sonatyperepoV397.ServiceAccountCreatedXO, *http.Response, error)
	// ListServiceAccountTokens lists all service account tokens. Secret material is never included.
	ListServiceAccountTokens(ctx context.Context) ([]sonatyperepoV397.ServiceAccountXO, *http.Response, error)
	// DeleteServiceAccountToken deletes (revokes) the token with the given id.
	DeleteServiceAccountToken(ctx context.Context, id string) (*http.Response, error)
}

// serviceAccountTokenServiceUnsupported implements ServiceAccountTokenService for NXRM servers
// older than 3.97.0, where the service account token REST API does not exist.
type serviceAccountTokenServiceUnsupported struct{}

// NewServiceAccountTokenServiceUnsupported creates a ServiceAccountTokenService adapter for NXRM
// servers older than 3.97.0, which have no service account token API to call.
func NewServiceAccountTokenServiceUnsupported() ServiceAccountTokenService {
	return &serviceAccountTokenServiceUnsupported{}
}

func (s *serviceAccountTokenServiceUnsupported) CreateServiceAccountToken(_ context.Context, _ sonatyperepoV397.ServiceAccountCreateXO) (*sonatyperepoV397.ServiceAccountCreatedXO, *http.Response, error) {
	return nil, nil, errServiceAccountTokensUnsupported
}

func (s *serviceAccountTokenServiceUnsupported) ListServiceAccountTokens(_ context.Context) ([]sonatyperepoV397.ServiceAccountXO, *http.Response, error) {
	return nil, nil, errServiceAccountTokensUnsupported
}

func (s *serviceAccountTokenServiceUnsupported) DeleteServiceAccountToken(_ context.Context, _ string) (*http.Response, error) {
	return nil, errServiceAccountTokensUnsupported
}
