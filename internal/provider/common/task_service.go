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

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// TaskFrequencyApiModel is a generation-agnostic representation of a Task's Frequency,
// shared by every NXRM API client generation.
type TaskFrequencyApiModel struct {
	Schedule       string
	StartDate      *int64
	TimeZoneOffset *string
	RecurringDays  []int32
	CronExpression *string
}

// TaskApiModel is a generation-agnostic representation of a Task, as returned by the
// Tasks API. TaskService implementations translate the raw, per-generation client type
// into this shape so that callers outside the adapter never depend on a specific
// client generation.
type TaskApiModel struct {
	Id                    *string
	Name                  *string
	Type                  *string
	Enabled               *bool
	AlertEmail            *string
	NotificationCondition *string
	Schedule              *string
	StartDate             *int64
	TimeZoneOffset        *string
	RecurringDays         []int32
	CronExpression        *string
	Properties            *map[string]string
}

// TaskCreateApiModel is a generation-agnostic representation of the request body used to
// create a Task.
type TaskCreateApiModel struct {
	Name                  string
	Enabled               bool
	Frequency             TaskFrequencyApiModel
	AlertEmail            *string
	NotificationCondition string
	Type                  string
	Properties            *map[string]string
}

// TaskUpdateApiModel is a generation-agnostic representation of the request body used to
// update a Task.
type TaskUpdateApiModel struct {
	Name                  string
	Enabled               bool
	Frequency             TaskFrequencyApiModel
	AlertEmail            *string
	NotificationCondition string
	Type                  *string
	Properties            *map[string]string
}

// TaskService abstracts the Tasks API across NXRM API client generations.
type TaskService interface {
	ListTasks(ctx context.Context) ([]TaskApiModel, *http.Response, error)
	GetTaskById(ctx context.Context, id string) (*TaskApiModel, *http.Response, error)
	CreateTask(ctx context.Context, plan *TaskCreateApiModel) (*TaskApiModel, *http.Response, error)
	UpdateTask(ctx context.Context, id string, plan *TaskUpdateApiModel) (*http.Response, error)
	DeleteTaskById(ctx context.Context, id string) (*http.Response, error)
}

func taskFrequencyToApiModelV382(f TaskFrequencyApiModel) sonatyperepoV382.FrequencyXO {
	return sonatyperepoV382.FrequencyXO{
		Schedule:       f.Schedule,
		StartDate:      f.StartDate,
		TimeZoneOffset: f.TimeZoneOffset,
		RecurringDays:  f.RecurringDays,
		CronExpression: f.CronExpression,
	}
}

func taskFrequencyToApiModelV395(f TaskFrequencyApiModel) sonatyperepoV395.FrequencyXO {
	return sonatyperepoV395.FrequencyXO{
		Schedule:       f.Schedule,
		StartDate:      f.StartDate,
		TimeZoneOffset: f.TimeZoneOffset,
		RecurringDays:  f.RecurringDays,
		CronExpression: f.CronExpression,
	}
}

