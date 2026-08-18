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

// StatusService abstracts the Status API across NXRM API client generations.
type StatusService interface {
	// ClusterNodeCount returns the number of nodes reported by the cluster status checks endpoint.
	ClusterNodeCount(ctx context.Context) (int, *http.Response, error)
	// IsWritable calls the writable status check endpoint. Callers only care about the
	// HTTP response status code, so no response body is returned.
	IsWritable(ctx context.Context) (*http.Response, error)
}

// statusServiceV382 implements StatusService against NXRM API client V382 (targets NXRM < 3.94.0).
type statusServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *statusServiceV382) ClusterNodeCount(ctx context.Context) (int, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.StatusAPI.GetClusterSystemStatusChecks(ctx).Execute()
	if err != nil {
		return 0, httpResponse, err
	}
	return len(apiResponse), httpResponse, nil
}

func (s *statusServiceV382) IsWritable(ctx context.Context) (*http.Response, error) {
	return s.client.StatusAPI.IsWritable(ctx).Execute()
}

// statusServiceV395 implements StatusService against NXRM API client V395 (targets NXRM 3.94.0+).
type statusServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *statusServiceV395) ClusterNodeCount(ctx context.Context) (int, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.StatusAPI.ListStatusCheckCluster(ctx).Execute()
	if err != nil {
		return 0, httpResponse, err
	}
	return len(apiResponse), httpResponse, nil
}

func (s *statusServiceV395) IsWritable(ctx context.Context) (*http.Response, error) {
	return s.client.StatusAPI.ListStatusWritable(ctx).Execute()
}
