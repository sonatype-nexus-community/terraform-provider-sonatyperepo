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
	"maps"
	"net/http"
	"reflect"
	"slices"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/repository/format"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	REPOSITORY_ERROR_RESPONSE_PREFIX           = "Error response: "
	REPOSITORY_GENERAL_ERROR_RESPONSE_GENERAL  = REPOSITORY_ERROR_RESPONSE_PREFIX + " %s"
	REPOSITORY_GENERAL_ERROR_RESPONSE_WITH_ERR = REPOSITORY_ERROR_RESPONSE_PREFIX + " %s - %s"
	REPOSITORY_ERROR_DID_NOT_EXIST             = "%s %s Repository did not exist to %s"
)

// Generic to all Repository Resources
type repositoryResource struct {
	common.BaseResource
	RepositoryFormat format.RepositoryFormat
	RepositoryType   format.RepositoryType
}

// Metadata returns the resource type name.
func (r *repositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_%s", req.ProviderTypeName, r.RepositoryFormat.GetResourceName(r.RepositoryType))
}

// Set Schema for this Resource
func (r *repositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := getHostedStandardSchema(r.RepositoryFormat.GetKey(), r.RepositoryType)
	maps.Copy(schema.Attributes, r.RepositoryFormat.GetFormatSchemaAttributes())
	resp.Schema = schema
}

