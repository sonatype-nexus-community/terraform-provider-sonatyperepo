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

// configurationServiceV397 implements ConfigurationService against NXRM API client V397
type configurationServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

// Email API implementation - V397

func (s *configurationServiceV397) SetEmailConfiguration(ctx context.Context, body sonatyperepoV382.ApiEmailConfiguration) (*http.Response, error) {
	var v397Body sonatyperepoV397.ApiEmailConfiguration
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.EmailAPI.UpdateEmail(ctx).ApiEmailConfiguration(v397Body).Execute()
}

func (s *configurationServiceV397) GetEmailConfiguration(ctx context.Context) (*sonatyperepoV382.ApiEmailConfiguration, *http.Response, error) {
	apiV397, httpResponse, err := s.client.EmailAPI.ListEmail(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.ApiEmailConfiguration
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *configurationServiceV397) DeleteEmailConfiguration(ctx context.Context) (*http.Response, error) {
	return s.client.EmailAPI.DeleteEmail(ctx).Execute()
}

// LDAP API implementation - V397

func (s *configurationServiceV397) CreateLdapServer(ctx context.Context, body sonatyperepoV382.CreateLdapServerXo) (*http.Response, error) {
	var v397Body sonatyperepoV397.CreateLdapServerXo
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.SecurityManagementLDAPAPI.CreateSecurityLdap(ctx).CreateLdapServerXo(v397Body).Execute()
}

func (s *configurationServiceV397) GetLdapServer(ctx context.Context, name string) (*sonatyperepoV382.ReadLdapServerXo, *http.Response, error) {
	apiV397, httpResponse, err := s.client.SecurityManagementLDAPAPI.GetSecurityLdap(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.ReadLdapServerXo
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *configurationServiceV397) UpdateLdapServer(ctx context.Context, name string, body sonatyperepoV382.UpdateLdapServerXo) (*http.Response, error) {
	var v397Body sonatyperepoV397.UpdateLdapServerXo
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.SecurityManagementLDAPAPI.UpdateSecurityLdap(ctx, name).UpdateLdapServerXo(v397Body).Execute()
}

func (s *configurationServiceV397) DeleteLdapServer(ctx context.Context, name string) (*http.Response, error) {
	return s.client.SecurityManagementLDAPAPI.DeleteSecurityLdap(ctx, name).Execute()
}

// IQ Connection API implementation - V397

func (s *configurationServiceV397) GetIqConnectionConfiguration(ctx context.Context) (*sonatyperepoV382.IqConnectionXo, *http.Response, error) {
	apiV397, httpResponse, err := s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.ListIq(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.IqConnectionXo
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *configurationServiceV397) UpdateIqConnectionConfiguration(ctx context.Context, body sonatyperepoV382.IqConnectionXo) (*sonatyperepoV382.IqConnectionXo, *http.Response, error) {
	var v397Body sonatyperepoV397.IqConnectionXo
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, nil, err
	}
	apiV397, httpResponse, err := s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.UpdateIq(ctx).IqConnectionXo(v397Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.IqConnectionXo
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *configurationServiceV397) DisableIq(ctx context.Context) (*http.Response, error) {
	return s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.CreateIqDisable(ctx).Execute()
}

func (s *configurationServiceV397) VerifyIqConnection(ctx context.Context) (*sonatyperepoV382.IqConnectionVerificationXo, *http.Response, error) {
	apiV397, httpResponse, err := s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.VerifyIqConnection(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.IqConnectionVerificationXo
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}
