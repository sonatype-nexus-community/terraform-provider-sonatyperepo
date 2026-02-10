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
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-sonatyperepo/internal/provider/repository/format"
)

// NewRepositoryTerraformProxyResource is a helper function to simplify the provider implementation.
func NewRepositoryTerraformProxyResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.TerraformRepositoryFormatProxy{},
		RepositoryType:   format.REPO_TYPE_PROXY,
	}
}

// NewRepositoryTerraformHostedResource is a helper function to simplify the provider implementation.
func NewRepositoryTerraformHostedResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.TerraformRepositoryFormatHosted{},
		RepositoryType:   format.REPO_TYPE_HOSTED,
	}
}
