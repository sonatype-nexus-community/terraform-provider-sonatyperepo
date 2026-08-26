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
)

// OCI (Open Container Initiative) repositories were introduced in NXRM 3.94.0 and have no
// equivalent type in the V382 ("v3") generated client generation at all -- unlike every other
// format in this package, there is no vendored V382 struct to express RepositoryManagementService's
// OCI methods in terms of. Rather than patch the vendored client (never done in this codebase) or
// type these methods directly against the V395 client (which would break every shared model helper
// in internal/provider/model that's hard-typed against V382's generic attribute types, e.g.
// repositoryProxyModel.MapFromApi(*sonatyperepoV382.ProxyAttributes)), these structs locally declare
// the same JSON shape NXRM's own OpenAPI schema describes for OCI, built out of the generic V382
// attribute types that already exist (StorageAttributes, ProxyAttributes, etc. are shared by every
// format, not OCI-specific) plus small OCI-only types mirroring sonatyperepoV395's OciAttributes/
// OciProxyAttributes/OciCosignConfiguration. jsonBridge can then convert between these and the real
// V395 client types exactly as it does for every V382-typed format, and the model package can reuse
// its existing V382-typed mapping helpers unchanged.

// OciAttributes mirrors sonatyperepoV395.OciAttributes.
type OciAttributes struct {
	ForceBasicAuth bool    `json:"forceBasicAuth"`
	HttpPort       *int32  `json:"httpPort,omitempty"`
	HttpsPort      *int32  `json:"httpsPort,omitempty"`
	PathEnabled    *bool   `json:"pathEnabled,omitempty"`
	Subdomain      *string `json:"subdomain,omitempty"`
	V1Enabled      bool    `json:"v1Enabled"`
}

// OciProxyAttributes mirrors sonatyperepoV395.OciProxyAttributes.
type OciProxyAttributes struct {
	CacheForeignLayers       *bool    `json:"cacheForeignLayers,omitempty"`
	ForeignLayerUrlWhitelist []string `json:"foreignLayerUrlWhitelist,omitempty"`
	IndexType                *string  `json:"indexType,omitempty"`
	IndexUrl                 *string  `json:"indexUrl,omitempty"`
}

// OciCosignConfiguration mirrors sonatyperepoV395.OciCosignConfiguration.
type OciCosignConfiguration struct {
	Enforcement   string  `json:"enforcement"`
	IdentityRegex *string `json:"identityRegex,omitempty"`
	IssuerRegex   *string `json:"issuerRegex,omitempty"`
}

// OciHostedApiRepository mirrors sonatyperepoV395.OciHostedApiRepository.
type OciHostedApiRepository struct {
	Cleanup   *sonatyperepoV382.CleanupPolicyAttributes      `json:"cleanup,omitempty"`
	Component *sonatyperepoV382.ComponentAttributes          `json:"component,omitempty"`
	Cosign    *OciCosignConfiguration                        `json:"cosign,omitempty"`
	Name      *string                                        `json:"name,omitempty"`
	Oci       OciAttributes                                  `json:"oci"`
	Online    bool                                           `json:"online"`
	Storage   sonatyperepoV382.DockerHostedStorageAttributes `json:"storage"`
	Format    *string                                        `json:"format,omitempty"`
	Type      *string                                        `json:"type,omitempty"`
	Url       *string                                        `json:"url,omitempty"`
}

// OciHostedRepositoryApiRequest mirrors sonatyperepoV395.OciHostedRepositoryApiRequest.
type OciHostedRepositoryApiRequest struct {
	Cleanup   *sonatyperepoV382.CleanupPolicyAttributes      `json:"cleanup,omitempty"`
	Component *sonatyperepoV382.ComponentAttributes          `json:"component,omitempty"`
	Cosign    *OciCosignConfiguration                        `json:"cosign,omitempty"`
	Name      string                                         `json:"name"`
	Oci       OciAttributes                                  `json:"oci"`
	Online    bool                                           `json:"online"`
	Storage   sonatyperepoV382.DockerHostedStorageAttributes `json:"storage"`
}

