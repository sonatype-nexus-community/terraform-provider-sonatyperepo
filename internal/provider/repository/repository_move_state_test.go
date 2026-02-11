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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"terraform-provider-sonatyperepo/internal/provider/repository/format"
)

func TestRepositoryResourceWithMoveStateMavenHosted(t *testing.T) {
	// Create the new Maven2 hosted resource
	res := NewRepositoryMavenHostedResource()

	// Assert it implements ResourceWithMoveState
	moveStateResource, ok := res.(resource.ResourceWithMoveState)
	require.True(t, ok, "Maven2 hosted resource should implement ResourceWithMoveState")

	// Get the state movers
	ctx := context.Background()
	movers := moveStateResource.MoveState(ctx)

	// Should have exactly one mover for the deprecated maven_hosted resource
	assert.Len(t, movers, 1, "Should have one state mover")

	// Verify the mover has a source schema
	assert.NotNil(t, movers[0].SourceSchema, "State mover should have a source schema")

	// Verify the mover function is not nil
	assert.NotNil(t, movers[0].StateMover, "State mover should have a mover function")
}

func TestRepositoryResourceWithMoveStateMavenProxy(t *testing.T) {
	// Create the new Maven2 proxy resource
	res := NewRepositoryMavenProxyResource()

	// Assert it implements ResourceWithMoveState
	moveStateResource, ok := res.(resource.ResourceWithMoveState)
	require.True(t, ok, "Maven2 proxy resource should implement ResourceWithMoveState")

	// Get the state movers
	ctx := context.Background()
	movers := moveStateResource.MoveState(ctx)

	// Should have exactly one mover for the deprecated maven_proxy resource
	assert.Len(t, movers, 1, "Should have one state mover")
}

func TestRepositoryResourceWithMoveStateMavenGroup(t *testing.T) {
	// Create the new Maven2 group resource
	res := NewRepositoryMavenGroupResource()

	// Assert it implements ResourceWithMoveState
	moveStateResource, ok := res.(resource.ResourceWithMoveState)
	require.True(t, ok, "Maven2 group resource should implement ResourceWithMoveState")

	// Get the state movers
	ctx := context.Background()
	movers := moveStateResource.MoveState(ctx)

	// Should have exactly one mover for the deprecated maven_group resource
	assert.Len(t, movers, 1, "Should have one state mover")
}

func TestRepositoryResourceWithMoveStateRubyGemsHosted(t *testing.T) {
	// Create the new RubyGems hosted resource
	res := NewRepositoryRubyGemsHostedResource()

	// Assert it implements ResourceWithMoveState
	moveStateResource, ok := res.(resource.ResourceWithMoveState)
	require.True(t, ok, "RubyGems hosted resource should implement ResourceWithMoveState")

	// Get the state movers
	ctx := context.Background()
	movers := moveStateResource.MoveState(ctx)

	// Should have exactly one mover for the deprecated ruby_gems_hosted resource
	assert.Len(t, movers, 1, "Should have one state mover")
}

func TestRepositoryResourceWithMoveStateRubyGemsProxy(t *testing.T) {
	// Create the new RubyGems proxy resource
	res := NewRepositoryRubyGemsProxyResource()

	// Assert it implements ResourceWithMoveState
	moveStateResource, ok := res.(resource.ResourceWithMoveState)
	require.True(t, ok, "RubyGems proxy resource should implement ResourceWithMoveState")

	// Get the state movers
	ctx := context.Background()
	movers := moveStateResource.MoveState(ctx)

	// Should have exactly one mover for the deprecated ruby_gems_proxy resource
	assert.Len(t, movers, 1, "Should have one state mover")
}

func TestRepositoryResourceWithMoveStateRubyGemsGroup(t *testing.T) {
	// Create the new RubyGems group resource
	res := NewRepositoryRubyGemsGroupResource()

	// Assert it implements ResourceWithMoveState
	moveStateResource, ok := res.(resource.ResourceWithMoveState)
	require.True(t, ok, "RubyGems group resource should implement ResourceWithMoveState")

	// Get the state movers
	ctx := context.Background()
	movers := moveStateResource.MoveState(ctx)

	// Should have exactly one mover for the deprecated ruby_gems_group resource
	assert.Len(t, movers, 1, "Should have one state mover")
}

