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

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure repositoryResourceWithMoveState implements ResourceWithMoveState
var _ resource.ResourceWithMoveState = &repositoryResourceWithMoveState{}

// repositoryResourceWithMoveState wraps repositoryResource to add MoveState functionality
// This allows migration from deprecated resource names to new resource names
type repositoryResourceWithMoveState struct {
	repositoryResource
	sourceResourceNames []string // List of deprecated resource names this resource can migrate from
}

// MoveState returns the list of state movers for migrating from deprecated resource names
func (r *repositoryResourceWithMoveState) MoveState(ctx context.Context) []resource.StateMover {
	sourceSchema := r.getSourceSchema(ctx)
	movers := make([]resource.StateMover, 0, len(r.sourceResourceNames))

	for _, sourceResourceName := range r.sourceResourceNames {
		movers = append(movers, resource.StateMover{
			SourceSchema: sourceSchema,
			StateMover:   r.createStateMover(sourceResourceName),
		})
	}

	return movers
}

// getSourceSchema retrieves the schema for state migration
func (r *repositoryResourceWithMoveState) getSourceSchema(ctx context.Context) *schema.Schema {
	var schemaResp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	return &schemaResp.Schema
}

type moveStateFunc = func(context.Context, resource.MoveStateRequest, *resource.MoveStateResponse)

// createStateMover creates a state mover function for a specific source resource name
func (r *repositoryResourceWithMoveState) createStateMover(sourceName string) moveStateFunc {
	return func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
		if req.SourceTypeName != sourceName {
			tflog.Debug(ctx, "Skipping state mover", map[string]interface{}{
				"source_type_name": req.SourceTypeName,
				"expected":         sourceName,
			})
			return
		}

		r.migrateState(ctx, req, resp)
	}
}

// migrateState performs the actual state migration from source to target
func (r *repositoryResourceWithMoveState) migrateState(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
	tflog.Info(ctx, "Migrating state from deprecated resource", map[string]interface{}{
		"from": req.SourceTypeName,
		"to":   r.RepositoryFormat.ResourceName(r.RepositoryType),
	})

	if req.SourceState == nil {
		tflog.Error(ctx, "Source state is nil, cannot migrate")
		resp.Diagnostics.AddError(
			"State Migration Failed",
			"Source state is nil, cannot migrate state from deprecated resource",
		)
		return
	}

	stateModel, diags := r.RepositoryFormat.StateAsModel(ctx, *req.SourceState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to read source state during migration")
		return
	}

	resp.Diagnostics.Append(resp.TargetState.Set(ctx, stateModel)...)
	if !resp.Diagnostics.HasError() {
		tflog.Info(ctx, "Successfully migrated state from deprecated resource", map[string]interface{}{
			"from": req.SourceTypeName,
		})
	}
}
