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

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DefaultValue return a string plan modifier that sets the specified value if the planned value is Null.
func DefaulInt32Value(i int32) planmodifier.Int32 {
	return defaultInt32Value{
		val: i,
	}
}

// defaultInt32Value holds our default value and allows us to implement the `planmodifier.String` interface
type defaultInt32Value struct {
	val int32
}

// Description implements the `planmodifier.String` interface
func (m defaultInt32Value) Description(context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %d", m.val)
}

// MarkdownDescription implements the `planmodifier.String` interface
func (m defaultInt32Value) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx) // reuse our plaintext Description
}

// PlanModifyString implements the `planmodifier.String` interface
func (m defaultInt32Value) PlanModifyInt32(ctx context.Context, req planmodifier.Int32Request, resp *planmodifier.Int32Response) {
	// If the attribute configuration is not null it is explicit; we should apply the planned value.
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	println(fmt.Sprintf("**** RETURNING DEFAULT INT32 VALUE: %d", m.val))

	// Otherwise, the configuration is null, so apply the default value to the response.
	resp.PlanValue = types.Int32Value(m.val)
}