func TestRepositoryResourceDeprecatedMavenHostedMetadata(t *testing.T) {
	// Create the deprecated Maven hosted resource
	res := NewRepositoryMavenHostedDeprecated()

	// Get metadata
	ctx := context.Background()
	req := resource.MetadataRequest{
		ProviderTypeName: "sonatyperepo",
	}
	var resp resource.MetadataResponse

	res.Metadata(ctx, req, &resp)

	// Verify the type name is set to the deprecated name
	assert.Equal(t, "sonatyperepo_repository_maven_hosted", resp.TypeName, "Deprecated resource should use old name")
}

func TestRepositoryResourceDeprecatedMavenProxyMetadata(t *testing.T) {
	// Create the deprecated Maven proxy resource
	res := NewRepositoryMavenProxyDeprecated()

	// Get metadata
	ctx := context.Background()
	req := resource.MetadataRequest{
		ProviderTypeName: "sonatyperepo",
	}
	var resp resource.MetadataResponse

	res.Metadata(ctx, req, &resp)

	// Verify the type name is set to the deprecated name
	assert.Equal(t, "sonatyperepo_repository_maven_proxy", resp.TypeName, "Deprecated resource should use old name")
}

func TestRepositoryResourceDeprecatedMavenGroupMetadata(t *testing.T) {
	// Create the deprecated Maven group resource
	res := NewRepositoryMavenGroupDeprecated()

	// Get metadata
	ctx := context.Background()
	req := resource.MetadataRequest{
		ProviderTypeName: "sonatyperepo",
	}
	var resp resource.MetadataResponse

	res.Metadata(ctx, req, &resp)

	// Verify the type name is set to the deprecated name
	assert.Equal(t, "sonatyperepo_repository_maven_group", resp.TypeName, "Deprecated resource should use old name")
}

func TestRepositoryResourceDeprecatedRubyGemsHostedMetadata(t *testing.T) {
	// Create the deprecated RubyGems hosted resource
	res := NewRepositoryRubyGemsHostedDeprecated()

	// Get metadata
	ctx := context.Background()
	req := resource.MetadataRequest{
		ProviderTypeName: "sonatyperepo",
	}
	var resp resource.MetadataResponse

	res.Metadata(ctx, req, &resp)

	// Verify the type name is set to the deprecated name
	assert.Equal(t, "sonatyperepo_repository_ruby_gems_hosted", resp.TypeName, "Deprecated resource should use old name")
}

func TestRepositoryResourceDeprecatedRubyGemsProxyMetadata(t *testing.T) {
	// Create the deprecated RubyGems proxy resource
	res := NewRepositoryRubyGemsProxyDeprecated()

	// Get metadata
	ctx := context.Background()
	req := resource.MetadataRequest{
		ProviderTypeName: "sonatyperepo",
	}
	var resp resource.MetadataResponse

	res.Metadata(ctx, req, &resp)

	// Verify the type name is set to the deprecated name
	assert.Equal(t, "sonatyperepo_repository_ruby_gems_proxy", resp.TypeName, "Deprecated resource should use old name")
}

func TestRepositoryResourceDeprecatedRubyGemsGroupMetadata(t *testing.T) {
	// Create the deprecated RubyGems group resource
	res := NewRepositoryRubyGemsGroupDeprecated()

	// Get metadata
	ctx := context.Background()
	req := resource.MetadataRequest{
		ProviderTypeName: "sonatyperepo",
	}
	var resp resource.MetadataResponse

	res.Metadata(ctx, req, &resp)

	// Verify the type name is set to the deprecated name
	assert.Equal(t, "sonatyperepo_repository_ruby_gems_group", resp.TypeName, "Deprecated resource should use old name")
}

func TestRepositoryResourceDeprecatedDelegatesSchema(t *testing.T) {
	// Create both deprecated and new resources
	deprecatedRes := NewRepositoryMavenHostedDeprecated()
	newRes := &repositoryResource{
		RepositoryFormat: &format.MavenRepositoryFormatHosted{},
		RepositoryType:   format.REPO_TYPE_HOSTED,
	}

	ctx := context.Background()
	req := resource.SchemaRequest{}

	// Get schemas from both
	var deprecatedResp resource.SchemaResponse
	deprecatedRes.Schema(ctx, req, &deprecatedResp)

	var newResp resource.SchemaResponse
	newRes.Schema(ctx, req, &newResp)

	// Both should have schemas with attributes
	assert.NotEmpty(t, deprecatedResp.Schema.Attributes, "Deprecated resource should have schema attributes")
	assert.NotEmpty(t, newResp.Schema.Attributes, "New resource should have schema attributes")

	// Both should have the same attributes (schema structure is identical)
	assert.Equal(t, len(newResp.Schema.Attributes), len(deprecatedResp.Schema.Attributes),
		"Deprecated and new resources should have same number of attributes")
}
