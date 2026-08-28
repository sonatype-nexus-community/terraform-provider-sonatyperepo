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

package format

import (
	"context"
	"maps"
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"terraform-provider-sonatyperepo/internal/provider/model"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

const (
	ociVersionRequiredError string = "OCI repositories require Sonatype Nexus Repository Manager 3.94.0 or later"
)

type OciRepositoryFormat struct {
	BaseRepositoryFormat
}

type OciRepositoryFormatHosted struct {
	OciRepositoryFormat
}

type OciRepositoryFormatProxy struct {
	OciRepositoryFormat
}

type OciRepositoryFormatGroup struct {
	OciRepositoryFormat
}

// --------------------------------------------
// Generic OCI Format Functions
// --------------------------------------------
func (f *OciRepositoryFormat) Key() string {
	return common.REPO_FORMAT_OCI
}

func (f *OciRepositoryFormat) ResourceName(repoType RepositoryType) string {
	return resourceName(f.Key(), repoType)
}

func (f *OciRepositoryFormat) AdditionalSchemaDescription() string {
	return `

**NOTE:** This resource requires Sonatype Nexus Repository 3.94.0 or later - see
[here](https://help.sonatype.com/en/sonatype-nexus-repository-3-94-0-release-notes.html)
for details.`
}

func validatePlanForOciRepository(version common.SystemVersion) []string {
	// SystemVersion.Patch/Build are int8 (max 127); comparing OlderThan(3, 93, 127, 127) rather
	// than OlderThan(3, 94, 0, 0) avoids the boundary bug where version == 3.94.0-0 would
	// otherwise be excluded (NewerThan is strict: 3.94.0-0 is not "newer than" 3.94.0-0). This
	// mirrors the same idiom common.NewServices uses to select the V395 client generation.
	if version.OlderThan(3, 93, 127, 127) {
		return []string{ociVersionRequiredError}
	}
	return nil
}

// --------------------------------------------
// Hosted OCI Format Functions
// --------------------------------------------
func (f *OciRepositoryFormatHosted) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryOciHostedModel)

	// Call API to Create
	return apiClient.CreateOciHostedRepository(ctx, planModel.ToApiCreateModel())
}

func (f *OciRepositoryFormatHosted) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryOciHostedModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetOciHostedRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}

	// Temporary Workaround (same as Docker Hosted, which shares the same storage attribute
	// shape): latest_policy is not returned from the READ API for OCI Hosted.
	if stateModel.Storage.LatestPolicy.IsNull() {
		apiResponse.Storage.LatestPolicy = common.NewFalse()
	} else {
		apiResponse.Storage.LatestPolicy = stateModel.Storage.LatestPolicy.ValueBoolPointer()
	}

	return *apiResponse, httpResponse, err
}

func (f *OciRepositoryFormatHosted) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryOciHostedModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryOciHostedModel)

	// Call API to Update
	return apiClient.UpdateOciHostedRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

// DoImportRequest implements the import functionality for OCI Hosted repositories
func (f *OciRepositoryFormatHosted) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetOciHostedRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *OciRepositoryFormatHosted) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonHostedSchemaAttributes()
	maps.Copy(additionalAttributes, ociSchemaAttributes())
	maps.Copy(additionalAttributes, ociCosignSchemaAttributes())
	return additionalAttributes
}

func (f *OciRepositoryFormatHosted) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryOciHostedModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *OciRepositoryFormatHosted) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryOciHostedModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *OciRepositoryFormatHosted) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryOciHostedModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *OciRepositoryFormatHosted) UpdateStateFromApi(state, api any) any {
	var stateModel model.RepositoryOciHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciHostedModel)
	}
	stateModel.FromApiModel((api).(common.OciHostedApiRepository))
	return stateModel
}

func (f *OciRepositoryFormatHosted) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryOciHostedModel)
	var stateModel model.RepositoryOciHostedModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciHostedModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	return stateModel
}

func (f *OciRepositoryFormatHosted) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	return validatePlanForOciRepository(version)
}

// --------------------------------------------
// PROXY OCI Format Functions
// --------------------------------------------
func (f *OciRepositoryFormatProxy) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryOciProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); this is the only mode OCI supports
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Create
	return apiClient.CreateOciProxyRepository(ctx, planModel.ToApiCreateModel(), &firewallMode)
}

