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
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOciValidatePlanForNxrmVersion asserts that OCI repositories (hosted, proxy and group) are
// rejected against any NXRM version older than 3.94.0 -- the version OCI was introduced in -- and
// accepted from 3.94.0 onward.
func TestOciValidatePlanForNxrmVersion(t *testing.T) {
	tooOld := common.SystemVersion{Major: 3, Minor: 93, Patch: 127, Build: 127}
	exactlySupported := common.SystemVersion{Major: 3, Minor: 94, Patch: 0, Build: 0}
	newer := common.SystemVersion{Major: 3, Minor: 95, Patch: 0, Build: 7}

	hosted := &OciRepositoryFormatHosted{}
	proxy := &OciRepositoryFormatProxy{}
	group := &OciRepositoryFormatGroup{}

	for _, version := range []common.SystemVersion{tooOld} {
		assert.Equal(t, []string{ociVersionRequiredError}, hosted.ValidatePlanForNxrmVersion(nil, version))
		assert.Equal(t, []string{ociVersionRequiredError}, proxy.ValidatePlanForNxrmVersion(nil, version))
		assert.Equal(t, []string{ociVersionRequiredError}, group.ValidatePlanForNxrmVersion(nil, version))
	}

	for _, version := range []common.SystemVersion{exactlySupported, newer} {
		assert.Nil(t, hosted.ValidatePlanForNxrmVersion(nil, version))
		assert.Nil(t, proxy.ValidatePlanForNxrmVersion(nil, version))
		assert.Nil(t, group.ValidatePlanForNxrmVersion(nil, version))
	}
}
