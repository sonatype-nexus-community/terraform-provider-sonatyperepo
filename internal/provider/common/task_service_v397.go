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
	"net/http"

	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

func taskFrequencyToApiModelV397(f TaskFrequencyApiModel) sonatyperepoV397.FrequencyXO {
	return sonatyperepoV397.FrequencyXO{
		Schedule:       f.Schedule,
		StartDate:      f.StartDate,
		TimeZoneOffset: f.TimeZoneOffset,
		RecurringDays:  f.RecurringDays,
		CronExpression: f.CronExpression,
	}
}

func taskApiModelFromV397(api *sonatyperepoV397.TaskXO) *TaskApiModel {
	if api == nil {
		return nil
	}
	m := &TaskApiModel{
		Id:                    api.Id,
		Name:                  api.Name,
		Type:                  api.Type,
		Enabled:               api.Enabled,
		AlertEmail:            api.AlertEmail,
		NotificationCondition: api.NotificationCondition,
		Schedule:              api.Schedule,
		TimeZoneOffset:        api.TimeZoneOffset,
		RecurringDays:         api.RecurringDays,
		CronExpression:        api.CronExpression,
		Properties:            api.Properties,
	}
	if api.StartDate != nil {
		unix := api.StartDate.Unix()
		m.StartDate = &unix
	}
	return m
}

// taskServiceV397 implements TaskService against NXRM API client V397 (targets NXRM 3.97.0+).
type taskServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

func (s *taskServiceV397) ListTasks(ctx context.Context) ([]TaskApiModel, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.TasksAPI.ListTasks(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	tasks := make([]TaskApiModel, 0, len(apiResponse.Items))
	for i := range apiResponse.Items {
		tasks = append(tasks, *taskApiModelFromV397(&apiResponse.Items[i]))
	}
	return tasks, httpResponse, nil
}

func (s *taskServiceV397) GetTaskById(ctx context.Context, id string) (*TaskApiModel, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.TasksAPI.GetTasks(ctx, id).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return taskApiModelFromV397(apiResponse), httpResponse, nil
}

func (s *taskServiceV397) CreateTask(ctx context.Context, plan *TaskCreateApiModel) (*TaskApiModel, *http.Response, error) {
	body := sonatyperepoV397.TaskTemplateXO{
		Name:                  plan.Name,
		Enabled:               plan.Enabled,
		Frequency:             taskFrequencyToApiModelV397(plan.Frequency),
		AlertEmail:            plan.AlertEmail,
		NotificationCondition: plan.NotificationCondition,
		Type:                  plan.Type,
		Properties:            plan.Properties,
	}
	apiResponse, httpResponse, err := s.client.TasksAPI.CreateTasks(ctx).TaskTemplateXO(body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return taskApiModelFromV397(apiResponse), httpResponse, nil
}

func (s *taskServiceV397) UpdateTask(ctx context.Context, id string, plan *TaskUpdateApiModel) (*http.Response, error) {
	body := sonatyperepoV397.UpdateTasksRequest{
		Name:                  plan.Name,
		Enabled:               plan.Enabled,
		Frequency:             taskFrequencyToApiModelV397(plan.Frequency),
		AlertEmail:            plan.AlertEmail,
		NotificationCondition: plan.NotificationCondition,
		Type:                  plan.Type,
		Properties:            plan.Properties,
	}
	return s.client.TasksAPI.UpdateTasks(ctx, id).UpdateTasksRequest(body).Execute()
}

func (s *taskServiceV397) DeleteTaskById(ctx context.Context, id string) (*http.Response, error) {
	return s.client.TasksAPI.DeleteTasks(ctx, id).Execute()
}
