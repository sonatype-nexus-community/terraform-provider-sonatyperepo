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
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
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
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, t.TaskType.ResourceName())
}

// Set Schema for this Resource
func (t *taskResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = taskSchema(t.TaskType)
}

// This allows users to import existing Tasks into Terraform state.
func (r *taskResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the Task ID as the import identifier
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (t *taskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan, diags := t.TaskType.PlanAsModel(ctx, req.Plan)
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
		errors.HandleAPIError(
			fmt.Sprintf("Error creating %s Task", t.TaskType.Type().String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}
	if !slices.Contains(t.TaskType.ApiCreateSuccessResponseCodes(), httpResponse.StatusCode) {
		errors.HandleAPIError(
			fmt.Sprintf("Creation of %s Task was not successful", t.TaskType.Type().String()),
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
	planModel, diags := t.TaskType.PlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	// Retrieve values from state
	stateModel, diags := t.TaskType.StateAsModel(ctx, req.State)
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
			errors.HandleAPIWarning(
				fmt.Sprintf(TASK_ERROR_DID_NOT_EXIST, t.TaskType.Type().String(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
				fmt.Sprintf(TASK_ERROR_DID_NOT_EXIST, t.TaskType.Type().String(), "update"),
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
	state, diags := t.TaskType.StateAsModel(ctx, req.State)
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
			tflog.Info(ctx, fmt.Sprintf("Unexpected response when deleting Task %s (attempt %d)", t.TaskType.Type().String(), attempts))
			attempts++
			continue
		}

		if err != nil {
			if httpResponse.StatusCode == http.StatusNotFound {
				resp.State.RemoveResource(ctx)
				errors.HandleAPIWarning(
					fmt.Sprintf(TASK_ERROR_DID_NOT_EXIST, t.TaskType.Type().String(), "delete"),
					&err,
					httpResponse,
					&resp.Diagnostics,
				)
			} else {
				errors.HandleAPIError(
					fmt.Sprintf(TASK_ERROR_DID_NOT_EXIST, t.TaskType.Type().String(), "delete"),
					&err,
					httpResponse,
					&resp.Diagnostics,
				)
			}
			return
		}
		if httpResponse.StatusCode != http.StatusNoContent {
			errors.HandleAPIError(
				fmt.Sprintf("Unexpected response when deleting %s Task (attempt %d)", t.TaskType.Type().String(), attempts),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)

			time.Sleep(1 * time.Second)
			attempts++
		} else {
			success = true
		}
	}
}

func taskSchema(tt tasktype.TaskTypeI) tfschema.Schema {
	attributes := map[string]tfschema.Attribute{
		"id":          schema.ResourceComputedString("The internal ID of the Task."),
		"name":        schema.ResourceRequiredString("The name of the Task."),
		"enabled":     schema.ResourceRequiredBool("Indicates if the task is enabled."),
		"alert_email": schema.ResourceOptionalString("E-mail address for task notifications."),
		"notification_condition": schema.ResourceRequiredStringEnum(
			"The type of Task.",
			common.NOTIFICATION_CONDITION_FAILURE,
			common.NOTIFICATION_CONDITION_SUCCESS_OR_FAILURE,
		),
		"frequency": schema.ResourceRequiredSingleNestedAttribute("Frequency Schedule for this Task.",
			map[string]tfschema.Attribute{
				"schedule": schema.ResourceRequiredStringEnum(
					"Type of Schedule.",
					common.FREQUENCY_SCHEDULE_MANUAL,
					common.FREQUENCY_SCHEDULE_ONCE,
					common.FREQUENCY_SCHEDULE_HOURLY,
					common.FREQUENCY_SCHEDULE_DAILY,
					common.FREQUENCY_SCHEDULE_WEEKLY,
					common.FREQUENCY_SCHEDULE_MONTHLY,
					common.FREQUENCY_SCHEDULE_CRON,
				),
				"start_date": schema.ResourceOptionalInt32(
					"Start date of the task represented in unix timestamp. Does not apply for \"manual\" schedule.",
				),
				"timezone_offset": schema.ResourceOptionalString("The offset time zone of the client. Example: -05:00"),
				"recurring_days": func() tfschema.ListAttribute {
					thisAttr := schema.ResourceOptionalInt32List(
						`Array with the number of the days the task must run.

- For "weekly" schedule allowed values, 1 to 7.
- For "monthly" schedule allowed values, 1 to 31.`,
					)
					thisAttr.Validators = []validator.List{
						listvalidator.SizeAtLeast(1),
					}
					return thisAttr
				}(),
				"cron_expression": schema.ResourceOptionalString("Cron expression for the task. Only applies for for \"cron\" schedule."),
			},
		),
		"last_updated": schema.ResourceLastUpdated(),
	}

	propertiesSchema := tt.PropertiesSchema()
	if len(propertiesSchema) > 0 {
		attributes["properties"] = schema.ResourceRequiredSingleNestedAttribute("Properties specific to this Task type", propertiesSchema)
	}

	return tfschema.Schema{
		MarkdownDescription: tt.MarkdownDescription(),
		Attributes:          attributes,
	}
}
