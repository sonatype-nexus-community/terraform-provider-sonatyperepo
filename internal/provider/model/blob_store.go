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

type BlobStoreSoftQuota struct {
	Type  types.String `tfsdk:"type"`
	Limit types.Int64  `tfsdk:"limit"`
}

type BlobStoreModel struct {
	Name                  types.String       `tfsdk:"name"`
	Type                  types.String       `tfsdk:"type"`
	Unavailable           types.Bool         `tfsdk:"unavailable"`
	BlobCount             types.Int64        `tfsdk:"blob_count"`
	TotalSizeInBytes      types.Int64        `tfsdk:"total_size_in_bytes"`
	AvailableSpaceInBytes types.Int64        `tfsdk:"available_space_in_bytes"`
	SoftQuota             BlobStoreSoftQuota `tfsdk:"soft_quota"`
}

type BlobStoresModel struct {
	BlobStores []BlobStoreModel `tfsdk:"blob_stores"`
}
