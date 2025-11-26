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

package task

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &tasksDataSource{}
	_ datasource.DataSourceWithConfigure = &tasksDataSource{}
)

// TasksDataSource is a helper function to simplify the provider implementation.
func TasksDataSource() datasource.DataSource {
	return &tasksDataSource{}
}

// tasksDataSource is the data source implementation.
type tasksDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *tasksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tasks"
}

// Schema defines the schema for the data source.
func (d *tasksDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = tfschema.Schema{
		Description: "Use this data source to get all Tasks",
		Attributes: map[string]tfschema.Attribute{
			"tasks": schema.DataSourceComputedListNestedAttribute(
				"List of Tasks",
				tfschema.NestedAttributeObject{
					Attributes: map[string]tfschema.Attribute{
						"id":   schema.DataSourceRequiredString("The ID of the Task."),
						"name": schema.DataSourceRequiredString("The name of the Task."),
						"type": schema.DataSourceRequiredString("The type of Task."),
					},
				},
			),
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *tasksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.TasksModel

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	tasksResponse, httpResponse, err := d.Client.TasksAPI.GetTasks(ctx).Execute()
	if err != nil {
		errors.HandleAPIError(
			"Unable to list tasks",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Iterating %d Tasks", len(tasksResponse.Items)))

	state.Tasks = make([]model.TaskModelSimple, 0)
	for _, task := range tasksResponse.Items {
		tflog.Debug(ctx, fmt.Sprintf("    Processing %s Task", *task.Id))
		taskModel := model.TaskModelSimple{}
		taskModel.MapFromApi(&task)

		state.Tasks = append(state.Tasks, taskModel)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
