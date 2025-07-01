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

package system

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

const defaultConnectionTimeoutSeconds int32 = 30

// systemConfigLdapResource is the resource implementation.
type systemConfigLdapResource struct {
	common.BaseResource
}

// NewSystemConfigLdapResource is a helper function to simplify the provider implementation.
func NewSystemConfigLdapResource() resource.Resource {
	return &systemConfigLdapResource{}
}

// Metadata returns the resource type name.
func (r *systemConfigLdapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_config_ldap_connection"
}

// Schema defines the schema for the resource.
func (r *systemConfigLdapResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Configure and LDAP connection",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Internal LDAP server ID",
				Computed:    true,
				Optional:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "LDAP connection name",
				Required:    true,
				Optional:    false,
			},
			"protocol": schema.StringAttribute{
				Description: "The user's status",
				Required:    true,
				Optional:    false,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						common.PROTOCOL_LDAP,
						common.PROTOCOL_LDAPS,
					}...),
				},
			},
			"nexus_trust_store_enabled": schema.BoolAttribute{
				Description: "Whether to use certificates stored in Nexus Repository Manager's truststore",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"hostname": schema.StringAttribute{
				Description: "LDAP server hostname",
				Required:    true,
				Optional:    false,
			},
			"port": schema.Int32Attribute{
				Description: "LDAP server port",
				Required:    true,
				Optional:    false,
			},
			"search_base": schema.StringAttribute{
				Description: "LDAP location to be added to the connection URL",
				Required:    true,
				Optional:    false,
			},
			"auth_scheme": schema.StringAttribute{
				Description: "Authentication scheme used for connecting to LDAP server",
				Required:    true,
				Optional:    false,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						common.AUTH_SCHEME_NONE,
						common.AUTH_SCHEME_SIMPLE,
						common.AUTH_SCHEME_DIGEST_MD5,
						common.AUTH_SCHEME_CRAM_MD5,
					}...),
				},
			},
			"auth_username": schema.StringAttribute{
				Description: "This must be a fully qualified username if simple authentication is used. Required if authScheme other than NONE.",
				Optional:    true,
				// TODO: Make this Requuired when auth_scheme != NONE
			},
			"auth_password": schema.StringAttribute{
				Description: "The password to bind with. Required if authScheme other than NONE.",
				Optional:    true,
				Sensitive:   true,
				// TODO: Make this Requuired when auth_scheme != NONE
			},
			"auth_realm": schema.StringAttribute{
				Description: "The SASL realm to bind to. Required if authScheme is CRAM_MD5 or DIGEST_MD5.",
				Optional:    true,
				// TODO: Make this Requuired when auth_scheme is CRAM_MD5 or DIGEST_MD5
			},
			"connection_timeout": schema.Int32Attribute{
				Description: "How many seconds to wait before timeout",
				// MarkdownDescription: "Recommended default value is 30",
				// Required: true,
				Optional: true,
				Computed: true,
				Default:  int32default.StaticInt32(defaultConnectionTimeoutSeconds),
				Validators: []validator.Int32{
					int32validator.Between(1, 3600),
				},
				// PlanModifiers: []planmodifier.Int32{
				// 	common.DefaulInt32Value(30),
				// },
			},
			"connection_retry_delay": schema.Int32Attribute{
				Description: "How many seconds to wait before retrying",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(300),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			"max_connection_attempts": schema.Int32Attribute{
				Description: "How many connection attempts before giving up",
				Required:    false,
				Optional:    true,
				Computed:    true,
				Default:     int32default.StaticInt32(3),
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.UseStateForUnknown(),
				},
			},
			// User Mapping
			"user_base_dn": schema.StringAttribute{
				Description: "The relative DN where user objects are found (e.g. ou=people). This value will have the Search base DN value appended to form the full User search base DN.",
				Optional:    true,
			},
			"user_subtree": schema.BoolAttribute{
				Description: "Are users located in structures below the user base DN?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"user_object_class": schema.StringAttribute{
				Description: "LDAP class for user objects - e.g. inetOrgPerson",
				Required:    true,
			},
			"user_ldap_filter": schema.StringAttribute{
				Description: "LDAP search filter to limit user search - e.g. (|(mail=*@example.com)(uid=dom*))",
				Optional:    true,
			},
			"user_id_attribute": schema.StringAttribute{
				Description: "This is used to find a user given its user ID - e.g. uid",
				Required:    true,
			},
			"user_real_name_attribute": schema.StringAttribute{
				Description: "This is used to find a real name given the user ID - e.g. cn",
				Required:    true,
			},
			"user_email_name_attribute": schema.StringAttribute{
				Description: "This is used to find an email address given the user ID - e.g. mail",
				Required:    true,
			},
			"user_password_attribute": schema.StringAttribute{
				Description: "If this field is blank the user will be authenticated against a bind with the LDAP server",
				Optional:    true,
			},
			"map_ldap_groups_to_roles": schema.BoolAttribute{
				Description: "Denotes whether LDAP assigned roles are used as Nexus Repository Manager roles",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			// Group Mapping
			"group_type": schema.StringAttribute{
				Description: "Defines a type of groups used: static (a group contains a list of users) or dynamic (a user contains a list of groups). Required if ldapGroupsAsRoles is true.",
				Required:    false,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{
						common.LDAP_GROUP_MAPPING_STATIC,
						common.LDAP_GROUP_MAPPING_DYNAMIC,
					}...),
				},
			},
			"user_member_of_attribute": schema.StringAttribute{
				Description: "Set this to the attribute used to store the attribute which holds groups DN in the user object. Required if group_type is DYNAMIC",
				Optional:    true,
				// TODO: Validate this is supplied if group_type == DYNAMIC
			},
			"group_base_dn": schema.StringAttribute{
				Description: "The relative DN where group objects are found (e.g. ou=Group). This value will have the Search base DN value appended to form the full Group search base DN. e.g. ou=Group",
				Optional:    true,
			},
			"group_subtree": schema.BoolAttribute{
				Description: "Are groups located in structures below the group base DN?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"group_object_class": schema.StringAttribute{
				Description: "LDAP class for group objects. Required if groupType is STATIC - e.g. posixGroup",
				Optional:    true,
			},
			"group_id_attribute": schema.StringAttribute{
				Description: "This field specifies the attribute of the Object class that defines the Group ID. Required if groupType is STATIC - e.g. cn",
				Optional:    true,
			},
			"group_member_attribute": schema.StringAttribute{
				Description: "LDAP attribute containing the usernames for the group. Required if groupType is STATIC - e.g. memberUid",
				Optional:    true,
			},
			"group_member_format": schema.StringAttribute{
				Description: "The format of user ID stored in the group member attribute. Required if groupType is STATIC - e.g. uid=${username},ou=people,dc=example,dc=com",
				Optional:    true,
			},
			"order": schema.Int32Attribute{
				Description: "Order number in which the server is being used when looking for a user - cannot be set during CREATE",
				Required:    false,
				Optional:    true,
			},
			// Meta
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
func (r *systemConfigLdapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan model.LdapServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

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
	apiResponse, err := r.Client.SecurityManagementLDAPAPI.CreateLdapServer(ctx).Body(*plan.ToApiCreateModel()).Execute()

	// Handle Error
	if err != nil || apiResponse.StatusCode != http.StatusCreated {
		common.HandleApiError(
			"Error creating LDAP Connection",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
		return
	}

	// Id & Order are not known until Create Request - we need to call GET now to obtain that
	ldapResonse, httpResponse, err := r.Client.SecurityManagementLDAPAPI.GetLdapServer(ctx, plan.Name.ValueString()).Execute()
	if err != nil || httpResponse.StatusCode != http.StatusOK {
		common.HandleApiError(
			"Error creating LDAP Connection - connection may be partially created",
			&err,
			apiResponse,
			&resp.Diagnostics,
		)
		return
	}

	plan.Id = types.StringPointerValue(ldapResonse.Id)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *systemConfigLdapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state model.LdapServerModel

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
	apiResponse, httpResponse, err := r.Client.SecurityManagementLDAPAPI.GetLdapServer(ctx, state.Name.ValueString()).Execute()

	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(ctx)
			common.HandleApiWarning(
				"LDAP Connection does not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			common.HandleApiError(
				"Error Reading LDAP Connection",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		}
		return
	} else {
		// Update State
		state.FromApiModel(apiResponse)
		state.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *systemConfigLdapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan & state
	var plan model.LdapServerModel
	var state model.LdapServerModel

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

	// Call API to Update
	body := plan.ToApiUpdateModel()
	httpResponse, err := r.Client.SecurityManagementLDAPAPI.UpdateLdapServer(ctx, state.Name.ValueString()).Body(*body).Execute()

	// Handle Error
	if err != nil || httpResponse.StatusCode != http.StatusNoContent {
		common.HandleApiError(
			"Error updating LDAP Connection",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *systemConfigLdapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state model.LdapServerModel

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

	httpResponse, err := r.Client.SecurityManagementLDAPAPI.DeleteLdapServer(ctx, state.Name.ValueString()).Execute()

	// Handle Error
	if err != nil || httpResponse.StatusCode != http.StatusNoContent {
		common.HandleApiError(
			"Error removing LDAP Connection",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}
}
