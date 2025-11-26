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
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &BaseResource{}
	_ resource.ResourceWithConfigure   = &BaseResource{}
	_ resource.ResourceWithImportState = &BaseResource{}
)

// BaseResource is the resource implementation for Sonatype Nexus Repository resources.
// It extends basic resource functionality with Sonatype-specific configuration.
type BaseResource struct {
	Auth         sonatyperepo.BasicAuth
	BaseUrl      string
	Client       *sonatyperepo.APIClient
	NxrmVersion  SystemVersion
	NxrmWritable bool
}

// ImportState implements resource.ResourceWithImportState.
func (r *BaseResource) ImportState(context.Context, resource.ImportStateRequest, *resource.ImportStateResponse) {
	panic("unimplemented")
}

// Create implements resource.Resource.
func (*BaseResource) Create(context.Context, resource.CreateRequest, *resource.CreateResponse) {
	panic("unimplemented")
}

// Delete implements resource.Resource.
func (*BaseResource) Delete(context.Context, resource.DeleteRequest, *resource.DeleteResponse) {
	panic("unimplemented")
}

// Read implements resource.Resource.
func (*BaseResource) Read(context.Context, resource.ReadRequest, *resource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements resource.Resource.
func (*BaseResource) Schema(context.Context, resource.SchemaRequest, *resource.SchemaResponse) {
	panic("unimplemented")
}

// Update implements resource.Resource.
func (*BaseResource) Update(context.Context, resource.UpdateRequest, *resource.UpdateResponse) {
	panic("unimplemented")
}

// Metadata implements resource.Resource.
func (*BaseResource) Metadata(context.Context, resource.MetadataRequest, *resource.MetadataResponse) {
	panic("unimplemented")
}

// Configure implements resource.ResourceWithConfigure.
func (r *BaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(SonatypeDataSourceData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Type",
			fmt.Sprintf("Expected provider.SonatypeDataSourceData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.Auth = config.Auth
	r.BaseUrl = config.BaseUrl
	r.Client = config.Client
	r.NxrmVersion = config.NxrmVersion
	r.NxrmWritable = config.NxrmWritable
}

// AuthContext returns a new context with authentication set up for API calls
func (r *BaseResource) AuthContext(ctx context.Context) context.Context {
	return WithAuth(ctx, r.Auth)
}

// AuthConfig returns the authentication configuration
func (r *BaseResource) AuthConfig() sonatyperepo.BasicAuth {
	return r.Auth
}

// URL returns the API base URL
func (r *BaseResource) URL() string {
	return r.BaseUrl
}

// APIClient returns the API client
func (r *BaseResource) APIClient() *sonatyperepo.APIClient {
	return r.Client
}

// IsConfigured checks if the resource has been properly configured
func (r *BaseResource) IsConfigured() bool {
	return r.Client != nil
}
