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

// HttpSettingsService abstracts the HTTP System Settings API across NXRM API client generations.
type HttpSettingsService interface {
	// GetHttpSettings retrieves current HTTP system settings.
	GetHttpSettings(ctx context.Context) (*sonatyperepoV382.HttpSettingsXo, *http.Response, error)
	// ResetHttpSettings resets HTTP system settings to defaults.
	ResetHttpSettings(ctx context.Context) (*http.Response, error)
	// UpdateHttpSettings updates HTTP system settings.
	UpdateHttpSettings(ctx context.Context, body sonatyperepoV382.HttpSettingsXo) (*http.Response, error)
}

// httpSettingsServiceV382 implements HttpSettingsService against NXRM API client V382 (targets NXRM < 3.94.0).
type httpSettingsServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *httpSettingsServiceV382) GetHttpSettings(ctx context.Context) (*sonatyperepoV382.HttpSettingsXo, *http.Response, error) {
	return s.client.ManageSonatypeHTTPSystemSettingsAPI.GetHttpSettings(ctx).Execute()
}

func (s *httpSettingsServiceV382) ResetHttpSettings(ctx context.Context) (*http.Response, error) {
	return s.client.ManageSonatypeHTTPSystemSettingsAPI.ResetHttpSettings(ctx).Execute()
}

func (s *httpSettingsServiceV382) UpdateHttpSettings(ctx context.Context, body sonatyperepoV382.HttpSettingsXo) (*http.Response, error) {
	return s.client.ManageSonatypeHTTPSystemSettingsAPI.UpdateHttpSettings(ctx).Body(body).Execute()
}

// httpSettingsServiceV395 implements HttpSettingsService against NXRM API client V395 (targets NXRM 3.94.0+).
type httpSettingsServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *httpSettingsServiceV395) GetHttpSettings(ctx context.Context) (*sonatyperepoV382.HttpSettingsXo, *http.Response, error) {
	v395Response, httpResponse, err := s.client.ManageSonatypeHTTPSystemSettingsAPI.ListHttp(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}

	// Bridge V395 response to V382 shape
	var v382Response sonatyperepoV382.HttpSettingsXo
	if err := jsonBridge(v395Response, &v382Response); err != nil {
		return nil, httpResponse, err
	}

	return &v382Response, httpResponse, nil
}

func (s *httpSettingsServiceV395) ResetHttpSettings(ctx context.Context) (*http.Response, error) {
	return s.client.ManageSonatypeHTTPSystemSettingsAPI.DeleteHttp(ctx).Execute()
}

func (s *httpSettingsServiceV395) UpdateHttpSettings(ctx context.Context, body sonatyperepoV382.HttpSettingsXo) (*http.Response, error) {
	// Bridge V382 request to V395 shape
	var v395Body sonatyperepoV395.HttpSettingsXo
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}

	return s.client.ManageSonatypeHTTPSystemSettingsAPI.UpdateHttp(ctx).HttpSettingsXo(v395Body).Execute()
}

// NewHttpSettingsServiceV382 creates a V382 HttpSettingsService adapter.
func NewHttpSettingsServiceV382(client *sonatyperepoV382.APIClient) HttpSettingsService {
	return &httpSettingsServiceV382{client: client}
}

// NewHttpSettingsServiceV395 creates a V395 HttpSettingsService adapter.
func NewHttpSettingsServiceV395(client *sonatyperepoV395.APIClient) HttpSettingsService {
	return &httpSettingsServiceV395{client: client}
}
