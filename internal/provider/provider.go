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

package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/blob_store"
	"terraform-provider-sonatyperepo/internal/provider/capability"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/content_selector"
	"terraform-provider-sonatyperepo/internal/provider/privilege"
	"terraform-provider-sonatyperepo/internal/provider/repository"
	"terraform-provider-sonatyperepo/internal/provider/role"
	"terraform-provider-sonatyperepo/internal/provider/system"
	"terraform-provider-sonatyperepo/internal/provider/task"
	"terraform-provider-sonatyperepo/internal/provider/user"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure SonatypeRepoProvider satisfies various provider interfaces.
var _ provider.Provider = &SonatypeRepoProvider{}

// SonatypeRepoProvider defines the provider implementation.
type SonatypeRepoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SonatypeRepoProviderModel describes the provider data model.
type SonatypeRepoProviderModel struct {
	Url         types.String `tfsdk:"url"`
	Username    types.String `tfsdk:"username"`
	Password    types.String `tfsdk:"password"`
	ApiBasePath types.String `tfsdk:"api_base_path"`
	VersionHint types.String `tfsdk:"version_hint"`
}

func (p *SonatypeRepoProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sonatyperepo"
	resp.Version = p.version
}

func (p *SonatypeRepoProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "Sonatype Nexus Repository Server URL",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for Sonatype Nexus Repository Server, requires role/permissions scoped to the resources you wish to manage",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for your user for Sonatype Nexus Repository Server",
				Required:            true,
				Sensitive:           true,
			},
			"api_base_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Base Path at which the API is present - defaults to `/service/rest`. This only needs to be set if you run Sonatype Nexus Repository at a Base Path that is not `/`.",
			},
			"version_hint": schema.StringAttribute{
				MarkdownDescription: `You can set this to the full version string (e.g. "3.85.0-03 (PRO)" or "3.80.0-06 (OSS)") of Sonatype Nexus Repository that you are connecting to.

> [!NOTE] 
> You can find the full version string in _Admin -> Support -> System Information_.
>				
> By default, this provider will attempt to automatically determine the version of Sonatype Nexus Repository you are connected to - but in some 
> real world cases, a Load Balancer or such may strip the HTTP Header that contians this information (the _Server_ header).

> [!TIP]
> If you receive an error such as ` + "`Plan is not supported for Sonatype Nexus Repository Manager: 0.0.0-0 (PRO=false)`" + ` then you should set 
> this attribute - otherwise, do not supply this attribute.
			`,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(\d+\.\d+\.\d+-\d+\s\((PRO|OSS)\))?$`),
						`Leave empty, or provide a version string in the format "3.85.0-03 (PRO)".`,
					),
				},
			},
		},
		MarkdownDescription: `Sonatype Nexus Repository must not be in read-only mode in order to use this Provider. This will be checked. 
		
Some resources and features depend on the version of Sonatype Nexus Repository you are running. See individual Data Source and Resource documentaiton for details.`,
	}
}

func (p *SonatypeRepoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SonatypeRepoProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	nxrmUrl, username, password, apiBasePath, versionHint := p.parseConfig(&config)

	p.validateConfig(resp, nxrmUrl, &config)
	if resp.Diagnostics.HasError() {
		return
	}

	ds := p.createClient(nxrmUrl, username, password, apiBasePath)

	p.checkVersion(ctx, &ds, resp, versionHint)

	resp.DataSourceData = ds
	resp.ResourceData = ds
}

func (p *SonatypeRepoProvider) parseConfig(config *SonatypeRepoProviderModel) (string, string, string, string, *string) {
	nxrmUrl := os.Getenv("NXRM_SERVER_URL")
	username := os.Getenv("NXRM_SERVER_USERNAME")
	password := os.Getenv("NXRM_SERVER_PASSWORD")
	apiBasePath := "/service/rest"
	var versionHint *string

	if !config.Url.IsNull() && len(config.Url.ValueString()) > 0 {
		nxrmUrl = config.Url.ValueString()
	}

	if !config.Username.IsNull() && len(config.Username.ValueString()) > 0 {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() && len(config.Password.ValueString()) > 0 {
		password = config.Password.ValueString()
	}

	if !config.ApiBasePath.IsNull() && len(config.ApiBasePath.ValueString()) > 0 {
		apiBasePath = config.ApiBasePath.ValueString()
	}

	if !config.VersionHint.IsNull() && len(config.VersionHint.ValueString()) > 0 {
		v := fmt.Sprintf("Nexus/%s", config.VersionHint.ValueString())
		versionHint = &v
	}

	return nxrmUrl, username, password, apiBasePath, versionHint
}

func (p *SonatypeRepoProvider) validateConfig(resp *provider.ConfigureResponse, nxrmUrl string, config *SonatypeRepoProviderModel) {
	if len(nxrmUrl) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown Sonatype Nexus Repository Server URL",
			"The provider is unable to work without a Sonatype Nexus Repository Server URL which should begin http:// or https://",
		)
	}

	if _, e := url.ParseRequestURI(nxrmUrl); e != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Invalid Sonatype Nexus Repository Server URL",
			"The provider is unable to work without a valid Sonatype Nexus Repository Server URL",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Username not supplied",
			"Administratrive credentials for your Sonatype Nexus Repository Server are required",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Username not supplied",
			"Administratrive credentials for your Sonatype Nexus Repository Server are required",
		)
	}
}

func (p *SonatypeRepoProvider) createClient(nxrmUrl, username, password, apiBasePath string) common.SonatypeDataSourceData {
	configuration := sonatyperepo.NewConfiguration()
	configuration.UserAgent = "sonatyperepo-terraform/" + p.version
	configuration.Servers = []sonatyperepo.ServerConfiguration{
		{
			URL:         fmt.Sprintf("%s%s", strings.TrimRight(nxrmUrl, "/"), strings.TrimRight(apiBasePath, "/")),
			Description: "Sonatype Nexus Repository Server",
		},
	}

	client := sonatyperepo.NewAPIClient(configuration)
	return common.SonatypeDataSourceData{
		Auth:    sonatyperepo.BasicAuth{UserName: username, Password: password},
		BaseUrl: strings.TrimRight(nxrmUrl, "/"),
		Client:  client,
	}
}

func (p *SonatypeRepoProvider) checkVersion(ctx context.Context, ds *common.SonatypeDataSourceData, resp *provider.ConfigureResponse, versionHint *string) {
	ds.CheckWritableAndGetVersion(ctx, &resp.Diagnostics, versionHint)
	tflog.Info(ctx, fmt.Sprintf("Detected Sonatype Nexus Repository to be version %s", ds.NxrmVersion.String()))

	if ds.NxrmVersion.OlderThan(3, 79, 1, 0) {
		resp.Diagnostics.AddWarning(
			`You are running against Sonatype Nexus Repository version older than 3.79.1`,
			`This provide has not been validated against versions older than 3.79.1 - things will probably work fine, but proceed with caution.`,
		)
	}
}

func (p *SonatypeRepoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		blob_store.NewBlobStoreFileResource,
		blob_store.NewBlobStoreGroupResource,
		blob_store.NewBlobStoreS3Resource,
		blob_store.NewBlobStoreGoogleCloudResource,
		capability.NewCapabilityAuditResource,
		capability.NewCapabilityCoreBaseUrlResource,
		capability.NewCapabilityCoreStorageSettingsResource,
		capability.NewCapabilityCustomS3RegionsResource,
		capability.NewCapabilityDefaultRoleResource,
		capability.NewCapabilityFirewallAuditQuarantineResource,
		capability.NewCapabilityHealthcheckResource,
		capability.NewCapabilityOutreachResource,
		capability.NewCapabilitySecurityRutAuthResource,
		capability.NewCapabilityUiBrandingResource,
		capability.NewCapabilityUiSettingsResource,
		capability.NewCapabilityWebhookGlobalResource,
		capability.NewCapabilityWebhookRepositoryResource,
		content_selector.NewContentSelectorResource,
		privilege.NewApplicationPrivilegeResource,
		privilege.NewRepositoryAdminPrivilegeResource,
		privilege.NewRepositoryContentSelectorPrivilegeResource,
		privilege.NewRepositoryViewPrivilegeResource,
		privilege.NewScriptPrivilegeResource,
		privilege.NewWildcardPrivilegeResource,
		repository.NewRepositoryAptHostedResource,
		repository.NewRepositoryAptProxyResource,
		repository.NewRepositoryCargoGroupResource,
		repository.NewRepositoryCargoHostedResource,
		repository.NewRepositoryCargoProxyResource,
		repository.NewRepositoryConanGroupResource,
		repository.NewRepositoryConanHostedResource,
		repository.NewRepositoryConanProxyResource,
		repository.NewRepositoryCocoaPodsProxyResource,
		repository.NewRepositoryComposerProxyResource,
		repository.NewRepositoryCondaProxyResource,
		repository.NewRepositoryDockerGroupResource,
		repository.NewRepositoryDockerHostedResource,
		repository.NewRepositoryDockerProxyResource,
		repository.NewRepositoryGitLfsHostedResource,
		repository.NewRepositoryGoGroupResource,
		repository.NewRepositoryGoProxyResource,
		repository.NewRepositoryHelmHostedResource,
		repository.NewRepositoryHelmProxyResource,
		repository.NewRepositoryHuggingFaceProxyResource,
		repository.NewRepositoryMavenGroupResource,
		repository.NewRepositoryMavenHostedResource,
		repository.NewRepositoryMavenProxyResource,
		repository.NewRepositoryNpmGroupResource,
		repository.NewRepositoryNpmHostedResource,
		repository.NewRepositoryNpmProxyResource,
		repository.NewRepositoryNugetGroupResource,
		repository.NewRepositoryNugetHostedResource,
		repository.NewRepositoryNugetProxyResource,
		repository.NewRepositoryP2ProxyResource,
		repository.NewRepositoryPyPiGroupResource,
		repository.NewRepositoryPyPiHostedResource,
		repository.NewRepositoryPyPiProxyResource,
		repository.NewRepositoryRGroupResource,
		repository.NewRepositoryRHostedResource,
		repository.NewRepositoryRProxyResource,
		repository.NewRepositoryRawGroupResource,
		repository.NewRepositoryRawHostedResource,
		repository.NewRepositoryRawProxyResource,
		repository.NewRepositoryRubyGemsHostedResource,
		repository.NewRepositoryRubyGemsGroupResource,
		repository.NewRepositoryRubyGemsProxyResource,
		repository.NewRepositoryYumGroupResource,
		repository.NewRepositoryYumHostedResource,
		repository.NewRepositoryYumProxyResource,
		repository.NewCleanupPolicyResource,
		repository.NewRoutingRuleResource,
		role.NewRoleResource,
		system.NewAnonymousAccessSystemResource,
		system.NewSystemConfigProductLicenseResource,
		system.NewSystemConfigLdapResource,
		system.NewSystemConfigMailResource,
		system.NewSystemConfigIqConnectionResource,
		system.NewSecurityRealmsResource,
		system.NewSecuritySamlResource,
		system.NewSecurityUserTokenResource,
		task.NewTaskBlobstoreCompactResource,
		task.NewTaskLicenseExpirationNotificationResource,
		task.NewTaskMalwareRemediatorResource,
		task.NewTaskRepairRebuildBrowseNodesResource,
		task.NewTaskRepositoryDockerGcResource,
		task.NewTaskRepositoryDockerUploadPurgeResource,
		task.NewTaskRepositoryMavenRemoveSnapshotsResource,
		user.NewUserResource,
	}
}

func (p *SonatypeRepoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		blob_store.BlobStoresDataSource,
		blob_store.BlobStoreFileDataSource,
		blob_store.BlobStoreGroupDataSource,
		blob_store.BlobStoreS3DataSource,
		capability.CapabilitiesDataSource,
		content_selector.ContentSelectorDataSource,
		content_selector.ContentSelectorsDataSource,
		privilege.PrivilegesDataSource,
		repository.RepositoriesDataSource,
		repository.RoutingRuleDataSource,
		repository.RoutingRulesDataSource,
		role.RolesDataSource,
		system.SecurityUserTokenDataSource,
		task.TaskDataSource,
		task.TasksDataSource,
		user.UsersDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SonatypeRepoProvider{
			version: version,
		}
	}
}
