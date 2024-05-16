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
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/blob_store"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/repository"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go"
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
				MarkdownDescription: "Base Path at which the API is present - defaults to /service/rest. This only needs to be set if you run Sonatype Nexus Repository at a Base Path that is not `/`.",
			},
		},
	}
}

func (p *SonatypeRepoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SonatypeRepoProviderModel

	diags := req.Config.Get(ctx, &config)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	nxrmUrl := os.Getenv("NXRM_SERVER_URL")
	username := os.Getenv("NXRM_SERVER_USERNAME")
	password := os.Getenv("NXRM_SERVER_PASSWORD")
	apiBasePath := "/service/rest"

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

	// Validate Provider Configuration
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

	if resp.Diagnostics.HasError() {
		return
	}

	// Example client configuration for data sources and resources
	configuration := sonatyperepo.NewConfiguration()
	configuration.UserAgent = "sonatyperepo-terraform/" + p.version
	configuration.Servers = []sonatyperepo.ServerConfiguration{
		{
			URL:         fmt.Sprintf("%s%s", strings.TrimRight(nxrmUrl, "/"), strings.TrimRight(apiBasePath, "/")),
			Description: "Sonatype Nexus Repository Server",
		},
	}

	client := sonatyperepo.NewAPIClient(configuration)
	resp.DataSourceData = common.SonatypeDataSourceData{
		Auth:    sonatyperepo.BasicAuth{UserName: username, Password: password},
		BaseUrl: strings.TrimRight(nxrmUrl, "/"),
		Client:  client,
	}
	resp.ResourceData = common.SonatypeDataSourceData{
		Auth:    sonatyperepo.BasicAuth{UserName: username, Password: password},
		BaseUrl: strings.TrimRight(nxrmUrl, "/"),
		Client:  client,
	}
}

func (p *SonatypeRepoProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		blob_store.NewBlobStoreFileResource,
		repository.NewRepositoryMavenGroupResource,
		repository.NewRepositoryMavenHostedResource,
		repository.NewRepositoryMavenProxyResource,
	}
}

func (p *SonatypeRepoProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		blob_store.BlobStoresDataSource,
		blob_store.BlobStoreFileDataSource,
		repository.RepositoriesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SonatypeRepoProvider{
			version: version,
		}
	}
}
