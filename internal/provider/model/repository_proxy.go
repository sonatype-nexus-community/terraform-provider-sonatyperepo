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
	"strconv"
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// RepositoryProxyModel
// --------------------------------------------------------
type RepositoryProxyModel struct {
	BasicRepositoryModel
	Storage       repositoryStorageModel       `tfsdk:"storage"`
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
	if api == nil {
		return
	}

	m.AutoBlock = types.BoolPointerValue(api.AutoBlock)
	m.Blocked = types.BoolPointerValue(api.Blocked)

	// Always populate into a brand new Connection/Authentication rather than mutating
	// m.Connection/m.Authentication in place: callers derive this model from the Plan
	// with a shallow copy (e.g. `stateModel := plan`), which copies these pointers
	// without cloning what they point to. Mutating through the existing pointer would
	// therefore also silently overwrite the original Plan's values (which
	// MapMissingApiFieldsFromPlan later needs to restore fields the API doesn't return,
	// such as password/preemptive).
	connection := &RepositoryHttpClientConnectionModel{}
	// Pass api.Connection which might be nil - the MapFromApi method will handle it
	connection.MapFromApi(api.Connection)
	m.Connection = connection

	if api.Authentication != nil {
		authentication := &RepositoryHttpClientAuthenticationModel{}
		// NXRM never returns password/bearerToken, and several formats' endpoints
		// (e.g. npm) omit preemptive from the response entirely too, so carry forward
		// whatever this model already held for them (e.g. from persisted state on a
		// plain Read, where there is no Plan to later restore fields from) before
		// overwriting the rest from the API response. MapFromApiHttpClientConnection-
		// AuthenticationAttributes only overwrites Preemptive when the API actually
		// returned one, so the carried-forward value survives formats that omit it.
		if m.Authentication != nil {
			authentication.Password = m.Authentication.Password
			authentication.BearerToken = m.Authentication.BearerToken
			authentication.Preemptive = m.Authentication.Preemptive
		}
		authentication.MapFromApiHttpClientConnectionAuthenticationAttributes(api.Authentication)
		// Default Preemptive to false when nothing was carried forward above and NXRM
		// didn't return it either (either because this format's API has no such field at
		// all - e.g. Helm, npm - or because it was never explicitly set on a format that
		// does support it). This only matters for terraform import: Create/Update always
		// have a Plan to carry the value forward from via MapMissingApiFieldsFromPlan
		// (called after this function returns), but ImportState calls this with a nil
		// prior model and no Plan at all, so without this default Preemptive would be left
		// null in freshly-imported state - and the schema's `false` default would then
		// reappear as a permanent phantom diff on the very next plan (see GH-493).
		if authentication.Preemptive.IsNull() {
			authentication.Preemptive = types.BoolValue(false)
		}
		m.Authentication = authentication
	} else {
		m.Authentication = nil
	}
}

func (m *repositoryHttpClientModel) MapToApiHttpClientAttributes(api *sonatyperepo.HttpClientAttributes) {
	api.AutoBlock = m.AutoBlock.ValueBoolPointer()
	api.Blocked = m.Blocked.ValueBoolPointer()

	if m.Connection != nil {
		api.Connection = &sonatyperepo.HttpClientConnectionAttributes{}
		m.Connection.MapToApi(api.Connection)
	}

	if m.Authentication.hasRequiredSecret() {
		api.Authentication = &sonatyperepo.HttpClientConnectionAuthenticationAttributes{}
		m.Authentication.MapToApiHttpClientConnectionAuthenticationAttributes(api.Authentication)
	}
}

func (m *repositoryHttpClientModel) MapToApiHttpClientAttributesWithPreemptiveAuth(api *sonatyperepo.HttpClientAttributesWithPreemptiveAuth) {
	api.AutoBlock = m.AutoBlock.ValueBoolPointer()
	api.Blocked = m.Blocked.ValueBoolPointer()

	if m.Connection != nil {
		api.Connection = &sonatyperepo.HttpClientConnectionAttributes{}
		m.Connection.MapToApi(api.Connection)
	}

	if m.Authentication.hasRequiredSecret() {
		api.Authentication = &sonatyperepo.HttpClientConnectionAuthenticationAttributesWithPreemptive{}
		m.Authentication.MapToApiHttpClientConnectionAuthenticationAttributesWithPreemptive(api.Authentication)
	}
}

