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
)

type TaskFrequency struct {
	Schedule       types.String   `tfsdk:"schedule"`
	StartDate      *types.Int32   `tfsdk:"start_date"`
	TimezoneOffset *types.String  `tfsdk:"timezone_offset"`
	RecurringDays  *[]types.Int32 `tfsdk:"recurring_days"`
	CronExpression *types.String  `tfsdk:"cron_expression"`
}

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

type TaskModel struct {
	TaskModelSimple
	Enabled               types.Bool    `tfsdk:"enabled"`
	AlertEmail            *types.String `tfsdk:"alert_email"`
	NotificationCondition types.String  `tfsdk:"notification_condition"`
	Frequency             TaskFrequency `tfsdk:"frequency"`
}

func (m *TaskModel) MapFromApi(api *sonatyperepo.TaskXO) {
	m.Id = types.StringPointerValue(api.Id)
	m.Name = types.StringPointerValue(api.Name)
	m.Type = types.StringPointerValue(api.Type)

}

// TasksModel
type TasksModel struct {
	Tasks []TaskModelSimple `tfsdk:"tasks"`
}
