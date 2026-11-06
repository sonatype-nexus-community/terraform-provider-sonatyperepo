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

// capabilityServiceV397 implements CapabilityService against NXRM API client V397 (targets NXRM 3.97.0+).
type capabilityServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *capabilityServiceV397) Create(ctx context.Context, body sonatyperepoV382.CapabilityDTO) (*sonatyperepoV382.CapabilityDTO, *http.Response, error) {
	// Bridge the V382-shaped request into V397 shape
	var v397Body sonatyperepoV397.CapabilityDTO
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, nil, err
	}

	apiV397, httpResponse, err := s.client.CapabilitiesAPI.CreateCapabilities(ctx).CapabilityDTO(v397Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.CapabilityDTO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *capabilityServiceV397) Update(ctx context.Context, capabilityId string, body sonatyperepoV382.CapabilityDTO) (*http.Response, error) {
	// Bridge the V382-shaped request into V397 shape
	var v397Body sonatyperepoV397.CapabilityDTO
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}

	return s.client.CapabilitiesAPI.UpdateCapabilities(ctx, capabilityId).CapabilityDTO(v397Body).Execute()
}

func (s *capabilityServiceV397) Delete(ctx context.Context, capabilityId string) (*http.Response, error) {
	return s.client.CapabilitiesAPI.DeleteCapabilities(ctx, capabilityId).Execute()
}

func (s *capabilityServiceV397) List(ctx context.Context) ([]sonatyperepoV382.CapabilityDTO, *http.Response, error) {
	apiV397, httpResponse, err := s.client.CapabilitiesAPI.ListCapabilities(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.CapabilityDTO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
