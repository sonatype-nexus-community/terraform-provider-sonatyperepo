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

// NewRepositoryConanHostedResource is a helper function to simplify the provider implementation.
func NewRepositoryConanHostedResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.ConanRepositoryFormatHosted{},
		RepositoryType:   format.REPO_TYPE_HOSTED,
	}
}

// NewRepositoryConanProxyResource is a helper function to simplify the provider implementation.
func NewRepositoryConanProxyResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.ConanRepositoryFormatProxy{},
		RepositoryType:   format.REPO_TYPE_PROXY,
	}
}

// NewRepositoryConanGroupResource is a helper function to simplify the provider implementation.
func NewRepositoryConanGroupResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.ConanRepositoryFormatGroup{},
		RepositoryType:   format.REPO_TYPE_GROUP,
	}
}
