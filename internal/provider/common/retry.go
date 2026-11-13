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
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// retryOnServerError calls fn and, while it fails with an HTTP 5xx response, calls it again after
// each delay in delays - so it makes at most len(delays)+1 calls. The result of the last call is
// returned. Anything else is returned immediately without a retry: success, a 4xx response, or an
// error with no HTTP response (a network failure, a decoding failure, ...). It stops waiting, and
// returns the last result, if ctx is cancelled.
//
// ONLY use this for calls that are safe to repeat (reads). Repeating a create or an update after
// a 5xx can apply it twice.
func retryOnServerError[T any](ctx context.Context, what string, delays []time.Duration, fn func() (T, *http.Response, error)) (T, *http.Response, error) {
	result, httpResponse, err := fn()

	for attempt := 0; attempt < len(delays) && isServerError(err, httpResponse); attempt++ {
		tflog.Warn(ctx, fmt.Sprintf("%s failed with HTTP %d, retrying in %s (attempt %d of %d)",
			what, httpResponse.StatusCode, delays[attempt], attempt+2, len(delays)+1))

		timer := time.NewTimer(delays[attempt])
		select {
		case <-ctx.Done():
			timer.Stop()
			return result, httpResponse, err
		case <-timer.C:
		}

		result, httpResponse, err = fn()
	}

	return result, httpResponse, err
}

// isServerError reports whether a call failed with an HTTP 5xx response.
func isServerError(err error, httpResponse *http.Response) bool {
	return err != nil && httpResponse != nil && httpResponse.StatusCode >= http.StatusInternalServerError
}
