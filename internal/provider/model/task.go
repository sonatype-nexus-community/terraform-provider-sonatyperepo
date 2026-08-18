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
	"terraform-provider-sonatyperepo/internal/provider/common"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Task Frequency
// ----------------------------------------
type taskFrequency struct {
	Schedule       types.String  `tfsdk:"schedule"`
	StartDate      types.Int64   `tfsdk:"start_date"`
	TimezoneOffset types.String  `tfsdk:"timezone_offset"`
	RecurringDays  []types.Int32 `tfsdk:"recurring_days"`
	CronExpression types.String  `tfsdk:"cron_expression"`
}

func (f *taskFrequency) ToApiModel(api *common.TaskFrequencyApiModel) {
	api.CronExpression = f.CronExpression.ValueStringPointer()
	api.RecurringDays = make([]int32, 0)
	for _, rd := range f.RecurringDays {
		api.RecurringDays = append(api.RecurringDays, rd.ValueInt32())
	}
	api.Schedule = f.Schedule.ValueString()
	api.StartDate = f.StartDate.ValueInt64Pointer()
	api.TimeZoneOffset = f.TimezoneOffset.ValueStringPointer()
}

// Tasks Model
// ----------------------------------------
type TasksModel struct {
	Tasks []TaskModelSimple `tfsdk:"tasks"`
}

// Task Model Base
// ----------------------------------------
type baseTaskModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

// Task Model (Simple)
// ----------------------------------------
type TaskModelSimple struct {
	baseTaskModel
	Type types.String `tfsdk:"type"`
}

func (m *TaskModelSimple) MapFromApi(api *common.TaskApiModel) {
	m.Id = types.StringPointerValue(api.Id)
	m.Name = types.StringPointerValue(api.Name)
	m.Type = types.StringPointerValue(api.Type)
}

// Base Task Model (Complete) - used for create and update
// ----------------------------------------
type BaseTaskModel struct {
	baseTaskModel
	Enabled               types.Bool     `tfsdk:"enabled"`
	AlertEmail            types.String   `tfsdk:"alert_email"`
	NotificationCondition types.String   `tfsdk:"notification_condition"`
	Frequency             *taskFrequency `tfsdk:"frequency"`
	LastUpdated           types.String   `tfsdk:"last_updated"`
}

func (m *BaseTaskModel) MapFromApi(api *common.TaskApiModel) {
	m.Id = types.StringPointerValue(api.Id)
	if api.Name != nil {
		m.Name = types.StringPointerValue(api.Name)
	}
}

func (m *BaseTaskModel) toApiCreateModel() *common.TaskCreateApiModel {
	api := &common.TaskCreateApiModel{}
	api.Name = m.Name.ValueString()
	api.Enabled = m.Enabled.ValueBool()
	if m.Frequency != nil {
		api.Frequency = common.TaskFrequencyApiModel{Schedule: m.Frequency.Schedule.ValueString()}
		m.Frequency.ToApiModel(&api.Frequency)
	}
	api.AlertEmail = m.AlertEmail.ValueStringPointer()
	api.NotificationCondition = m.NotificationCondition.ValueString()
	return api
}

func (m *BaseTaskModel) toApiUpdateModel() *common.TaskUpdateApiModel {
	api := &common.TaskUpdateApiModel{}
	api.Name = m.Name.ValueString()
	api.Enabled = m.Enabled.ValueBool()
	if m.Frequency != nil {
		api.Frequency = common.TaskFrequencyApiModel{Schedule: m.Frequency.Schedule.ValueString()}
		m.Frequency.ToApiModel(&api.Frequency)
	}
	api.AlertEmail = m.AlertEmail.ValueStringPointer()
	api.NotificationCondition = m.NotificationCondition.ValueString()
	return api
}
