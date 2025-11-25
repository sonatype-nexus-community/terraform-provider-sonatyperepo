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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	tfschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/sonatype-nexus-community/terraform-provider-shared/schema"
)

func getCommonGroupSchemaAttributes(includeDeploy bool) map[string]tfschema.Attribute {
	attributes := map[string]tfschema.Attribute{
		"member_names": func() tfschema.ListAttribute {
			thisAttr := schema.ResourceOptionalStringList("Member repositories' names")
			thisAttr.Validators = []validator.List{
				listvalidator.SizeAtLeast(1),
			}
			return thisAttr
		}(),
	}
	if includeDeploy {
		attributes["writable_member"] = schema.ResourceOptionalString(
			"This field is for the Group Deployment feature available in Sonatype Nexus Repository Pro.",
		)
	}
	return map[string]tfschema.Attribute{
		"group": schema.ResourceRequiredSingleNestedAttribute(
			"Group specific configuration for this Repository",
			attributes,
		),
	}
}
