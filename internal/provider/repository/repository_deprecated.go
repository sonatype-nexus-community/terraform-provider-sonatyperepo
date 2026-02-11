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
	"context"
	"fmt"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"terraform-provider-sonatyperepo/internal/provider/repository/format"
)

// repositoryResourceDeprecated wraps repositoryResource to provide deprecated resource aliases
// that delegate all operations to the underlying resource while maintaining backward compatibility
type repositoryResourceDeprecated struct {
	repositoryResource
	deprecatedName string
	newName        string
}

// Metadata returns the resource type name using the deprecated name
func (r *repositoryResourceDeprecated) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	// Set the type name to the deprecated name
	resp.TypeName = req.ProviderTypeName + "_repository_" + r.getShortName()
}

// Set Schema for this Resource
func (r *repositoryResourceDeprecated) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := standardRepositorySchema(r.RepositoryFormat.Key(), r.RepositoryType, r.RepositoryFormat.AdditionalSchemaDescription())
	maps.Copy(schema.Attributes, r.RepositoryFormat.FormatSchemaAttributes())
	schema.DeprecationMessage = fmt.Sprintf("This resource is deprecated - use instead `%s`", r.newName)
	schema.MarkdownDescription = fmt.Sprintf("~> This resource is deprecated and will be removed in the next major version (v2.x.x) - see %s", r.newName)
	resp.Schema = schema
}

// getShortName extracts the short name from the deprecated name
// e.g., "sonatyperepo_repository_maven_hosted" -> "maven_hosted"
func (r *repositoryResourceDeprecated) getShortName() string {
	// Strip the "sonatyperepo_repository_" prefix
	const prefix = "sonatyperepo_repository_"
	if len(r.deprecatedName) > len(prefix) {
		return r.deprecatedName[len(prefix):]
	}
	return r.deprecatedName
}

// Deprecated Maven resources (old names) - these are aliases to the new maven2 resources
// Users should migrate to the new resource names using moved blocks

func NewRepositoryMavenHostedDeprecated() resource.Resource {
	return &repositoryResourceDeprecated{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.MavenRepositoryFormatHosted{},
			RepositoryType:   format.REPO_TYPE_HOSTED,
		},
		deprecatedName: "sonatyperepo_repository_maven_hosted",
		newName:        "sonatyperepo_repository_maven2_hosted",
	}
}

func NewRepositoryMavenProxyDeprecated() resource.Resource {
	return &repositoryResourceDeprecated{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.MavenRepositoryFormatProxy{},
			RepositoryType:   format.REPO_TYPE_PROXY,
		},
		deprecatedName: "sonatyperepo_repository_maven_proxy",
		newName:        "sonatyperepo_repository_maven2_proxy",
	}
}

func NewRepositoryMavenGroupDeprecated() resource.Resource {
	return &repositoryResourceDeprecated{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.MavenRepositoryFormatGroup{},
			RepositoryType:   format.REPO_TYPE_GROUP,
		},
		deprecatedName: "sonatyperepo_repository_maven_group",
		newName:        "sonatyperepo_repository_maven2_group",
	}
}

// Deprecated RubyGems resources (old names) - these are aliases to the new rubygems resources

func NewRepositoryRubyGemsHostedDeprecated() resource.Resource {
	return &repositoryResourceDeprecated{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.RubyGemsRepositoryFormatHosted{},
			RepositoryType:   format.REPO_TYPE_HOSTED,
		},
		deprecatedName: "sonatyperepo_repository_ruby_gems_hosted",
		newName:        "sonatyperepo_repository_rubygems_hosted",
	}
}

func NewRepositoryRubyGemsProxyDeprecated() resource.Resource {
	return &repositoryResourceDeprecated{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.RubyGemsRepositoryFormatProxy{},
			RepositoryType:   format.REPO_TYPE_PROXY,
		},
		deprecatedName: "sonatyperepo_repository_ruby_gems_proxy",
		newName:        "sonatyperepo_repository_rubygems_proxy",
	}
}

func NewRepositoryRubyGemsGroupDeprecated() resource.Resource {
	return &repositoryResourceDeprecated{
		repositoryResource: repositoryResource{
			RepositoryFormat: &format.RubyGemsRepositoryFormatGroup{},
			RepositoryType:   format.REPO_TYPE_GROUP,
		},
		deprecatedName: "sonatyperepo_repository_ruby_gems_group",
		newName:        "sonatyperepo_repository_rubygems_group",
	}
}
