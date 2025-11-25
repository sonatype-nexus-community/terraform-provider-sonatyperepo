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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

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
	resp.Schema = tfschema.Schema{
		Description: "Configure and LDAP connection",
		Attributes: map[string]tfschema.Attribute{
			"id":   schema.ResourceComputedString("Internal LDAP server ID"),
			"name": schema.ResourceRequiredString("LDAP connection name"),
			"protocol": schema.ResourceStringEnum(
				"LDAP protocol to use",
				common.PROTOCOL_LDAP,
				common.PROTOCOL_LDAPS,
			),
			"nexus_trust_store_enabled": schema.ResourceOptionalBoolWithDefault(
				"Whether to use certificates stored in Nexus Repository Manager's truststore", false),
			"hostname":    schema.ResourceRequiredString("LDAP server hostname"),
			"port":        schema.ResourceRequiredInt32("LDAP server port"),
			"search_base": schema.ResourceRequiredString("LDAP location to be added to the connection URL"),
			"auth_scheme": schema.ResourceStringEnum(
				"Authentication scheme used for connecting to LDAP server",
				common.AUTH_SCHEME_NONE,
				common.AUTH_SCHEME_SIMPLE,
				common.AUTH_SCHEME_DIGEST_MD5,
				common.AUTH_SCHEME_CRAM_MD5,
			),
			"auth_username": schema.ResourceOptionalString("This must be a fully qualified username if simple authentication is used. Required if authScheme other than NONE."),
			"auth_password": schema.ResourceSensitiveString("The password to bind with. Required if authScheme other than NONE."),
			"auth_realm":    schema.ResourceOptionalString("The SASL realm to bind to. Required if authScheme is CRAM_MD5 or DIGEST_MD5."),
			"connection_timeout": schema.ResourceOptionalInt32WithDefaultAndValidator(
				"How many seconds to wait before timeout",
				common.LDAP_DEFAULT_CONNECTION_TIMEOUT_SECONDS,
				int32validator.Between(
					common.LDAP_MIN_CONNECTION_TIMEOUT_SECONDS,
					common.LDAP_MAX_CONNECTION_TIMEOUT_SECONDS,
				),
			),
			"connection_retry_delay": schema.ResourceOptionalInt32WithDefaultAndPlanModifier(
				"How many seconds to wait before retrying",
				common.LDAP_DEFAULT_CONNECTION_RETRY_SECONDS,
				int32planmodifier.UseStateForUnknown(),
			),
			"max_connection_attempts": schema.ResourceOptionalInt32WithDefaultAndPlanModifier(
				"How many connection attempts before giving up",
				common.LDAP_DEFAULT_MAX_CONENCTION_ATTEMPTS,
				int32planmodifier.UseStateForUnknown(),
			),
			// User Mapping
			"user_base_dn":              schema.ResourceOptionalString("The relative DN where user objects are found (e.g. ou=people). This value will have the Search base DN value appended to form the full User search base DN."),
			"user_subtree":              schema.ResourceOptionalBoolWithDefault("Are users located in structures below the user base DN?", false),
			"user_object_class":         schema.ResourceRequiredString("LDAP class for user objects - e.g. inetOrgPerson"),
			"user_ldap_filter":          schema.ResourceOptionalString("LDAP search filter to limit user search - e.g. (|(mail=*@example.com)(uid=dom*))"),
			"user_id_attribute":         schema.ResourceRequiredString("This is used to find a user given its user ID - e.g. uid"),
			"user_real_name_attribute":  schema.ResourceRequiredString("This is used to find a real name given the user ID - e.g. cn"),
			"user_email_name_attribute": schema.ResourceRequiredString("This is used to find an email address given the user ID - e.g. mail"),
			"user_password_attribute":   schema.ResourceOptionalString("If this field is blank the user will be authenticated against a bind with the LDAP server"),
			"map_ldap_groups_to_roles":  schema.ResourceOptionalBoolWithDefault("Denotes whether LDAP assigned roles are used as Nexus Repository Manager roles", false),
			// Group Mapping
			"group_type": schema.ResourceStringEnum(
				"Defines a type of groups used: static (a group contains a list of users) or dynamic (a user contains a list of groups). Required if ldapGroupsAsRoles is true.",
				common.LDAP_GROUP_MAPPING_STATIC,
				common.LDAP_GROUP_MAPPING_DYNAMIC,
			),
			"user_member_of_attribute": schema.ResourceOptionalString("Set this to the attribute used to store the attribute which holds groups DN in the user object. Required if group_type is DYNAMIC"),
			"group_base_dn":            schema.ResourceOptionalString("The relative DN where group objects are found (e.g. ou=Group). This value will have the Search base DN value appended to form the full Group search base DN. e.g. ou=Group"),
			"group_subtree":            schema.ResourceOptionalBoolWithDefault("Are groups located in structures below the group base DN?", false),
			"group_object_class":       schema.ResourceOptionalString("LDAP class for group objects. Required if groupType is STATIC - e.g. posixGroup"),
			"group_id_attribute":       schema.ResourceOptionalString("This field specifies the attribute of the Object class that defines the Group ID. Required if groupType is STATIC - e.g. cn"),
			"group_member_attribute":   schema.ResourceOptionalString("LDAP attribute containing the usernames for the group. Required if groupType is STATIC - e.g. memberUid"),
			"group_member_format":      schema.ResourceOptionalString("The format of user ID stored in the group member attribute. Required if groupType is STATIC - e.g. uid=${username},ou=people,dc=example,dc=com"),
			"order": schema.ResourceOptionalInt32(
				"Order number in which the server is being used when looking for a user - cannot be set during CREATE",
			),
			// Meta
			"last_updated": schema.ResourceLastUpdated(),
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
		errors.HandleAPIError(
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
		errors.HandleAPIError(
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
			errors.HandleAPIWarning(
				"LDAP Connection does not exist",
				&err,
				httpResponse,
				&resp.Diagnostics,
			)
		} else {
			errors.HandleAPIError(
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
		errors.HandleAPIError(
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
		errors.HandleAPIError(
			"Error removing LDAP Connection",
			&err,
			httpResponse,
			&resp.Diagnostics,
		)
		return
	}
}
