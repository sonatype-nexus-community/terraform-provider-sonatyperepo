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

package privilege_type

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateActionsAttribute runs every validator configured on the "actions"
// attribute of the given schema against a single-element set and returns the
// resulting diagnostics.
func validateActionsAttribute(t *testing.T, attrs map[string]resourceschema.Attribute, action string) validator.SetResponse {
	t.Helper()

	actionsAttr, ok := attrs["actions"].(resourceschema.SetAttribute)
	if !ok {
		t.Fatalf("expected 'actions' attribute to be a resourceschema.SetAttribute, got %T", attrs["actions"])
	}

	req := validator.SetRequest{
		Path:        path.Root("actions"),
		ConfigValue: types.SetValueMust(types.StringType, []attr.Value{types.StringValue(action)}),
	}
	resp := &validator.SetResponse{}

	for _, v := range actionsAttr.Validators {
		v.ValidateSet(context.Background(), req, resp)
	}

	return *resp
}

// GH-475: repository-view, repository-admin and repository-content-selector
// privileges document that "ALL" is a valid BREAD action, and the NXRM API
// accepts it, but the shared schema validator only allowed the explicit
// BREAD values (ADD, BROWSE, DELETE, EDIT, READ).
func TestRepositoryScopedPrivilegeActionsAllowsAll(t *testing.T) {
	privilegeTypes := map[string]PrivilegeType{
		"repository-view":             &RepositoryViewPrivilegeType{},
		"repository-admin":            &RepositoryAdminPrivilegeType{},
		"repository-content-selector": &RepositoryContentSelectorPrivilegeType{},
	}

	for name, pt := range privilegeTypes {
		t.Run(name, func(t *testing.T) {
			resp := validateActionsAttribute(t, pt.PrivilegeTypeSchemaAttributes(), "ALL")

			if resp.Diagnostics.HasError() {
				t.Fatalf(
					"actions validator rejected \"ALL\" for %s privileges, but the schema description "+
						"and NXRM API both document it as a supported BREAD action: %s",
					name, resp.Diagnostics,
				)
			}
		})
	}
}

// Regression guard: the explicit BREAD values must keep validating
// regardless of how "ALL" is handled.
func TestRepositoryScopedPrivilegeActionsAllowsExplicitBread(t *testing.T) {
	pt := &RepositoryViewPrivilegeType{}

	for _, action := range BreadActions() {
		t.Run(action, func(t *testing.T) {
			resp := validateActionsAttribute(t, pt.PrivilegeTypeSchemaAttributes(), action)

			if resp.Diagnostics.HasError() {
				t.Fatalf("actions validator unexpectedly rejected explicit BREAD action %q: %s", action, resp.Diagnostics)
			}
		})
	}
}

// Regression guard: nonsense values must still be rejected.
func TestRepositoryScopedPrivilegeActionsRejectsUnknownValue(t *testing.T) {
	pt := &RepositoryViewPrivilegeType{}

	resp := validateActionsAttribute(t, pt.PrivilegeTypeSchemaAttributes(), "NOT-A-REAL-ACTION")

	if !resp.Diagnostics.HasError() {
		t.Fatal("actions validator should reject an unrecognised action value")
	}
}
