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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

func TestModelRepositoryNpmHostedToApiCreateModelEmpty(t *testing.T) {
	m := &RepositoryNpmHostedModel{}
	m.Name = types.StringValue(TEST_HOSTED_REPO_NAME)

	// Map to API Model
	apiModel := m.ToApiCreateModel()

	assert.Equal(t, TEST_HOSTED_REPO_NAME, apiModel.Name)
	assert.False(t, apiModel.Online)
	assert.Equal(t, "", apiModel.Storage.BlobStoreName)
	assert.False(t, apiModel.Storage.StrictContentTypeValidation)
	assert.Equal(t, "", apiModel.Storage.WritePolicy)
	assert.Equal(t, len(apiModel.Cleanup.GetPolicyNames()), 0)
	assert.False(t, *apiModel.Component.ProprietaryComponents)
}