func (f *OciRepositoryFormatProxy) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryOciProxyModel)

	// Call to API to Read
	apiResponse, firewallMode, httpResponse, err := apiClient.GetOciProxyRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, err
}

func (f *OciRepositoryFormatProxy) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryOciProxyModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryOciProxyModel)

	// Compute inline firewall mode (NXRM 3.94+); this is the only mode OCI supports
	firewallMode := ComputeFirewallMode(f, planModel)

	// Call API to Update
	return apiClient.UpdateOciProxyRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel(), &firewallMode)
}

// DoImportRequest implements the import functionality for OCI Proxy repositories
func (f *OciRepositoryFormatProxy) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, firewallMode, httpResponse, err := apiClient.GetOciProxyRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return ProxyApiResponseWithFirewall{Repository: *apiResponse, FirewallMode: firewallMode}, httpResponse, nil
}

func (f *OciRepositoryFormatProxy) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonProxySchemaAttributes(f.SupportsRepositoryFirewall(), f.SupportsRepositoryFirewallPccs())
	maps.Copy(additionalAttributes, ociSchemaAttributes())
	maps.Copy(additionalAttributes, ociProxySchemaAttributes())
	maps.Copy(additionalAttributes, ociCosignSchemaAttributes())
	return additionalAttributes
}

func (f *OciRepositoryFormatProxy) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryOciProxyModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *OciRepositoryFormatProxy) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryOciProxyModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *OciRepositoryFormatProxy) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryOciProxyModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *OciRepositoryFormatProxy) UpdateStateFromApi(state, api any) any {
	var stateModel model.RepositoryOciProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciProxyModel)
	}

	// NXRM 3.94+ (the only version OCI supports) returns the repository wrapped with its
	// inline firewall mode; use that directly instead of the Capability-based UpateStateWithCapability path.
	if wrapped, ok := api.(ProxyApiResponseWithFirewall); ok {
		stateModel.FromApiModel((wrapped.Repository).(common.OciProxyApiRepository))
		if wrapped.FirewallMode != nil {
			keep, enabled, quarantine, _ := ResolveFirewallBlockFlags(stateModel.FirewallAuditAndQuarantine != nil, *wrapped.FirewallMode)
			if !keep {
				stateModel.FirewallAuditAndQuarantine = nil
			} else {
				if stateModel.FirewallAuditAndQuarantine == nil {
					stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineModelWithDefaults()
				}
				stateModel.FirewallAuditAndQuarantine.CapabilityId = types.StringNull()
				stateModel.FirewallAuditAndQuarantine.Enabled = types.BoolValue(enabled)
				stateModel.FirewallAuditAndQuarantine.Quarantine = types.BoolValue(quarantine)
			}
		}
		return stateModel
	}

	stateModel.FromApiModel((api).(common.OciProxyApiRepository))
	return stateModel
}

func (f *OciRepositoryFormatProxy) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryOciProxyModel)
	var stateModel model.RepositoryOciProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciProxyModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	return stateModel
}

func (f *OciRepositoryFormatProxy) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	return validatePlanForOciRepository(version)
}

func (f *OciRepositoryFormatProxy) GetRepositoryId(state any) string {
	var stateModel model.RepositoryOciProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciProxyModel)
	}
	return stateModel.Name.ValueString()
}

func (f *OciRepositoryFormatProxy) UpateStateWithCapability(state any, capability *sonatyperepo.CapabilityDTO) any {
	var stateModel = (state).(model.RepositoryOciProxyModel)
	if capability != nil {
		if stateModel.FirewallAuditAndQuarantine == nil {
			stateModel.FirewallAuditAndQuarantine = model.NewFirewallAuditAndQuarantineModelWithDefaults()
		}
		stateModel.FirewallAuditAndQuarantine.MapFromCapabilityDTO(capability)
	} else {
		stateModel.FirewallAuditAndQuarantine = nil
	}
	return stateModel
}

// Returns true only if `repository_firewall` block is supplied
func (f *OciRepositoryFormatProxy) HasFirewallConfig(state any) bool {
	var stateModel model.RepositoryOciProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine != nil {
		return true
	}
	return false
}

