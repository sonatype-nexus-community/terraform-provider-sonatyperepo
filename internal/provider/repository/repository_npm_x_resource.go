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

// NewRepositoryNpmHostedResource is a helper function to simplify the provider implementation.
func NewRepositoryNpmHostedResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.NpmRepositoryFormatHosted{},
		RepositoryType:   format.REPO_TYPE_HOSTED,
	}
}

// NewRepositoryNpmProxyResource is a helper function to simplify the provider implementation.
func NewRepositoryNpmProxyResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.NpmRepositoryFormatProxy{},
		RepositoryType:   format.REPO_TYPE_PROXY,
	}
}

// NewRepositoryNpmGroupResource is a helper function to simplify the provider implementation.
func NewRepositoryNpmGroupResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.NpmRepositoryFormatGroup{},
		RepositoryType:   format.REPO_TYPE_GROUP,
	}
}
