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