func (f *OciRepositoryFormatProxy) GetRepositoryFirewallEnabled(state any) bool {
	var stateModel model.RepositoryOciProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine == nil {
		return false
	}
	return stateModel.FirewallAuditAndQuarantine.Enabled.ValueBool()
}

func (f *OciRepositoryFormatProxy) GetRepositoryFirewallQuarantineEnabled(state any) bool {
	var stateModel model.RepositoryOciProxyModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciProxyModel)
	}
	if stateModel.FirewallAuditAndQuarantine != nil {
		return stateModel.FirewallAuditAndQuarantine.Quarantine.ValueBool()
	}
	return false
}

// --------------------------------------------
// GROUP OCI Format Functions
// --------------------------------------------
func (f *OciRepositoryFormatGroup) DoCreateRequest(plan any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryOciGroupModel)

	// Call API to Create
	return apiClient.CreateOciGroupRepository(ctx, planModel.ToApiCreateModel())
}

func (f *OciRepositoryFormatGroup) DoReadRequest(state any, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryOciGroupModel)

	// Call to API to Read
	apiResponse, httpResponse, err := apiClient.GetOciGroupRepository(ctx, stateModel.Name.ValueString())
	if apiResponse == nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, err
}

func (f *OciRepositoryFormatGroup) DoUpdateRequest(plan any, state any, apiClient common.RepositoryManagementService, ctx context.Context) (*http.Response, error) {
	// Cast to correct Plan Model Type
	planModel := (plan).(model.RepositoryOciGroupModel)

	// Cast to correct State Model Type
	stateModel := (state).(model.RepositoryOciGroupModel)

	// Call API to Update
	return apiClient.UpdateOciGroupRepository(ctx, stateModel.Name.ValueString(), planModel.ToApiUpdateModel())
}

// DoImportRequest implements the import functionality for OCI Group repositories
func (f *OciRepositoryFormatGroup) DoImportRequest(repositoryName string, apiClient common.RepositoryManagementService, ctx context.Context) (any, *http.Response, error) {
	// Call to API to Read repository for import
	apiResponse, httpResponse, err := apiClient.GetOciGroupRepository(ctx, repositoryName)
	if err != nil {
		return nil, httpResponse, err
	}
	return *apiResponse, httpResponse, nil
}

func (f *OciRepositoryFormatGroup) FormatSchemaAttributes() map[string]tfschema.Attribute {
	additionalAttributes := commonGroupSchemaAttributes(true)
	maps.Copy(additionalAttributes, ociSchemaAttributes())
	maps.Copy(additionalAttributes, ociCosignSchemaAttributes())
	return additionalAttributes
}

func (f *OciRepositoryFormatGroup) PlanAsModel(ctx context.Context, plan tfsdk.Plan) (any, diag.Diagnostics) {
	var planModel model.RepositoryOciGroupModel
	return planModel, plan.Get(ctx, &planModel)
}

func (f *OciRepositoryFormatGroup) StateAsModel(ctx context.Context, state tfsdk.State) (any, diag.Diagnostics) {
	var stateModel model.RepositoryOciGroupModel
	return stateModel, state.Get(ctx, &stateModel)
}

func (f *OciRepositoryFormatGroup) UpdatePlanForState(plan any) any {
	var planModel = (plan).(model.RepositoryOciGroupModel)
	planModel.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	return planModel
}

func (f *OciRepositoryFormatGroup) UpdateStateFromApi(state, api any) any {
	var stateModel model.RepositoryOciGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciGroupModel)
	}
	stateModel.FromApiModel((api).(common.OciGroupApiRepository))
	return stateModel
}

func (f *OciRepositoryFormatGroup) UpdateStateFromPlanForNonApiFields(plan, state any) any {
	var planModel = (plan).(model.RepositoryOciGroupModel)
	var stateModel model.RepositoryOciGroupModel
	// During import, state might be nil, so we create a new model
	if state != nil {
		stateModel = (state).(model.RepositoryOciGroupModel)
	}

	stateModel.MapMissingApiFieldsFromPlan(planModel)
	return stateModel
}

func (f *OciRepositoryFormatGroup) ValidatePlanForNxrmVersion(plan any, version common.SystemVersion) []string {
	return validatePlanForOciRepository(version)
}

