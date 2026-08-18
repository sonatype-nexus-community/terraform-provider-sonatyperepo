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
	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// Services bundles every domain service adapter, resolved once per provider
// configuration to the NXRM API client generation matching the connected server.
type Services struct {
	Status          StatusService
	Task            TaskService
	Repository      RepositoryManagementService
	BlobStore       BlobStoreService
	Capability      CapabilityService
	Privilege       PrivilegeService
	ContentSelector ContentSelectorService
	Role            RoleService
	User            UserService
	Configuration   ConfigurationService
	RoutingRule     RoutingRuleService
	CleanupPolicy   CleanupPolicyService
	Saml            SamlService
	Realms          RealmsService
	AnonymousAccess AnonymousAccessService
	SsrfProtection  SsrfProtectionService
	UserTokens      UserTokensService
	Certificates    CertificatesService
	License         LicenseService
	HttpSettings    HttpSettingsService
	OAuth2          OAuth2Service
}

// NewServices resolves every domain service to the client generation appropriate for
// the connected NXRM server's version. NXRM 3.94.0 and later use V395; anything older
// uses V382. This is the only place a version is compared or a concrete client type is named.
func NewServices(version SystemVersion, clientV382 *sonatyperepoV382.APIClient, clientV395 *sonatyperepoV395.APIClient) Services {
	// SystemVersion.Patch/Build are int8 (max 127); 127 is used rather than 999 (which
	// would overflow int8 and wrap negative, defeating the comparison) to ensure any
	// 3.93.x.y patch/build stays on V382 and only 3.94.0+ selects V395.
	if version.NewerThan(3, 93, 127, 127) {
		return Services{
			Status:          &statusServiceV395{client: clientV395},
			Task:            &taskServiceV395{client: clientV395},
			Repository:      &repositoryManagementServiceV395{client: clientV395},
			BlobStore:       &blobStoreServiceV395{client: clientV395},
			Capability:      &capabilityServiceV395{client: clientV395},
			Privilege:       &privilegeServiceV395{client: clientV395},
			ContentSelector: &contentSelectorServiceV395{client: clientV395},
			Role:            &roleServiceV395{client: clientV395},
			User:            &userServiceV395{client: clientV395},
			Configuration:   &configurationServiceV395{client: clientV395},
			RoutingRule:     &routingRuleServiceV395{client: clientV395},
			CleanupPolicy:   NewCleanupPolicyServiceV395(clientV395),
			Saml:            NewSamlServiceV395(clientV395),
			Realms:          NewRealmsServiceV395(clientV395),
			AnonymousAccess: NewAnonymousAccessServiceV395(clientV395),
			SsrfProtection:  NewSsrfProtectionServiceV395(clientV395),
			UserTokens:      NewUserTokensServiceV395(clientV395),
			Certificates:    NewCertificatesServiceV395(clientV395),
			License:         NewLicenseServiceV395(clientV395),
			HttpSettings:    NewHttpSettingsServiceV395(clientV395),
			OAuth2:          NewOAuth2ServiceV395(clientV395),
		}
	}
	return Services{
		Status:          &statusServiceV382{client: clientV382},
		Task:            &taskServiceV382{client: clientV382},
		Repository:      &repositoryManagementServiceV382{client: clientV382},
		BlobStore:       &blobStoreServiceV382{client: clientV382},
		Capability:      &capabilityServiceV382{client: clientV382},
		Privilege:       &privilegeServiceV382{client: clientV382},
		ContentSelector: &contentSelectorServiceV382{client: clientV382},
		Role:            &roleServiceV382{client: clientV382},
		User:            &userServiceV382{client: clientV382},
		Configuration:   &configurationServiceV382{client: clientV382},
		RoutingRule:     &routingRuleServiceV382{client: clientV382},
		CleanupPolicy:   NewCleanupPolicyServiceV382(clientV382),
		Saml:            NewSamlServiceV382(clientV382),
		Realms:          NewRealmsServiceV382(clientV382),
		AnonymousAccess: NewAnonymousAccessServiceV382(clientV382),
		SsrfProtection:  NewSsrfProtectionServiceV382(clientV382),
		UserTokens:      NewUserTokensServiceV382(clientV382),
		Certificates:    NewCertificatesServiceV382(clientV382),
		License:         NewLicenseServiceV382(clientV382),
		HttpSettings:    NewHttpSettingsServiceV382(clientV382),
		OAuth2:          NewOAuth2ServiceUnsupported(),
	}
}
