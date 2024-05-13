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

package provider

func NewFalse() *bool {
	b := false
	return &b
}

func NewTrue() *bool {
	b := true
	return &b
}

// Repository Types
const (
	REPOSITORY_TYPE_HOSTED string = "hosted"
	REPOSITORY_TYPE_PROXY  string = "proxy"
	REPOSITORY_TYPE_GROUP  string = "group"
)

// Repository Formats
const (
	REPOSITORY_FORMAT_MAVEN string = "maven2"
)
