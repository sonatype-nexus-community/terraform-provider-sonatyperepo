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

package repository

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go"
)

// repositoryMavenGroupResource is the resource implementation.
type repositoryMavenGroupResource struct {
	common.BaseResource
}

// NewRepositoryMavenGroupResource is a helper function to simplify the provider implementation.
func NewRepositoryMavenGroupResource() resource.Resource {
	return &repositoryMavenGroupResource{}
}

// Metadata returns the resource type name.
func (r *repositoryMavenGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_maven_group"
}

// Schema defines the schema for the resource.
func (r *repositoryMavenGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Group Maven Repositories",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Name of the Repository",
				Required:    true,
			},
			"url": schema.StringAttribute{
				Description: "URL to access the Repository",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"format": schema.StringAttribute{
				Description: fmt.Sprintf("Format of this Repository - will always be '%s'", REPOSITORY_FORMAT_MAVEN),
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(REPOSITORY_FORMAT_MAVEN),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: fmt.Sprintf("Type of this Repository - will always be '%s'", REPOSITORY_TYPE_PROXY),
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(REPOSITORY_TYPE_PROXY),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"online": schema.BoolAttribute{
				Description: "Whether this Repository is online and accepting incoming requests",
				Required:    true,
			},
			"storage": schema.SingleNestedAttribute{
				Description: "Storage configuration for this Repository",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"blob_store_name": schema.StringAttribute{
						Description: "Name of the Blob Store to use",
						Required:    true,
						Optional:    false,
					},
					"strict_content_type_validation": schema.BoolAttribute{
						Description: "Whether this Repository validates that all content uploaded to this repository is of a MIME type appropriate for the repository format",
						Required:    true,
					},
				},
			},
			"group": schema.SingleNestedAttribute{
				Description: "Repository Group configuration",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"member_names": schema.ListAttribute{
						Description: "Member repositories' names",
						ElementType: types.StringType,
						Required:    false,
						Optional:    true,
						Validators: []validator.List{
							listvalidator.SizeAtLeast(1),
							listvalidator.UniqueValues(),
							listvalidator.IsRequired(),
						},
					},
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *repositoryMavenGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.RepositoryMavenGroupModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Call API to Create
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	requestPayload := r.makeApiRequest(&plan)
	createRequest := r.Client.RepositoryManagementAPI.CreateMavenGroupRepository(ctx).Body(requestPayload)
	httpResponse, err := createRequest.Execute()

	// Handle Error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Maven Group Repository",
			fmt.Sprintf("Error creating Maven Group Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	} else if httpResponse.StatusCode != http.StatusCreated && httpResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error creating Maven Group Repository",
			fmt.Sprintf("Unexpected Response Code whilst creating Maven Group Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}

	// Crank in some defaults that whilst send in request, do not appear in response
	if plan.Format.IsNull() {
		plan.Format = types.StringValue(REPOSITORY_FORMAT_MAVEN)
	}
	if plan.Type.IsNull() {
		plan.Type = types.StringValue(REPOSITORY_TYPE_PROXY)
	}
	// E.g. http://localhost:8081/repository/maven-proxy-repo-test - this is not included in response to CREATE
	plan.Url = types.StringValue(fmt.Sprintf("%s/repository/%s", r.BaseUrl, plan.Name.ValueString()))

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *repositoryMavenGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.RepositoryMavenGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Read API Call
	repositoryApiResponse, httpResponse, err := r.Client.RepositoryManagementAPI.GetMavenGroupRepository(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Requested Maven Group Repository does not exist",
				fmt.Sprintf("Unable to read Maven Group Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Maven Group Repository",
				fmt.Sprintf("Unable to read Maven Group Repository: %s: %s", httpResponse.Status, err),
			)
		}
		return
	} else {
		// Update State
		state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		state.Name = types.StringValue(*repositoryApiResponse.Name)
		state.Format = types.StringValue(REPOSITORY_FORMAT_MAVEN)
		state.Type = types.StringValue(REPOSITORY_TYPE_PROXY)
		state.Url = types.StringValue(*repositoryApiResponse.Url)
		state.Online = types.BoolValue(repositoryApiResponse.Online)
		state.Storage.BlobStoreName = types.StringValue(repositoryApiResponse.Storage.BlobStoreName)
		state.Storage.StrictContentTypeValidation = types.BoolValue(repositoryApiResponse.Storage.StrictContentTypeValidation)
		if len(repositoryApiResponse.Group.MemberNames) > 0 {
			groupMemberNames := make([]types.String, 0, len(repositoryApiResponse.Group.MemberNames))
			for _, p := range repositoryApiResponse.Group.MemberNames {
				groupMemberNames = append(groupMemberNames, types.StringValue(p))
			}
			state.Group = model.RepositoryGroupModel{
				MemberNames: groupMemberNames,
			}
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repositoryMavenGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.RepositoryMavenGroupModel
	var state model.RepositoryMavenGroupModel

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

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Update API Call
	requestPayload := r.makeApiRequest(&plan)
	apiUpdateRequest := r.Client.RepositoryManagementAPI.UpdateMavenGroupRepository(ctx, state.Name.ValueString()).Body(requestPayload)

	// Call API
	httpResponse, err := apiUpdateRequest.Execute()

	// Handle Error(s)
	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Maven Proxy Repository to update did not exist",
				fmt.Sprintf("Unable to update Maven Proxy Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Updating Maven Proxy Repository",
				fmt.Sprintf("Unable to update Maven Proxy Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		}
		return
	} else if httpResponse.StatusCode == http.StatusNoContent {
		// Map response body to schema and populate Computed attribute values
		plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		// Set state to fully populated data
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError(
			"Unknown Error Updating Maven Proxy Repository",
			fmt.Sprintf("Unable to update Maven Proxy Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repositoryMavenGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.RepositoryMavenGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting state data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	DeleteRepository(r.Client, &ctx, state.Name.ValueString(), resp)
}
func (r *repositoryMavenGroupResource) makeApiRequest(plan *model.RepositoryMavenGroupModel) sonatyperepo.MavenGroupRepositoryApiRequest {
	requestPayload := sonatyperepo.MavenGroupRepositoryApiRequest{
		Name:   plan.Name.ValueString(),
		Online: plan.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{
			BlobStoreName:               plan.Storage.BlobStoreName.ValueString(),
			StrictContentTypeValidation: plan.Storage.StrictContentTypeValidation.ValueBool(),
		},
	}
	if len(plan.Group.MemberNames) > 0 {
		groupMembers := make([]string, 0, len(plan.Group.MemberNames))
		for _, p := range plan.Group.MemberNames {
			groupMembers = append(groupMembers, p.ValueString())
		}
		requestPayload.Group = sonatyperepo.GroupAttributes{
			MemberNames: groupMembers,
		}
	}

	return requestPayload
}
