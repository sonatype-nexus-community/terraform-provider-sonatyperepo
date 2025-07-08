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
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// RepositoryProxyModel
// --------------------------------------------------------
type RepositoryProxyModel struct {
	BasicRepositoryModel
	Proxy         repositoryProxyModel         `tfsdk:"proxy"`
	NegativeCache repositoryNegativeCacheModel `tfsdk:"negative_cache"`
	HttpClient    repositoryHttpClientModel    `tfsdk:"http_client"`
	RoutingRule   types.String                 `tfsdk:"routing_rule"`
	Replication   *RepositoryReplicationModel  `tfsdk:"replication"`
}

// repositoryProxyModel
// --------------------------------------------------------
type repositoryProxyModel struct {
	RemoteUrl      types.String `tfsdk:"remote_url"`
	ContentMaxAge  types.Int64  `tfsdk:"content_max_age"`
	MetadataMaxAge types.Int64  `tfsdk:"metadata_max_age"`
}

func (m *repositoryProxyModel) MapFromApi(api *sonatyperepo.ProxyAttributes) {
	m.ContentMaxAge = types.Int64Value(int64(api.ContentMaxAge))
	m.MetadataMaxAge = types.Int64Value(int64(api.MetadataMaxAge))
	m.RemoteUrl = types.StringPointerValue(api.RemoteUrl)
}

func (m *repositoryProxyModel) MapToApi(api *sonatyperepo.ProxyAttributes) {
	api.ContentMaxAge = int32(m.ContentMaxAge.ValueInt64())
	api.MetadataMaxAge = int32(m.MetadataMaxAge.ValueInt64())
	api.RemoteUrl = m.RemoteUrl.ValueStringPointer()
}

func (m *repositoryHttpClientModel) MapFromApiHttpClientAttributes(api *sonatyperepo.HttpClientAttributes) {
	m.AutoBlock = types.BoolValue(api.AutoBlock)
	m.Blocked = types.BoolValue(api.Blocked)

	if api.Connection != nil {
		m.Connection.MapFromApi(api.Connection)
	}
	if api.Authentication != nil {
		m.Authentication.MapFromApiHttpClientConnectionAuthenticationAttributes(api.Authentication)
	}
}

func (m *repositoryHttpClientModel) MapToApiHttpClientAttributes(api *sonatyperepo.HttpClientAttributes) {
	api.AutoBlock = m.AutoBlock.ValueBool()
	api.Blocked = m.Blocked.ValueBool()

	if m.Connection != nil {
		api.Connection = &sonatyperepo.HttpClientConnectionAttributes{}
		m.Connection.MapToApi(api.Connection)
	}

	if m.Authentication != nil {
		api.Authentication = &sonatyperepo.HttpClientConnectionAuthenticationAttributes{}
		m.Authentication.MapToApiHttpClientConnectionAuthenticationAttributes(api.Authentication)
	}
}

func (m *repositoryHttpClientModel) MapToApiHttpClientAttributesWithPreemptiveAuth(api *sonatyperepo.HttpClientAttributesWithPreemptiveAuth) {
	api.AutoBlock = m.AutoBlock.ValueBool()
	api.Blocked = m.Blocked.ValueBool()

	if m.Connection != nil {
		api.Connection = &sonatyperepo.HttpClientConnectionAttributes{}
		m.Connection.MapToApi(api.Connection)
	}

	if m.Authentication != nil {
		api.Authentication = &sonatyperepo.HttpClientConnectionAuthenticationAttributesWithPreemptive{}
		m.Authentication.MapToApiHttpClientConnectionAuthenticationAttributesWithPreemptive(api.Authentication)
	}
}

type repositoryNegativeCacheModel struct {
	Enabled    types.Bool  `tfsdk:"enabled"`
	TimeToLive types.Int64 `tfsdk:"time_to_live"`
}

