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
	"encoding/json"
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type AnonymousAccessModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	RealmName   types.String `tfsdk:"realm_name"`
	UserId      types.String `tfsdk:"user_id"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

type EmailConfigurationModel struct {
	Enabled                       types.Bool   `tfsdk:"enabled"`
	Host                          types.String `tfsdk:"host"`
	Port                          types.Int64  `tfsdk:"port"`
	Username                      types.String `tfsdk:"username"`
	Password                      types.String `tfsdk:"password"`
	FromAddress                   types.String `tfsdk:"from_address"`
	SubjectPrefix                 types.String `tfsdk:"subject_prefix"`
	StartTLSEnabled               types.Bool   `tfsdk:"start_tls_enabled"`
	StartTLSRequired              types.Bool   `tfsdk:"start_tls_required"`
	SSLOnConnectEnabled           types.Bool   `tfsdk:"ssl_on_connect_enabled"`
	SSLServerIdentityCheckEnabled types.Bool   `tfsdk:"ssl_server_identity_check_enabled"`
	NexusTrustStoreEnabled        types.Bool   `tfsdk:"nexus_trust_store_enabled"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
}

func (m *EmailConfigurationModel) MapFromApi(api *sonatyperepo.ApiEmailConfiguration) {
	m.Enabled = types.BoolPointerValue(api.Enabled)
	m.Host = types.StringPointerValue(api.Host)
	m.Port = types.Int64Value(int64(api.Port))
	if api.Username != nil && *api.Username != "" {
		m.Username = types.StringPointerValue(api.Username)
	}
	// Skip password mapping for security reasons
	m.FromAddress = types.StringPointerValue(api.FromAddress)
	m.SubjectPrefix = types.StringPointerValue(api.SubjectPrefix)
	m.StartTLSEnabled = types.BoolPointerValue(api.StartTlsEnabled)
	m.StartTLSRequired = types.BoolPointerValue(api.StartTlsRequired)
	m.SSLOnConnectEnabled = types.BoolPointerValue(api.SslOnConnectEnabled)
	m.SSLServerIdentityCheckEnabled = types.BoolPointerValue(api.SslServerIdentityCheckEnabled)
	m.NexusTrustStoreEnabled = types.BoolPointerValue(api.NexusTrustStoreEnabled)
}

type IqConnectionModel struct {
	Enabled                types.Bool   `tfsdk:"enabled"`
	Url                    types.String `tfsdk:"url"`
	NexusTrustStoreEnabled types.Bool   `tfsdk:"nexus_trust_store_enabled"`
	AuthenticationMethod   types.String `tfsdk:"authentication_method"`
	Username               types.String `tfsdk:"username"`
	Password               types.String `tfsdk:"password"`
	ConnectionTimeout      types.Int32  `tfsdk:"connection_timeout"`
	Properties             types.String `tfsdk:"properties"`
	ShowIQServerLink       types.Bool   `tfsdk:"show_iq_server_link"`
	FailOpenModeEnabled    types.Bool   `tfsdk:"fail_open_mode_enabled"`
	LastUpdated            types.String `tfsdk:"last_updated"`
}

type SecurityRealmsModel struct {
	Active types.List   `tfsdk:"active"`
	ID     types.String `tfsdk:"id"`
}

type SecuritySamlModel struct {
	IdpMetadata                types.String `tfsdk:"idp_metadata"`
	UsernameAttribute          types.String `tfsdk:"username_attribute"`
	FirstNameAttribute         types.String `tfsdk:"first_name_attribute"`
	LastNameAttribute          types.String `tfsdk:"last_name_attribute"`
	EmailAttribute             types.String `tfsdk:"email_attribute"`
	GroupsAttribute            types.String `tfsdk:"groups_attribute"`
	ValidateResponseSignature  types.Bool   `tfsdk:"validate_response_signature"`
	ValidateAssertionSignature types.Bool   `tfsdk:"validate_assertion_signature"`
	EntityId                   types.String `tfsdk:"entity_id"`
}