// OciProxyApiRepository mirrors sonatyperepoV395.OciProxyApiRepository, minus the inline
// `firewall` field -- that is threaded separately as a *FirewallMode return value, exactly as
// every other proxy format's V382-shaped Get method does.
type OciProxyApiRepository struct {
	Cleanup         *sonatyperepoV382.CleanupPolicyAttributes `json:"cleanup,omitempty"`
	Cosign          *OciCosignConfiguration                   `json:"cosign,omitempty"`
	HttpClient      sonatyperepoV382.HttpClientAttributes     `json:"httpClient"`
	Name            *string                                   `json:"name,omitempty"`
	NegativeCache   sonatyperepoV382.NegativeCacheAttributes  `json:"negativeCache"`
	Oci             OciAttributes                             `json:"oci"`
	OciProxy        *OciProxyAttributes                       `json:"ociProxy,omitempty"`
	Online          bool                                      `json:"online"`
	Proxy           sonatyperepoV382.ProxyAttributes          `json:"proxy"`
	Replication     *sonatyperepoV382.ReplicationAttributes   `json:"replication,omitempty"`
	RoutingRuleName *string                                   `json:"routingRuleName,omitempty"`
	Storage         sonatyperepoV382.StorageAttributes        `json:"storage"`
	Format          *string                                   `json:"format,omitempty"`
	Type            *string                                   `json:"type,omitempty"`
	Url             *string                                   `json:"url,omitempty"`
}

// OciProxyRepositoryApiRequest mirrors sonatyperepoV395.OciProxyRepositoryApiRequest, minus the
// inline `firewall` field -- passed separately as a *FirewallMode parameter, as with every other
// proxy format's V382-shaped Create/Update methods.
type OciProxyRepositoryApiRequest struct {
	Cleanup       *sonatyperepoV382.CleanupPolicyAttributes `json:"cleanup,omitempty"`
	Cosign        *OciCosignConfiguration                   `json:"cosign,omitempty"`
	HttpClient    sonatyperepoV382.HttpClientAttributes     `json:"httpClient"`
	Name          string                                    `json:"name"`
	NegativeCache sonatyperepoV382.NegativeCacheAttributes  `json:"negativeCache"`
	Oci           OciAttributes                             `json:"oci"`
	OciProxy      *OciProxyAttributes                       `json:"ociProxy,omitempty"`
	Online        bool                                      `json:"online"`
	Proxy         sonatyperepoV382.ProxyAttributes          `json:"proxy"`
	Replication   *sonatyperepoV382.ReplicationAttributes   `json:"replication,omitempty"`
	RoutingRule   *string                                   `json:"routingRule,omitempty"`
	Storage       sonatyperepoV382.StorageAttributes        `json:"storage"`
}

// OciGroupApiRepository mirrors sonatyperepoV395.OciGroupApiRepository.
type OciGroupApiRepository struct {
	Cosign  *OciCosignConfiguration                `json:"cosign,omitempty"`
	Group   sonatyperepoV382.GroupDeployAttributes `json:"group"`
	Name    *string                                `json:"name,omitempty"`
	Oci     OciAttributes                          `json:"oci"`
	Online  bool                                   `json:"online"`
	Storage sonatyperepoV382.StorageAttributes     `json:"storage"`
	Format  *string                                `json:"format,omitempty"`
	Type    *string                                `json:"type,omitempty"`
	Url     *string                                `json:"url,omitempty"`
}

// OciGroupRepositoryApiRequest mirrors sonatyperepoV395.OciGroupRepositoryApiRequest.
type OciGroupRepositoryApiRequest struct {
	Cosign  *OciCosignConfiguration                `json:"cosign,omitempty"`
	Group   sonatyperepoV382.GroupDeployAttributes `json:"group"`
	Name    string                                 `json:"name"`
	Oci     OciAttributes                          `json:"oci"`
	Online  bool                                   `json:"online"`
	Storage sonatyperepoV382.StorageAttributes     `json:"storage"`
}