// --------------------------------------------
// Common Functions
// --------------------------------------------
func ociSchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"oci": schema.ResourceRequiredSingleNestedAttribute(
			"OCI specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"force_basic_auth": schema.ResourceRequiredBool("Whether to force authentication (OCI Bearer Token Realm required if false)"),
				"http_port":        schema.ResourceOptionalInt32("Create an HTTP connector at specified port"),
				"https_port":       schema.ResourceOptionalInt32("Create an HTTPS connector at specified port"),
				"path_enabled":     schema.ResourceOptionalBool("Allows to use repository name in OCI image paths"),
				"subdomain":        schema.ResourceOptionalString("Allows to use subdomain"),
				"v1_enabled":       schema.ResourceRequiredBool("Whether to allow clients to use the V1 API to interact with this repository"),
			},
		),
	}
}

func ociProxySchemaAttributes() map[string]tfschema.Attribute {
	return map[string]tfschema.Attribute{
		"oci_proxy": schema.ResourceRequiredSingleNestedAttribute(
			"OCI Proxy specific configuration for this Repository",
			map[string]tfschema.Attribute{
				"cache_foreign_layers": schema.ResourceComputedOptionalBoolWithDefault(
					"Allow Nexus Repository Manager to download and cache foreign layers",
					false,
				),
				"foreign_layer_url_whitelist": func() tfschema.SetAttribute {
					thisAttr := schema.ResourceOptionalStringSet("Foreign Layer URL Whitelist")
					thisAttr.Computed = true
					thisAttr.Default = setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{}))
					return thisAttr
				}(),
				"index_type": schema.ResourceStringEnumWithDefault(
					"Type of OCI Index",
					common.OCI_PROXY_INDEX_TYPE_REGISTRY,
					common.OCI_PROXY_INDEX_TYPE_HUB,
					common.OCI_PROXY_INDEX_TYPE_REGISTRY,
					common.OCI_PROXY_INDEX_TYPE_CUSTOM,
				),
				"index_url": schema.ResourceOptionalString("Url of OCI Index to use"),
			},
		),
	}
}

// ociCosignDefaultObjectType/Value describe the object NXRM itself defaults `cosign` to
// (`{enforcement: "NONE", identity_regex: null, issuer_regex: null}`) when a request omits it --
// confirmed against a live NXRM 3.95.0-07 server, which always echoes this shape back on GET even
// when the create/update request had no `cosign` at all. The schema attribute must be Computed
// (not just Optional) and carry a matching default so Terraform doesn't see this as the provider
// conjuring a value out of a planned null, exactly like the "component" attribute in hosted.go.
var ociCosignDefaultObjectType = map[string]attr.Type{
	"enforcement":    types.StringType,
	"identity_regex": types.StringType,
	"issuer_regex":   types.StringType,
}

var ociCosignDefaultObjectValue = map[string]attr.Value{
	"enforcement":    types.StringValue(common.OCI_COSIGN_ENFORCEMENT_NONE),
	"identity_regex": types.StringNull(),
	"issuer_regex":   types.StringNull(),
}

func ociCosignSchemaAttributes() map[string]tfschema.Attribute {
	thisAttr := schema.ResourceComputedOptionalSingleNestedAttribute(
		"Cosign keyless signature enforcement configuration for this Repository",
		map[string]tfschema.Attribute{
			"enforcement": schema.ResourceStringEnumWithDefault(
				"Cosign enforcement mode. NONE disables enforcement; KEYLESS requires a valid cosign signature attached as an OCI referrer for every pull",
				common.OCI_COSIGN_ENFORCEMENT_NONE,
				common.OCI_COSIGN_ENFORCEMENT_NONE,
				common.OCI_COSIGN_ENFORCEMENT_KEYLESS,
			),
			"identity_regex": schema.ResourceOptionalString(
				"Regex matched against the Fulcio certificate Subject Alternative Name. Ignored when enforcement is NONE",
			),
			"issuer_regex": schema.ResourceOptionalString(
				"Regex matched against the OIDC issuer extension on the Fulcio certificate. Ignored when enforcement is NONE",
			),
		},
	)
	thisAttr.Default = objectdefault.StaticValue(types.ObjectValueMust(ociCosignDefaultObjectType, ociCosignDefaultObjectValue))
	thisAttr.PlanModifiers = []planmodifier.Object{
		objectplanmodifier.UseStateForUnknown(),
	}
	return map[string]tfschema.Attribute{
		"cosign": thisAttr,
	}
}
