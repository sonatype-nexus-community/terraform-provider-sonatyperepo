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

package common

const (
	NXRM_DEFAULT_ANONYMOUS_USERNAME string = "anonymous"
	NXRM_DEFAULT_REALM_NAME         string = "NexusAuthorizingRealm"
	NXRM_PLACEHOLDER_PASSWORD       string = "#~NXRM~PLACEHOLDER~PASSWORD~#"
)

func NewFalse() *bool {
	b := false
	return &b
}

func NewTrue() *bool {
	b := true
	return &b
}

func StringPointer(s string) *string { return &s }
