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
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &taskDataSource{}
	_ datasource.DataSourceWithConfigure = &taskDataSource{}
)

// TaskDataSource is a helper function to simplify the provider implementation.
func TaskDataSource() datasource.DataSource {
	return &taskDataSource{}
}

// taskDataSource is the data source implementation.
type taskDataSource struct {
	common.BaseDataSource
}

// Metadata returns the data source type name.
func (d *taskDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

// Schema defines the schema for the data source.
func (d *taskDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get a single Task by ID",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The name of the Task.",
				Required:    true,
				Optional:    false,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Task.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of Task.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			// 			"enabled": schema.BoolAttribute{
			// 				Description: "Indicates if the task is enabled.",
			// 				Required:    false,
			// 				Optional:    false,
			// 				Computed:    true,
			// 			},
			// 			"alert_email": schema.StringAttribute{
			// 				Description: "E-mail address for task notifications.",
			// 				Required:    false,
			// 				Optional:    false,
			// 				Computed:    true,
			// 			},
			// 			"notification_condition": schema.StringAttribute{
			// 				Description: "The type of Task.",
			// 				Required:    false,
			// 				Optional:    false,
			// 				Computed:    true,
			// 				Validators: []validator.String{
			// 					stringvalidator.OneOf(
			// 						common.NOTIFICATION_CONDITION_FAILURE, common.NOTIFICATION_CONDITION_SUCCESS_OR_FAILURE,
			// 					),
			// 				},
			// 			},
			// 			"frequency": schema.SingleNestedAttribute{
			// 				Description: "Frequency Schedule for this Task.",
			// 				Required:    false,
			// 				Optional:    false,
			// 				Computed:    true,
			// 				Attributes: map[string]schema.Attribute{
			// 					"schedule": schema.StringAttribute{
			// 						Description: "Type of Schedule.",
			// 						Required:    true,
			// 						Optional:    false,
			// 						Validators: []validator.String{
			// 							stringvalidator.OneOf(
			// 								common.FREQUENCY_SCHEDULE_MANUAL,
			// 								common.FREQUENCY_SCHEDULE_ONCE,
			// 								common.FREQUENCY_SCHEDULE_HOURLY,
			// 								common.FREQUENCY_SCHEDULE_DAILY,
			// 								common.FREQUENCY_SCHEDULE_WEEKLY,
			// 								common.FREQUENCY_SCHEDULE_MONTHLY,
			// 								common.FREQUENCY_SCHEDULE_CRON,
			// 							),
			// 						},
			// 					},
			// 					"start_date": schema.Int64Attribute{
			// 						Description: "Start date of the task represented in unix timestamp. Does not apply for \"manual\" schedule.",
			// 						Required:    false,
			// 						Optional:    true,
			// 					},
			// 					"timezone_offset": schema.StringAttribute{
			// 						Description: "The offset time zone of the client. Example: -05:00",
			// 						Required:    false,
			// 						Optional:    true,
			// 					},
			// 					"recurring_days": schema.ListAttribute{
			// 						MarkdownDescription: `Array with the number of the days the task must run.

			// - For "weekly" schedule allowed values, 1 to 7.
			// - For "monthly" schedule allowed values, 1 to 31.`,
			// 						ElementType: types.Int32Type,
			// 						Required:    false,
			// 						Optional:    true,
			// 						Validators: []validator.List{
			// 							listvalidator.SizeAtLeast(1),
			// 						},
			// 					},
			// 					"cron_expression": schema.StringAttribute{
			// 						Description: "Cron expression for the task. Only applies for for \"cron\" schedule.",
			// 						Required:    false,
			// 						Optional:    true,
			// 					},
			// 				},
			// 			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *taskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data model.TaskModelSimple

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	if data.Id.IsNull() {
		resp.Diagnostics.AddError("Name must not be empty.", "Name must be provided.")
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		d.Auth,
	)

	taskResponse, httpResponse, err := d.Client.TasksAPI.GetTaskById(ctx, data.Id.ValueString()).Execute()

	state := model.TaskModelSimple{}
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.Diagnostics.AddWarning(
				"No Task with supplied ID",
				fmt.Sprintf("No Task with supplied ID: %s", httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error finding Task",
				fmt.Sprintf("Error finding Task: %s", httpResponse.Status),
			)
			return
		}
	} else if httpResponse.StatusCode == http.StatusOK {
		state.MapFromApi(taskResponse)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
