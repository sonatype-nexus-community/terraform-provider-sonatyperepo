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

package capability

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"

	"github.com/sonatype-nexus-community/terraform-provider-shared/errors"
)

// CapabilityHelper provides reusable capability management functions for repository resources
type CapabilityHelper struct {
	client         *v3.APIClient
	capabilityType common.CapabilityType
}

// NewCapabilityHelper creates a new capability helper
func NewCapabilityHelper(client *v3.APIClient, capabilityType common.CapabilityType) *CapabilityHelper {
	return &CapabilityHelper{
		client:         client,
		capabilityType: capabilityType,
	}
}

// FindCapabilityByRepositoryId searches for a firewall audit and quarantine capability for a given repository
func (ch *CapabilityHelper) FindCapabilityByRepositoryId(ctx context.Context, repositoryId string, diags *diag.Diagnostics) *v3.CapabilityDTO {
	capabilities, httpResponse, err := ch.client.CapabilitiesAPI.List(ctx).Execute()
	if err != nil {
		errors.HandleAPIError(
			fmt.Sprintf("Error listing capabilities to find capability for repository %s", repositoryId),
			&err,
			httpResponse,
			diags,
		)
		return nil
	}

	return ch.findMatchingCapability(capabilities, repositoryId)
}

// findMatchingCapability searches for a capability matching the given repository ID
func (ch *CapabilityHelper) findMatchingCapability(capabilities []v3.CapabilityDTO, repositoryId string) *v3.CapabilityDTO {
	for _, cap := range capabilities {
		if ch.isCapabilityMatch(&cap, repositoryId) {
			return &cap
		}
	}
	return nil
}

// isCapabilityMatch checks if a capability matches the expected type and repository ID
func (ch *CapabilityHelper) isCapabilityMatch(cap *v3.CapabilityDTO, repositoryId string) bool {
	if cap.Type == nil {
		return false
	}
	if *cap.Type != ch.capabilityType.String() {
		return false
	}
	if cap.Properties == nil {
		return false
	}

	repoId, ok := (*cap.Properties)["repository"]
	return ok && repoId == repositoryId
}

// CapabilityExists checks if a capability exists for a repository
func (ch *CapabilityHelper) CapabilityExists(ctx context.Context, repositoryId string, diags *diag.Diagnostics) bool {
	return ch.FindCapabilityByRepositoryId(ctx, repositoryId, diags) != nil
}

// CreateCapability creates a capability for a repository
func (ch *CapabilityHelper) CreateCapability(ctx context.Context, repositoryId string, quarantineEnabled bool, diags *diag.Diagnostics) *v3.CapabilityDTO {
	// Create the capability request
	enabled := true
	properties := map[string]string{
		"repository": repositoryId,
		"quarantine": fmt.Sprintf("%t", quarantineEnabled),
	}
	capabilityRequest := v3.CapabilityDTO{
		Type:       ch.capabilityType.StringPointer(),
		Enabled:    &enabled,
		Properties: &properties,
	}

	apiResponse, httpResponse, err := ch.client.CapabilitiesAPI.Create3(ctx).Body(capabilityRequest).Execute()
	if err != nil {
		errors.HandleAPIError(
			fmt.Sprintf("Error creating %s capability for repository %s", ch.capabilityType.String(), repositoryId),
			&err,
			httpResponse,
			diags,
		)
		return nil
	}

	return apiResponse
}

// UpdateCapability updates an existing capability
func (ch *CapabilityHelper) UpdateCapability(ctx context.Context, capabilityId string, repositoryId string, capabilityEnabled bool, quarantineEnabled bool, diags *diag.Diagnostics) (*v3.CapabilityDTO, error) {
	properties := map[string]string{
		"repository": repositoryId,
		"quarantine": fmt.Sprintf("%t", quarantineEnabled),
	}
	capabilityRequest := v3.CapabilityDTO{
		Id:         &capabilityId,
		Type:       ch.capabilityType.StringPointer(),
		Enabled:    &capabilityEnabled,
		Properties: &properties,
	}

	httpResponse, err := ch.client.CapabilitiesAPI.Update3(ctx, capabilityId).Body(capabilityRequest).Execute()
	if err != nil {
		errors.HandleAPIError(
			fmt.Sprintf("Error updating %s capability (ID=%s)", ch.capabilityType.String(), capabilityId),
			&err,
			httpResponse,
			diags,
		)
		return nil, fmt.Errorf("failed to update exisiting Capability")
	}

	return ch.FindCapabilityByRepositoryId(ctx, repositoryId, diags), nil
}

// DeleteCapability deletes an existng capability with retry logic
func (ch *CapabilityHelper) DeleteCapability(ctx context.Context, capabilityId string, diags *diag.Diagnostics) bool {
	const maxAttempts = 3

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if ch.attemptDeleteCapability(ctx, capabilityId, attempt, maxAttempts, diags) {
			return true
		}
	}
	return false
}

// attemptDeleteCapability performs a single delete attempt and returns true if successful
func (ch *CapabilityHelper) attemptDeleteCapability(ctx context.Context, capabilityId string, attempt int, maxAttempts int, diags *diag.Diagnostics) bool {
	httpResponse, err := ch.client.CapabilitiesAPI.Delete4(ctx, capabilityId).Execute()

	// Trap 500 Error as they occur when Repo is not in appropriate internal state
	if httpResponse.StatusCode == http.StatusInternalServerError {
		tflog.Info(ctx, fmt.Sprintf("Unexpected response when deleting capability %s (attempt %d)", capabilityId, attempt))
		if attempt < maxAttempts {
			time.Sleep(1 * time.Second)
		}
		return false
	}

	// Handle errors other than success
	if err != nil {
		ch.handleDeleteError(err, httpResponse, capabilityId, diags)
		return false
	}

	// Check for success
	if httpResponse.StatusCode == http.StatusNoContent {
		return true
	}

	// Unexpected status code - retry if attempts remaining
	errors.HandleAPIError(
		fmt.Sprintf("Unexpected response when deleting firewall capability %s (attempt %d)", capabilityId, attempt),
		&err,
		httpResponse,
		diags,
	)
	if attempt < maxAttempts {
		time.Sleep(1 * time.Second)
	}
	return false
}

// handleDeleteError handles errors from delete operations
func (ch *CapabilityHelper) handleDeleteError(err error, httpResponse *http.Response, capabilityId string, diags *diag.Diagnostics) {
	if httpResponse.StatusCode == http.StatusNotFound {
		errors.HandleAPIWarning(
			fmt.Sprintf("Firewall capability (ID=%s) did not exist to delete", capabilityId),
			&err,
			httpResponse,
			diags,
		)
	} else {
		errors.HandleAPIError(
			fmt.Sprintf("Error deleting firewall capability (ID=%s)", capabilityId),
			&err,
			httpResponse,
			diags,
		)
	}
}
