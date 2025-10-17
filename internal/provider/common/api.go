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
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"syscall"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func HandleApiError(message string, err *error, httpResponse *http.Response, respDiags *diag.Diagnostics) {
	networkError, errorMessage := handleNetworkError(*err)
	if networkError {
		respDiags.AddError(
			errorMessage,
			fmt.Sprintf("Networking Error: %s (%v)", errorMessage, *err),
		)
	} else {
		respDiags.AddError(
			message,
			fmt.Sprintf("%s: %s: %s", *err, httpResponse.Status, getResponseBody(httpResponse)),
		)
	}
}

func HandleApiWarning(message string, err *error, httpResponse *http.Response, respDiags *diag.Diagnostics) {
	respDiags.AddWarning(
		message,
		fmt.Sprintf("%s: %s: %s", *err, httpResponse.Status, getResponseBody(httpResponse)),
	)
}

func handleNetworkError(err error) (bool, string) {
	// Check for specific error types
	if opErr, ok := err.(*net.OpError); ok {
		// Network operation error (dial, read, write, etc.)
		return true, fmt.Sprintf("OpError: %s, %s", opErr.Op, opErr.Net)
	}

	if dnsErr, ok := err.(*net.DNSError); ok {
		// DNS resolution error
		return true, fmt.Sprintf("DNS Error: %v", dnsErr)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		// Timeout error
		return true, fmt.Sprintf("Connection timed out: %v", err)
	}

	// General network error check
	if errors.Is(err, syscall.ECONNREFUSED) {
		return true, fmt.Sprintf("Connection refused: %v", err)
	}

	return false, ""
}

func getResponseBody(httpResponse *http.Response) []byte {
	body, _ := io.ReadAll(httpResponse.Body)
	err := httpResponse.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	return body
}
