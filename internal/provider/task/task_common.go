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
	"reflect"
	"slices"
	"terraform-provider-sonatyperepo/internal/provider/common"
	tasktype "terraform-provider-sonatyperepo/internal/provider/task/task_type"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

const (
	TASK_ERROR_RESPONSE_PREFIX           = "Error response: "
	TASK_GENERAL_ERROR_RESPONSE_GENERAL  = TASK_ERROR_RESPONSE_PREFIX + " %s"
	TASK_GENERAL_ERROR_RESPONSE_WITH_ERR = TASK_ERROR_RESPONSE_PREFIX + " %s - %s"
	TASK_ERROR_DID_NOT_EXIST             = "%s Task did not exist to %s"
)

// Generic to all Task Resources
type taskResource struct {
	common.BaseResource
	TaskType tasktype.TaskTypeI
}

// Metadata returns the resource type name.
func (t *taskResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, t.TaskType.GetResourceName())
}

// Set Schema for this Resource
func (t *taskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = getTaskSchema(t.TaskType)
}

// This allows users to import existing Tasks into Terraform state.
func (r *taskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the Task ID as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (t *taskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan, diags := t.TaskType.GetPlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting Plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		t.Auth,
	)

	// Make API requet
	taskCreateResponse, httpResponse, err := t.TaskType.DoCreateRequest(plan, t.Client, ctx, t.NxrmVersion)

	// Handle Errors
	if err != nil {
		common.HandleApiError(
			fmt.Sprintf("Error creating %s Task", t.TaskType.GetType().String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}
	if !slices.Contains(t.TaskType.GetApiCreateSuccessResponseCodes(), httpResponse.StatusCode) {
		common.HandleApiError(
			fmt.Sprintf("Creation of %s Task was not successful", t.TaskType.GetType().String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	plan = t.TaskType.UpdateStateFromApi(plan, *taskCreateResponse)
	plan = t.TaskType.UpdatePlanForState(plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (t *taskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// EMPTY
}

func (t *taskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	planModel, diags := t.TaskType.GetPlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	// Retrieve values from state
	stateModel, diags := t.TaskType.GetStateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		t.Auth,
	)

	// Make API requet
	httpResponse, err := t.TaskType.DoUpdateRequest(planModel, stateModel, t.Client, ctx, t.NxrmVersion)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			common.HandleApiWarning(
				fmt.Sprintf(TASK_ERROR_DID_NOT_EXIST, t.TaskType.GetType().String(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			common.HandleApiError(
				fmt.Sprintf(TASK_ERROR_DID_NOT_EXIST, t.TaskType.GetType().String(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	planModel = t.TaskType.UpdateStateFromPlanForUpdate(planModel, stateModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, planModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (t *taskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	state, diags := t.TaskType.GetStateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Handle any errors
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		t.Auth,
	)

	// Make API request
	taskIdStructField := reflect.Indirect(reflect.ValueOf(state)).FieldByName("Id").Interface()
	taskId, ok := taskIdStructField.(basetypes.StringValue)
	if !ok {
		resp.Diagnostics.AddError(
			"Failed to determine Task ID to delete from State",
			fmt.Sprintf("%s %s", TASK_ERROR_RESPONSE_PREFIX, taskIdStructField),
		)
		return
	}

	attempts := 1
	maxAttempts := 3
	success := false

	for !success && attempts < maxAttempts {
		httpResponse, err := t.Client.TasksAPI.DeleteTaskById(ctx, taskId.ValueString()).Execute()

		// Trap 500 Error as they occur when Repo is not in appropriate internal state
		if httpResponse.StatusCode == http.StatusInternalServerError {
			tflog.Info(ctx, fmt.Sprintf("Unexpected response when deleting Task %s (attempt %d)", t.TaskType.GetType().String(), attempts))
			attempts++
			continue
		}

		if err != nil {
			if httpResponse.StatusCode == http.StatusNotFound {
				resp.State.RemoveResource(ctx)
				resp.Diagnostics.AddWarning(
					fmt.Sprintf(TASK_ERROR_DID_NOT_EXIST, t.TaskType.GetType().String(), "delete"),
					fmt.Sprintf(TASK_GENERAL_ERROR_RESPONSE_GENERAL, httpResponse.Status),
				)
			} else {
				resp.Diagnostics.AddError(
					fmt.Sprintf(TASK_ERROR_DID_NOT_EXIST, t.TaskType.GetType().String(), "delete"),
					fmt.Sprintf(TASK_GENERAL_ERROR_RESPONSE_WITH_ERR, httpResponse.Status, err),
				)
			}
			return
		}
		if httpResponse.StatusCode != http.StatusNoContent {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unexpected response when deleting %s Task (attempt %d)", t.TaskType.GetType().String(), attempts),
				fmt.Sprintf("Error response: %s", httpResponse.Status),
			)

			time.Sleep(1 * time.Second)
			attempts++
		} else {
			success = true
		}
	}
}

func getTaskSchema(tt tasktype.TaskTypeI) schema.Schema {
	return schema.Schema{
		Description: fmt.Sprintf("Manage Task '%s' (%s)", tt.GetPublicName(), tt.GetType().String()),
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The internal ID of the Task.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Task.",
				Required:    true,
				Optional:    false,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates if the task is enabled.",
				Required:    true,
				Optional:    false,
			},
			"alert_email": schema.StringAttribute{
				Description: "E-mail address for task notifications.",
				Required:    false,
				Optional:    true,
			},
			"notification_condition": schema.StringAttribute{
				Description: "The type of Task.",
				Required:    true,
				Optional:    false,
				Validators: []validator.String{
					stringvalidator.OneOf(
						common.NOTIFICATION_CONDITION_FAILURE, common.NOTIFICATION_CONDITION_SUCCESS_OR_FAILURE,
					),
				},
			},
			"frequency": schema.SingleNestedAttribute{
				Description: "Frequency Schedule for this Task.",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"schedule": schema.StringAttribute{
						Description: "Type of Schedule.",
						Required:    true,
						Optional:    false,
						Validators: []validator.String{
							stringvalidator.OneOf(
								common.FREQUENCY_SCHEDULE_MANUAL,
								common.FREQUENCY_SCHEDULE_ONCE,
								common.FREQUENCY_SCHEDULE_HOURLY,
								common.FREQUENCY_SCHEDULE_DAILY,
								common.FREQUENCY_SCHEDULE_WEEKLY,
								common.FREQUENCY_SCHEDULE_MONTHLY,
								common.FREQUENCY_SCHEDULE_CRON,
							),
						},
					},
					"start_date": schema.Int32Attribute{
						Description: "Start date of the task represented in unix timestamp. Does not apply for \"manual\" schedule.",
						Required:    false,
						Optional:    true,
					},
					"timezone_offset": schema.StringAttribute{
						Description: "The offset time zone of the client. Example: -05:00",
						Required:    false,
						Optional:    true,
					},
					"recurring_days": schema.ListAttribute{
						MarkdownDescription: `Array with the number of the days the task must run.

- For "weekly" schedule allowed values, 1 to 7.
- For "monthly" schedule allowed values, 1 to 31.`,
						ElementType: types.Int32Type,
						Required:    false,
						Optional:    true,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
						},
					},
					"cron_expression": schema.StringAttribute{
						Description: "Cron expression for the task. Only applies for for \"cron\" schedule.",
						Required:    false,
						Optional:    true,
					},
				},
			},
			"properties": schema.SingleNestedAttribute{
				Description: "Properties specific to this Task type",
				Required:    true,
				Optional:    false,
				Attributes:  tt.GetPropertiesSchema(),
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