func (m *repositoryHttpClientModel) MapMissingApiFieldsFromPlan(planModel repositoryHttpClientModel) {
	if planModel.Authentication == nil {
		return
	}

	if !planModel.Authentication.hasRequiredSecret() {
		// MapToApiHttpClientAttributes/MapToApiHttpClientAttributesWithPreemptiveAuth omitted
		// Authentication from the outbound request entirely, because the plan's copy was
		// missing the secret NXRM requires - so NXRM was never asked to change anything (see
		// GH-491). `authentication` isn't a Computed schema attribute, so Terraform requires
		// state to match the plan exactly; a fresh Read() would otherwise report whatever NXRM
		// now has (typically nil - omitting the field on NXRM's full-replace PUT clears it),
		// which would diverge from the plan and fail Terraform's post-apply consistency check.
		// Mirroring the plan's value verbatim - a value copy, not the same pointer, to avoid
		// aliasing the Plan's own struct (see GH-489) - is therefore the only state that's both
		// truthful (nothing changed) and consistent with what Terraform already committed to.
		authentication := *planModel.Authentication
		m.Authentication = &authentication
		return
	}

	if m.Authentication == nil {
		m.Authentication = &RepositoryHttpClientAuthenticationModel{}
	}
	m.Authentication.MapMissingApiFieldsFromPlan(planModel.Authentication)
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
	// Check if api is nil to prevent nil pointer dereference
	if api == nil {
		// Set default values
		m.EnableCircularRedirects = types.BoolValue(false)
		m.EnableCookies = types.BoolValue(false)
		m.UseTrustStore = types.BoolValue(false)
		m.UserAgentSuffix = types.StringNull()
		m.Retries = types.Int64Value(common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES)
		m.Timeout = types.Int64Value(common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT)
		return
	}

	m.EnableCircularRedirects = types.BoolPointerValue(api.EnableCircularRedirects)
	m.EnableCookies = types.BoolPointerValue(api.EnableCookies)
	m.UseTrustStore = types.BoolPointerValue(api.UseTrustStore)
	m.UserAgentSuffix = types.StringPointerValue(api.UserAgentSuffix)

	if api.Retries != nil {
		m.Retries = types.Int64Value(int64(*api.Retries))
	} else {
		m.Retries = types.Int64Value(common.DEFAULT_HTTP_CLIENT_CONNECTION_RETRIES)
	}

	if api.Timeout != nil {
		m.Timeout = types.Int64Value(int64(*api.Timeout))
	} else {
		m.Timeout = types.Int64Value(common.DEFAULT_HTTP_CLIENT_CONNECTION_TIMEOUT)
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

// hasRequiredSecret reports whether m carries the secret NXRM requires for its configured
// Type - Password for "username"/"ntlm", BearerToken for "bearerToken" - so callers building an
// outbound API request can tell a genuinely usable authentication block apart from one that's
// missing its secret and would be rejected server-side.
//
// This matters because NXRM never returns password/bearerToken from a GET (they're write-only in
// its OpenAPI schema - see MapFromApiHttpClientConnectionAuthenticationAttributes), and repository
// updates are a full PUT replace with no partial-update/PATCH endpoint at all: there is no way to
// tell NXRM "leave authentication as-is". So whenever this model's Type/Username were populated
// from a prior Read (state) but its secret was never restored from a Plan - e.g. upstream
// authentication configured directly against NXRM (UI/API) rather than through this provider,
// combined with `lifecycle { ignore_changes = [http_client.authentication] }` freezing that
// state-derived value into the plan sent to Update - sending the object through anyway produces a
// 400 from NXRM's `UsernameAuthenticationConfiguration`/`NtlmAuthenticationConfiguration`/
// `BearerTokenAuthenticationConfiguration` validators ("password must not be null" /
// "bearerToken must not be null"), surfaced by the provider as a confusing
// "Repository did not exist to update" error (see GH-491). Omitting the whole Authentication
// object in that case - rather than sending a half-populated one - degrades to the same "no
// authentication configured" request NXRM already accepts when the attribute is unset entirely.
func (m *RepositoryHttpClientAuthenticationModel) hasRequiredSecret() bool {
	if m == nil {
		return false
	}
	switch m.Type.ValueString() {
	case common.HTTP_AUTH_TYPE_BEARER_TOKEN:
		return !m.BearerToken.IsNull() && !m.BearerToken.IsUnknown()
	case common.HTTP_AUTH_TYPE_USERNAME, common.HTTP_AUTH_TYPE_NTLM:
		return !m.Password.IsNull() && !m.Password.IsUnknown()
	default:
		// Unknown/empty Type: nothing for NXRM to validate a secret against, so there's
		// no missing-secret condition to guard here - let the existing schema validation
		// (type must be one of the known enum values) catch anything actually invalid.
		return true
	}
}

func (m *RepositoryHttpClientAuthenticationModel) MapFromApiHttpClientConnectionAuthenticationAttributes(api *sonatyperepo.HttpClientConnectionAuthenticationAttributes) {
	m.Type = types.StringPointerValue(api.Type)
	if api.Preemptive != nil {
		m.Preemptive = types.BoolPointerValue(api.Preemptive)
	}

	if api.Type != nil {
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

func (m *RepositoryHttpClientAuthenticationModel) MapMissingApiFieldsFromPlan(planModel *RepositoryHttpClientAuthenticationModel) {
	// Most NXRM repository format endpoints (e.g. npm, maven, r, raw, ...) never return
	// password/bearerToken, and several also omit preemptive from their response
	// entirely (verified for npm; the Terraform format is the odd one out and always
	// echoes back an explicit `preemptive: false` when unset). All three must therefore
	// always be restored from the Plan rather than trusted from the API response.
	m.Password = planModel.Password
	m.BearerToken = planModel.BearerToken
	m.Preemptive = planModel.Preemptive
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

// FirewallAuditAndQuarantineModel
// --------------------------------------------------------
type FirewallAuditAndQuarantineModel struct {
	CapabilityId types.String `tfsdk:"capability_id"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	Quarantine   types.Bool   `tfsdk:"quarantine"`
}

func NewFirewallAuditAndQuarantineModelWithDefaults() *FirewallAuditAndQuarantineModel {
	return &FirewallAuditAndQuarantineModel{
		CapabilityId: types.StringNull(),
		Enabled:      types.BoolValue(false),
		Quarantine:   types.BoolValue(false),
	}
}

// MapFromCapabilityDTO populates the model from a CapabilityDTO returned from the API
func (m *FirewallAuditAndQuarantineModel) MapFromCapabilityDTO(api *sonatyperepo.CapabilityDTO) {
	if api == nil {
		return
	}
	m.CapabilityId = types.StringPointerValue(api.Id)
	m.Enabled = types.BoolPointerValue(api.Enabled)
	if api.Properties != nil {
		if quarantineStr, ok := (*api.Properties)["quarantine"]; ok {
			if quarantine, err := strconv.ParseBool(quarantineStr); err == nil {
				m.Quarantine = types.BoolValue(quarantine)
			}
		}
	}
}

// FirewallAuditAndQuarantineWithPccsModel
// --------------------------------------------------------
type FirewallAuditAndQuarantineWithPccsModel struct {
	FirewallAuditAndQuarantineModel
	PccsEnabled types.Bool `tfsdk:"pccs_enabled"`
}

func NewFirewallAuditAndQuarantineWithPccsModelWithDefaults() *FirewallAuditAndQuarantineWithPccsModel {
	return &FirewallAuditAndQuarantineWithPccsModel{
		FirewallAuditAndQuarantineModel: *NewFirewallAuditAndQuarantineModelWithDefaults(),
		PccsEnabled:                     types.BoolValue(false),
	}
}

// MapFromCapabilityDTO populates the model from a CapabilityDTO returned from the API
func (m *FirewallAuditAndQuarantineWithPccsModel) MapFromCapabilityDTO(api *sonatyperepo.CapabilityDTO) {
	if api == nil {
		return
	}
	m.CapabilityId = types.StringPointerValue(api.Id)
	m.Enabled = types.BoolPointerValue(api.Enabled)
	if api.Properties != nil {
		if quarantineStr, ok := (*api.Properties)["quarantine"]; ok {
			if quarantine, err := strconv.ParseBool(quarantineStr); err == nil {
				m.Quarantine = types.BoolValue(quarantine)
			}
		}
	}
}

// RepositoryReplicationModel
// --------------------------------------------------------
type ProxyRemoveQuarrantiedModel struct {
	RemoveQuarrantined types.Bool `tfsdk:"remove_quarrantined"`
}

func (m *ProxyRemoveQuarrantiedModel) MapFromNpmApi(api *sonatyperepo.NpmAttributes) {
	m.RemoveQuarrantined = types.BoolValue(api.RemoveQuarantined)
}

func (m *ProxyRemoveQuarrantiedModel) MapToNpmApi(api *sonatyperepo.NpmAttributes) {
	api.RemoveQuarantined = m.RemoveQuarrantined.ValueBool()
}

func (m *ProxyRemoveQuarrantiedModel) MapFromPyPiApi(api *sonatyperepo.PyPiProxyAttributes) {
	m.RemoveQuarrantined = types.BoolValue(api.RemoveQuarantined)
}

func (m *ProxyRemoveQuarrantiedModel) MapToPyPiApi(api *sonatyperepo.PyPiProxyAttributes) {
	api.RemoveQuarantined = m.RemoveQuarrantined.ValueBool()
}