func (m *SecuritySamlModel) MapFromApi(api *sonatyperepo.SamlConfigurationXO) {
	if m.EntityId.IsNull() || m.EntityId.IsUnknown() {
		m.EntityId = types.StringNull()
	} else {
		m.EntityId = types.StringPointerValue(api.EntityId)
	}
	m.IdpMetadata = types.StringValue(api.IdpMetadata)
	m.UsernameAttribute = types.StringValue(api.UsernameAttribute)
	m.FirstNameAttribute = types.StringPointerValue(api.FirstNameAttribute)
	m.LastNameAttribute = types.StringPointerValue(api.LastNameAttribute)
	m.EmailAttribute = types.StringPointerValue(api.EmailAttribute)
	m.GroupsAttribute = types.StringPointerValue(api.GroupsAttribute)
	m.ValidateResponseSignature = types.BoolPointerValue(api.ValidateResponseSignature)
	m.ValidateAssertionSignature = types.BoolPointerValue(api.ValidateAssertionSignature)
}

func (m *IqConnectionModel) MapFromApi(api *sonatyperepo.IqConnectionXo) {
	m.AuthenticationMethod = types.StringValue(api.AuthenticationType)
	m.Enabled = types.BoolPointerValue(api.Enabled)
	m.NexusTrustStoreEnabled = types.BoolPointerValue(api.UseTrustStoreForUrl)
	m.Url = types.StringPointerValue(api.Url)
	m.Username = types.StringPointerValue(api.Username)
	// Skip password
	m.ConnectionTimeout = types.Int32PointerValue(api.TimeoutSeconds)
	m.Properties = types.StringPointerValue(api.Properties)
	m.ShowIQServerLink = types.BoolPointerValue(api.ShowLink)
	m.FailOpenModeEnabled = types.BoolPointerValue(api.FailOpenModeEnabled)
}

func (m *IqConnectionModel) MapToApi(api *sonatyperepo.IqConnectionXo) {
	api.Enabled = m.Enabled.ValueBoolPointer()
	api.Url = m.Url.ValueStringPointer()
	api.AuthenticationType = m.AuthenticationMethod.ValueString()
	api.Username = m.Username.ValueStringPointer()
	api.Password = m.Password.ValueStringPointer()
	api.TimeoutSeconds = m.ConnectionTimeout.ValueInt32Pointer()
	api.Properties = m.Properties.ValueStringPointer()
	api.ShowLink = m.ShowIQServerLink.ValueBoolPointer()
	api.FailOpenModeEnabled = m.FailOpenModeEnabled.ValueBoolPointer()
}

type LdapServerModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Protocol               types.String `tfsdk:"protocol"`
	NexusTrustStoreEnabled types.Bool   `tfsdk:"nexus_trust_store_enabled"`
	Hostname               types.String `tfsdk:"hostname"`
	Port                   types.Int32  `tfsdk:"port"`
	SearchBase             types.String `tfsdk:"search_base"`
	AuthScheme             types.String `tfsdk:"auth_scheme"`
	AuthUsername           types.String `tfsdk:"auth_username"`
	AuthPassword           types.String `tfsdk:"auth_password"`
	AuthRealm              types.String `tfsdk:"auth_realm"`
	ConnectionTimeout      types.Int32  `tfsdk:"connection_timeout"`
	ConnectionRetryDelay   types.Int32  `tfsdk:"connection_retry_delay"`
	MaxConnectionAttempts  types.Int32  `tfsdk:"max_connection_attempts"`
	Order                  types.Int32  `tfsdk:"order"`
	// User Mapping
	UserBaseDn            types.String `tfsdk:"user_base_dn"`
	UserSubtree           types.Bool   `tfsdk:"user_subtree"`
	UserObjectClass       types.String `tfsdk:"user_object_class"`
	UserLdapFilter        types.String `tfsdk:"user_ldap_filter"`
	UserIdAttribute       types.String `tfsdk:"user_id_attribute"`
	UserRealNameAttribute types.String `tfsdk:"user_real_name_attribute"`
	UserEmailAttribute    types.String `tfsdk:"user_email_name_attribute"`
	UserPasswordAttribute types.String `tfsdk:"user_password_attribute"`
	// Group Mapping
	MapLdapGroupsAsRoles  types.Bool   `tfsdk:"map_ldap_groups_to_roles"`
	UserMemberOfAttribute types.String `tfsdk:"user_member_of_attribute"`
	GroupType             types.String `tfsdk:"group_type"`
	GroupBaseDn           types.String `tfsdk:"group_base_dn"`
	GroupSubtree          types.Bool   `tfsdk:"group_subtree"`
	GroupObjectClass      types.String `tfsdk:"group_object_class"`
	GroupIdAttribute      types.String `tfsdk:"group_id_attribute"`
	GroupMemberAttribute  types.String `tfsdk:"group_member_attribute"`
	GroupMemberFormat     types.String `tfsdk:"group_member_format"`
	// Meta for Terraform
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (model *LdapServerModel) FromApiModel(api *sonatyperepo.ReadLdapServerXo) {
	model.Id = types.StringPointerValue(api.Id)
	model.Name = types.StringPointerValue(&api.Name)
	model.Protocol = types.StringPointerValue(&api.Protocol)
	model.NexusTrustStoreEnabled = types.BoolPointerValue(api.UseTrustStore)
	model.Hostname = types.StringPointerValue(&api.Host)
	model.Port = types.Int32PointerValue(&api.Port)
	model.SearchBase = types.StringPointerValue(&api.SearchBase)
	model.AuthScheme = types.StringPointerValue(&api.AuthScheme)
	model.AuthUsername = types.StringPointerValue(api.AuthUsername)
	// model.AuthPassword = types.StringPointerValue(api.???)
	model.AuthRealm = types.StringPointerValue(api.AuthRealm)
	model.ConnectionTimeout = types.Int32PointerValue(&api.ConnectionTimeoutSeconds)
	model.ConnectionRetryDelay = types.Int32PointerValue(&api.ConnectionRetryDelaySeconds)
	model.MaxConnectionAttempts = types.Int32PointerValue(&api.MaxIncidentsCount)
	model.Order = types.Int32PointerValue(api.Order)
	model.UserBaseDn = types.StringPointerValue(api.UserBaseDn)
	model.UserSubtree = types.BoolPointerValue(api.UserSubtree)
	model.UserObjectClass = types.StringPointerValue(api.UserObjectClass)
	model.UserLdapFilter = types.StringPointerValue(api.UserLdapFilter)
	model.UserIdAttribute = types.StringPointerValue(api.UserIdAttribute)
	model.UserPasswordAttribute = types.StringPointerValue(api.UserPasswordAttribute)
	model.UserMemberOfAttribute = types.StringPointerValue(api.UserMemberOfAttribute)
	model.MapLdapGroupsAsRoles = types.BoolPointerValue(api.LdapGroupsAsRoles)
	model.GroupType = types.StringPointerValue(api.GroupType)
	model.GroupBaseDn = types.StringPointerValue(api.GroupBaseDn)
	model.GroupSubtree = types.BoolPointerValue(api.GroupSubtree)
	model.GroupObjectClass = types.StringPointerValue(api.GroupObjectClass)
	model.GroupIdAttribute = types.StringPointerValue(api.GroupIdAttribute)
	model.GroupMemberAttribute = types.StringPointerValue(api.GroupMemberAttribute)
	model.GroupMemberFormat = types.StringPointerValue(api.GroupMemberFormat)
}

