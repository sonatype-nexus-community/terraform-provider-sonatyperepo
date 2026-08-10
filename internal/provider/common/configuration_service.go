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

// ConfigurationService abstracts Email, LDAP, and IQ Connection configuration
// APIs across NXRM API client generations. Every method's request/response shape
// is expressed in terms of the V382 generated types, matching the existing model
// mapping layer's expectations.
type ConfigurationService interface {
	// Email
	SetEmailConfiguration(ctx context.Context, body sonatyperepoV382.ApiEmailConfiguration) (*http.Response, error)
	GetEmailConfiguration(ctx context.Context) (*sonatyperepoV382.ApiEmailConfiguration, *http.Response, error)
	DeleteEmailConfiguration(ctx context.Context) (*http.Response, error)

	// LDAP
	CreateLdapServer(ctx context.Context, body sonatyperepoV382.CreateLdapServerXo) (*http.Response, error)
	GetLdapServer(ctx context.Context, name string) (*sonatyperepoV382.ReadLdapServerXo, *http.Response, error)
	UpdateLdapServer(ctx context.Context, name string, body sonatyperepoV382.UpdateLdapServerXo) (*http.Response, error)
	DeleteLdapServer(ctx context.Context, name string) (*http.Response, error)

	// IQ Connection
	GetIqConnectionConfiguration(ctx context.Context) (*sonatyperepoV382.IqConnectionXo, *http.Response, error)
	UpdateIqConnectionConfiguration(ctx context.Context, body sonatyperepoV382.IqConnectionXo) (*sonatyperepoV382.IqConnectionXo, *http.Response, error)
	DisableIq(ctx context.Context) (*http.Response, error)
	VerifyIqConnection(ctx context.Context) (*sonatyperepoV382.IqConnectionVerificationXo, *http.Response, error)
}

// configurationServiceV382 implements ConfigurationService against NXRM API client V382
// (targets NXRM < 3.94.0).
type configurationServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

// configurationServiceV395 implements ConfigurationService against NXRM API client V395
// (targets NXRM 3.94.0+). Requests/responses are bridged to/from V382 shapes via JSON,
// since the internal/provider/model package's mapping methods are written in terms of
// V382 types.
type configurationServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

// Email API implementation - V382

func (s *configurationServiceV382) SetEmailConfiguration(ctx context.Context, body sonatyperepoV382.ApiEmailConfiguration) (*http.Response, error) {
	return s.client.EmailAPI.SetEmailConfiguration(ctx).Body(body).Execute()
}

func (s *configurationServiceV382) GetEmailConfiguration(ctx context.Context) (*sonatyperepoV382.ApiEmailConfiguration, *http.Response, error) {
	return s.client.EmailAPI.GetEmailConfiguration(ctx).Execute()
}

func (s *configurationServiceV382) DeleteEmailConfiguration(ctx context.Context) (*http.Response, error) {
	return s.client.EmailAPI.DeleteEmailConfiguration(ctx).Execute()
}

// Email API implementation - V395

func (s *configurationServiceV395) SetEmailConfiguration(ctx context.Context, body sonatyperepoV382.ApiEmailConfiguration) (*http.Response, error) {
	var v395Body sonatyperepoV395.ApiEmailConfiguration
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.EmailAPI.UpdateEmail(ctx).ApiEmailConfiguration(v395Body).Execute()
}

