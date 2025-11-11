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

package capabilitytype

import (
	"context"
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// --------------------------------------------
// Capabiltiy Type: UI: Settings
// --------------------------------------------
type UiSettingsCapability struct {
	BaseCapabilityType
}

func NewUiSettingsCapability() *UiSettingsCapability {
	return &UiSettingsCapability{
		BaseCapabilityType: BaseCapabilityType{
			capabilityType: common.CAPABILITY_TYPE_UI_SETTINGS,
			publicName:     "UI Settings",
		},
	}
}

// --------------------------------------------
// Capabiltiy Type: UI: Branding Functions
// --------------------------------------------
func (f *UiSettingsCapability) DoCreateRequest(plan any, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*v3.CapabilityDTO, *http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.UiSettingsCapabilityModel)

	// Call API to Create
	return apiClient.CapabilitiesAPI.Create3(ctx).Body(*planModel.ToApiCreateModel(version)).Execute()
}

func (f *UiSettingsCapability) DoUpdateRequest(plan any, capabilityId string, apiClient *v3.APIClient, ctx context.Context, version common.SystemVersion) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.UiSettingsCapabilityModel)

	// Call API to Update
	return apiClient.CapabilitiesAPI.Update3(ctx, capabilityId).Body(*planModel.ToApiUpdateModel(version)).Execute()
}

func (f *UiSettingsCapability) GetPlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.UiSettingsCapabilityModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *UiSettingsCapability) GetPropertiesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"title": schema.StringAttribute{
			Description: "Browser page title.",
			Optional:    true,
			Computed:    true,
			Default: stringdefault.StaticString(
				common.CAPABILITY_UI_SETTINGS_DEFAULT_TITLE,
			),
		},
		"debug_allowed": schema.BoolAttribute{
			Description: "Allow developer debugging.",
			Optional:    true,
			Computed:    true,
			Default: booldefault.StaticBool(
				common.CAPABILITY_UI_SETTINGS_DEFAULT_DEBUG_ALLOWED,
			),
		},
		"status_interval_authenticated": schema.Int32Attribute{
			Description: "Interval between status requests for authenticated users (seconds).",
			Optional:    true,
			Computed:    true,
			Default: int32default.StaticInt32(
				common.CAPABILITY_UI_SETTINGS_DEFAULT_STATUS_INTERVAL_AUTHENTICATED,
			),
		},
		"status_interval_anonymous": schema.Int32Attribute{
			Description: "Interval between status requests for anonymous user (seconds).",
			Optional:    true,
			Computed:    true,
			Default: int32default.StaticInt32(
				common.CAPABILITY_UI_SETTINGS_DEFAULT_STATUS_INTERVAL_ANONYMOUS,
			),
		},
		"session_timeout": schema.Int32Attribute{
			Description: "Period of inactivity before session times out (minutes). A value of 0 will mean that a session never expires.",
			Optional:    true,
			Computed:    true,
			Default: int32default.StaticInt32(
				common.CAPABILITY_UI_SETTINGS_DEFAULT_SESSION_TIMEOUT,
			),
		},
		"request_timeout": schema.Int32Attribute{
			Description: "Period of time to keep the connection alive for requests expected to take a normal period of time (seconds).",
			Optional:    true,
			Computed:    true,
			Default: int32default.StaticInt32(
				common.CAPABILITY_UI_SETTINGS_DEFAULT_REQUEST_TIMEOUT,
			),
		},
		"long_request_timeout": schema.Int32Attribute{
			Description: "Period of time to keep the connection alive for requests expected to take an extended period of time (seconds).",
			Optional:    true,
			Computed:    true,
			Default: int32default.StaticInt32(
				common.CAPABILITY_UI_SETTINGS_DEFAULT_LONG_REQUEST_TIMEOUT,
			),
		},
		// "search_request_timeout": schema.Int32Attribute{
		// 	Description: "Search request timeout in milliseconds.",
		// 	Optional:    true,
		// 	Computed:    true,
		// 	Default: int32default.StaticInt32(
		// 		common.CAPABILITY_UI_SETTINGS_DEFAULT_SEARCH_REQUEST_TIMEOUT,
		// 	),
		// },
	}
}

func (f *UiSettingsCapability) GetStateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.UiSettingsCapabilityModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *UiSettingsCapability) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.UiSettingsCapabilityModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *UiSettingsCapability) UpdateStateFromApi(state any, api any) any {
	stateModel := (state).(model.UiSettingsCapabilityModel)
	apiModel := (api).(*v3.CapabilityDTO)
	stateModel.FromApiModel(apiModel)
	stateModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return stateModel
}

func (f *UiSettingsCapability) UpdateStateFromPlanForUpdate(plan any, state any) any {
	planModel := (plan).(model.UiSettingsCapabilityModel)
	stateModel := (state).(model.UiSettingsCapabilityModel)

	planModel.Id = stateModel.Id
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	return planModel
}
