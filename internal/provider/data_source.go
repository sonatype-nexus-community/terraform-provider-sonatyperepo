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

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &baseDataSource{}
	_ datasource.DataSourceWithConfigure = &baseDataSource{}
)

// applicationsDataSource is the data source implementation.
type baseDataSource struct {
	client *sonatyperepo.APIClient
	auth   sonatyperepo.BasicAuth
}

// Configure implements datasource.DataSourceWithConfigure.
func (d *baseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = config.client
	d.auth = config.auth
}

// Metadata implements datasource.DataSource.
func (*baseDataSource) Metadata(context.Context, datasource.MetadataRequest, *datasource.MetadataResponse) {
	panic("unimplemented")
}

// Read implements datasource.DataSource.
func (*baseDataSource) Read(context.Context, datasource.ReadRequest, *datasource.ReadResponse) {
	panic("unimplemented")
}

// Schema implements datasource.DataSource.
func (*baseDataSource) Schema(context.Context, datasource.SchemaRequest, *datasource.SchemaResponse) {
	panic("unimplemented")
}
