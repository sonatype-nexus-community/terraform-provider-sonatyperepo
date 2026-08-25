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

// RepositoryManagementService abstracts repository CRUD across NXRM API client
// generations. Every method's request/response shape is expressed in terms of the V382
// ("github.com/.../nexus-repo-api-client-go/v3") generated types -- the vocabulary the
// existing internal/provider/model package's per-format FromApiModel/ToApiCreateModel/
// ToApiUpdateModel methods are already written against. The V395 adapter bridges to/from
// its own generated types internally so that callers (internal/provider/repository/format)
// never need to know which generation is behind the interface.
type RepositoryManagementService interface {
	DeleteRepository(ctx context.Context, repositoryName string) (*http.Response, error)
	ListRepositories(ctx context.Context) ([]sonatyperepoV382.RepositoryXO, *http.Response, error)
	CreateAlpineGroupRepository(ctx context.Context, body sonatyperepoV382.AlpineGroupRepositoryApiRequest) (*http.Response, error)
	GetAlpineGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineGroupApiRepository, *http.Response, error)
	UpdateAlpineGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineGroupRepositoryApiRequest) (*http.Response, error)
	CreateAlpineHostedRepository(ctx context.Context, body sonatyperepoV382.AlpineHostedRepositoryApiRequest) (*http.Response, error)
	GetAlpineHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineHostedApiRepository, *http.Response, error)
	UpdateAlpineHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineHostedRepositoryApiRequest) (*http.Response, error)
	CreateAlpineProxyRepository(ctx context.Context, body sonatyperepoV382.AlpineProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetAlpineProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateAlpineProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateAnsiblegalaxyGroupRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyGroupRepositoryApiRequest) (*http.Response, error)
	GetAnsiblegalaxyGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyGroupApiRepository, *http.Response, error)
	UpdateAnsiblegalaxyGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyGroupRepositoryApiRequest) (*http.Response, error)
	CreateAnsiblegalaxyHostedRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyHostedRepositoryApiRequest) (*http.Response, error)
	GetAnsiblegalaxyHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyHostedApiRepository, *http.Response, error)
	UpdateAnsiblegalaxyHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyHostedRepositoryApiRequest) (*http.Response, error)
	CreateAnsiblegalaxyProxyRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyProxyRepositoryApiRequest) (*http.Response, error)
	GetAnsiblegalaxyProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyProxyApiRepository, *http.Response, error)
	UpdateAnsiblegalaxyProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyProxyRepositoryApiRequest) (*http.Response, error)
	CreateAptHostedRepository(ctx context.Context, body sonatyperepoV382.AptHostedRepositoryApiRequest) (*http.Response, error)
	GetAptHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AptHostedApiRepository, *http.Response, error)
	UpdateAptHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AptHostedRepositoryApiRequest) (*http.Response, error)
	CreateAptProxyRepository(ctx context.Context, body sonatyperepoV382.AptProxyRepositoryApiRequest) (*http.Response, error)
	GetAptProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AptProxyApiRepository, *http.Response, error)
	UpdateAptProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AptProxyRepositoryApiRequest) (*http.Response, error)
	CreateCargoGroupRepository(ctx context.Context, body sonatyperepoV382.CargoGroupRepositoryApiRequest) (*http.Response, error)
	GetCargoGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.CargoGroupApiRepository, *http.Response, error)
	UpdateCargoGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoGroupRepositoryApiRequest) (*http.Response, error)
	CreateCargoHostedRepository(ctx context.Context, body sonatyperepoV382.CargoHostedRepositoryApiRequest) (*http.Response, error)
	GetCargoHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateCargoHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoHostedRepositoryApiRequest) (*http.Response, error)
	CreateCargoProxyRepository(ctx context.Context, body sonatyperepoV382.CargoProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetCargoProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.CargoProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateCargoProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateCocoapodsProxyRepository(ctx context.Context, body sonatyperepoV382.CocoapodsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetCocoapodsProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error)
	UpdateCocoapodsProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CocoapodsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateComposerProxyRepository(ctx context.Context, body sonatyperepoV382.ComposerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetComposerProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error)
	UpdateComposerProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ComposerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateConanGroupRepository(ctx context.Context, body sonatyperepoV382.ConanGroupRepositoryApiRequest) (*http.Response, error)
	GetConanGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error)
	UpdateConanGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanGroupRepositoryApiRequest) (*http.Response, error)
	CreateConanHostedRepository(ctx context.Context, body sonatyperepoV382.ConanHostedRepositoryApiRequest) (*http.Response, error)
	GetConanHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateConanHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanHostedRepositoryApiRequest) (*http.Response, error)
	CreateConanProxyRepository(ctx context.Context, body sonatyperepoV382.ConanProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetConanProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.ConanProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateConanProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateCondaProxyRepository(ctx context.Context, body sonatyperepoV382.CondaProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetCondaProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error)
	UpdateCondaProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CondaProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateDockerGroupRepository(ctx context.Context, body sonatyperepoV382.DockerGroupRepositoryApiRequest) (*http.Response, error)
	GetDockerGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerGroupApiRepository, *http.Response, error)
	UpdateDockerGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerGroupRepositoryApiRequest) (*http.Response, error)
	CreateDockerHostedRepository(ctx context.Context, body sonatyperepoV382.DockerHostedRepositoryApiRequest) (*http.Response, error)
	GetDockerHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerHostedApiRepository, *http.Response, error)
	UpdateDockerHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerHostedRepositoryApiRequest) (*http.Response, error)
	CreateDockerProxyRepository(ctx context.Context, body sonatyperepoV382.DockerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetDockerProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateDockerProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateGitlfsHostedRepository(ctx context.Context, body sonatyperepoV382.GitLfsHostedRepositoryApiRequest) (*http.Response, error)
	GetGitlfsHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateGitlfsHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GitLfsHostedRepositoryApiRequest) (*http.Response, error)
	CreateGoGroupRepository(ctx context.Context, body sonatyperepoV382.GolangGroupRepositoryApiRequest) (*http.Response, error)
	GetGoGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error)
	UpdateGoGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangGroupRepositoryApiRequest) (*http.Response, error)
	CreateGoHostedRepository(ctx context.Context, body sonatyperepoV382.GolangHostedRepositoryApiRequest) (*http.Response, error)
	GetGoHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateGoHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangHostedRepositoryApiRequest) (*http.Response, error)
	CreateGoProxyRepository(ctx context.Context, body sonatyperepoV382.GolangProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetGoProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error)
	UpdateGoProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateHelmGroupRepository(ctx context.Context, body sonatyperepoV382.HelmGroupRepositoryApiRequest) (*http.Response, error)
	GetHelmGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error)
	UpdateHelmGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmGroupRepositoryApiRequest) (*http.Response, error)
	CreateHelmHostedRepository(ctx context.Context, body sonatyperepoV382.HelmHostedRepositoryApiRequest) (*http.Response, error)
	GetHelmHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateHelmHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmHostedRepositoryApiRequest) (*http.Response, error)
	CreateHelmProxyRepository(ctx context.Context, body sonatyperepoV382.HelmProxyRepositoryApiRequest) (*http.Response, error)
	GetHelmProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *http.Response, error)
	UpdateHelmProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmProxyRepositoryApiRequest) (*http.Response, error)
	CreateHuggingfaceProxyRepository(ctx context.Context, body sonatyperepoV382.HuggingFaceProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetHuggingfaceProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error)
	UpdateHuggingfaceProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HuggingFaceProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateMavenGroupRepository(ctx context.Context, body sonatyperepoV382.MavenGroupRepositoryApiRequest) (*http.Response, error)
	GetMavenGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error)
	UpdateMavenGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenGroupRepositoryApiRequest) (*http.Response, error)
	CreateMavenHostedRepository(ctx context.Context, body sonatyperepoV382.MavenHostedRepositoryApiRequest) (*http.Response, error)
	GetMavenHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.MavenHostedApiRepository, *http.Response, error)
	UpdateMavenHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenHostedRepositoryApiRequest) (*http.Response, error)
	CreateMavenProxyRepository(ctx context.Context, body sonatyperepoV382.MavenProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetMavenProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.MavenProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateMavenProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateNpmGroupRepository(ctx context.Context, body sonatyperepoV382.NpmGroupRepositoryApiRequest) (*http.Response, error)
	GetNpmGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error)
	UpdateNpmGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmGroupRepositoryApiRequest) (*http.Response, error)
	CreateNpmHostedRepository(ctx context.Context, body sonatyperepoV382.NpmHostedRepositoryApiRequest) (*http.Response, error)
	GetNpmHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateNpmHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmHostedRepositoryApiRequest) (*http.Response, error)
	CreateNpmProxyRepository(ctx context.Context, body sonatyperepoV382.NpmProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetNpmProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.NpmProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateNpmProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateNugetGroupRepository(ctx context.Context, body sonatyperepoV382.NugetGroupRepositoryApiRequest) (*http.Response, error)
	GetNugetGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error)
	UpdateNugetGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetGroupRepositoryApiRequest) (*http.Response, error)
	CreateNugetHostedRepository(ctx context.Context, body sonatyperepoV382.NugetHostedRepositoryApiRequest) (*http.Response, error)
	GetNugetHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateNugetHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetHostedRepositoryApiRequest) (*http.Response, error)
	CreateNugetProxyRepository(ctx context.Context, body sonatyperepoV382.NugetProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetNugetProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.NugetProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateNugetProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateP2ProxyRepository(ctx context.Context, body sonatyperepoV382.P2ProxyRepositoryApiRequest) (*http.Response, error)
	GetP2ProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *http.Response, error)
	UpdateP2ProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.P2ProxyRepositoryApiRequest) (*http.Response, error)
	CreatePubGroupRepository(ctx context.Context, body sonatyperepoV382.PubGroupRepositoryApiRequest) (*http.Response, error)
	GetPubGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error)
	UpdatePubGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubGroupRepositoryApiRequest) (*http.Response, error)
	CreatePubHostedRepository(ctx context.Context, body sonatyperepoV382.PubHostedRepositoryApiRequest) (*http.Response, error)
	GetPubHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdatePubHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubHostedRepositoryApiRequest) (*http.Response, error)
	CreatePubProxyRepository(ctx context.Context, body sonatyperepoV382.PubProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetPubProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error)
	UpdatePubProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreatePypiGroupRepository(ctx context.Context, body sonatyperepoV382.PypiGroupRepositoryApiRequest) (*http.Response, error)
	GetPypiGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error)
	UpdatePypiGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiGroupRepositoryApiRequest) (*http.Response, error)
	CreatePypiHostedRepository(ctx context.Context, body sonatyperepoV382.PypiHostedRepositoryApiRequest) (*http.Response, error)
	GetPypiHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdatePypiHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiHostedRepositoryApiRequest) (*http.Response, error)
	CreatePypiProxyRepository(ctx context.Context, body sonatyperepoV382.PypiProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetPypiProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.PyPiProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdatePypiProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateRGroupRepository(ctx context.Context, body sonatyperepoV382.RGroupRepositoryApiRequest) (*http.Response, error)
	GetRGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error)
	UpdateRGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RGroupRepositoryApiRequest) (*http.Response, error)
	CreateRHostedRepository(ctx context.Context, body sonatyperepoV382.RHostedRepositoryApiRequest) (*http.Response, error)
	GetRHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateRHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RHostedRepositoryApiRequest) (*http.Response, error)
	CreateRProxyRepository(ctx context.Context, body sonatyperepoV382.RProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetRProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error)
	UpdateRProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateRawGroupRepository(ctx context.Context, body sonatyperepoV382.RawGroupRepositoryApiRequest) (*http.Response, error)
	GetRawGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawGroupApiRepository, *http.Response, error)
	UpdateRawGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawGroupRepositoryApiRequest) (*http.Response, error)
	CreateRawHostedRepository(ctx context.Context, body sonatyperepoV382.RawHostedRepositoryApiRequest) (*http.Response, error)
	GetRawHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawHostedApiRepository, *http.Response, error)
	UpdateRawHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawHostedRepositoryApiRequest) (*http.Response, error)
	CreateRawProxyRepository(ctx context.Context, body sonatyperepoV382.RawProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetRawProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateRawProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateRubygemsGroupRepository(ctx context.Context, body sonatyperepoV382.RubyGemsGroupRepositoryApiRequest) (*http.Response, error)
	GetRubygemsGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error)
	UpdateRubygemsGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsGroupRepositoryApiRequest) (*http.Response, error)
	CreateRubygemsHostedRepository(ctx context.Context, body sonatyperepoV382.RubyGemsHostedRepositoryApiRequest) (*http.Response, error)
	GetRubygemsHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error)
	UpdateRubygemsHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsHostedRepositoryApiRequest) (*http.Response, error)
	CreateRubygemsProxyRepository(ctx context.Context, body sonatyperepoV382.RubyGemsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetRubygemsProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error)
	UpdateRubygemsProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	CreateSwiftGroupRepository(ctx context.Context, body sonatyperepoV382.SwiftGroupRepositoryApiRequest) (*http.Response, error)
	GetSwiftGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SwiftGroupApiRepository, *http.Response, error)
	UpdateSwiftGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.SwiftGroupRepositoryApiRequest) (*http.Response, error)
	CreateSwiftProxyRepository(ctx context.Context, body sonatyperepoV382.SwiftProxyRepositoryApiRequest) (*http.Response, error)
	GetSwiftProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SwiftProxyApiRepository, *http.Response, error)
	UpdateSwiftProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.SwiftProxyRepositoryApiRequest) (*http.Response, error)
	CreateTerraformGroupRepository(ctx context.Context, body sonatyperepoV382.TerraformGroupRepositoryApiRequest) (*http.Response, error)
	GetTerraformGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformGroupApiRepository, *http.Response, error)
	UpdateTerraformGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformGroupRepositoryApiRequest) (*http.Response, error)
	CreateTerraformHostedRepository(ctx context.Context, body sonatyperepoV382.TerraformHostedRepositoryApiRequest) (*http.Response, error)
	GetTerraformHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformHostedRepositoryApiRequest, *http.Response, error)
	UpdateTerraformHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformHostedRepositoryApiRequest) (*http.Response, error)
	CreateTerraformProxyRepository(ctx context.Context, body sonatyperepoV382.TerraformProxyRepositoryApiRequest) (*http.Response, error)
	GetTerraformProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformProxyApiRepository, *http.Response, error)
	UpdateTerraformProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformProxyRepositoryApiRequest) (*http.Response, error)
	CreateYumGroupRepository(ctx context.Context, body sonatyperepoV382.YumGroupRepositoryApiRequest) (*http.Response, error)
	GetYumGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumGroupApiRepository, *http.Response, error)
	UpdateYumGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumGroupRepositoryApiRequest) (*http.Response, error)
	CreateYumHostedRepository(ctx context.Context, body sonatyperepoV382.YumHostedRepositoryApiRequest) (*http.Response, error)
	GetYumHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumHostedApiRepository, *http.Response, error)
	UpdateYumHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumHostedRepositoryApiRequest) (*http.Response, error)
	CreateYumProxyRepository(ctx context.Context, body sonatyperepoV382.YumProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
	GetYumProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumProxyApiRepository, *FirewallMode, *http.Response, error)
	UpdateYumProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error)
}

