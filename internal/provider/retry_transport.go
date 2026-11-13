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

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// repositoryRetryDelays are the waits before each retry of a repository update or delete that
// Nexus rejected with a transient lifecycle error - see isTransientRepositoryStateError. So a
// request is sent at most len(repositoryRetryDelays)+1 times.
//
// Nexus stops and starts a repository's facets when it is updated, and when it is removed from
// a group on delete. A request that arrives while another change to the same repository (or its
// group) is still settling is rejected with 400 "Invalid state: STARTED; allowed: [STOPPED]" (or
// the reverse), and succeeds if sent again a moment later (#536). The same request, with the
// same body, is safe to repeat.
var repositoryRetryDelays = []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}

// transientRepositoryStateRx matches only the two lifecycle states a repository passes through
// while it is being restarted. Deliberately NOT FAILED (or anything else): that means the
// repository is broken, which sending the request again does not fix.
var transientRepositoryStateRx = regexp.MustCompile(`Invalid state: (STARTED|STOPPED);`)

// isRepositoryLifecycleRequest reports whether req updates or deletes a repository. Other calls,
// creates (POST) in particular, are never retried here.
func isRepositoryLifecycleRequest(req *http.Request) bool {
	return (req.Method == http.MethodPut || req.Method == http.MethodDelete) &&
		strings.Contains(req.URL.Path, "/v1/repositories/")
}

// isTransientRepositoryStateError reports whether resp is Nexus's transient "Invalid state"
// rejection. It reads the response body and puts it back, so the caller can still read it.
func isTransientRepositoryStateError(resp *http.Response) bool {
	if resp == nil || resp.StatusCode != http.StatusBadRequest || resp.Body == nil {
		return false
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	_ = resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	return err == nil && transientRepositoryStateRx.Match(body)
}

// canReplay reports whether req can be sent a second time: it has no body, or it can recreate
// its body. Requests built by the generated API clients always can.
func canReplay(req *http.Request) bool {
	return req.Body == nil || req.Body == http.NoBody || req.GetBody != nil
}

// replayable returns a copy of req with a fresh body.
func replayable(req *http.Request) (*http.Request, error) {
	retry := req.Clone(req.Context())
	if req.Body != nil && req.Body != http.NoBody {
		body, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		retry.Body = body
	}
	return retry, nil
}

// roundTripWithRepositoryRetry sends req and, if it is a repository update or delete that Nexus
// rejected with a transient "Invalid state" error, sends it again after each delay in
// repositoryRetryDelays. Whatever the last attempt returned is returned, with its body readable.
// It stops waiting if the request's context is cancelled.
func (t *MiddlewareTransport) roundTripWithRepositoryRetry(req *http.Request) (*http.Response, error) {
	resp, err := t.Base.RoundTrip(req)
	if !isRepositoryLifecycleRequest(req) || !canReplay(req) {
		return resp, err
	}

	for attempt := 0; attempt < len(repositoryRetryDelays); attempt++ {
		if err != nil || !isTransientRepositoryStateError(resp) {
			break
		}

		tflog.Warn(req.Context(), fmt.Sprintf("%s %s was rejected with a transient repository state error, retrying in %s (attempt %d of %d)",
			req.Method, req.URL.Path, repositoryRetryDelays[attempt], attempt+2, len(repositoryRetryDelays)+1))

		timer := time.NewTimer(repositoryRetryDelays[attempt])
		select {
		case <-req.Context().Done():
			timer.Stop()
			return resp, err
		case <-timer.C:
		}

		retry, cloneErr := replayable(req)
		if cloneErr != nil {
			return resp, err
		}
		_ = resp.Body.Close()
		resp, err = t.Base.RoundTrip(retry)
	}

	return resp, err
}