func (m *repositoryNegativeCacheModel) MapFromApi(api *sonatyperepo.NegativeCacheAttributes) {
	m.Enabled = types.BoolValue(api.Enabled)
	m.TimeToLive = types.Int64Value(int64(api.TimeToLive))
}

func (m *repositoryNegativeCacheModel) MapToApi(api *sonatyperepo.NegativeCacheAttributes) {
	api.Enabled = m.Enabled.ValueBool()
	api.TimeToLive = int32(m.TimeToLive.ValueInt64())
}

type repositoryHttpClientModel struct {
	Blocked        types.Bool                               `tfsdk:"blocked"`
	AutoBlock      types.Bool                               `tfsdk:"auto_block"`
	Connection     *RepositoryHttpClientConnectionModel     `tfsdk:"connection"`
	Authentication *RepositoryHttpClientAuthenticationModel `tfsdk:"authentication"`
}

// RepositoryHttpClientConnectionModel
// --------------------------------------------------------
type RepositoryHttpClientConnectionModel struct {
	Retries                 types.Int64  `tfsdk:"retries"`
	UserAgentSuffix         types.String `tfsdk:"user_agent_suffix"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
	EnableCircularRedirects types.Bool   `tfsdk:"enable_circular_redirects"`
	EnableCookies           types.Bool   `tfsdk:"enable_cookies"`
	UseTrustStore           types.Bool   `tfsdk:"use_trust_store"`
}

func (m *RepositoryHttpClientConnectionModel) MapFromApi(api *sonatyperepo.HttpClientConnectionAttributes) {
	m.EnableCircularRedirects = types.BoolPointerValue(api.EnableCircularRedirects)
	m.EnableCookies = types.BoolPointerValue(api.EnableCookies)
	m.UseTrustStore = types.BoolPointerValue(api.UseTrustStore)
	m.UserAgentSuffix = types.StringPointerValue(api.UserAgentSuffix)

	if api.Retries != nil {
		m.Retries = types.Int64Value(int64(*api.Retries))
	}
	if api.Timeout != nil {
		m.Timeout = types.Int64Value(int64(*api.Timeout))
	}
}

func (m *RepositoryHttpClientConnectionModel) MapToApi(api *sonatyperepo.HttpClientConnectionAttributes) {
	api.EnableCircularRedirects = m.EnableCircularRedirects.ValueBoolPointer()
	api.EnableCookies = m.EnableCookies.ValueBoolPointer()
	api.UseTrustStore = m.UseTrustStore.ValueBoolPointer()
	api.UserAgentSuffix = m.UserAgentSuffix.ValueStringPointer()

	if !m.Retries.IsNull() {
		retries := int32(m.Retries.ValueInt64())
		api.Retries = &retries
	}

	if !m.Timeout.IsNull() {
		timeout := int32(m.Timeout.ValueInt64())
		api.Timeout = &timeout
	}
}

// RepositoryHttpClientAuthenticationModel
// --------------------------------------------------------
type RepositoryHttpClientAuthenticationModel struct {
	Type        types.String `tfsdk:"type"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	NtlmHost    types.String `tfsdk:"ntlm_host"`
	NtlmDomain  types.String `tfsdk:"ntlm_domain"`
	Preemptive  types.Bool   `tfsdk:"preemptive"`
	BearerToken types.String `tfsdk:"bearer_token"`
}

func (m *RepositoryHttpClientAuthenticationModel) MapFromApiHttpClientConnectionAuthenticationAttributes(api *sonatyperepo.HttpClientConnectionAuthenticationAttributes) {
	m.Type = types.StringPointerValue(api.Type)
	if api.Preemptive != nil {
		m.Preemptive = types.BoolPointerValue(api.Preemptive)
	}

	if *api.Type == common.HTTP_AUTH_TYPE_BEARER_TOKEN {
		m.BearerToken = types.StringPointerValue(api.BearerToken)
	} else if api.Type != nil {
		m.Username = types.StringPointerValue(api.Username)
		// m.Password = types.StringPointerValue(api.Password)

		if *api.Type == common.HTTP_AUTH_TYPE_NTLM {
			m.NtlmDomain = types.StringPointerValue(api.NtlmDomain)
			m.NtlmHost = types.StringPointerValue(api.NtlmHost)
		}
	}
}

