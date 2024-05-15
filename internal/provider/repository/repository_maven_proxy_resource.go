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

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
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

// repositoryMavenProxyResource is the resource implementation.
type repositoryMavenProxyResource struct {
	common.BaseResource
}

// NewRepositoryMavenProxyResource is a helper function to simplify the provider implementation.
func NewRepositoryMavenProxyResource() resource.Resource {
	return &repositoryMavenProxyResource{}
}

// Metadata returns the resource type name.
func (r *repositoryMavenProxyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_maven_proxy"
}

// Schema defines the schema for the resource.
func (r *repositoryMavenProxyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage Proxy Maven Repositories",
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
					"write_policy": schema.StringAttribute{
						Description: "Controls if deployments of and updates to assets are allowed",
						Required:    false,
						Optional:    false,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
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
			"proxy": schema.SingleNestedAttribute{
				Description: "Proxy specific configuration for this Repository",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"remote_url": schema.StringAttribute{
						Description: "Location of the remote repository being proxied",
						Required:    false,
						Optional:    true,
					},
					"content_max_age": schema.Int64Attribute{
						Description: "How long to cache artifacts before rechecking the remote repository (in minutes)",
						Required:    true,
						Optional:    false,
					},
					"metadata_max_age": schema.Int64Attribute{
						Description: "How long to cache metadata before rechecking the remote repository (in minutes)",
						Required:    true,
						Optional:    false,
					},
				},
			},
			"negative_cache": schema.SingleNestedAttribute{
				Description: "Negative Cache configuration for this Repository",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Whether to cache responses for content not present in the proxied repository",
						Required:    true,
						Optional:    false,
					},
					"time_to_live": schema.Int64Attribute{
						Description: "How long to cache the fact that a file was not found in the repository (in minutes)",
						Required:    true,
						Optional:    false,
					},
				},
			},
			"http_client": schema.SingleNestedAttribute{
				Description: "HTTP Client configuration for this Repository",
				Required:    true,
				Optional:    false,
				Attributes: map[string]schema.Attribute{
					"blocked": schema.BoolAttribute{
						Description: "Whether to block outbound connections on the repository",
						Required:    true,
						Optional:    false,
					},
					"auto_block": schema.BoolAttribute{
						Description: "Whether to auto-block outbound connections if remote peer is detected as unreachable/unresponsive",
						Required:    true,
						Optional:    false,
					},
					"connection": schema.SingleNestedAttribute{
						Description: "HTTP Client Connection configuration for this Repository",
						Required:    false,
						Optional:    true,
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"retries": schema.Int64Attribute{
								Description: "Total retries if the initial connection attempt suffers a timeout",
								Required:    false,
								Optional:    true,
								Computed:    true,
								Default:     int64default.StaticInt64(0),
								Validators: []validator.Int64{
									int64validator.Between(0, 10),
								},
							},
							"user_agent_suffix": schema.StringAttribute{
								Description: "Custom fragment to append to User-Agent header in HTTP requests",
								Required:    false,
								Optional:    true,
							},
							"timeout": schema.Int64Attribute{
								Description: "Seconds to wait for activity before stopping and retrying the connection",
								Required:    false,
								Optional:    true,
								Computed:    true,
								Default:     int64default.StaticInt64(60),
								Validators: []validator.Int64{
									int64validator.Between(1, 3600),
								},
							},
							"enable_circular_redirects": schema.BoolAttribute{
								Description: "Whether to enable redirects to the same location (may be required by some servers)",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"enable_cookies": schema.BoolAttribute{
								Description: "Whether to allow cookies to be stored and used",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
							"use_trust_store": schema.BoolAttribute{
								Description: "Use certificates stored in the Nexus Repository Manager truststore to connect to external systems",
								Optional:    true,
								Computed:    true,
								Default:     booldefault.StaticBool(false),
							},
						},
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
					},
					"authentication": schema.SingleNestedAttribute{
						Description: "Maven specific configuration for this Repository",
						Required:    false,
						Optional:    true,
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Description: "Authentication type - either 'username' or 'ntlm'",
								Required:    false,
								Optional:    true,
								Validators: []validator.String{
									stringvalidator.OneOf("username", "ntlm"),
								},
							},
							"username": schema.StringAttribute{
								Description: "Username",
								Required:    false,
								Optional:    true,
							},
							"password": schema.StringAttribute{
								Description: "Password",
								Required:    false,
								Optional:    true,
								Sensitive:   true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"ntlm_host": schema.StringAttribute{
								Description: "NTLM Host",
								Required:    false,
								Optional:    true,
							},
							"ntlm_domain": schema.StringAttribute{
								Description: "NTLM Domain",
								Required:    false,
								Optional:    true,
							},
							"preemptive": schema.BoolAttribute{
								Description: "Whether to use pre-emptive authentication. Use with caution. Defaults to false.",
								Required:    false,
								Optional:    true,
								// Computed:    true,
								// Default:     booldefault.StaticBool(false),
							},
						},
					},
				},
			},
			"routing_rule": schema.StringAttribute{
				Description: "Routing Rule",
				Required:    false,
				Optional:    true,
			},
			"replication": schema.SingleNestedAttribute{
				Description: "Replication configuration for this Repository",
				Required:    false,
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"preemptive_pull_enabled": schema.BoolAttribute{
						Description: "Whether pre-emptive pull is enabled",
						Required:    true,
						Optional:    false,
					},
					"asset_path_regex": schema.StringAttribute{
						Description: "Regular Expression of Asset Paths to pull pre-emptively pull",
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
func (r *repositoryMavenProxyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.RepositoryMavenProxyModel

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

	requestPayload := makeApiRequest(&plan)
	createRequest := r.Client.RepositoryManagementAPI.CreateMavenProxyRepository(ctx).Body(requestPayload)
	httpResponse, err := createRequest.Execute()

	// Handle Error
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Maven Proxy Repository",
			fmt.Sprintf("Error creating Maven Proxy Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
		return
	} else if httpResponse.StatusCode != http.StatusCreated && httpResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Error creating Maven Proxy Repository",
			fmt.Sprintf("Unexpected Response Code whilst creating Maven Proxy Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}

	// Crank in some defaults that whilst send in request, do not appear in response
	if plan.Format.IsNull() {
		plan.Format = types.StringValue(REPOSITORY_FORMAT_MAVEN)
	}
	if plan.Type.IsNull() {
		plan.Type = types.StringValue(REPOSITORY_TYPE_PROXY)
	}
	if plan.Storage.WritePolicy.IsNull() {
		plan.Storage.WritePolicy = types.StringValue("ALLOW")
	}
	// E.g. http://localhost:8081/repository/maven-proxy-repo-test - this is not included in response to CREATE
	plan.Url = types.StringValue(fmt.Sprintf("%s/repository/%s", r.BaseUrl, plan.Name.ValueString()))
	if plan.HttpClient.Connection.EnableCircularRedirects.IsNull() {
		plan.HttpClient.Connection.EnableCircularRedirects = types.BoolValue(false)
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *repositoryMavenProxyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.RepositoryMavenProxyModel

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
	repositoryApiResponse, httpResponse, err := r.Client.RepositoryManagementAPI.GetMavenProxyRepository(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			resp.Diagnostics.AddWarning(
				"Requested Maven Proxy Repository does not exist",
				fmt.Sprintf("Unable to read Maven Proxy Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error Reading Maven Proxy Repository",
				fmt.Sprintf("Unable to read Maven Proxy Repository: %s: %s", httpResponse.Status, err),
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
		state.Storage.WritePolicy = types.StringValue(*repositoryApiResponse.Storage.WritePolicy)
		if repositoryApiResponse.Cleanup != nil {
			policies := make([]types.String, len(repositoryApiResponse.Cleanup.PolicyNames), 0)
			for i, p := range repositoryApiResponse.Cleanup.PolicyNames {
				policies[i] = types.StringValue(p)
			}
			state.Cleanup = &model.RepositoryCleanupModel{
				PolicyNames: policies,
			}
		}
		state.Proxy.ContentMaxAge = types.Int64Value(int64(repositoryApiResponse.Proxy.ContentMaxAge))
		state.Proxy.MetadataMaxAge = types.Int64Value(int64(repositoryApiResponse.Proxy.MetadataMaxAge))
		if repositoryApiResponse.Proxy.RemoteUrl != nil {
			state.Proxy.RemoteUrl = types.StringValue(*repositoryApiResponse.Proxy.RemoteUrl)
		}
		state.NegativeCache.Enabled = types.BoolValue(repositoryApiResponse.NegativeCache.Enabled)
		state.NegativeCache.TimeToLive = types.Int64Value(int64(repositoryApiResponse.NegativeCache.TimeToLive))
		state.HttpClient.Blocked = types.BoolValue(repositoryApiResponse.HttpClient.Blocked)
		state.HttpClient.AutoBlock = types.BoolValue(repositoryApiResponse.HttpClient.AutoBlock)
		if repositoryApiResponse.HttpClient.Connection != nil {
			state.HttpClient.Connection = &model.RepositoryHttpClientConnectionModel{}
			if repositoryApiResponse.HttpClient.Connection.Retries != nil {
				state.HttpClient.Connection.Retries = types.Int64Value(int64(*repositoryApiResponse.HttpClient.Connection.Retries))
			}
			if repositoryApiResponse.HttpClient.Connection.UserAgentSuffix != nil {
				state.HttpClient.Connection.UserAgentSuffix = types.StringValue(*repositoryApiResponse.HttpClient.Connection.UserAgentSuffix)
			}
			if repositoryApiResponse.HttpClient.Connection.Timeout != nil {
				state.HttpClient.Connection.Timeout = types.Int64Value(int64(*repositoryApiResponse.HttpClient.Connection.Timeout))
			}
			if repositoryApiResponse.HttpClient.Connection.EnableCircularRedirects != nil {
				state.HttpClient.Connection.EnableCircularRedirects = types.BoolValue(*repositoryApiResponse.HttpClient.Connection.EnableCircularRedirects)
			}
			if repositoryApiResponse.HttpClient.Connection.EnableCookies != nil {
				state.HttpClient.Connection.EnableCookies = types.BoolValue(*repositoryApiResponse.HttpClient.Connection.EnableCookies)
			}
			if repositoryApiResponse.HttpClient.Connection.UseTrustStore != nil {
				state.HttpClient.Connection.UseTrustStore = types.BoolValue(*repositoryApiResponse.HttpClient.Connection.UseTrustStore)
			}
		}
		if repositoryApiResponse.HttpClient.Authentication != nil {
			if state.HttpClient.Authentication == nil {
				state.HttpClient.Authentication = &model.RepositoryHttpClientAuthenticationModel{}
			}
			if repositoryApiResponse.HttpClient.Authentication.Type != nil {
				state.HttpClient.Authentication.Type = types.StringValue(*repositoryApiResponse.HttpClient.Authentication.Type)
			}
			if repositoryApiResponse.HttpClient.Authentication.Username != nil {
				state.HttpClient.Authentication.Username = types.StringValue(*repositoryApiResponse.HttpClient.Authentication.Username)
			}
			if repositoryApiResponse.HttpClient.Authentication.NtlmHost != nil {
				state.HttpClient.Authentication.NtlmHost = types.StringValue(*repositoryApiResponse.HttpClient.Authentication.NtlmHost)
			}
			if repositoryApiResponse.HttpClient.Authentication.NtlmDomain != nil {
				state.HttpClient.Authentication.NtlmDomain = types.StringValue(*repositoryApiResponse.HttpClient.Authentication.NtlmDomain)
			}
			if repositoryApiResponse.HttpClient.Authentication.Preemptive != nil {
				state.HttpClient.Authentication.Preemptive = types.BoolValue(*repositoryApiResponse.HttpClient.Authentication.Preemptive)
			}
		}
		if repositoryApiResponse.RoutingRuleName != nil {
			state.RoutingRule = types.StringValue(*repositoryApiResponse.RoutingRuleName)
		}
		if repositoryApiResponse.Replication != nil {
			state.Replication = &model.RepositoryReplicationModel{
				PreemptivePullEnabled: types.BoolValue(repositoryApiResponse.Replication.PreemptivePullEnabled),
			}
			if repositoryApiResponse.Replication.AssetPathRegex != nil {
				state.Replication.AssetPathRegex = types.StringValue(*repositoryApiResponse.Replication.AssetPathRegex)
			}
		}
		state.Maven.ContentDisposition = types.StringValue(*repositoryApiResponse.Maven.ContentDisposition)
		state.Maven.LayoutPolicy = types.StringValue(*repositoryApiResponse.Maven.LayoutPolicy)
		state.Maven.VersionPolicy = types.StringValue(*repositoryApiResponse.Maven.VersionPolicy)

		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repositoryMavenProxyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.RepositoryMavenProxyModel
	var state model.RepositoryMavenProxyModel

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
	requestPayload := makeApiRequest(&plan)
	apiUpdateRequest := r.Client.RepositoryManagementAPI.UpdateMavenProxyRepository(ctx, state.Name.ValueString()).Body(requestPayload)

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
func (r *repositoryMavenProxyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.RepositoryMavenProxyModel

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

func makeApiRequest(plan *model.RepositoryMavenProxyModel) sonatyperepo.MavenProxyRepositoryApiRequest {
	requestPayload := sonatyperepo.MavenProxyRepositoryApiRequest{
		Name:   plan.Name.ValueString(),
		Maven:  sonatyperepo.MavenAttributes{},
		Online: plan.Online.ValueBool(),
		Storage: sonatyperepo.StorageAttributes{
			BlobStoreName:               plan.Storage.BlobStoreName.ValueString(),
			StrictContentTypeValidation: plan.Storage.StrictContentTypeValidation.ValueBool(),
		},
		Proxy: sonatyperepo.ProxyAttributes{
			ContentMaxAge:  int32(plan.Proxy.ContentMaxAge.ValueInt64()),
			MetadataMaxAge: int32(plan.Proxy.MetadataMaxAge.ValueInt64()),
		},
		NegativeCache: sonatyperepo.NegativeCacheAttributes{
			Enabled:    plan.NegativeCache.Enabled.ValueBool(),
			TimeToLive: int32(plan.NegativeCache.TimeToLive.ValueInt64()),
		},
		HttpClient: sonatyperepo.HttpClientAttributesWithPreemptiveAuth{
			Blocked:   plan.HttpClient.Blocked.ValueBool(),
			AutoBlock: plan.HttpClient.AutoBlock.ValueBool(),
		},
	}
	if !plan.Proxy.RemoteUrl.IsNull() {
		requestPayload.Proxy.RemoteUrl = plan.Proxy.RemoteUrl.ValueStringPointer()
	}
	if plan.HttpClient.Connection != nil {
		requestPayload.HttpClient.Connection = &sonatyperepo.HttpClientConnectionAttributes{}
		if !plan.HttpClient.Connection.Retries.IsNull() {
			retries := int32(plan.HttpClient.Connection.Retries.ValueInt64())
			requestPayload.HttpClient.Connection.Retries = &retries
		}
		if !plan.HttpClient.Connection.UserAgentSuffix.IsNull() {
			requestPayload.HttpClient.Connection.UserAgentSuffix = plan.HttpClient.Connection.UserAgentSuffix.ValueStringPointer()
		}
		if !plan.HttpClient.Connection.Timeout.IsNull() {
			timeout := int32(plan.HttpClient.Connection.Timeout.ValueInt64())
			requestPayload.HttpClient.Connection.Timeout = &timeout
		}
		if !plan.HttpClient.Connection.EnableCircularRedirects.IsNull() {
			requestPayload.HttpClient.Connection.EnableCircularRedirects = plan.HttpClient.Connection.EnableCircularRedirects.ValueBoolPointer()
		}
		if !plan.HttpClient.Connection.EnableCookies.IsNull() {
			requestPayload.HttpClient.Connection.EnableCookies = plan.HttpClient.Connection.EnableCookies.ValueBoolPointer()
		}
		if !plan.HttpClient.Connection.UseTrustStore.IsNull() {
			requestPayload.HttpClient.Connection.UseTrustStore = plan.HttpClient.Connection.UseTrustStore.ValueBoolPointer()
		}
		if plan.HttpClient.Authentication != nil {
			requestPayload.HttpClient.Authentication = &sonatyperepo.HttpClientConnectionAuthenticationAttributesWithPreemptive{}
			if !plan.HttpClient.Authentication.Type.IsNull() {
				requestPayload.HttpClient.Authentication.Type = plan.HttpClient.Authentication.Type.ValueStringPointer()
			}
			if !plan.HttpClient.Authentication.Username.IsNull() {
				requestPayload.HttpClient.Authentication.Username = plan.HttpClient.Authentication.Username.ValueStringPointer()
			}
			if !plan.HttpClient.Authentication.Password.IsNull() {
				requestPayload.HttpClient.Authentication.Password = plan.HttpClient.Authentication.Password.ValueStringPointer()
			}
			if !plan.HttpClient.Authentication.NtlmHost.IsNull() {
				requestPayload.HttpClient.Authentication.NtlmHost = plan.HttpClient.Authentication.NtlmHost.ValueStringPointer()
			}
			if !plan.HttpClient.Authentication.NtlmDomain.IsNull() {
				requestPayload.HttpClient.Authentication.NtlmDomain = plan.HttpClient.Authentication.NtlmDomain.ValueStringPointer()
			}
			if !plan.HttpClient.Authentication.Preemptive.IsNull() {
				requestPayload.HttpClient.Authentication.Preemptive = plan.HttpClient.Authentication.Preemptive.ValueBoolPointer()
			}
		}
	} else {
		plan.HttpClient.Connection = &model.RepositoryHttpClientConnectionModel{
			EnableCircularRedirects: types.BoolValue(false),
			EnableCookies:           types.BoolValue(false),
			UseTrustStore:           types.BoolValue(false),
		}
	}
	if !plan.RoutingRule.IsNull() {
		requestPayload.RoutingRule = plan.RoutingRule.ValueStringPointer()
	}
	if plan.Replication != nil {
		requestPayload.Replication = &sonatyperepo.ReplicationAttributes{
			PreemptivePullEnabled: plan.Replication.PreemptivePullEnabled.ValueBool(),
		}
		if !plan.Replication.AssetPathRegex.IsNull() {
			requestPayload.Replication.AssetPathRegex = plan.Replication.AssetPathRegex.ValueStringPointer()
		}
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

	return requestPayload
}
