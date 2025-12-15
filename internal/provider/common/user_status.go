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

type UserStatusType string

func (ust UserStatusType) String() string {
	return string(ust)
}

const (
	USER_STATUS_ACTIVE          UserStatusType = "active"
	USER_STATUS_LOCKED          UserStatusType = "locked"
	USER_STATUS_DISABLED        UserStatusType = "disabled"
	USER_STATUS_CHANGE_PASSWORD UserStatusType = "changepassword"
)

func AllUserStatusTypes() []string {
	return []string{
		USER_STATUS_ACTIVE.String(),
		USER_STATUS_LOCKED.String(),
		USER_STATUS_DISABLED.String(),
		USER_STATUS_CHANGE_PASSWORD.String(),
	}
}
