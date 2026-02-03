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

// Generic Capability DTO for when Firewall is not configured
func FirewallCapabilityNotDefinedDTO() *v3.CapabilityDTO {
	return &v3.CapabilityDTO{
		Id:      nil,
		Type:    common.CAPABILITY_TYPE_FIREWALL_AUDIT_QUARANTINE.StringPointer(),
		Enabled: nil,
	}
}

// CapabilityHelper provides reusable capability management functions for repository resources
type CapabilityHelper struct {
	client         *v3.APIClient
	ctx            *context.Context
	capabilityType common.CapabilityType
}

// NewCapabilityHelper creates a new capability helper
func NewCapabilityHelper(client *v3.APIClient, ctx *context.Context, capabilityType common.CapabilityType) *CapabilityHelper {
	return &CapabilityHelper{
		client:         client,
		ctx:            ctx,
		capabilityType: capabilityType,
	}
}

// FindCapabilityByRepositoryId searches for a firewall audit and quarantine capability for a given repository
func (ch *CapabilityHelper) FindCapabilityByRepositoryId(repositoryId string, diags *diag.Diagnostics) *v3.CapabilityDTO {
	capabilities, httpResponse, err := ch.client.CapabilitiesAPI.List(*ch.ctx).Execute()
	if err != nil {
		errors.HandleAPIError(
			fmt.Sprintf("Error listing capabilities to find capability for repository %s", repositoryId),
			&err,
			httpResponse,
			diags,
		)
		return nil
	}

	// Search for firewall audit and quarantine capability with matching repository ID
	for _, cap := range capabilities {
		if cap.Type != nil && *cap.Type == ch.capabilityType.String() {
			if cap.Properties != nil {
				if repoId, ok := (*cap.Properties)["repository"]; ok {
					if repoId == repositoryId {
						return &cap
					}
				}
			}
		}
	}

	return nil
}

// CapabilityExists checks if a capability exists for a repository
func (ch *CapabilityHelper) CapabilityExists(repositoryId string, diags *diag.Diagnostics) bool {
	return ch.FindCapabilityByRepositoryId(repositoryId, diags) != nil
}

// CreateCapability creates a capability for a repository
func (ch *CapabilityHelper) CreateCapability(repositoryId string, quarantineEnabled bool, diags *diag.Diagnostics) *v3.CapabilityDTO {
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

	apiResponse, httpResponse, err := ch.client.CapabilitiesAPI.Create3(*ch.ctx).Body(capabilityRequest).Execute()
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
func (ch *CapabilityHelper) UpdateCapability(capabilityId string, repositoryId string, quarantineEnabled bool, diags *diag.Diagnostics) (*v3.CapabilityDTO, error) {
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

	httpResponse, err := ch.client.CapabilitiesAPI.Update3(*ch.ctx, capabilityId).Body(capabilityRequest).Execute()
	if err != nil {
		errors.HandleAPIError(
			fmt.Sprintf("Error updating %s capability (ID=%s)", ch.capabilityType.String(), capabilityId),
			&err,
			httpResponse,
			diags,
		)
		return nil, fmt.Errorf("Failed to update exisiting Capability")
	}

	return ch.FindCapabilityByRepositoryId(repositoryId, diags), nil
}

// DeleteCapability deletes an existng capability with retry logic
func (ch *CapabilityHelper) DeleteCapability(capabilityId string, diags *diag.Diagnostics) bool {
	attempts := 1
	maxAttempts := 3

	for attempts <= maxAttempts {
		httpResponse, err := ch.client.CapabilitiesAPI.Delete4(*ch.ctx, capabilityId).Execute()

		// Trap 500 Error as they occur when Repo is not in appropriate internal state
		if httpResponse.StatusCode == http.StatusInternalServerError {
			tflog.Info(*ch.ctx, fmt.Sprintf("Unexpected response when deleting capability %s (attempt %d)", capabilityId, attempts))
			attempts++
			if attempts <= maxAttempts {
				time.Sleep(1 * time.Second)
				continue
			}
		}

		if err != nil {
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
			return false
		}

		if httpResponse.StatusCode != http.StatusNoContent {
			errors.HandleAPIError(
				fmt.Sprintf("Unexpected response when deleting firewall capability %s (attempt %d)", capabilityId, attempts),
				&err,
				httpResponse,
				diags,
			)
			attempts++
			if attempts <= maxAttempts {
				time.Sleep(1 * time.Second)
				continue
			}
			return false
		}

		return true
	}

	return false
}