func (model *LdapServerModel) ToApiCreateModel() *sonatyperepo.CreateLdapServerXo {
	apiModel := sonatyperepo.CreateLdapServerXo{
		Name:                        model.Name.ValueString(),
		Protocol:                    model.Protocol.ValueString(),
		Host:                        model.Hostname.ValueString(),
		Port:                        model.Port.ValueInt32(),
		SearchBase:                  model.SearchBase.ValueString(),
		AuthScheme:                  model.AuthScheme.ValueString(),
		ConnectionTimeoutSeconds:    model.ConnectionTimeout.ValueInt32(),
		ConnectionRetryDelaySeconds: model.ConnectionRetryDelay.ValueInt32(),
		MaxIncidentsCount:           model.MaxConnectionAttempts.ValueInt32(),
		UserBaseDn:                  model.UserBaseDn.ValueStringPointer(),
		UserSubtree:                 model.UserSubtree.ValueBoolPointer(),
		UserObjectClass:             model.UserObjectClass.ValueStringPointer(),
		UserLdapFilter:              model.UserLdapFilter.ValueStringPointer(),
		UserIdAttribute:             model.UserIdAttribute.ValueStringPointer(),
		UserRealNameAttribute:       model.UserRealNameAttribute.ValueStringPointer(),
		UserEmailAddressAttribute:   model.UserEmailAttribute.ValueStringPointer(),
		UserPasswordAttribute:       model.UserPasswordAttribute.ValueStringPointer(),
		UserMemberOfAttribute:       model.UserMemberOfAttribute.ValueStringPointer(),
		LdapGroupsAsRoles:           model.MapLdapGroupsAsRoles.ValueBoolPointer(),
	}
	if apiModel.Protocol == common.PROTOCOL_LDAPS {
		apiModel.UseTrustStore = model.NexusTrustStoreEnabled.ValueBoolPointer()
	}
	if apiModel.AuthScheme != common.AUTH_SCHEME_NONE {
		apiModel.AuthUsername = model.AuthUsername.ValueStringPointer()
		apiModel.AuthPassword = *model.AuthPassword.ValueStringPointer()
	}
	if apiModel.AuthScheme == common.AUTH_SCHEME_DIGEST_MD5 || apiModel.AuthScheme == common.AUTH_SCHEME_CRAM_MD5 {
		apiModel.AuthRealm = model.AuthRealm.ValueStringPointer()
	}
	if model.MapLdapGroupsAsRoles.ValueBool() {
		apiModel.LdapGroupsAsRoles = common.NewTrue()
		apiModel.GroupType = model.GroupType.ValueStringPointer()
		if *apiModel.GroupType == common.LDAP_GROUP_MAPPING_DYNAMIC {
			apiModel.GroupMemberAttribute = model.GroupMemberAttribute.ValueStringPointer()
		}
		if *apiModel.GroupType == common.LDAP_GROUP_MAPPING_STATIC {
			apiModel.GroupBaseDn = model.GroupBaseDn.ValueStringPointer()
			apiModel.GroupSubtree = model.GroupSubtree.ValueBoolPointer()
			apiModel.GroupObjectClass = model.GroupObjectClass.ValueStringPointer()
			apiModel.GroupIdAttribute = model.GroupIdAttribute.ValueStringPointer()
			apiModel.GroupMemberAttribute = model.GroupMemberAttribute.ValueStringPointer()
			apiModel.GroupMemberFormat = model.GroupMemberFormat.ValueStringPointer()
		}
	} else {
		apiModel.LdapGroupsAsRoles = common.NewFalse()
	}

	return &apiModel
}

func (model *LdapServerModel) ToApiUpdateModel() *sonatyperepo.UpdateLdapServerXo {
	creatModel := model.ToApiCreateModel()
	temp, _ := json.Marshal(creatModel)
	var updateModel sonatyperepo.UpdateLdapServerXo
	_ = json.Unmarshal(temp, &updateModel)
	updateModel.Id = model.Id.ValueStringPointer()
	return &updateModel
}

type SecurityUserTokenModel struct {
	ID                types.String `tfsdk:"id"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	ExpirationDays    types.Int32  `tfsdk:"expiration_days"`
	ExpirationEnabled types.Bool   `tfsdk:"expiration_enabled"`
	ProtectContent    types.Bool   `tfsdk:"protect_content"`
	LastUpdated       types.String `tfsdk:"last_updated"`
}

func (m *SecurityUserTokenModel) MapFromApi(api *sonatyperepo.UserTokensApiModel) {
	m.Enabled = types.BoolPointerValue(api.Enabled)
	m.ExpirationDays = types.Int32PointerValue(api.ExpirationDays)
	m.ExpirationEnabled = types.BoolPointerValue(api.ExpirationEnabled)
	m.ProtectContent = types.BoolPointerValue(api.ProtectContent)
}

func (m *SecurityUserTokenModel) MapToApi(api *sonatyperepo.UserTokensApiModel) {
	api.Enabled = m.Enabled.ValueBoolPointer()
	// Set ExpirationDays to default value if not specified, to satisfy API's minimum value requirement
	if !m.ExpirationDays.IsNull() && !m.ExpirationDays.IsUnknown() {
		api.ExpirationDays = m.ExpirationDays.ValueInt32Pointer()
	} else {
		defaultValue := common.SECURITY_USER_TOKEN_DEFAULT_EXPIRATION_DAYS
		api.ExpirationDays = &defaultValue
	}
	api.ExpirationEnabled = m.ExpirationEnabled.ValueBoolPointer()
	api.ProtectContent = m.ProtectContent.ValueBoolPointer()
}
