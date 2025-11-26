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

package format

import (
	"fmt"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getResourceNameTestCase struct {
	repoFormat     string
	repoType       RepositoryType
	expectedResult string
}

func TestRepositoryFormatCommonGetResourceName(t *testing.T) {
	testCases := []getResourceNameTestCase{
		{
			common.REPO_FORMAT_NPM,
			REPO_TYPE_HOSTED,
			"repository_npm_hosted",
		},
		{
			common.REPO_FORMAT_NPM,
			REPO_TYPE_PROXY,
			"repository_npm_proxy",
		},
	}

	for i, testCase := range testCases {
		assert.Equal(
			t,
			testCase.expectedResult,
			getResourceName(testCase.repoFormat, testCase.repoType),
			fmt.Sprintf("%d: Resource Name not as expected: %s", i, testCase.expectedResult),
		)
	}
}
