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
)

var CurrenTestNxrmVersion = common.ParseServerHeaderToVersion(fmt.Sprintf("Nexus/%s (PRO)", os.Getenv("NXRM_VERSION")))

func SkipIfNxrmVersionEq(t *testing.T, v *common.SystemVersion) {
	t.Helper()

	t.Logf("Checking if Current NXRM Test Version (%s) is equal to %s", CurrenTestNxrmVersion.String(), v.String())

	if v.Major == CurrenTestNxrmVersion.Major && v.Minor == CurrenTestNxrmVersion.Minor && v.Patch == CurrenTestNxrmVersion.Patch {
		t.Skipf("NXRM Version is == %s - skipping", v.String())
	}
}
