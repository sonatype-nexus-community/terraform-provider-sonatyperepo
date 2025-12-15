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

package blob_store

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

// Common Delete Blob Store implementation.
func DeleteBlobStore(client *sonatyperepo.APIClient, ctx *context.Context, blobStoreName string, resp *resource.DeleteResponse) {
	// Delete API Call
	api_requeest := client.BlobStoreAPI.DeleteBlobStore(*ctx, blobStoreName)

	// Call API
	api_response, err := api_requeest.Execute()

	// Handle Error(s)
	if err != nil {
		if api_response.StatusCode == 404 {
			resp.State.RemoveResource(*ctx)
			resp.Diagnostics.AddWarning(
				"Blob Store to delete did not exist",
				fmt.Sprintf("Unable to delete Blob Store: %d: %s", api_response.StatusCode, api_response.Status),
			)
		} else {
			resp.Diagnostics.AddError(
				"Error deleting Blob Store",
				fmt.Sprintf("Unable to delete Blob Store: %d: %s", api_response.StatusCode, api_response.Status),
			)
		}
		return
	} else if api_response.StatusCode == http.StatusNoContent {
		resp.State.RemoveResource(*ctx)
	} else {
		resp.Diagnostics.AddError(
			"Error deleting Blob Store",
			fmt.Sprintf("Unexpected Response Code whilst deleting Blob Store: %d: %s", api_response.StatusCode, api_response.Status),
		)
	}
}
