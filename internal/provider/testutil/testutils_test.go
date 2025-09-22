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

package testutil

import (
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccVersionInRangeTrue384001(t *testing.T) {
	var testVer = common.ParseServerHeaderToVersion("Nexus/3.84.0-01 (PRO)")
	inRange, err := VersionInRange(
		&testVer,
		&common.SystemVersion{
			Major: 2,
			Minor: 0,
			Patch: 0,
			Build: 0,
		}, &common.SystemVersion{
			Major: 4,
			Minor: 0,
			Patch: 0,
			Build: 0,
		},
	)
	assert.Nil(t, err)
	assert.True(t, inRange)

	inRange, err = VersionInRange(
		&testVer,
		&common.SystemVersion{
			Major: 3,
			Minor: 0,
			Patch: 0,
			Build: 0,
		}, &common.SystemVersion{
			Major: 4,
			Minor: 0,
			Patch: 0,
			Build: 0,
		},
	)
	assert.Nil(t, err)
	assert.True(t, inRange)

	inRange, err = VersionInRange(
		&testVer,
		&common.SystemVersion{
			Major: 3,
			Minor: 0,
			Patch: 0,
			Build: 0,
		}, &common.SystemVersion{
			Major: 3,
			Minor: 85,
			Patch: 0,
			Build: 0,
		},
	)
	assert.Nil(t, err)
	assert.True(t, inRange)

	inRange, err = VersionInRange(
		&testVer,
		&common.SystemVersion{
			Major: 3,
			Minor: 84,
			Patch: 0,
			Build: 0,
		}, &common.SystemVersion{
			Major: 3,
			Minor: 85,
			Patch: 0,
			Build: 0,
		},
	)
	assert.Nil(t, err)
	assert.True(t, inRange)
}

func TestAccVersionInRangeTrue382108(t *testing.T) {
	var testVer = common.ParseServerHeaderToVersion("Nexus/3.82.1-08 (PRO)")

	inRange, err := VersionInRange(
		&testVer,
		&common.SystemVersion{
			Major: 3,
			Minor: 82,
			Patch: 0,
		},
		&common.SystemVersion{
			Major: 3,
			Minor: 84,
			Patch: 99,
		},
	)
	assert.Nil(t, err)
	assert.True(t, inRange)
}

func TestAccVersionInRangeTrue384101(t *testing.T) {
	var testVer = common.ParseServerHeaderToVersion("Nexus/3.84.1-01 (PRO)")

	inRange, err := VersionInRange(
		&testVer,
		&common.SystemVersion{
			Major: 3,
			Minor: 84,
			Patch: 0,
			Build: 1,
		}, &common.SystemVersion{
			Major: 3,
			Minor: 85,
			Patch: 0,
			Build: 0,
		},
	)
	assert.Nil(t, err)
	assert.True(t, inRange)

	inRange, err = VersionInRange(
		&testVer,
		&common.SystemVersion{
			Major: 3,
			Minor: 84,
			Patch: 1,
			Build: 1,
		}, &common.SystemVersion{
			Major: 3,
			Minor: 85,
			Patch: 0,
			Build: 0,
		},
	)
	assert.Nil(t, err)
	assert.True(t, inRange)
}

func TestAccVersionInRangeFalse384101(t *testing.T) {
	var testVer = common.ParseServerHeaderToVersion("Nexus/3.84.1-01 (PRO)")
	inRange, err := VersionInRange(
		&testVer,
		&common.SystemVersion{
			Major: 3,
			Minor: 84,
			Patch: 1,
			Build: 2,
		}, &common.SystemVersion{
			Major: 3,
			Minor: 85,
			Patch: 0,
			Build: 0,
		},
	)
	assert.Nil(t, err)
	assert.False(t, inRange)
}
