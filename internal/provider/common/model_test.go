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

package common_test

import (
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

type systemVersionTestCase struct {
	input         string
	expectedMajor int
	expectedMinor int
	expectedPatch int
	expectedBuild int
	expectedIsPro bool
}

func TestModelSystemVersionParse(t *testing.T) {
	testCases := []systemVersionTestCase{
		{
			input:         "Nexus/3.0.0-03 (OSS)",
			expectedMajor: 3,
			expectedMinor: 0,
			expectedPatch: 0,
			expectedBuild: 3,
			expectedIsPro: false,
		},
		{
			input:         "Nexus/3.80.0-06 (PRO)",
			expectedMajor: 3,
			expectedMinor: 80,
			expectedPatch: 0,
			expectedBuild: 6,
			expectedIsPro: true,
		},
		{
			input:         "Nexus/3.82.0-08 (PRO)",
			expectedMajor: 3,
			expectedMinor: 82,
			expectedPatch: 0,
			expectedBuild: 8,
			expectedIsPro: true,
		},
		{
			input:         "NEXUS/3.82.0-08 (PRO)",
			expectedMajor: 3,
			expectedMinor: 82,
			expectedPatch: 0,
			expectedBuild: 8,
			expectedIsPro: true,
		},
	}

	for _, tc := range testCases {
		sv := common.ParseServerHeaderToVersion(tc.input)
		assert.Equal(t, int8(tc.expectedMajor), sv.Major)
		assert.Equal(t, int8(tc.expectedMinor), sv.Minor)
		assert.Equal(t, int8(tc.expectedPatch), sv.Patch)
		assert.Equal(t, int8(tc.expectedBuild), sv.Build)
		assert.Equal(t, tc.expectedIsPro, sv.ProVersion)
	}
}

func TestModelSystemVersionNewerThan(t *testing.T) {
	sv := common.SystemVersion{
		Major:      int8(3),
		Minor:      int8(70),
		Patch:      int8(2),
		Build:      int8(8),
		ProVersion: true,
	}

	assert.True(t, sv.NewerThan(2, 0, 0, 0))
	assert.True(t, sv.NewerThan(3, 69, 0, 0))
	assert.True(t, sv.NewerThan(3, 69, 3, 0))
	assert.True(t, sv.NewerThan(3, 70, 2, 7))
	assert.False(t, sv.NewerThan(3, 70, 3, 0))
	assert.False(t, sv.NewerThan(3, 70, 2, 9))
}
