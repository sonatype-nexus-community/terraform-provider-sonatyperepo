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

// ToApiUpdateModel methods are already written against. The V397 adapter bridges to/from
// repositoryManagementServiceV397 implements RepositoryManagementService against NXRM API
// client V397 (targets NXRM 3.97.0+). Requests/responses are bridged to/from V382 shapes via
//   - RoutingRule (V382) -> RoutingRuleName (V397), on the write side of proxy formats.
//   - npm/pypi proxy "RemoveQuarantined" (PCCS) has no equivalent field on the V397 repository
type repositoryManagementServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *repositoryManagementServiceV397) DeleteRepository(ctx context.Context, repositoryName string) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.DeleteRepositories(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV397) ListRepositories(ctx context.Context) ([]sonatyperepoV382.RepositoryXO, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAllRepositories(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.RepositoryXO
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) CreateAlpineGroupRepository(ctx context.Context, body sonatyperepoV382.AlpineGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AlpineGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAlpineGroupRepository(ctx).AlpineGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetAlpineGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAlpineGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AlpineGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateAlpineGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AlpineGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAlpineGroupRepository(ctx, repositoryName).AlpineGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateAlpineHostedRepository(ctx context.Context, body sonatyperepoV382.AlpineHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AlpineHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAlpineHostedRepository(ctx).AlpineHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetAlpineHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineHostedApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAlpineHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AlpineHostedApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateAlpineHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AlpineHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAlpineHostedRepository(ctx, repositoryName).AlpineHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateAlpineProxyRepository(ctx context.Context, body sonatyperepoV382.AlpineProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.AlpineProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateAlpineProxyRepository(ctx).AlpineProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetAlpineProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAlpineProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.AlpineProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateAlpineProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.AlpineProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateAlpineProxyRepository(ctx, repositoryName).AlpineProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateAnsiblegalaxyGroupRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AnsibleGalaxyGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyGroupRepository(ctx).AnsibleGalaxyGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetAnsiblegalaxyGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAnsiblegalaxyGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnsibleGalaxyGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateAnsiblegalaxyGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AnsibleGalaxyGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyGroupRepository(ctx, repositoryName).AnsibleGalaxyGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateAnsiblegalaxyHostedRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AnsibleGalaxyHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyHostedRepository(ctx).AnsibleGalaxyHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetAnsiblegalaxyHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyHostedApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAnsiblegalaxyHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnsibleGalaxyHostedApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateAnsiblegalaxyHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AnsibleGalaxyHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyHostedRepository(ctx, repositoryName).AnsibleGalaxyHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateAnsiblegalaxyProxyRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AnsibleGalaxyProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyProxyRepository(ctx).AnsibleGalaxyProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetAnsiblegalaxyProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyProxyApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAnsiblegalaxyProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnsibleGalaxyProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateAnsiblegalaxyProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AnsibleGalaxyProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyProxyRepository(ctx, repositoryName).AnsibleGalaxyProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateAptHostedRepository(ctx context.Context, body sonatyperepoV382.AptHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AptHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAptHostedRepository(ctx).AptHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetAptHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AptHostedApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAptHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AptHostedApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateAptHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AptHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AptHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAptHostedRepository(ctx, repositoryName).AptHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateAptProxyRepository(ctx context.Context, body sonatyperepoV382.AptProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AptProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateAptProxyRepository(ctx).AptProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetAptProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AptProxyApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetAptProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AptProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateAptProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AptProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.AptProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateAptProxyRepository(ctx, repositoryName).AptProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateCargoGroupRepository(ctx context.Context, body sonatyperepoV382.CargoGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.CargoGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateCargoGroupRepository(ctx).CargoGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetCargoGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.CargoGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetCargoGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.CargoGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateCargoGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.CargoGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateCargoGroupRepository(ctx, repositoryName).CargoGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateCargoHostedRepository(ctx context.Context, body sonatyperepoV382.CargoHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.CargoHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateCargoHostedRepository(ctx).CargoHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetCargoHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetCargoHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateCargoHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.CargoHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateCargoHostedRepository(ctx, repositoryName).CargoHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateCargoProxyRepository(ctx context.Context, body sonatyperepoV382.CargoProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.CargoProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateCargoProxyRepository(ctx).CargoProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetCargoProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.CargoProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetCargoProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil && apiV397.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.CargoProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateCargoProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.CargoProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateCargoProxyRepository(ctx, repositoryName).CargoProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateCocoapodsProxyRepository(ctx context.Context, body sonatyperepoV382.CocoapodsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.CocoapodsProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateCocoapodsProxyRepository(ctx).CocoapodsProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetCocoapodsProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetCocoapodsProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateCocoapodsProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CocoapodsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.CocoapodsProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateCocoapodsProxyRepository(ctx, repositoryName).CocoapodsProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateComposerProxyRepository(ctx context.Context, body sonatyperepoV382.ComposerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.ComposerProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateComposerProxyRepository(ctx).ComposerProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetComposerProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetComposerProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateComposerProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ComposerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.ComposerProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateComposerProxyRepository(ctx, repositoryName).ComposerProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateConanGroupRepository(ctx context.Context, body sonatyperepoV382.ConanGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.ConanGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateConanGroupRepository(ctx).ConanGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetConanGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetConanGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupDeployRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateConanGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.ConanGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateConanGroupRepository(ctx, repositoryName).ConanGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateConanHostedRepository(ctx context.Context, body sonatyperepoV382.ConanHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.ConanHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateConanHostedRepository(ctx).ConanHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetConanHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetConanHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateConanHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.ConanHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateConanHostedRepository(ctx, repositoryName).ConanHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateConanProxyRepository(ctx context.Context, body sonatyperepoV382.ConanProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.ConanProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateConanProxyRepository(ctx).ConanProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetConanProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.ConanProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetConanProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.ConanProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateConanProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.ConanProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateConanProxyRepository(ctx, repositoryName).ConanProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateCondaProxyRepository(ctx context.Context, body sonatyperepoV382.CondaProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.CondaProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateCondaProxyRepository(ctx).CondaProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetCondaProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetCondaProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateCondaProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CondaProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.CondaProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateCondaProxyRepository(ctx, repositoryName).CondaProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateDockerGroupRepository(ctx context.Context, body sonatyperepoV382.DockerGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.DockerGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateDockerGroupRepository(ctx).DockerGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetDockerGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetDockerGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.DockerGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateDockerGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.DockerGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateDockerGroupRepository(ctx, repositoryName).DockerGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateDockerHostedRepository(ctx context.Context, body sonatyperepoV382.DockerHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.DockerHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateDockerHostedRepository(ctx).DockerHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetDockerHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerHostedApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetDockerHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.DockerHostedApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateDockerHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.DockerHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateDockerHostedRepository(ctx, repositoryName).DockerHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateDockerProxyRepository(ctx context.Context, body sonatyperepoV382.DockerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.DockerProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateDockerProxyRepository(ctx).DockerProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetDockerProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetDockerProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil && apiV397.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.DockerProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateDockerProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.DockerProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateDockerProxyRepository(ctx, repositoryName).DockerProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateGitlfsHostedRepository(ctx context.Context, body sonatyperepoV382.GitLfsHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.GitLfsHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateGitlfsHostedRepository(ctx).GitLfsHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetGitlfsHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetGitlfsHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateGitlfsHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GitLfsHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.GitLfsHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateGitlfsHostedRepository(ctx, repositoryName).GitLfsHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateGoGroupRepository(ctx context.Context, body sonatyperepoV382.GolangGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.GolangGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateGoGroupRepository(ctx).GolangGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetGoGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetGoGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateGoGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.GolangGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateGoGroupRepository(ctx, repositoryName).GolangGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateGoHostedRepository(ctx context.Context, body sonatyperepoV382.GolangHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.GolangHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateGoHostedRepository(ctx).GolangHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetGoHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetGoHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateGoHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.GolangHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateGoHostedRepository(ctx, repositoryName).GolangHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateGoProxyRepository(ctx context.Context, body sonatyperepoV382.GolangProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.GolangProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateGoProxyRepository(ctx).GolangProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetGoProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetGoProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil && apiV397.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateGoProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.GolangProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateGoProxyRepository(ctx, repositoryName).GolangProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateHelmGroupRepository(ctx context.Context, body sonatyperepoV382.HelmGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.HelmGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateHelmGroupRepository(ctx).HelmGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetHelmGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetHelmGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateHelmGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.HelmGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateHelmGroupRepository(ctx, repositoryName).HelmGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateHelmHostedRepository(ctx context.Context, body sonatyperepoV382.HelmHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.HelmHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateHelmHostedRepository(ctx).HelmHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetHelmHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetHelmHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateHelmHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.HelmHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateHelmHostedRepository(ctx, repositoryName).HelmHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateHelmProxyRepository(ctx context.Context, body sonatyperepoV382.HelmProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.HelmProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateHelmProxyRepository(ctx).HelmProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetHelmProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetHelmProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateHelmProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.HelmProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateHelmProxyRepository(ctx, repositoryName).HelmProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateHuggingfaceProxyRepository(ctx context.Context, body sonatyperepoV382.HuggingFaceProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.HuggingFaceProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateHuggingfaceProxyRepository(ctx).HuggingFaceProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetHuggingfaceProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetHuggingfaceProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil && apiV397.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateHuggingfaceProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HuggingFaceProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.HuggingFaceProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateHuggingfaceProxyRepository(ctx, repositoryName).HuggingFaceProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateMavenGroupRepository(ctx context.Context, body sonatyperepoV382.MavenGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.MavenGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateMavenGroupRepository(ctx).MavenGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetMavenGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetMavenGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateMavenGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.MavenGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateMavenGroupRepository(ctx, repositoryName).MavenGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateMavenHostedRepository(ctx context.Context, body sonatyperepoV382.MavenHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.MavenHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateMavenHostedRepository(ctx).MavenHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetMavenHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.MavenHostedApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetMavenHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.MavenHostedApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateMavenHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.MavenHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateMavenHostedRepository(ctx, repositoryName).MavenHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateMavenProxyRepository(ctx context.Context, body sonatyperepoV382.MavenProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.MavenProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateMavenProxyRepository(ctx).MavenProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetMavenProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.MavenProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetMavenProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.MavenProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateMavenProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.MavenProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateMavenProxyRepository(ctx, repositoryName).MavenProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateNpmGroupRepository(ctx context.Context, body sonatyperepoV382.NpmGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.NpmGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateNpmGroupRepository(ctx).NpmGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetNpmGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetNpmGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupDeployRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateNpmGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.NpmGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateNpmGroupRepository(ctx, repositoryName).NpmGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateNpmHostedRepository(ctx context.Context, body sonatyperepoV382.NpmHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.NpmHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateNpmHostedRepository(ctx).NpmHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetNpmHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetNpmHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateNpmHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.NpmHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateNpmHostedRepository(ctx, repositoryName).NpmHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateNpmProxyRepository(ctx context.Context, body sonatyperepoV382.NpmProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.NpmProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateNpmProxyRepository(ctx).NpmProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetNpmProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.NpmProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetNpmProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var result sonatyperepoV382.NpmProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	var mode *FirewallMode
	if apiV397.Firewall != nil && apiV397.Firewall.Mode != nil {
		m := FirewallMode(*apiV397.Firewall.Mode)
		mode = &m
	}
	return &result, mode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateNpmProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.NpmProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateNpmProxyRepository(ctx, repositoryName).NpmProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateNugetGroupRepository(ctx context.Context, body sonatyperepoV382.NugetGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.NugetGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateNugetGroupRepository(ctx).NugetGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetNugetGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetNugetGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateNugetGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.NugetGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateNugetGroupRepository(ctx, repositoryName).NugetGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateNugetHostedRepository(ctx context.Context, body sonatyperepoV382.NugetHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.NugetHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateNugetHostedRepository(ctx).NugetHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetNugetHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetNugetHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateNugetHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.NugetHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateNugetHostedRepository(ctx, repositoryName).NugetHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateNugetProxyRepository(ctx context.Context, body sonatyperepoV382.NugetProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.NugetProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateNugetProxyRepository(ctx).NugetProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetNugetProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.NugetProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetNugetProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.NugetProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateNugetProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.NugetProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateNugetProxyRepository(ctx, repositoryName).NugetProxyRepositoryApiRequest(v397Body).Execute()
}

// other proxy/hosted/group format's V397 adapter, just without a real vendored V382 counterpart.
func (s *repositoryManagementServiceV397) CreateOciGroupRepository(ctx context.Context, body OciGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.OciGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateOciGroupRepository(ctx).OciGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetOciGroupRepository(ctx context.Context, repositoryName string) (*OciGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetOciGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result OciGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateOciGroupRepository(ctx context.Context, repositoryName string, body OciGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.OciGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateOciGroupRepository(ctx, repositoryName).OciGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateOciHostedRepository(ctx context.Context, body OciHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.OciHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateOciHostedRepository(ctx).OciHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetOciHostedRepository(ctx context.Context, repositoryName string) (*OciHostedApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetOciHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result OciHostedApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateOciHostedRepository(ctx context.Context, repositoryName string, body OciHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.OciHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateOciHostedRepository(ctx, repositoryName).OciHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateOciProxyRepository(ctx context.Context, body OciProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.OciProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateOciProxyRepository(ctx).OciProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetOciProxyRepository(ctx context.Context, repositoryName string) (*OciProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetOciProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result OciProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateOciProxyRepository(ctx context.Context, repositoryName string, body OciProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.OciProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateOciProxyRepository(ctx, repositoryName).OciProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateP2ProxyRepository(ctx context.Context, body sonatyperepoV382.P2ProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.P2ProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateP2ProxyRepository(ctx).P2ProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetP2ProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetP2ProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateP2ProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.P2ProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.P2ProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateP2ProxyRepository(ctx, repositoryName).P2ProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreatePubGroupRepository(ctx context.Context, body sonatyperepoV382.PubGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.PubGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreatePubGroupRepository(ctx).PubGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetPubGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetPubGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdatePubGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.PubGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdatePubGroupRepository(ctx, repositoryName).PubGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreatePubHostedRepository(ctx context.Context, body sonatyperepoV382.PubHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.PubHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreatePubHostedRepository(ctx).PubHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetPubHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetPubHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdatePubHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.PubHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdatePubHostedRepository(ctx, repositoryName).PubHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreatePubProxyRepository(ctx context.Context, body sonatyperepoV382.PubProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.PubProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreatePubProxyRepository(ctx).PubProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetPubProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetPubProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil && apiV397.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdatePubProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.PubProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdatePubProxyRepository(ctx, repositoryName).PubProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreatePypiGroupRepository(ctx context.Context, body sonatyperepoV382.PypiGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.PypiGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreatePypiGroupRepository(ctx).PypiGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetPypiGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetPypiGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupDeployRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdatePypiGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.PypiGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdatePypiGroupRepository(ctx, repositoryName).PypiGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreatePypiHostedRepository(ctx context.Context, body sonatyperepoV382.PypiHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.PypiHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreatePypiHostedRepository(ctx).PypiHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetPypiHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetPypiHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdatePypiHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.PypiHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdatePypiHostedRepository(ctx, repositoryName).PypiHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreatePypiProxyRepository(ctx context.Context, body sonatyperepoV382.PypiProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.PypiProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreatePypiProxyRepository(ctx).PypiProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetPypiProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.PyPiProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetPypiProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.PyPiProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdatePypiProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.PypiProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdatePypiProxyRepository(ctx, repositoryName).PypiProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRGroupRepository(ctx context.Context, body sonatyperepoV382.RGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRGroupRepository(ctx).RGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRGroupRepository(ctx, repositoryName).RGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRHostedRepository(ctx context.Context, body sonatyperepoV382.RHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRHostedRepository(ctx).RHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRHostedRepository(ctx, repositoryName).RHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRProxyRepository(ctx context.Context, body sonatyperepoV382.RProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.RProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateRProxyRepository(ctx).RProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.RProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateRProxyRepository(ctx, repositoryName).RProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRawGroupRepository(ctx context.Context, body sonatyperepoV382.RawGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RawGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRawGroupRepository(ctx).RawGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRawGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRawGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RawGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRawGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RawGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRawGroupRepository(ctx, repositoryName).RawGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRawHostedRepository(ctx context.Context, body sonatyperepoV382.RawHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RawHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRawHostedRepository(ctx).RawHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRawHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawHostedApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRawHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RawHostedApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRawHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RawHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRawHostedRepository(ctx, repositoryName).RawHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRawProxyRepository(ctx context.Context, body sonatyperepoV382.RawProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.RawProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateRawProxyRepository(ctx).RawProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRawProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRawProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.RawProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRawProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.RawProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateRawProxyRepository(ctx, repositoryName).RawProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRubygemsGroupRepository(ctx context.Context, body sonatyperepoV382.RubyGemsGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RubyGemsGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRubygemsGroupRepository(ctx).RubyGemsGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRubygemsGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRubygemsGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRubygemsGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RubyGemsGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRubygemsGroupRepository(ctx, repositoryName).RubyGemsGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRubygemsHostedRepository(ctx context.Context, body sonatyperepoV382.RubyGemsHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RubyGemsHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRubygemsHostedRepository(ctx).RubyGemsHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRubygemsHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRubygemsHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRubygemsHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.RubyGemsHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRubygemsHostedRepository(ctx, repositoryName).RubyGemsHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateRubygemsProxyRepository(ctx context.Context, body sonatyperepoV382.RubyGemsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.RubyGemsProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateRubygemsProxyRepository(ctx).RubyGemsProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetRubygemsProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetRubygemsProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateRubygemsProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.RubyGemsProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateRubygemsProxyRepository(ctx, repositoryName).RubyGemsProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateSwiftGroupRepository(ctx context.Context, body sonatyperepoV382.SwiftGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.SwiftGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateSwiftGroupRepository(ctx).SwiftGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetSwiftGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SwiftGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetSwiftGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SwiftGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateSwiftGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.SwiftGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.SwiftGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateSwiftGroupRepository(ctx, repositoryName).SwiftGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateSwiftProxyRepository(ctx context.Context, body sonatyperepoV382.SwiftProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.SwiftProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateSwiftProxyRepository(ctx).SwiftProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetSwiftProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SwiftProxyApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetSwiftProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SwiftProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateSwiftProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.SwiftProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.SwiftProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateSwiftProxyRepository(ctx, repositoryName).SwiftProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateTerraformGroupRepository(ctx context.Context, body sonatyperepoV382.TerraformGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.TerraformGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateTerraformGroupRepository(ctx).TerraformGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetTerraformGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetTerraformGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.TerraformGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateTerraformGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.TerraformGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateTerraformGroupRepository(ctx, repositoryName).TerraformGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateTerraformHostedRepository(ctx context.Context, body sonatyperepoV382.TerraformHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.TerraformHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateTerraformHostedRepository(ctx).TerraformHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetTerraformHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformHostedRepositoryApiRequest, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetTerraformHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.TerraformHostedRepositoryApiRequest
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateTerraformHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.TerraformHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateTerraformHostedRepository(ctx, repositoryName).TerraformHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateTerraformProxyRepository(ctx context.Context, body sonatyperepoV382.TerraformProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.TerraformProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateTerraformProxyRepository(ctx).TerraformProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetTerraformProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformProxyApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetTerraformProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.TerraformProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateTerraformProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformProxyRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.TerraformProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateTerraformProxyRepository(ctx, repositoryName).TerraformProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateYumGroupRepository(ctx context.Context, body sonatyperepoV382.YumGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.YumGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateYumGroupRepository(ctx).YumGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetYumGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumGroupApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetYumGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.YumGroupApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateYumGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumGroupRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.YumGroupRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateYumGroupRepository(ctx, repositoryName).YumGroupRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateYumHostedRepository(ctx context.Context, body sonatyperepoV382.YumHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.YumHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateYumHostedRepository(ctx).YumHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetYumHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumHostedApiRepository, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetYumHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.YumHostedApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateYumHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumHostedRepositoryApiRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.YumHostedRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateYumHostedRepository(ctx, repositoryName).YumHostedRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) CreateYumProxyRepository(ctx context.Context, body sonatyperepoV382.YumProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.YumProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateYumProxyRepository(ctx).YumProxyRepositoryApiRequest(v397Body).Execute()
}

func (s *repositoryManagementServiceV397) GetYumProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV397, httpResponse, err := s.client.RepositoryManagementAPI.GetYumProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV397.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV397.Firewall.Mode)
	}
	var result sonatyperepoV382.YumProxyApiRepository
	if err := jsonBridge(apiV397, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV397) UpdateYumProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v397Body sonatyperepoV397.YumProxyRepositoryApiRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	v397Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v397Body.Firewall = &sonatyperepoV397.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateYumProxyRepository(ctx, repositoryName).YumProxyRepositoryApiRequest(v397Body).Execute()
}
