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

package model

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CleanupPolicyModel represents the Terraform model for a cleanup policy
type CleanupPolicyModel struct {
	Name        types.String                `tfsdk:"name"`
	Notes       types.String                `tfsdk:"notes"`
	Format      types.String                `tfsdk:"format"`
	Criteria    *CleanupPolicyCriteriaModel `tfsdk:"criteria"`
	Retain      types.Int64                 `tfsdk:"retain"`
	LastUpdated types.String                `tfsdk:"last_updated"`
}

// CleanupPolicyCriteriaModel represents the criteria for a cleanup policy
type CleanupPolicyCriteriaModel struct {
	LastBlobUpdated types.Int64  `tfsdk:"last_blob_updated"`
	LastDownloaded  types.Int64  `tfsdk:"last_downloaded"`
	ReleaseType     types.String `tfsdk:"release_type"`
	AssetRegex      types.String `tfsdk:"asset_regex"`
}
