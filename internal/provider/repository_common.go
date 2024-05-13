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

package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go"
)

// Common Delete Repository implementation
func DeleteRepository(client *sonatyperepo.APIClient, ctx *context.Context, repositoryName string, resp *resource.DeleteResponse) {
	// Delete API Call
	apiDeleteRequest := client.RepositoryManagementAPI.DeleteRepository(*ctx, repositoryName)

	// Call API
	httpResponse, err := apiDeleteRequest.Execute()

	// Handle Error(s)
	if err != nil {
		if httpResponse.StatusCode == 404 {
			resp.State.RemoveResource(*ctx)
			resp.Diagnostics.AddWarning(
				"Repository to delete did not exist",
				fmt.Sprintf("Unable to delete Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting Repository",
				fmt.Sprintf("Unable to delete Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
			)
		}
		return
	} else if httpResponse.StatusCode != http.StatusNoContent {
		resp.Diagnostics.AddError(
			"Error deleting Repository",
			fmt.Sprintf("Unexpected Response Code whilst deleting Repository: %d: %s", httpResponse.StatusCode, httpResponse.Status),
		)
	}
}