func (m *RepositoryHttpClientAuthenticationModel) MapFromApiHttpClientConnectionAuthenticationAttributesWithPreemptive(api *sonatyperepo.HttpClientConnectionAuthenticationAttributesWithPreemptive) {
	m.Type = types.StringPointerValue(api.Type)
	m.Preemptive = types.BoolPointerValue(api.Preemptive)

	if api.Type != nil {
		if *api.Type == common.HTTP_AUTH_TYPE_BEARER_TOKEN {
			m.BearerToken = types.StringPointerValue(api.BearerToken)
		} else if api.Type != nil {
			m.Username = types.StringPointerValue(api.Username)
			m.Password = types.StringPointerValue(api.Password)

			if *api.Type == common.HTTP_AUTH_TYPE_NTLM {
				m.NtlmDomain = types.StringPointerValue(api.NtlmDomain)
				m.NtlmHost = types.StringPointerValue(api.NtlmHost)
			}
		}
	}
}

func (m *RepositoryHttpClientAuthenticationModel) MapToApiHttpClientConnectionAuthenticationAttributesWithPreemptive(api *sonatyperepo.HttpClientConnectionAuthenticationAttributesWithPreemptive) {
	api.Type = m.Type.ValueStringPointer()
	api.Preemptive = m.Preemptive.ValueBoolPointer()

	if m.Type.ValueString() == common.HTTP_AUTH_TYPE_BEARER_TOKEN {
		api.BearerToken = m.BearerToken.ValueStringPointer()
	} else if !m.Type.IsNull() {
		api.Username = m.Username.ValueStringPointer()
		api.Password = m.Password.ValueStringPointer()

		if m.Type.ValueString() == common.HTTP_AUTH_TYPE_NTLM {
			api.NtlmDomain = m.NtlmDomain.ValueStringPointer()
			api.NtlmHost = m.NtlmHost.ValueStringPointer()
		}
	}
}

func (m *RepositoryHttpClientAuthenticationModel) MapToApiHttpClientConnectionAuthenticationAttributes(api *sonatyperepo.HttpClientConnectionAuthenticationAttributes) {
	api.Type = m.Type.ValueStringPointer()
	api.Preemptive = m.Preemptive.ValueBoolPointer()

	if m.Type.ValueString() == common.HTTP_AUTH_TYPE_BEARER_TOKEN {
		api.BearerToken = m.BearerToken.ValueStringPointer()
	} else if !m.Type.IsNull() {
		api.Username = m.Username.ValueStringPointer()
		api.Password = m.Password.ValueStringPointer()

		if m.Type.ValueString() == common.HTTP_AUTH_TYPE_NTLM {
			api.NtlmDomain = m.NtlmDomain.ValueStringPointer()
			api.NtlmHost = m.NtlmHost.ValueStringPointer()
		}
	}
}

// RepositoryReplicationModel
// --------------------------------------------------------
type RepositoryReplicationModel struct {
	PreemptivePullEnabled types.Bool   `tfsdk:"preemptive_pull_enabled"`
	AssetPathRegex        types.String `tfsdk:"asset_path_regex"`
}

func (m *RepositoryReplicationModel) MapFromApi(api *sonatyperepo.ReplicationAttributes) {
	m.PreemptivePullEnabled = types.BoolValue(api.PreemptivePullEnabled)
	m.AssetPathRegex = types.StringPointerValue(api.AssetPathRegex)
}

func (m *RepositoryReplicationModel) MapToApi(api *sonatyperepo.ReplicationAttributes) {
	api.PreemptivePullEnabled = m.PreemptivePullEnabled.ValueBool()
	if m.AssetPathRegex.String() != "" {
		api.AssetPathRegex = m.AssetPathRegex.ValueStringPointer()
	}
}
