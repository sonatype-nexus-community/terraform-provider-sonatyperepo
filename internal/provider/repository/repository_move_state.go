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
	movers := make([]resource.StateMover, 0, len(r.sourceResourceNames))

	// Get the schema for this resource to use as the source schema
	// Since deprecated resources use the same schema, we can use the current resource's schema
	var schemaResp resource.SchemaResponse
	r.repositoryResource.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	sourceSchema := &schemaResp.Schema

	for _, sourceResourceName := range r.sourceResourceNames {
		// Capture the source name in the closure
		sourceName := sourceResourceName
		movers = append(movers, resource.StateMover{
			// Provide the schema so SourceState will be populated
			SourceSchema: sourceSchema,

			StateMover: func(ctx context.Context, req resource.MoveStateRequest, resp *resource.MoveStateResponse) {
				// Check if this mover should handle the request by verifying the source type name
				if req.SourceTypeName != sourceName {
					// Not for this source type, skip
					tflog.Debug(ctx, "Skipping state mover", map[string]interface{}{
						"source_type_name": req.SourceTypeName,
						"expected":         sourceName,
					})
					return
				}

				tflog.Info(ctx, "Migrating state from deprecated resource", map[string]interface{}{
					"from": req.SourceTypeName,
					"to":   r.repositoryResource.RepositoryFormat.ResourceName(r.repositoryResource.RepositoryType),
				})

				// The state structure is identical between deprecated and new resources
				// Copy the source state to the target state
				if req.SourceState != nil {
					// Get the state model from the source
					stateModel, diags := r.repositoryResource.RepositoryFormat.StateAsModel(ctx, *req.SourceState)
					resp.Diagnostics.Append(diags...)
					if resp.Diagnostics.HasError() {
						tflog.Error(ctx, "Failed to read source state during migration")
						return
					}

					// Set the state model into the target state
					resp.Diagnostics.Append(resp.TargetState.Set(ctx, stateModel)...)
					if !resp.Diagnostics.HasError() {
						tflog.Info(ctx, "Successfully migrated state from deprecated resource", map[string]interface{}{
							"from": req.SourceTypeName,
						})
					}
				} else {
					tflog.Error(ctx, "Source state is nil, cannot migrate")
					resp.Diagnostics.AddError(
						"State Migration Failed",
						"Source state is nil, cannot migrate state from deprecated resource",
					)
				}
			},
		})
	}

	return movers
}