func (r *repositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	plan, diags := r.RepositoryFormat.GetPlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting Plan data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Make API requet
	httpResponse, err := r.RepositoryFormat.DoCreateRequest(plan, r.Client, ctx)

	// Handle Errors
	if err != nil {
		common.HandleApiError(
			fmt.Sprintf("Error creating %s %s Repository", r.RepositoryFormat.GetKey(), r.RepositoryType.String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}
	if !slices.Contains(r.RepositoryFormat.GetApiCreateSuccessResposneCodes(), httpResponse.StatusCode) {
		common.HandleApiError(
			fmt.Sprintf("Creation of %s %s Repository was not successful", r.RepositoryFormat.GetKey(), r.RepositoryType.String()),
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
	}

	// Call Read API as that contains more complete information for mapping to State
	apiResponse, httpResponse, err := r.RepositoryFormat.DoReadRequest(plan, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			common.HandleApiWarning(
				fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			common.HandleApiError(
				fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	stateModel := r.RepositoryFormat.UpdateStateFromApi(plan, apiResponse)
	stateModel = r.RepositoryFormat.UpdatePlanForState(stateModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	stateModel, diags := r.RepositoryFormat.GetStateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Handle any errors
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Getting request data has errors: %v", resp.Diagnostics.Errors()))
		return
	}

	// Set API Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Make API Request
	apiResponse, httpResponse, err := r.RepositoryFormat.DoReadRequest(stateModel, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			common.HandleApiWarning(
				fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			common.HandleApiError(
				fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Update State from Response
	r.RepositoryFormat.UpdateStateFromApi(stateModel, apiResponse)
	resp.Diagnostics.Append(resp.State.Set(ctx, &stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	planModel, diags := r.RepositoryFormat.GetPlanAsModel(ctx, req.Plan)
	resp.Diagnostics.Append(diags...)

	// Retrieve values from state
	stateModel, diags := r.RepositoryFormat.GetStateAsModel(ctx, req.State)
	resp.Diagnostics.Append(diags...)

	// Request Context
	ctx = context.WithValue(
		ctx,
		sonatyperepo.ContextBasicAuth,
		r.Auth,
	)

	// Make API requet
	httpResponse, err := r.RepositoryFormat.DoUpdateRequest(planModel, stateModel, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			common.HandleApiWarning(
				fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			common.HandleApiError(
				fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "update"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	// Call Read API as that contains more complete information for mapping to State
	apiResponse, httpResponse, err := r.RepositoryFormat.DoReadRequest(planModel, r.Client, ctx)

	// Handle any errors
	if err != nil {
		if httpResponse.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			common.HandleApiWarning(
				fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			common.HandleApiError(
				fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "read"),
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	}

	stateModel = r.RepositoryFormat.UpdateStateFromApi(planModel, apiResponse)
	// stateModel = (r.RepositoryFormat.UpdatePlanForState(stateModel)).(model.RepositoryNpmHostedModel)
	stateModel = r.RepositoryFormat.UpdatePlanForState(stateModel)
	resp.Diagnostics.Append(resp.State.Set(ctx, stateModel)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *repositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	state, diags := r.RepositoryFormat.GetStateAsModel(ctx, req.State)
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
		r.Auth,
	)

	// Make API request
	repoNameStructField := reflect.Indirect(reflect.ValueOf(state)).FieldByName("Name").Interface()
	repositoryName, ok := repoNameStructField.(basetypes.StringValue)
	if !ok {
		resp.Diagnostics.AddError(
			"Failed to determine Repository Name to delete from State",
			fmt.Sprintf("%s %s", REPOSITORY_ERROR_RESPONSE_PREFIX, repoNameStructField),
		)
		return
	}

	attempts := 1
	maxAttempts := 3
	success := false

	for !success && attempts < maxAttempts {
		httpResponse, err := r.RepositoryFormat.DoDeleteRequest(repositoryName.ValueString(), r.Client, ctx)

		// Trap 500 Error as they occur when Repo is not in appropriate internal state
		if httpResponse.StatusCode == http.StatusInternalServerError {
			tflog.Info(ctx, fmt.Sprintf("Unexpected response when deleting %s %s Repository (attempt %d)", r.RepositoryFormat.GetKey(), r.RepositoryFormat, attempts))
			attempts++
			continue
		}

		if err != nil {
			if httpResponse.StatusCode == http.StatusNotFound {
				resp.State.RemoveResource(ctx)
				resp.Diagnostics.AddWarning(
					fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryType.String(), r.RepositoryFormat.GetKey(), "delete"),
					fmt.Sprintf(REPOSITORY_GENERAL_ERROR_RESPONSE_GENERAL, httpResponse.Status),
				)
			} else {
				resp.Diagnostics.AddError(
					fmt.Sprintf(REPOSITORY_ERROR_DID_NOT_EXIST, r.RepositoryFormat.GetKey(), r.RepositoryFormat, "delete"),
					fmt.Sprintf(REPOSITORY_GENERAL_ERROR_RESPONSE_WITH_ERR, httpResponse.Status, err),
				)
			}
			return
		}
		if httpResponse.StatusCode != http.StatusNoContent {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unexpected response when deleting %s %s Repository (attempt %d)", r.RepositoryFormat.GetKey(), r.RepositoryFormat, attempts),
				fmt.Sprintf("Error response: %s", httpResponse.Status),
			)

			time.Sleep(1 * time.Second)
			attempts++
		} else {
			success = true
		}
	}
}

func getHostedStandardSchema(repoFormat string, repoType format.RepositoryType) schema.Schema {
	storageAttributes := map[string]schema.Attribute{
		"blob_store_name": schema.StringAttribute{
			Description: "Name of the Blob Store to use",
			Required:    true,
			Optional:    false,
		},
		"strict_content_type_validation": schema.BoolAttribute{
			Description: "Whether this Repository validates that all content uploaded to this repository is of a MIME type appropriate for the repository format",
			Required:    true,
		},
	}

	// Write Policy is only for Hosted Repositories
	if repoType == format.REPO_TYPE_HOSTED {
		storageAttributes["write_policy"] = schema.StringAttribute{
			Description: "Controls if deployments of and updates to assets are allowed",
			Required:    true,
			Optional:    false,
			Validators: []validator.String{
				stringvalidator.OneOf(
					common.WRITE_POLICY_ALLOW,
					common.WRITE_POLICY_ALLOW_ONCE,
					common.WRITE_POLICY_DENY,
				),
			},
		}
	}

	// LatestPolicy is only for Docker Hosted Repositories
	if repoFormat == common.REPO_FORMAT_DOCKER && repoType == format.REPO_TYPE_HOSTED {
		storageAttributes["latest_policy"] = schema.BoolAttribute{
			Description: "Whether to allow redeploying the 'latest' tag but defer to the Deployment Policy for all other tags. Only applicable for Hosted Docker Repositories when Deployment Policy is set to Disable.",
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(false),
		}
	}

	return schema.Schema{
		Description: fmt.Sprintf("Manage %s %s Repositories", cases.Title(language.Und).String(repoType.String()), repoFormat),
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
			"online": schema.BoolAttribute{
				Description: "Whether this Repository is online and accepting incoming requests",
				Required:    true,
			},
			"storage": schema.SingleNestedAttribute{
				Description: "Storage configuration for this Repository",
				Required:    true,
				Optional:    false,
				Attributes:  storageAttributes,
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
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}
