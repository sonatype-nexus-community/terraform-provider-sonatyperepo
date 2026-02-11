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

// NewRepositoryMavenHostedResource is a helper function to simplify the provider implementation.
// This resource supports state migration from the deprecated maven_hosted resource name.
func NewRepositoryMavenHostedResource() resource.Resource {
	return &repositoryResourceWithMoveState{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.MavenRepositoryFormatHosted{},
			RepositoryType:   format.REPO_TYPE_HOSTED,
		},
		sourceResourceNames: []string{"sonatyperepo_repository_maven_hosted"},
	}
}

// NewRepositoryMavenProxyResource is a helper function to simplify the provider implementation.
// This resource supports state migration from the deprecated maven_proxy resource name.
func NewRepositoryMavenProxyResource() resource.Resource {
	return &repositoryResourceWithMoveState{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.MavenRepositoryFormatProxy{},
			RepositoryType:   format.REPO_TYPE_PROXY,
		},
		sourceResourceNames: []string{"sonatyperepo_repository_maven_proxy"},
	}
}

// NewRepositoryMavenGroupResource is a helper function to simplify the provider implementation.
// This resource supports state migration from the deprecated maven_group resource name.
func NewRepositoryMavenGroupResource() resource.Resource {
	return &repositoryResourceWithMoveState{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.MavenRepositoryFormatGroup{},
			RepositoryType:   format.REPO_TYPE_GROUP,
		},
		sourceResourceNames: []string{"sonatyperepo_repository_maven_group"},
	}
}
