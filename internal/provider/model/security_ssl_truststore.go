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
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

type SecuritySslTruststoreModel struct {
	Id                        types.String `tfsdk:"id"`
	Pem                       types.String `tfsdk:"pem"`
	Fingerprint               types.String `tfsdk:"fingerprint"`
	SerialNumber              types.String `tfsdk:"serial_number"`
	SubjectCommonName         types.String `tfsdk:"subject_common_name"`
	SubjectOrganization       types.String `tfsdk:"subject_organization"`
	SubjectOrganizationalUnit types.String `tfsdk:"subject_organizational_unit"`
	IssuerCommonName          types.String `tfsdk:"issuer_common_name"`
	IssuerOrganization        types.String `tfsdk:"issuer_organization"`
	IssuerOrganizationalUnit  types.String `tfsdk:"issuer_organizational_unit"`
	IssuedOn                  types.Int64  `tfsdk:"issued_on"`
	ExpiresOn                 types.Int64  `tfsdk:"expires_on"`
	LastUpdated               types.String `tfsdk:"last_updated"`
}

func (m *SecuritySslTruststoreModel) MapFromApi(api *sonatyperepo.ApiCertificate) {
	m.Id = types.StringPointerValue(api.Id)
	m.Pem = types.StringPointerValue(api.Pem)
	m.Fingerprint = types.StringPointerValue(api.Fingerprint)
	m.SerialNumber = types.StringPointerValue(api.SerialNumber)
	m.SubjectCommonName = types.StringPointerValue(api.SubjectCommonName)
	m.SubjectOrganization = types.StringPointerValue(api.SubjectOrganization)
	m.SubjectOrganizationalUnit = types.StringPointerValue(api.SubjectOrganizationalUnit)
	m.IssuerCommonName = types.StringPointerValue(api.IssuerCommonName)
	m.IssuerOrganization = types.StringPointerValue(api.IssuerOrganization)
	m.IssuerOrganizationalUnit = types.StringPointerValue(api.IssuerOrganizationalUnit)
	m.IssuedOn = types.Int64PointerValue(api.IssuedOn)
	m.ExpiresOn = types.Int64PointerValue(api.ExpiresOn)
}
