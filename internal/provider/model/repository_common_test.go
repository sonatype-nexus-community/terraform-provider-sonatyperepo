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
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"

	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	"github.com/stretchr/testify/assert"
)

var (
	TEST_HOSTED_REPO_NAME string = "test-hosted-repo"
	TEST_BLOB_STORE_NAME  string = "default"
)

func TestModelRepositoryCommonMapSimpleApiHostedRepositoryEmpty(t *testing.T) {
	m := &RepositoryNpmHostedModel{}

	// Set Name only
	m.mapSimpleApiHostedRepository(sonatyperepo.SimpleApiHostedRepository{
		Name: &TEST_HOSTED_REPO_NAME,
	})

	assert.Equal(t, TEST_HOSTED_REPO_NAME, m.Name.ValueString())
	assert.False(t, m.Online.ValueBool())
	assert.Equal(t, "", m.Storage.BlobStoreName.ValueString())
	assert.False(t, m.Storage.StrictContentTypeValidation.ValueBool())
	assert.Equal(t, "", m.Storage.WritePolicy.ValueString())
	assert.Equal(t, "", m.Url.ValueString())
	assert.Nil(t, m.Cleanup)
	assert.Nil(t, m.Component)
}

func TestModelRepositoryCommonMapSimpleApiHostedRepositoryComplete(t *testing.T) {
	m := &RepositoryNpmHostedModel{}

	// Set Name only
	m.mapSimpleApiHostedRepository(sonatyperepo.SimpleApiHostedRepository{
		Name:   &TEST_HOSTED_REPO_NAME,
		Online: true,
		Storage: sonatyperepo.HostedStorageAttributes{
			BlobStoreName:               TEST_BLOB_STORE_NAME,
			StrictContentTypeValidation: true,
			WritePolicy:                 common.WRITE_POLICY_ALLOW_ONCE,
		},
		Url: &TEST_HOSTED_REPO_NAME,
		Cleanup: &sonatyperepo.CleanupPolicyAttributes{
			PolicyNames: []string{"policy-1"},
		},
		Component: &sonatyperepo.ComponentAttributes{
			ProprietaryComponents: common.NewTrue(),
		},
	})

	assert.Equal(t, TEST_HOSTED_REPO_NAME, m.Name.ValueString())
	assert.True(t, m.Online.ValueBool())
	assert.Equal(t, TEST_BLOB_STORE_NAME, m.Storage.BlobStoreName.ValueString())
	assert.True(t, m.Storage.StrictContentTypeValidation.ValueBool())
	assert.Equal(t, common.WRITE_POLICY_ALLOW_ONCE, m.Storage.WritePolicy.ValueString())
	assert.Equal(t, TEST_HOSTED_REPO_NAME, m.Url.ValueString())
	assert.NotNil(t, m.Cleanup)
	assert.Equal(t, len(m.Cleanup.PolicyNames), 1)
	assert.NotNil(t, m.Component)
	assert.True(t, m.Component.ProprietaryComponents.ValueBool())
}
