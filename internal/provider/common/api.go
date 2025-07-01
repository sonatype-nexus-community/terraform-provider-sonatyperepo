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
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func HandleApiError(message string, err *error, httpResponse *http.Response, respDiags *diag.Diagnostics) {
	respDiags.AddError(
		message,
		fmt.Sprintf("%s: %s: %s", message, httpResponse.Status, getResponseBody(httpResponse)),
	)
}

func HandleApiWarning(message string, err *error, httpResponse *http.Response, respDiags *diag.Diagnostics) {
	respDiags.AddWarning(
		"LDAP Connection does not exist",
		fmt.Sprintf("%s: %s: %s", message, httpResponse.Status, getResponseBody(httpResponse)),
	)
}

func getResponseBody(httpResponse *http.Response) []byte {
	defer httpResponse.Body.Close()
	body, _ := io.ReadAll(httpResponse.Body)
	return body
}