// repositoryManagementServiceV382 implements RepositoryManagementService against NXRM API
// client V382 (targets NXRM < 3.94.0). Every method is a direct, unmodified delegation to
// the generated client -- V382 types are this interface's native vocabulary.
type repositoryManagementServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *repositoryManagementServiceV382) DeleteRepository(ctx context.Context, repositoryName string) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.DeleteRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) ListRepositories(ctx context.Context) ([]sonatyperepoV382.RepositoryXO, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetAllRepositories(ctx).Execute()
}

func (s *repositoryManagementServiceV382) CreateAlpineGroupRepository(ctx context.Context, body sonatyperepoV382.AlpineGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateAlpineGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetAlpineGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineGroupApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetAlpineGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateAlpineGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateAlpineGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateAlpineHostedRepository(ctx context.Context, body sonatyperepoV382.AlpineHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateAlpineHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetAlpineHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineHostedApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetAlpineHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateAlpineHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateAlpineHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateAlpineProxyRepository(ctx context.Context, body sonatyperepoV382.AlpineProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateAlpineProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetAlpineProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetAlpineProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateAlpineProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateAlpineProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateAnsiblegalaxyGroupRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetAnsiblegalaxyGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyGroupApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetAnsiblegalaxyGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateAnsiblegalaxyGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateAnsiblegalaxyHostedRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetAnsiblegalaxyHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyHostedApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetAnsiblegalaxyHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateAnsiblegalaxyHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateAnsiblegalaxyProxyRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetAnsiblegalaxyProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyProxyApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetAnsiblegalaxyProxyRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateAnsiblegalaxyProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateAptHostedRepository(ctx context.Context, body sonatyperepoV382.AptHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateAptHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetAptHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AptHostedApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetAptHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateAptHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AptHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateAptHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateAptProxyRepository(ctx context.Context, body sonatyperepoV382.AptProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateAptProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetAptProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AptProxyApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetAptProxyRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateAptProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AptProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateAptProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateCargoGroupRepository(ctx context.Context, body sonatyperepoV382.CargoGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateCargoGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetCargoGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.CargoGroupApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetCargoGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateCargoGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateCargoGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateCargoHostedRepository(ctx context.Context, body sonatyperepoV382.CargoHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateCargoHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetCargoHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetCargoHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateCargoHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateCargoHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateCargoProxyRepository(ctx context.Context, body sonatyperepoV382.CargoProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateCargoProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetCargoProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.CargoProxyApiRepository, *FirewallMode, *http.Response, error) {
	result, httpResponse, err := s.client.RepositoryManagementAPI.GetCargoProxyRepository(ctx, repositoryName).Execute()
	return result, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateCargoProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateCargoProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateCocoapodsProxyRepository(ctx context.Context, body sonatyperepoV382.CocoapodsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateCocoapodsProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetCocoapodsProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetCocoapodsProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateCocoapodsProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CocoapodsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateCocoapodsProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateComposerProxyRepository(ctx context.Context, body sonatyperepoV382.ComposerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateComposerProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetComposerProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetComposerProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateComposerProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ComposerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateComposerProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateConanGroupRepository(ctx context.Context, body sonatyperepoV382.ConanGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateConanGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetConanGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetConanGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateConanGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateConanGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateConanHostedRepository(ctx context.Context, body sonatyperepoV382.ConanHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateConanHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetConanHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetConanHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateConanHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateConanHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateConanProxyRepository(ctx context.Context, body sonatyperepoV382.ConanProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateConanProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetConanProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.ConanProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetConanProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateConanProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateConanProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateCondaProxyRepository(ctx context.Context, body sonatyperepoV382.CondaProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateCondaProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetCondaProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetCondaProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateCondaProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CondaProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateCondaProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateDockerGroupRepository(ctx context.Context, body sonatyperepoV382.DockerGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateDockerGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetDockerGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerGroupApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetDockerGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateDockerGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateDockerGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateDockerHostedRepository(ctx context.Context, body sonatyperepoV382.DockerHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateDockerHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetDockerHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerHostedApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetDockerHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateDockerHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateDockerHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateDockerProxyRepository(ctx context.Context, body sonatyperepoV382.DockerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateDockerProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetDockerProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetDockerProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateDockerProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateDockerProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateGitlfsHostedRepository(ctx context.Context, body sonatyperepoV382.GitLfsHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateGitlfsHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetGitlfsHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetGitlfsHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateGitlfsHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GitLfsHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateGitlfsHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateGoGroupRepository(ctx context.Context, body sonatyperepoV382.GolangGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateGoGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetGoGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetGoGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateGoGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateGoGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateGoHostedRepository(ctx context.Context, body sonatyperepoV382.GolangHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateGoHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetGoHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetGoHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateGoHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateGoHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateGoProxyRepository(ctx context.Context, body sonatyperepoV382.GolangProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateGoProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetGoProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	result, httpResponse, err := s.client.RepositoryManagementAPI.GetGoProxyRepository(ctx, repositoryName).Execute()
	return result, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateGoProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateGoProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateHelmGroupRepository(ctx context.Context, body sonatyperepoV382.HelmGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateHelmGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetHelmGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetHelmGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateHelmGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateHelmGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateHelmHostedRepository(ctx context.Context, body sonatyperepoV382.HelmHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateHelmHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetHelmHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetHelmHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateHelmHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateHelmHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateHelmProxyRepository(ctx context.Context, body sonatyperepoV382.HelmProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateHelmProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetHelmProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetHelmProxyRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateHelmProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateHelmProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateHuggingfaceProxyRepository(ctx context.Context, body sonatyperepoV382.HuggingFaceProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateHuggingfaceProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetHuggingfaceProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	result, httpResponse, err := s.client.RepositoryManagementAPI.GetHuggingfaceProxyRepository(ctx, repositoryName).Execute()
	return result, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateHuggingfaceProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HuggingFaceProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateHuggingfaceProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateMavenGroupRepository(ctx context.Context, body sonatyperepoV382.MavenGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateMavenGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetMavenGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetMavenGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateMavenGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateMavenGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateMavenHostedRepository(ctx context.Context, body sonatyperepoV382.MavenHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateMavenHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetMavenHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.MavenHostedApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetMavenHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateMavenHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateMavenHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateMavenProxyRepository(ctx context.Context, body sonatyperepoV382.MavenProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateMavenProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetMavenProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.MavenProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetMavenProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateMavenProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateMavenProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateNpmGroupRepository(ctx context.Context, body sonatyperepoV382.NpmGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateNpmGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetNpmGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetNpmGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateNpmGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateNpmGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateNpmHostedRepository(ctx context.Context, body sonatyperepoV382.NpmHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateNpmHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetNpmHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetNpmHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateNpmHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateNpmHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateNpmProxyRepository(ctx context.Context, body sonatyperepoV382.NpmProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateNpmProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetNpmProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.NpmProxyApiRepository, *FirewallMode, *http.Response, error) {
	api, httpResp, err := s.client.RepositoryManagementAPI.GetNpmProxyRepository(ctx, repositoryName).Execute()
	return api, nil, httpResp, err
}

func (s *repositoryManagementServiceV382) UpdateNpmProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateNpmProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateNugetGroupRepository(ctx context.Context, body sonatyperepoV382.NugetGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateNugetGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetNugetGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetNugetGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateNugetGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateNugetGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateNugetHostedRepository(ctx context.Context, body sonatyperepoV382.NugetHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateNugetHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetNugetHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetNugetHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateNugetHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateNugetHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateNugetProxyRepository(ctx context.Context, body sonatyperepoV382.NugetProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateNugetProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetNugetProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.NugetProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetNugetProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateNugetProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateNugetProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateP2ProxyRepository(ctx context.Context, body sonatyperepoV382.P2ProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateP2ProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetP2ProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetP2ProxyRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateP2ProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.P2ProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateP2ProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreatePubGroupRepository(ctx context.Context, body sonatyperepoV382.PubGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreatePubGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetPubGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetPubGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdatePubGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdatePubGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreatePubHostedRepository(ctx context.Context, body sonatyperepoV382.PubHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreatePubHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetPubHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetPubHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdatePubHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdatePubHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreatePubProxyRepository(ctx context.Context, body sonatyperepoV382.PubProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreatePubProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetPubProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	result, httpResponse, err := s.client.RepositoryManagementAPI.GetPubProxyRepository(ctx, repositoryName).Execute()
	return result, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdatePubProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdatePubProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreatePypiGroupRepository(ctx context.Context, body sonatyperepoV382.PypiGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreatePypiGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetPypiGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetPypiGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdatePypiGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdatePypiGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreatePypiHostedRepository(ctx context.Context, body sonatyperepoV382.PypiHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreatePypiHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetPypiHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetPypiHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdatePypiHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdatePypiHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreatePypiProxyRepository(ctx context.Context, body sonatyperepoV382.PypiProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreatePypiProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetPypiProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.PyPiProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetPypiProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdatePypiProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdatePypiProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRGroupRepository(ctx context.Context, body sonatyperepoV382.RGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetRGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateRGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRHostedRepository(ctx context.Context, body sonatyperepoV382.RHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetRHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateRHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRProxyRepository(ctx context.Context, body sonatyperepoV382.RProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetRProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateRProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRawGroupRepository(ctx context.Context, body sonatyperepoV382.RawGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRawGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRawGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawGroupApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetRawGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateRawGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRawGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRawHostedRepository(ctx context.Context, body sonatyperepoV382.RawHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRawHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRawHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawHostedApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetRawHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateRawHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRawHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRawProxyRepository(ctx context.Context, body sonatyperepoV382.RawProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRawProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRawProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetRawProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateRawProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRawProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRubygemsGroupRepository(ctx context.Context, body sonatyperepoV382.RubyGemsGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRubygemsGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRubygemsGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetRubygemsGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateRubygemsGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRubygemsGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRubygemsHostedRepository(ctx context.Context, body sonatyperepoV382.RubyGemsHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRubygemsHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRubygemsHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetRubygemsHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateRubygemsHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRubygemsHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateRubygemsProxyRepository(ctx context.Context, body sonatyperepoV382.RubyGemsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateRubygemsProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetRubygemsProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetRubygemsProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateRubygemsProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateRubygemsProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateSwiftGroupRepository(ctx context.Context, body sonatyperepoV382.SwiftGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateSwiftGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetSwiftGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SwiftGroupApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetSwiftGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateSwiftGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.SwiftGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateSwiftGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateSwiftProxyRepository(ctx context.Context, body sonatyperepoV382.SwiftProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateSwiftProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetSwiftProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SwiftProxyApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetSwiftProxyRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateSwiftProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.SwiftProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateSwiftProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateTerraformGroupRepository(ctx context.Context, body sonatyperepoV382.TerraformGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateTerraformGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetTerraformGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformGroupApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetTerraformGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateTerraformGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateTerraformGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateTerraformHostedRepository(ctx context.Context, body sonatyperepoV382.TerraformHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateTerraformHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetTerraformHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformHostedRepositoryApiRequest, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetTerraformHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateTerraformHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateTerraformHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateTerraformProxyRepository(ctx context.Context, body sonatyperepoV382.TerraformProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateTerraformProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetTerraformProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformProxyApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetTerraformProxyRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateTerraformProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformProxyRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateTerraformProxyRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateYumGroupRepository(ctx context.Context, body sonatyperepoV382.YumGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateYumGroupRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetYumGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumGroupApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetYumGroupRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateYumGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumGroupRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateYumGroupRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateYumHostedRepository(ctx context.Context, body sonatyperepoV382.YumHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateYumHostedRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetYumHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumHostedApiRepository, *http.Response, error) {
	return s.client.RepositoryManagementAPI.GetYumHostedRepository(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV382) UpdateYumHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumHostedRepositoryApiRequest) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateYumHostedRepository(ctx, repositoryName).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) CreateYumProxyRepository(ctx context.Context, body sonatyperepoV382.YumProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.CreateYumProxyRepository(ctx).Body(body).Execute()
}

func (s *repositoryManagementServiceV382) GetYumProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.RepositoryManagementAPI.GetYumProxyRepository(ctx, repositoryName).Execute()
	return apiResponse, nil, httpResponse, err
}

func (s *repositoryManagementServiceV382) UpdateYumProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.UpdateYumProxyRepository(ctx, repositoryName).Body(body).Execute()
}

// repositoryManagementServiceV395 implements RepositoryManagementService against NXRM API
// client V395 (targets NXRM 3.94.0+). Requests/responses are bridged to/from V382 shapes via
// JSON, since the internal/provider/model package's mapping methods are written in terms of
// V382 types. Known field renames that survive the JSON bridge incorrectly (because the JSON
// key itself changed) are patched explicitly:
//   - RoutingRule (V382) -> RoutingRuleName (V395), on the write side of proxy formats.
//   - npm/pypi proxy "RemoveQuarantined" (PCCS) has no equivalent field on the V395 repository
//     struct at all (confirmed removed, not renamed) -- reads of this flag against NXRM 3.94+
//     always come back false, and writes of it are silently dropped. Restoring this behavior
//     requires integrating a separate Capability-API-based PCCS toggle, which is out of scope
//     for this struct-mapping bridge; tracked as a follow-up.
type repositoryManagementServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *repositoryManagementServiceV395) DeleteRepository(ctx context.Context, repositoryName string) (*http.Response, error) {
	return s.client.RepositoryManagementAPI.DeleteRepositories(ctx, repositoryName).Execute()
}

func (s *repositoryManagementServiceV395) ListRepositories(ctx context.Context) ([]sonatyperepoV382.RepositoryXO, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAllRepositories(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.RepositoryXO
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) CreateAlpineGroupRepository(ctx context.Context, body sonatyperepoV382.AlpineGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AlpineGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAlpineGroupRepository(ctx).AlpineGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetAlpineGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineGroupApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAlpineGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AlpineGroupApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateAlpineGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AlpineGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAlpineGroupRepository(ctx, repositoryName).AlpineGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateAlpineHostedRepository(ctx context.Context, body sonatyperepoV382.AlpineHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AlpineHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAlpineHostedRepository(ctx).AlpineHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetAlpineHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineHostedApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAlpineHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AlpineHostedApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateAlpineHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AlpineHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAlpineHostedRepository(ctx, repositoryName).AlpineHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateAlpineProxyRepository(ctx context.Context, body sonatyperepoV382.AlpineProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.AlpineProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateAlpineProxyRepository(ctx).AlpineProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetAlpineProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AlpineProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAlpineProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.AlpineProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateAlpineProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AlpineProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.AlpineProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateAlpineProxyRepository(ctx, repositoryName).AlpineProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateAnsiblegalaxyGroupRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AnsibleGalaxyGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyGroupRepository(ctx).AnsibleGalaxyGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetAnsiblegalaxyGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyGroupApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAnsiblegalaxyGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnsibleGalaxyGroupApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateAnsiblegalaxyGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AnsibleGalaxyGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyGroupRepository(ctx, repositoryName).AnsibleGalaxyGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateAnsiblegalaxyHostedRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AnsibleGalaxyHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyHostedRepository(ctx).AnsibleGalaxyHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetAnsiblegalaxyHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyHostedApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAnsiblegalaxyHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnsibleGalaxyHostedApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateAnsiblegalaxyHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AnsibleGalaxyHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyHostedRepository(ctx, repositoryName).AnsibleGalaxyHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateAnsiblegalaxyProxyRepository(ctx context.Context, body sonatyperepoV382.AnsibleGalaxyProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AnsibleGalaxyProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateAnsiblegalaxyProxyRepository(ctx).AnsibleGalaxyProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetAnsiblegalaxyProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AnsibleGalaxyProxyApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAnsiblegalaxyProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AnsibleGalaxyProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateAnsiblegalaxyProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AnsibleGalaxyProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AnsibleGalaxyProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateAnsiblegalaxyProxyRepository(ctx, repositoryName).AnsibleGalaxyProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateAptHostedRepository(ctx context.Context, body sonatyperepoV382.AptHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AptHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateAptHostedRepository(ctx).AptHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetAptHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AptHostedApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAptHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AptHostedApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateAptHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AptHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AptHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateAptHostedRepository(ctx, repositoryName).AptHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateAptProxyRepository(ctx context.Context, body sonatyperepoV382.AptProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AptProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateAptProxyRepository(ctx).AptProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetAptProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.AptProxyApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetAptProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AptProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateAptProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.AptProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.AptProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateAptProxyRepository(ctx, repositoryName).AptProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateCargoGroupRepository(ctx context.Context, body sonatyperepoV382.CargoGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.CargoGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateCargoGroupRepository(ctx).CargoGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetCargoGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.CargoGroupApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetCargoGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.CargoGroupApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateCargoGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.CargoGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateCargoGroupRepository(ctx, repositoryName).CargoGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateCargoHostedRepository(ctx context.Context, body sonatyperepoV382.CargoHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.CargoHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateCargoHostedRepository(ctx).CargoHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetCargoHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetCargoHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateCargoHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.CargoHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateCargoHostedRepository(ctx, repositoryName).CargoHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateCargoProxyRepository(ctx context.Context, body sonatyperepoV382.CargoProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.CargoProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateCargoProxyRepository(ctx).CargoProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetCargoProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.CargoProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetCargoProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil && apiV395.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.CargoProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateCargoProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CargoProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.CargoProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateCargoProxyRepository(ctx, repositoryName).CargoProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateCocoapodsProxyRepository(ctx context.Context, body sonatyperepoV382.CocoapodsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.CocoapodsProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateCocoapodsProxyRepository(ctx).CocoapodsProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetCocoapodsProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetCocoapodsProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateCocoapodsProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CocoapodsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.CocoapodsProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateCocoapodsProxyRepository(ctx, repositoryName).CocoapodsProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateComposerProxyRepository(ctx context.Context, body sonatyperepoV382.ComposerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.ComposerProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateComposerProxyRepository(ctx).ComposerProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetComposerProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetComposerProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateComposerProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ComposerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.ComposerProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateComposerProxyRepository(ctx, repositoryName).ComposerProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateConanGroupRepository(ctx context.Context, body sonatyperepoV382.ConanGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.ConanGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateConanGroupRepository(ctx).ConanGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetConanGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetConanGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupDeployRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateConanGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.ConanGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateConanGroupRepository(ctx, repositoryName).ConanGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateConanHostedRepository(ctx context.Context, body sonatyperepoV382.ConanHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.ConanHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateConanHostedRepository(ctx).ConanHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetConanHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetConanHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateConanHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.ConanHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateConanHostedRepository(ctx, repositoryName).ConanHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateConanProxyRepository(ctx context.Context, body sonatyperepoV382.ConanProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.ConanProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateConanProxyRepository(ctx).ConanProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetConanProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.ConanProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetConanProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.ConanProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateConanProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.ConanProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.ConanProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateConanProxyRepository(ctx, repositoryName).ConanProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateCondaProxyRepository(ctx context.Context, body sonatyperepoV382.CondaProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.CondaProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateCondaProxyRepository(ctx).CondaProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetCondaProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetCondaProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateCondaProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.CondaProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.CondaProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateCondaProxyRepository(ctx, repositoryName).CondaProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateDockerGroupRepository(ctx context.Context, body sonatyperepoV382.DockerGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.DockerGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateDockerGroupRepository(ctx).DockerGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetDockerGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerGroupApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetDockerGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.DockerGroupApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateDockerGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.DockerGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateDockerGroupRepository(ctx, repositoryName).DockerGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateDockerHostedRepository(ctx context.Context, body sonatyperepoV382.DockerHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.DockerHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateDockerHostedRepository(ctx).DockerHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetDockerHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerHostedApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetDockerHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.DockerHostedApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateDockerHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.DockerHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateDockerHostedRepository(ctx, repositoryName).DockerHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateDockerProxyRepository(ctx context.Context, body sonatyperepoV382.DockerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.DockerProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateDockerProxyRepository(ctx).DockerProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetDockerProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.DockerProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetDockerProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil && apiV395.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.DockerProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateDockerProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.DockerProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.DockerProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateDockerProxyRepository(ctx, repositoryName).DockerProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateGitlfsHostedRepository(ctx context.Context, body sonatyperepoV382.GitLfsHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.GitLfsHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateGitlfsHostedRepository(ctx).GitLfsHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetGitlfsHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetGitlfsHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateGitlfsHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GitLfsHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.GitLfsHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateGitlfsHostedRepository(ctx, repositoryName).GitLfsHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateGoGroupRepository(ctx context.Context, body sonatyperepoV382.GolangGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.GolangGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateGoGroupRepository(ctx).GolangGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetGoGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetGoGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateGoGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.GolangGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateGoGroupRepository(ctx, repositoryName).GolangGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateGoHostedRepository(ctx context.Context, body sonatyperepoV382.GolangHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.GolangHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateGoHostedRepository(ctx).GolangHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetGoHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetGoHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateGoHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.GolangHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateGoHostedRepository(ctx, repositoryName).GolangHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateGoProxyRepository(ctx context.Context, body sonatyperepoV382.GolangProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.GolangProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateGoProxyRepository(ctx).GolangProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetGoProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetGoProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil && apiV395.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateGoProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.GolangProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.GolangProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateGoProxyRepository(ctx, repositoryName).GolangProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateHelmGroupRepository(ctx context.Context, body sonatyperepoV382.HelmGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.HelmGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateHelmGroupRepository(ctx).HelmGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetHelmGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetHelmGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateHelmGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.HelmGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateHelmGroupRepository(ctx, repositoryName).HelmGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateHelmHostedRepository(ctx context.Context, body sonatyperepoV382.HelmHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.HelmHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateHelmHostedRepository(ctx).HelmHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetHelmHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetHelmHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateHelmHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.HelmHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateHelmHostedRepository(ctx, repositoryName).HelmHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateHelmProxyRepository(ctx context.Context, body sonatyperepoV382.HelmProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.HelmProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateHelmProxyRepository(ctx).HelmProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetHelmProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetHelmProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateHelmProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HelmProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.HelmProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateHelmProxyRepository(ctx, repositoryName).HelmProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateHuggingfaceProxyRepository(ctx context.Context, body sonatyperepoV382.HuggingFaceProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.HuggingFaceProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateHuggingfaceProxyRepository(ctx).HuggingFaceProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetHuggingfaceProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetHuggingfaceProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil && apiV395.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateHuggingfaceProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.HuggingFaceProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.HuggingFaceProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateHuggingfaceProxyRepository(ctx, repositoryName).HuggingFaceProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateMavenGroupRepository(ctx context.Context, body sonatyperepoV382.MavenGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.MavenGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateMavenGroupRepository(ctx).MavenGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetMavenGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetMavenGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateMavenGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.MavenGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateMavenGroupRepository(ctx, repositoryName).MavenGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateMavenHostedRepository(ctx context.Context, body sonatyperepoV382.MavenHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.MavenHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateMavenHostedRepository(ctx).MavenHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetMavenHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.MavenHostedApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetMavenHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.MavenHostedApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateMavenHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.MavenHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateMavenHostedRepository(ctx, repositoryName).MavenHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateMavenProxyRepository(ctx context.Context, body sonatyperepoV382.MavenProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.MavenProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateMavenProxyRepository(ctx).MavenProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetMavenProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.MavenProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetMavenProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.MavenProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateMavenProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.MavenProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.MavenProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateMavenProxyRepository(ctx, repositoryName).MavenProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateNpmGroupRepository(ctx context.Context, body sonatyperepoV382.NpmGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.NpmGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateNpmGroupRepository(ctx).NpmGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetNpmGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetNpmGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupDeployRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateNpmGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.NpmGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateNpmGroupRepository(ctx, repositoryName).NpmGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateNpmHostedRepository(ctx context.Context, body sonatyperepoV382.NpmHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.NpmHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateNpmHostedRepository(ctx).NpmHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetNpmHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetNpmHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateNpmHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.NpmHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateNpmHostedRepository(ctx, repositoryName).NpmHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateNpmProxyRepository(ctx context.Context, body sonatyperepoV382.NpmProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.NpmProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateNpmProxyRepository(ctx).NpmProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetNpmProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.NpmProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetNpmProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var result sonatyperepoV382.NpmProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	var mode *FirewallMode
	if apiV395.Firewall != nil && apiV395.Firewall.Mode != nil {
		m := FirewallMode(*apiV395.Firewall.Mode)
		mode = &m
	}
	return &result, mode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateNpmProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NpmProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.NpmProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateNpmProxyRepository(ctx, repositoryName).NpmProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateNugetGroupRepository(ctx context.Context, body sonatyperepoV382.NugetGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.NugetGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateNugetGroupRepository(ctx).NugetGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetNugetGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetNugetGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateNugetGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.NugetGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateNugetGroupRepository(ctx, repositoryName).NugetGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateNugetHostedRepository(ctx context.Context, body sonatyperepoV382.NugetHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.NugetHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateNugetHostedRepository(ctx).NugetHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetNugetHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetNugetHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateNugetHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.NugetHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateNugetHostedRepository(ctx, repositoryName).NugetHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateNugetProxyRepository(ctx context.Context, body sonatyperepoV382.NugetProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.NugetProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateNugetProxyRepository(ctx).NugetProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetNugetProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.NugetProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetNugetProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.NugetProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateNugetProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.NugetProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.NugetProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateNugetProxyRepository(ctx, repositoryName).NugetProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateP2ProxyRepository(ctx context.Context, body sonatyperepoV382.P2ProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.P2ProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateP2ProxyRepository(ctx).P2ProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetP2ProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetP2ProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateP2ProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.P2ProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.P2ProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateP2ProxyRepository(ctx, repositoryName).P2ProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreatePubGroupRepository(ctx context.Context, body sonatyperepoV382.PubGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.PubGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreatePubGroupRepository(ctx).PubGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetPubGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetPubGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdatePubGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.PubGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdatePubGroupRepository(ctx, repositoryName).PubGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreatePubHostedRepository(ctx context.Context, body sonatyperepoV382.PubHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.PubHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreatePubHostedRepository(ctx).PubHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetPubHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetPubHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdatePubHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.PubHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdatePubHostedRepository(ctx, repositoryName).PubHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreatePubProxyRepository(ctx context.Context, body sonatyperepoV382.PubProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.PubProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreatePubProxyRepository(ctx).PubProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetPubProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetPubProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil && apiV395.Firewall.Mode != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdatePubProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PubProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.PubProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdatePubProxyRepository(ctx, repositoryName).PubProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreatePypiGroupRepository(ctx context.Context, body sonatyperepoV382.PypiGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.PypiGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreatePypiGroupRepository(ctx).PypiGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetPypiGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupDeployRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetPypiGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupDeployRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdatePypiGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.PypiGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdatePypiGroupRepository(ctx, repositoryName).PypiGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreatePypiHostedRepository(ctx context.Context, body sonatyperepoV382.PypiHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.PypiHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreatePypiHostedRepository(ctx).PypiHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetPypiHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetPypiHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdatePypiHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.PypiHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdatePypiHostedRepository(ctx, repositoryName).PypiHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreatePypiProxyRepository(ctx context.Context, body sonatyperepoV382.PypiProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.PypiProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreatePypiProxyRepository(ctx).PypiProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetPypiProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.PyPiProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetPypiProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.PyPiProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdatePypiProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.PypiProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.PypiProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdatePypiProxyRepository(ctx, repositoryName).PypiProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRGroupRepository(ctx context.Context, body sonatyperepoV382.RGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRGroupRepository(ctx).RGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRGroupRepository(ctx, repositoryName).RGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRHostedRepository(ctx context.Context, body sonatyperepoV382.RHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRHostedRepository(ctx).RHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRHostedRepository(ctx, repositoryName).RHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRProxyRepository(ctx context.Context, body sonatyperepoV382.RProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.RProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateRProxyRepository(ctx).RProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.RProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateRProxyRepository(ctx, repositoryName).RProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRawGroupRepository(ctx context.Context, body sonatyperepoV382.RawGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RawGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRawGroupRepository(ctx).RawGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRawGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawGroupApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRawGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RawGroupApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRawGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RawGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRawGroupRepository(ctx, repositoryName).RawGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRawHostedRepository(ctx context.Context, body sonatyperepoV382.RawHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RawHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRawHostedRepository(ctx).RawHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRawHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawHostedApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRawHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.RawHostedApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRawHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RawHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRawHostedRepository(ctx, repositoryName).RawHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRawProxyRepository(ctx context.Context, body sonatyperepoV382.RawProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.RawProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateRawProxyRepository(ctx).RawProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRawProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.RawProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRawProxyRepository(ctx, repositoryName).Execute()
	var result sonatyperepoV382.RawProxyApiRepository
	if err := bridgeFromResponse(apiV395, httpResponse, err, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	// v395.95.0's RawProxyApiRepository response type has no Firewall field (unlike its
	// request-type counterpart, and unlike most other proxy formats' response types), so the
	// mode cannot be read back for Raw via this endpoint.
	return &result, nil, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRawProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RawProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.RawProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateRawProxyRepository(ctx, repositoryName).RawProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRubygemsGroupRepository(ctx context.Context, body sonatyperepoV382.RubyGemsGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RubyGemsGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRubygemsGroupRepository(ctx).RubyGemsGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRubygemsGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiGroupRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRubygemsGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiGroupRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRubygemsGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RubyGemsGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRubygemsGroupRepository(ctx, repositoryName).RubyGemsGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRubygemsHostedRepository(ctx context.Context, body sonatyperepoV382.RubyGemsHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RubyGemsHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateRubygemsHostedRepository(ctx).RubyGemsHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRubygemsHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiHostedRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRubygemsHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SimpleApiHostedRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRubygemsHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.RubyGemsHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateRubygemsHostedRepository(ctx, repositoryName).RubyGemsHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateRubygemsProxyRepository(ctx context.Context, body sonatyperepoV382.RubyGemsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.RubyGemsProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateRubygemsProxyRepository(ctx).RubyGemsProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetRubygemsProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SimpleApiProxyRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetRubygemsProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.SimpleApiProxyRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateRubygemsProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.RubyGemsProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.RubyGemsProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateRubygemsProxyRepository(ctx, repositoryName).RubyGemsProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateSwiftGroupRepository(ctx context.Context, body sonatyperepoV382.SwiftGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.SwiftGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateSwiftGroupRepository(ctx).SwiftGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetSwiftGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SwiftGroupApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetSwiftGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SwiftGroupApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateSwiftGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.SwiftGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.SwiftGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateSwiftGroupRepository(ctx, repositoryName).SwiftGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateSwiftProxyRepository(ctx context.Context, body sonatyperepoV382.SwiftProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.SwiftProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateSwiftProxyRepository(ctx).SwiftProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetSwiftProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.SwiftProxyApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetSwiftProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.SwiftProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateSwiftProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.SwiftProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.SwiftProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateSwiftProxyRepository(ctx, repositoryName).SwiftProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateTerraformGroupRepository(ctx context.Context, body sonatyperepoV382.TerraformGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.TerraformGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateTerraformGroupRepository(ctx).TerraformGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetTerraformGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformGroupApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetTerraformGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.TerraformGroupApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateTerraformGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.TerraformGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateTerraformGroupRepository(ctx, repositoryName).TerraformGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateTerraformHostedRepository(ctx context.Context, body sonatyperepoV382.TerraformHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.TerraformHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateTerraformHostedRepository(ctx).TerraformHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetTerraformHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformHostedRepositoryApiRequest, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetTerraformHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.TerraformHostedRepositoryApiRequest
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateTerraformHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.TerraformHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateTerraformHostedRepository(ctx, repositoryName).TerraformHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateTerraformProxyRepository(ctx context.Context, body sonatyperepoV382.TerraformProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.TerraformProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.CreateTerraformProxyRepository(ctx).TerraformProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetTerraformProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.TerraformProxyApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetTerraformProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.TerraformProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateTerraformProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.TerraformProxyRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.TerraformProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	return s.client.RepositoryManagementAPI.UpdateTerraformProxyRepository(ctx, repositoryName).TerraformProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateYumGroupRepository(ctx context.Context, body sonatyperepoV382.YumGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.YumGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateYumGroupRepository(ctx).YumGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetYumGroupRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumGroupApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetYumGroupRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.YumGroupApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateYumGroupRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumGroupRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.YumGroupRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateYumGroupRepository(ctx, repositoryName).YumGroupRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateYumHostedRepository(ctx context.Context, body sonatyperepoV382.YumHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.YumHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.CreateYumHostedRepository(ctx).YumHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetYumHostedRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumHostedApiRepository, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetYumHostedRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.YumHostedApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateYumHostedRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumHostedRepositoryApiRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.YumHostedRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.RepositoryManagementAPI.UpdateYumHostedRepository(ctx, repositoryName).YumHostedRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) CreateYumProxyRepository(ctx context.Context, body sonatyperepoV382.YumProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.YumProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.CreateYumProxyRepository(ctx).YumProxyRepositoryApiRequest(v395Body).Execute()
}

func (s *repositoryManagementServiceV395) GetYumProxyRepository(ctx context.Context, repositoryName string) (*sonatyperepoV382.YumProxyApiRepository, *FirewallMode, *http.Response, error) {
	apiV395, httpResponse, err := s.client.RepositoryManagementAPI.GetYumProxyRepository(ctx, repositoryName).Execute()
	if err != nil {
		return nil, nil, httpResponse, err
	}
	var firewallMode *FirewallMode
	if apiV395.Firewall != nil {
		firewallMode = (*FirewallMode)(apiV395.Firewall.Mode)
	}
	var result sonatyperepoV382.YumProxyApiRepository
	if err := jsonBridge(apiV395, &result); err != nil {
		return nil, nil, httpResponse, err
	}
	return &result, firewallMode, httpResponse, nil
}

func (s *repositoryManagementServiceV395) UpdateYumProxyRepository(ctx context.Context, repositoryName string, body sonatyperepoV382.YumProxyRepositoryApiRequest, firewallMode *FirewallMode) (*http.Response, error) {
	var v395Body sonatyperepoV395.YumProxyRepositoryApiRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	v395Body.RoutingRuleName = body.RoutingRule
	if firewallMode != nil {
		v395Body.Firewall = &sonatyperepoV395.FirewallAttributes{Mode: (*string)(firewallMode)}
	}
	return s.client.RepositoryManagementAPI.UpdateYumProxyRepository(ctx, repositoryName).YumProxyRepositoryApiRequest(v395Body).Execute()
}
