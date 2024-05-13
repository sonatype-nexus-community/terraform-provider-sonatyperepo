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

type RepositoryModel struct {
	Name   types.String `tfsdk:"name"`
	Format types.String `tfsdk:"format"`
	Type   types.String `tfsdk:"type"`
	Url    types.String `tfsdk:"url"`
}

type RepositoriesModel struct {
	Repositories []RepositoryModel `tfsdk:"repositories"`
}

type RepositoryMavenModel struct {
	Name        types.String                 `tfsdk:"name"`
	Format      types.String                 `tfsdk:"format"`
	Type        types.String                 `tfsdk:"type"`
	Url         types.String                 `tfsdk:"url"`
	Online      types.Bool                   `tfsdk:"online"`
	Storage     repositoryStorageModel       `tfsdk:"storage"`
	Cleanup     *RepositoryCleanupModel      `tfsdk:"cleanup"`
	Maven       repositoryMavenSpecificModel `tfsdk:"maven"`
	Component   *RepositoryComponentModel    `tfsdk:"component"`
	LastUpdated types.String                 `tfsdk:"last_updated"`
}

type repositoryStorageModel struct {
	BlobStoreName               types.String `tfsdk:"blob_store_name"`
	StrictContentTypeValidation types.Bool   `tfsdk:"strict_content_type_validation"`
	WritePolicy                 types.String `tfsdk:"write_policy"`
}

type RepositoryCleanupModel struct {
	PolicyNames []types.String `tfsdk:"policy_names"`
}

type RepositoryComponentModel struct {
	ProprietaryComponents types.Bool `tfsdk:"proprietary_components"`
}

type repositoryMavenSpecificModel struct {
	VersionPolicy      types.String `tfsdk:"version_policy"`
	LayoutPolicy       types.String `tfsdk:"layout_policy"`
	ContentDisposition types.String `tfsdk:"content_disposition"`
}
