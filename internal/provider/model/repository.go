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

package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepositoryModel struct {
	Name   types.String `tfsdk:"name"`
	Format types.String `tfsdk:"format"`
	Type   types.String `tfsdk:"type"`
	Url    types.String `tfsdk:"url"`
}

type RepositoriesModel struct {
	Repositories []RepositoryModel `tfsdk:"repositories"`
}

type RepositoryMavenGroupModel struct {
	Name        types.String               `tfsdk:"name"`
	Format      types.String               `tfsdk:"format"`
	Type        types.String               `tfsdk:"type"`
	Url         types.String               `tfsdk:"url"`
	Online      types.Bool                 `tfsdk:"online"`
	Storage     repositoryStorageModeGroup `tfsdk:"storage"`
	Group       RepositoryGroupModel       `tfsdk:"group"`
	LastUpdated types.String               `tfsdk:"last_updated"`
}

type RepositoryMavenHostedModel struct {
	Name        types.String                   `tfsdk:"name"`
	Format      types.String                   `tfsdk:"format"`
	Type        types.String                   `tfsdk:"type"`
	Url         types.String                   `tfsdk:"url"`
	Online      types.Bool                     `tfsdk:"online"`
	Storage     repositoryStorageModelNonGroup `tfsdk:"storage"`
	Cleanup     *RepositoryCleanupModel        `tfsdk:"cleanup"`
	Maven       repositoryMavenSpecificModel   `tfsdk:"maven"`
	Component   *RepositoryComponentModel      `tfsdk:"component"`
	LastUpdated types.String                   `tfsdk:"last_updated"`
}

type RepositoryMavenProxyModel struct {
	Name          types.String                   `tfsdk:"name"`
	Format        types.String                   `tfsdk:"format"`
	Type          types.String                   `tfsdk:"type"`
	Url           types.String                   `tfsdk:"url"`
	Online        types.Bool                     `tfsdk:"online"`
	Storage       repositoryStorageModelNonGroup `tfsdk:"storage"`
	Cleanup       *RepositoryCleanupModel        `tfsdk:"cleanup"`
	Proxy         repositoryProxyModel           `tfsdk:"proxy"`
	NegativeCache repositoryNegativeCacheModel   `tfsdk:"negative_cache"`
	HttpClient    repositoryHttpClientModel      `tfsdk:"http_client"`
	RoutingRule   types.String                   `tfsdk:"routing_rule"`
	Replication   *RepositoryReplicationModel    `tfsdk:"replication"`
	Maven         repositoryMavenSpecificModel   `tfsdk:"maven"`
	LastUpdated   types.String                   `tfsdk:"last_updated"`
}

type repositoryStorageModelNonGroup struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
	WritePolicy                 types.String `tfsdk:"write_policy"`
}

type repositoryStorageModeGroup struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
}

type RepositoryCleanupModel struct {
	PolicyNames []types.String `tfsdk:"policy_names"`
}

type RepositoryComponentModel struct {
	ProprietaryComponents types.Bool `tfsdk:"proprietary_components"`
}

type repositoryMavenSpecificModel struct {
	VersionPolicy      types.String `tfsdk:"version_policy"`
	LayoutPolicy       types.String `tfsdk:"layout_policy"`
	ContentDisposition types.String `tfsdk:"content_disposition"`
}

type repositoryProxyModel struct {
	RemoteUrl      types.String `tfsdk:"remote_url"`
	ContentMaxAge  types.Int64  `tfsdk:"content_max_age"`
	MetadataMaxAge types.Int64  `tfsdk:"metadata_max_age"`
}

type repositoryNegativeCacheModel struct {
	Enabled    types.Bool  `tfsdk:"enabled"`
	TimeToLive types.Int64 `tfsdk:"time_to_live"`
}

type repositoryHttpClientModel struct {
	Blocked        types.Bool                               `tfsdk:"blocked"`
	AutoBlock      types.Bool                               `tfsdk:"auto_block"`
	Connection     *RepositoryHttpClientConnectionModel     `tfsdk:"connection"`
	Authentication *RepositoryHttpClientAuthenticationModel `tfsdk:"authentication"`
}

type RepositoryHttpClientConnectionModel struct {
	Retries                 types.Int64  `tfsdk:"retries"`
	UserAgentSuffix         types.String `tfsdk:"user_agent_suffix"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
	EnableCircularRedirects types.Bool   `tfsdk:"enable_circular_redirects"`
	EnableCookies           types.Bool   `tfsdk:"enable_cookies"`
	UseTrustStore           types.Bool   `tfsdk:"use_trust_store"`
}

type RepositoryHttpClientAuthenticationModel struct {
	Type       types.String `tfsdk:"type"`
	Username   types.String `tfsdk:"username"`
	Password   types.String `tfsdk:"password"`
	NtlmHost   types.String `tfsdk:"ntlm_host"`
	NtlmDomain types.String `tfsdk:"ntlm_domain"`
	Preemptive types.Bool   `tfsdk:"preemptive"`
}

type RepositoryReplicationModel struct {
	PreemptivePullEnabled types.Bool   `tfsdk:"preemptive_pull_enabled"`
	AssetPathRegex        types.String `tfsdk:"asset_path_regex"`
}

type RepositoryGroupModel struct {
	MemberNames []types.String `tfsdk:"member_names"`
}
