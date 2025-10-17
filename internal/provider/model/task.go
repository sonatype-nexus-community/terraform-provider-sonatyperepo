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
	"github.com/hashicorp/terraform-plugin-framework/types"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Task Frequency
// ----------------------------------------
type TaskFrequency struct {
	Schedule       types.String   `tfsdk:"schedule"`
	StartDate      *types.Int32   `tfsdk:"start_date"`
	TimezoneOffset *types.String  `tfsdk:"timezone_offset"`
	RecurringDays  *[]types.Int32 `tfsdk:"recurring_days"`
	CronExpression *types.String  `tfsdk:"cron_expression"`
}

func (f *TaskFrequency) ToApiModel(api *v3.FrequencyXO) {
	api.CronExpression = f.CronExpression.ValueStringPointer()
	api.RecurringDays = make([]int32, 0)
	for _, rd := range *f.RecurringDays {
		api.RecurringDays = append(api.RecurringDays, rd.ValueInt32())
	}
	api.Schedule = f.Schedule.ValueString()
	if f.StartDate.ValueInt32Pointer() != nil {
		val := int64(*f.StartDate.ValueInt32Pointer())
		api.StartDate = &val
	} else {
		api.StartDate = nil
	}
	api.TimeZoneOffset = f.TimezoneOffset.ValueStringPointer()
}

// Tasks Model
// ----------------------------------------
type TasksModel struct {
	Tasks []TaskModelSimple `tfsdk:"tasks"`
}

// Task Model (Simple)
// ----------------------------------------
type TaskModelSimple struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Type types.String `tfsdk:"type"`
}

func (m *TaskModelSimple) MapFromApi(api *sonatyperepo.TaskXO) {
	m.Id = types.StringPointerValue(api.Id)
	m.Name = types.StringPointerValue(api.Name)
	m.Type = types.StringPointerValue(api.Type)
}

// Base Task Model (Complete) - used for create and update
// ----------------------------------------
type BaseTaskModel struct {
	TaskModelSimple
	Enabled               types.Bool    `tfsdk:"enabled"`
	AlertEmail            *types.String `tfsdk:"alert_email"`
	NotificationCondition types.String  `tfsdk:"notification_condition"`
	Frequency             TaskFrequency `tfsdk:"frequency"`
}

func (m *BaseTaskModel) MapFromApi(api *sonatyperepo.TaskXO) {
	m.Id = types.StringPointerValue(api.Id)
	m.Name = types.StringPointerValue(api.Name)
	m.Type = types.StringPointerValue(api.Type)
}

func (m *BaseTaskModel) toApiCreateModel() *v3.TaskTemplateXO {
	api := v3.NewTaskTemplateXOWithDefaults()
	api.Name = m.Name.ValueString()
	api.Type = m.Type.ValueString()
	api.Enabled = m.Enabled.ValueBool()
	api.Frequency = *v3.NewFrequencyXO(m.Frequency.Schedule.ValueString())
	m.Frequency.ToApiModel(&api.Frequency)
	api.AlertEmail = m.AlertEmail.ValueStringPointer()
	api.NotificationCondition = m.NotificationCondition.ValueString()
	return api
}

// Base Task Properties
// ----------------------------------------
type BaseTaskProperties struct{}

func (p *BaseTaskProperties) AsMap() *map[string]string {
	return StructToMap(p)
}
