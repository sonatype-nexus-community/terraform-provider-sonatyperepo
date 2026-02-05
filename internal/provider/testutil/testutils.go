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
	"fmt"
	"os"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"

	semver "github.com/hashicorp/go-version"
)

var CurrenTestNxrmVersion = common.ParseServerHeaderToVersion(fmt.Sprintf("Nexus/%s (PRO)", os.Getenv("NXRM_VERSION")))

func SkipIfNxrmVersionEq(t *testing.T, v *common.SystemVersion) {
	t.Helper()

	if v.Major == CurrenTestNxrmVersion.Major && v.Minor == CurrenTestNxrmVersion.Minor && v.Patch == CurrenTestNxrmVersion.Patch {
		t.Skipf("NXRM Version is == %s - skipping", v.String())
	}
}

func SkipIfNxrmVersionInRange(t *testing.T, low *common.SystemVersion, high *common.SystemVersion) {
	t.Helper()

	inRange, err := VersionInRange(&CurrenTestNxrmVersion, low, high)

	if err != nil {
		t.Errorf("Error comparing versions: %v", err)
		t.FailNow()
	}

	if inRange {
		t.Skipf("NXRM Version within range %s and %s - skipping", low.String(), high.String())
	}
}

func VersionInRange(ver *common.SystemVersion, low *common.SystemVersion, high *common.SystemVersion) (bool, error) {
	thisVersion, err := semver.NewVersion(ver.SemVerString())
	if err != nil {
		return false, err
	}

	lowVersion, err := semver.NewVersion(low.SemVerString())
	if err != nil {
		return false, err
	}

	highVersion, err := semver.NewVersion(high.SemVerString())
	if err != nil {
		return false, err
	}

	if lowVersion.LessThanOrEqual(thisVersion) && highVersion.GreaterThanOrEqual(thisVersion) {
		return true, nil
	}

	return false, nil
}
