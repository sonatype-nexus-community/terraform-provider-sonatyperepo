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

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
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

// repositoryMavenHostedResource is the resource implementation.
type repositoryMavenHostedResource struct {
	common.BaseResource
}

// NewRepositoryMavenResource is a helper function to simplify the provider implementation.
func NewRepositoryMavenHostedResource() resource.Resource {
	return &repositoryMavenHostedResource{}
}

// Metadata returns the resource type name.
func (r *repositoryMavenHostedResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_maven_hosted"
}

// Schema defines the schema for the resource.
func (r *repositoryMavenHostedResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Hosted Maven Repositories",
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
				Description: fmt.Sprintf("Type of this Repository - will always be '%s'", REPOSITORY_TYPE_HOSTED),
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(REPOSITORY_TYPE_HOSTED),
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
					"write_policy": schema.StringAttribute{
						Description: "Controls if deployments of and updates to assets are allowed",
						Required:    true,
						Optional:    false,
						Validators: []validator.String{
							stringvalidator.OneOf("ALLOW", "ALLOW_ONCE", "DENY"),
						},
					},
				},
			},
			"cleanup": schema.SingleNestedAttribute{
				Description: "Repository Cleanup configuration",
				Required:    false,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"policy_names": schema.ListAttribute{
						Description: "Components that match any of the applied policies will be deleted",
						ElementType: types.StringType,
						Required:    false,
						Optional:    true,
					},
				},
			},
			"maven": schema.SingleNestedAttribute{
				Description: "Maven specific configuration for this Repository",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"version_policy": schema.StringAttribute{
						Description: "What type of artifacts does this repository store?",
						Required:    false,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("RELEASE", "SNAPSHOT", "MIXED"),
						},
					},
					"layout_policy": schema.StringAttribute{
						Description: "Validate that all paths are maven artifact or metadata paths",
						Required:    false,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("STRICT", "PERMISSIVE"),
						},
					},
					"content_disposition": schema.StringAttribute{
						Description: "Add Content-Disposition header as 'ATTACHMENT' to disable some content from being inline in a browser.",
						Required:    false,
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.OneOf("INLINE", "ATTACHMENT"),
						},
					},
				},
			},
			"component": schema.SingleNestedAttribute{
				Description: "Component configuration for this Repository",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"proprietary_components": schema.BoolAttribute{
						Description: "Components in this repository count as proprietary for namespace conflict attacks (requires Sonatype Nexus Firewall)",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
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
func (r *repositoryMavenHostedResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.RepositoryMavenHostedModel

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

	requestPayload := sonatyperepo.MavenHostedRepositoryApiRequest{
		Name:   plan.Name.ValueString(),
		Maven:  sonatyperepo.MavenAttributes{},
		Online: plan.Online.ValueBool(),
		Storage: sonatyperepo.HostedStorageAttributes{
			BlobStoreName:               plan.Storage.BlobStoreName.ValueString(),
			StrictContentTypeValidation: plan.Storage.StrictContentTypeValidation.ValueBool(),
			WritePolicy:                 plan.Storage.WritePolicy.ValueString(),
		},
	}
	if !plan.Maven.ContentDisposition.IsNull() {
		requestPayload.Maven.ContentDisposition = plan.Maven.ContentDisposition.ValueStringPointer()
	}
	if !plan.Maven.LayoutPolicy.IsNull() {
		requestPayload.Maven.LayoutPolicy = plan.Maven.LayoutPolicy.ValueStringPointer()
	}
	if !plan.Maven.VersionPolicy.IsNull() {
		requestPayload.Maven.VersionPolicy = plan.Maven.VersionPolicy.ValueStringPointer()
	}

	if plan.Cleanup != nil {
		if len(plan.Cleanup.PolicyNames) > 0 {
			policies := make([]string, len(plan.Cleanup.PolicyNames), 0)
			for _, p := range plan.Cleanup.PolicyNames {
				policies = append(policies, p.ValueString())
			}
			requestPayload.Cleanup = &sonatyperepo.CleanupPolicyAttributes{
				PolicyNames: policies,
			}
		}
	}

	if plan.Component != nil && !plan.Component.ProprietaryComponents.IsNull() {
		requestPayload.Component = &sonatyperepo.ComponentAttributes{
			ProprietaryComponents: plan.Component.ProprietaryComponents.ValueBoolPointer(),
		}
	}

	createRequest := r.Client.RepositoryManagementAPI.CreateMavenHostedRepository(ctx).Body(requestPayload)
	httpResponse, err := createRequest.Execute()

	// Handle Error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Maven Hosted Repository",
			fmt.Sprintf("Error creating Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	} else if httpResponse.StatusCode != http.StatusCreated && httpResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error creating Maven Hosted Repository",
			fmt.Sprintf("Unexpected Response Code whilst creating Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}

	// Crank in some defaults that whilst send in request, do not appear in response
	if plan.Format.IsNull() {
		plan.Format = types.StringValue(REPOSITORY_FORMAT_MAVEN)
	}
	if plan.Type.IsNull() {
		plan.Type = types.StringValue(REPOSITORY_TYPE_HOSTED)
	}
	if plan.Component == nil {
		plan.Component = &model.RepositoryComponentModel{
			ProprietaryComponents: types.BoolValue(false),
		}
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *repositoryMavenHostedResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.RepositoryMavenHostedModel

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
	repositoryApiResponse, httpResponse, err := r.Client.RepositoryManagementAPI.GetMavenHostedRepository(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Requested Maven Hosted Repository does not exist",
				fmt.Sprintf("Unable to read Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Maven Hosted Repository",
				fmt.Sprintf("Unable to read Maven Hosted Repository: %s: %s", httpResponse.Status, err),
			)
		}
		return
	} else {
		// Update State
		state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		state.Name = types.StringValue(*repositoryApiResponse.Name)
		state.Format = types.StringValue(*repositoryApiResponse.Format)
		state.Type = types.StringValue(*repositoryApiResponse.Type)
		state.Url = types.StringValue(*repositoryApiResponse.Url)
		state.Online = types.BoolValue(repositoryApiResponse.Online)
		state.Storage.BlobStoreName = types.StringValue(repositoryApiResponse.Storage.BlobStoreName)
		state.Storage.StrictContentTypeValidation = types.BoolValue(repositoryApiResponse.Storage.StrictContentTypeValidation)
		state.Storage.WritePolicy = types.StringValue(repositoryApiResponse.Storage.WritePolicy)
		if repositoryApiResponse.Cleanup != nil {
			policies := make([]types.String, len(repositoryApiResponse.Cleanup.PolicyNames), 0)
			for i, p := range repositoryApiResponse.Cleanup.PolicyNames {
				policies[i] = types.StringValue(p)
			}
			state.Cleanup = &model.RepositoryCleanupModel{
				PolicyNames: policies,
			}
		}
		state.Maven.ContentDisposition = types.StringValue(*repositoryApiResponse.Maven.ContentDisposition)
		state.Maven.LayoutPolicy = types.StringValue(*repositoryApiResponse.Maven.LayoutPolicy)
		state.Maven.VersionPolicy = types.StringValue(*repositoryApiResponse.Maven.VersionPolicy)
		if repositoryApiResponse.Component != nil && repositoryApiResponse.Component.ProprietaryComponents != nil {
			state.Component = &model.RepositoryComponentModel{
				ProprietaryComponents: types.BoolValue(*repositoryApiResponse.Component.ProprietaryComponents),
			}
		} else {
			state.Component = nil
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repositoryMavenHostedResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.RepositoryMavenHostedModel
	var state model.RepositoryMavenHostedModel

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
	requestPayload := sonatyperepo.MavenHostedRepositoryApiRequest{
		Name:   *plan.Name.ValueStringPointer(),
		Maven:  sonatyperepo.MavenAttributes{},
		Online: plan.Online.ValueBool(),
		Storage: sonatyperepo.HostedStorageAttributes{
			BlobStoreName:               plan.Storage.BlobStoreName.ValueString(),
			StrictContentTypeValidation: plan.Storage.StrictContentTypeValidation.ValueBool(),
			WritePolicy:                 plan.Storage.WritePolicy.ValueString(),
		},
	}
	if !plan.Maven.ContentDisposition.IsNull() {
		requestPayload.Maven.ContentDisposition = plan.Maven.ContentDisposition.ValueStringPointer()
	}
	if !plan.Maven.LayoutPolicy.IsNull() {
		requestPayload.Maven.LayoutPolicy = plan.Maven.LayoutPolicy.ValueStringPointer()
	}
	if !plan.Maven.VersionPolicy.IsNull() {
		requestPayload.Maven.VersionPolicy = plan.Maven.VersionPolicy.ValueStringPointer()
	}
	if len(plan.Cleanup.PolicyNames) > 0 {
		policies := make([]string, len(plan.Cleanup.PolicyNames), 0)
		for _, p := range plan.Cleanup.PolicyNames {
			policies = append(policies, p.ValueString())
		}
		requestPayload.Cleanup = &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: policies,
		}
	}
	if !plan.Component.ProprietaryComponents.IsNull() {
		requestPayload.Component = &sonatyperepo.ComponentAttributes{
			ProprietaryComponents: plan.Component.ProprietaryComponents.ValueBoolPointer(),
		}
	}
	apiUpdateRequest := r.Client.RepositoryManagementAPI.UpdateMavenHostedRepository(ctx, state.Name.ValueString()).Body(requestPayload)

	// Call API
	httpResponse, err := apiUpdateRequest.Execute()

	// Handle Error(s)
	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Maven Hosted Repository to update did not exist",
				fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Updating Maven Hosted Repository",
				fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
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
			"Unknown Error Updating Maven Hosted Repository",
			fmt.Sprintf("Unable to update Maven Hosted Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repositoryMavenHostedResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.RepositoryMavenHostedModel

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
