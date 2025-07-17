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

// NewRepositoryCargoHostedResource is a helper function to simplify the provider implementation.
func NewRepositoryCargoHostedResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.CargoRepositoryFormatHosted{},
		RepositoryType:   format.REPO_TYPE_HOSTED,
	}
}

// NewRepositoryCargoProxyResource is a helper function to simplify the provider implementation.
func NewRepositoryCargoProxyResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.CargoRepositoryFormatProxy{},
		RepositoryType:   format.REPO_TYPE_PROXY,
	}
}

// NewRepositoryCargoGroupResource is a helper function to simplify the provider implementation.
func NewRepositoryCargoGroupResource() resource.Resource {
	return &repositoryResource{
		RepositoryFormat: &format.CargoRepositoryFormatGroup{},
		RepositoryType:   format.REPO_TYPE_GROUP,
	}
}
