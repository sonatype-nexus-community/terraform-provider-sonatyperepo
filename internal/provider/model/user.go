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
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UserModel struct {
	UserId        types.String   `tfsdk:"user_id"`
	FirstName     types.String   `tfsdk:"first_name"`
	LastName      types.String   `tfsdk:"last_name"`
	EmailAddress  types.String   `tfsdk:"email_address"`
	ReadOnly      types.Bool     `tfsdk:"read_only"`
	Source        types.String   `tfsdk:"source"`
	Status        types.String   `tfsdk:"status"`
	Roles         []types.String `tfsdk:"roles"`
	ExternalRoles []types.String `tfsdk:"external_roles"`
}

type UsersModel struct {
	Users []UserModel `tfsdk:"users"`
}
