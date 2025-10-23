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
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getCommonGroupSchemaAttributes(includeDeploy bool) map[string]schema.Attribute {
	attributes := map[string]schema.Attribute{
		"member_names": schema.SetAttribute{
			Description: "Member repositories' names",
			ElementType: types.StringType,
			Required:    false,
			Optional:    true,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		},
	}
	if includeDeploy {
		attributes["writable_member"] = schema.StringAttribute{
			Description: "This field is for the Group Deployment feature available in Sonatype Nexus Repository Pro.",
			Required:    false,
			Optional:    true,
		}
	}
	return map[string]schema.Attribute{
		"group": schema.SingleNestedAttribute{
			Description: "Group specific configuration for this Repository",
			Required:    true,
			Optional:    false,
			Attributes:  attributes,
		},
	}
}
