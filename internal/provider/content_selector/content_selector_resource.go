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

package content_selector

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	tfschema "github.com/sonatype-nexus-community/terraform-provider-shared/schema"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
)

const contentSelectorNamePattern = `^[a-zA-Z0-9\-]{1}[a-zA-Z0-9_\-\.]*$`

// contentSelectorResource is the resource implementation.
type contentSelectorResource struct {
	common.BaseResource
}

// NewContentSelectorResource is a helper function to simplify the provider implementation.
func NewContentSelectorResource() resource.Resource {
	return &contentSelectorResource{}
}

// Metadata returns the resource type name.
func (r *contentSelectorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_content_selector"
}

// Schema defines the schema for the resource.
func (r *contentSelectorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Content Selectors in Sonatype Nexus Repository",
		Attributes: map[string]schema.Attribute{
			"name": tfschema.RequiredStringWithRegex(
				"The name of the Content Selector.",
				regexp.MustCompile(contentSelectorNamePattern),
				"Content Selector name must match pattern `^[a-zA-Z0-9\\-]{1}[a-zA-Z0-9_\\-\\.]*$`",
			),
			"description": tfschema.RequiredString("The description of this Content Selector."),
			"expression":  tfschema.RequiredString("The Content Selector expression used to identify content."),
			"last_updated": tfschema.Timestamp(),
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *contentSelectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan model.ContentSelectorModelResource
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	ctx = r.GetAuthContext(ctx)
	apiBody := sonatyperepo.NewContentSelectorApiCreateRequest()
	plan.MapToApiCreate(apiBody)
	httpResponse, err := r.Client.ContentSelectorsAPI.CreateContentSelector(ctx).Body(*apiBody).Execute()

	if err != nil {
		sharederr.HandleAPIError(
			"Error creating Content Selector",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		sharederr.HandleAPIError(
			"Creation of Content Selector was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	// Update State
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *contentSelectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.ContentSelectorModelResource
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.GetAuthContext(ctx)

	// Read API Call
	apiResponse, httpResponse, err := r.Client.ContentSelectorsAPI.GetContentSelector(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			sharederr.HandleAPIWarning(
				"Content Selector to read did not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			sharederr.HandleAPIError(
				"Error reading Content Selector",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Update State based on Response
	state.MapFromApi(apiResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *contentSelectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.ContentSelectorModelResource
	var state model.ContentSelectorModelResource

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Update
	ctx = r.GetAuthContext(ctx)
	apiBody := sonatyperepo.NewContentSelectorApiUpdateRequest()
	plan.MapToApiUpdate(apiBody)
	httpResponse, err := r.Client.ContentSelectorsAPI.UpdateContentSelector(ctx, state.Name.ValueString()).Body(*apiBody).Execute()

	if err != nil {
		sharederr.HandleAPIError(
			"Error updating Content Selector",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		sharederr.HandleAPIError(
			"Update of Content Selector was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *contentSelectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.ContentSelectorModelResource

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = r.GetAuthContext(ctx)

	httpResponse, err := r.Client.ContentSelectorsAPI.DeleteContentSelector(ctx, state.Name.ValueString()).Execute()

	// Handle Error
	if err != nil {
		sharederr.HandleAPIError(
			"Error removing Content Selector",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		sharederr.HandleAPIError(
			"Removal of Content Selector was not successful",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}
}
