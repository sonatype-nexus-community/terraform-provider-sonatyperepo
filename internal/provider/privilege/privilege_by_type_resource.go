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

package privilege

import (
	"terraform-provider-sonatyperepo/internal/provider/privilege/privilege_type"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// NewApplicationPrivilegeResource is a helper function to simplify the provider implementation.
func NewApplicationPrivilegeResource() resource.Resource {
	return &privilegeResource{
		PrivilegeType:     &privilege_type.ApplicationPrivilegeType{},
		PrivilegeTypeType: privilege_type.TypeApplication,
	}
}

// NewRepositoryAdminPrivilegeResource is a helper function to simplify the provider implementation.
func NewRepositoryAdminPrivilegeResource() resource.Resource {
	return &privilegeResource{
		PrivilegeType:     &privilege_type.RepositoryAdminPrivilegeType{},
		PrivilegeTypeType: privilege_type.TypeRepositoryAdmin,
	}
}