func (s *configurationServiceV395) GetEmailConfiguration(ctx context.Context) (*sonatyperepoV382.ApiEmailConfiguration, *http.Response, error) {
	apiV395, httpResponse, err := s.client.EmailAPI.ListEmail(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.ApiEmailConfiguration
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *configurationServiceV395) DeleteEmailConfiguration(ctx context.Context) (*http.Response, error) {
	return s.client.EmailAPI.DeleteEmail(ctx).Execute()
}

// LDAP API implementation - V382

func (s *configurationServiceV382) CreateLdapServer(ctx context.Context, body sonatyperepoV382.CreateLdapServerXo) (*http.Response, error) {
	return s.client.SecurityManagementLDAPAPI.CreateLdapServer(ctx).Body(body).Execute()
}

func (s *configurationServiceV382) GetLdapServer(ctx context.Context, name string) (*sonatyperepoV382.ReadLdapServerXo, *http.Response, error) {
	return s.client.SecurityManagementLDAPAPI.GetLdapServer(ctx, name).Execute()
}

func (s *configurationServiceV382) UpdateLdapServer(ctx context.Context, name string, body sonatyperepoV382.UpdateLdapServerXo) (*http.Response, error) {
	return s.client.SecurityManagementLDAPAPI.UpdateLdapServer(ctx, name).Body(body).Execute()
}

func (s *configurationServiceV382) DeleteLdapServer(ctx context.Context, name string) (*http.Response, error) {
	return s.client.SecurityManagementLDAPAPI.DeleteLdapServer(ctx, name).Execute()
}

// LDAP API implementation - V395

func (s *configurationServiceV395) CreateLdapServer(ctx context.Context, body sonatyperepoV382.CreateLdapServerXo) (*http.Response, error) {
	var v395Body sonatyperepoV395.CreateLdapServerXo
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.SecurityManagementLDAPAPI.CreateSecurityLdap(ctx).CreateLdapServerXo(v395Body).Execute()
}

func (s *configurationServiceV395) GetLdapServer(ctx context.Context, name string) (*sonatyperepoV382.ReadLdapServerXo, *http.Response, error) {
	apiV395, httpResponse, err := s.client.SecurityManagementLDAPAPI.GetSecurityLdap(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.ReadLdapServerXo
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *configurationServiceV395) UpdateLdapServer(ctx context.Context, name string, body sonatyperepoV382.UpdateLdapServerXo) (*http.Response, error) {
	var v395Body sonatyperepoV395.UpdateLdapServerXo
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.SecurityManagementLDAPAPI.UpdateSecurityLdap(ctx, name).UpdateLdapServerXo(v395Body).Execute()
}

func (s *configurationServiceV395) DeleteLdapServer(ctx context.Context, name string) (*http.Response, error) {
	return s.client.SecurityManagementLDAPAPI.DeleteSecurityLdap(ctx, name).Execute()
}

// IQ Connection API implementation - V382

func (s *configurationServiceV382) GetIqConnectionConfiguration(ctx context.Context) (*sonatyperepoV382.IqConnectionXo, *http.Response, error) {
	return s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.GetConfiguration1(ctx).Execute()
}

func (s *configurationServiceV382) UpdateIqConnectionConfiguration(ctx context.Context, body sonatyperepoV382.IqConnectionXo) (*sonatyperepoV382.IqConnectionXo, *http.Response, error) {
	return s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.UpdateConfiguration1(ctx).Body(body).Execute()
}

func (s *configurationServiceV382) DisableIq(ctx context.Context) (*http.Response, error) {
	return s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.DisableIq(ctx).Execute()
}

func (s *configurationServiceV382) VerifyIqConnection(ctx context.Context) (*sonatyperepoV382.IqConnectionVerificationXo, *http.Response, error) {
	return s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.VerifyIqConnection(ctx).Execute()
}

// IQ Connection API implementation - V395

func (s *configurationServiceV395) GetIqConnectionConfiguration(ctx context.Context) (*sonatyperepoV382.IqConnectionXo, *http.Response, error) {
	apiV395, httpResponse, err := s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.ListIq(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.IqConnectionXo
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *configurationServiceV395) UpdateIqConnectionConfiguration(ctx context.Context, body sonatyperepoV382.IqConnectionXo) (*sonatyperepoV382.IqConnectionXo, *http.Response, error) {
	var v395Body sonatyperepoV395.IqConnectionXo
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, nil, err
	}
	apiV395, httpResponse, err := s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.UpdateIq(ctx).IqConnectionXo(v395Body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.IqConnectionXo
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *configurationServiceV395) DisableIq(ctx context.Context) (*http.Response, error) {
	return s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.CreateIqDisable(ctx).Execute()
}

func (s *configurationServiceV395) VerifyIqConnection(ctx context.Context) (*sonatyperepoV382.IqConnectionVerificationXo, *http.Response, error) {
	apiV395, httpResponse, err := s.client.ManageSonatypeRepositoryFirewallConfigurationAPI.VerifyIqConnection(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.IqConnectionVerificationXo
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}