func taskApiModelFromV382(api *sonatyperepoV382.TaskXO) *TaskApiModel {
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

func taskApiModelFromV395(api *sonatyperepoV395.TaskXO) *TaskApiModel {
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

// taskServiceV382 implements TaskService against NXRM API client V382 (targets NXRM < 3.94.0).
type taskServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

func (s *taskServiceV382) ListTasks(ctx context.Context) ([]TaskApiModel, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.TasksAPI.GetTasks(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	tasks := make([]TaskApiModel, 0, len(apiResponse.Items))
	for i := range apiResponse.Items {
		tasks = append(tasks, *taskApiModelFromV382(&apiResponse.Items[i]))
	}
	return tasks, httpResponse, nil
}

func (s *taskServiceV382) GetTaskById(ctx context.Context, id string) (*TaskApiModel, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.TasksAPI.GetTaskById(ctx, id).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return taskApiModelFromV382(apiResponse), httpResponse, nil
}

func (s *taskServiceV382) CreateTask(ctx context.Context, plan *TaskCreateApiModel) (*TaskApiModel, *http.Response, error) {
	body := sonatyperepoV382.TaskTemplateXO{
		Name:                  plan.Name,
		Enabled:               plan.Enabled,
		Frequency:             taskFrequencyToApiModelV382(plan.Frequency),
		AlertEmail:            plan.AlertEmail,
		NotificationCondition: plan.NotificationCondition,
		Type:                  plan.Type,
		Properties:            plan.Properties,
	}
	apiResponse, httpResponse, err := s.client.TasksAPI.CreateTask(ctx).Body(body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return taskApiModelFromV382(apiResponse), httpResponse, nil
}

func (s *taskServiceV382) UpdateTask(ctx context.Context, id string, plan *TaskUpdateApiModel) (*http.Response, error) {
	body := sonatyperepoV382.UpdateTaskRequest{
		Name:                  plan.Name,
		Enabled:               plan.Enabled,
		Frequency:             taskFrequencyToApiModelV382(plan.Frequency),
		AlertEmail:            plan.AlertEmail,
		NotificationCondition: plan.NotificationCondition,
		Type:                  plan.Type,
		Properties:            plan.Properties,
	}
	return s.client.TasksAPI.UpdateTask(ctx, id).Body(body).Execute()
}

func (s *taskServiceV382) DeleteTaskById(ctx context.Context, id string) (*http.Response, error) {
	return s.client.TasksAPI.DeleteTaskById(ctx, id).Execute()
}

// taskServiceV395 implements TaskService against NXRM API client V395 (targets NXRM 3.94.0+).
type taskServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

func (s *taskServiceV395) ListTasks(ctx context.Context) ([]TaskApiModel, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.TasksAPI.ListTasks(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	tasks := make([]TaskApiModel, 0, len(apiResponse.Items))
	for i := range apiResponse.Items {
		tasks = append(tasks, *taskApiModelFromV395(&apiResponse.Items[i]))
	}
	return tasks, httpResponse, nil
}

func (s *taskServiceV395) GetTaskById(ctx context.Context, id string) (*TaskApiModel, *http.Response, error) {
	apiResponse, httpResponse, err := s.client.TasksAPI.GetTasks(ctx, id).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return taskApiModelFromV395(apiResponse), httpResponse, nil
}

func (s *taskServiceV395) CreateTask(ctx context.Context, plan *TaskCreateApiModel) (*TaskApiModel, *http.Response, error) {
	body := sonatyperepoV395.TaskTemplateXO{
		Name:                  plan.Name,
		Enabled:               plan.Enabled,
		Frequency:             taskFrequencyToApiModelV395(plan.Frequency),
		AlertEmail:            plan.AlertEmail,
		NotificationCondition: plan.NotificationCondition,
		Type:                  plan.Type,
		Properties:            plan.Properties,
	}
	apiResponse, httpResponse, err := s.client.TasksAPI.CreateTasks(ctx).TaskTemplateXO(body).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	return taskApiModelFromV395(apiResponse), httpResponse, nil
}

func (s *taskServiceV395) UpdateTask(ctx context.Context, id string, plan *TaskUpdateApiModel) (*http.Response, error) {
	body := sonatyperepoV395.UpdateTasksRequest{
		Name:                  plan.Name,
		Enabled:               plan.Enabled,
		Frequency:             taskFrequencyToApiModelV395(plan.Frequency),
		AlertEmail:            plan.AlertEmail,
		NotificationCondition: plan.NotificationCondition,
		Type:                  plan.Type,
		Properties:            plan.Properties,
	}
	return s.client.TasksAPI.UpdateTasks(ctx, id).UpdateTasksRequest(body).Execute()
}

func (s *taskServiceV395) DeleteTaskById(ctx context.Context, id string) (*http.Response, error) {
	return s.client.TasksAPI.DeleteTasks(ctx, id).Execute()
}
