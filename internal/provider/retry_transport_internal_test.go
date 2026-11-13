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
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	transientStarted = `{"id":"*","message":"\"Invalid state: STARTED; allowed: [STOPPED]\""}`
	transientStopped = `{"id":"*","message":"\"Invalid state: STOPPED; allowed: [STARTED]\""}`
	failedState      = `{"id":"*","message":"\"Failed to validate facets; 1 failure\nInvalid state: FAILED; allowed: [ATTACHED, STARTED, STOPPED]\""}`
	updatePath       = "/service/rest/v1/repositories/npm/proxy/my-repo"
	deletePath       = "/service/rest/v1/repositories/my-repo"
)

// noRetryDelays makes retries immediate; the number of attempts is unchanged.
func noRetryDelays(t *testing.T) {
	t.Helper()
	original := repositoryRetryDelays
	repositoryRetryDelays = make([]time.Duration, len(original))
	t.Cleanup(func() { repositoryRetryDelays = original })
}

// scriptedServer answers the first len(replies) requests with the given (status, body) pairs and
// every later request with 204. It records every request body it receives.
type scriptedServer struct {
	*httptest.Server
	requests atomic.Int32
	mu       sync.Mutex
	bodies   []string
}

type reply struct {
	status int
	body   string
}

func newScriptedServer(t *testing.T, replies ...reply) *scriptedServer {
	t.Helper()
	s := &scriptedServer{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := int(s.requests.Add(1))
		body, _ := io.ReadAll(r.Body)
		s.mu.Lock()
		s.bodies = append(s.bodies, string(body))
		s.mu.Unlock()
		if n <= len(replies) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(replies[n-1].status)
			_, _ = w.Write([]byte(replies[n-1].body))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(s.Close)
	return s
}

func do(t *testing.T, ctx context.Context, method, url, body string) (*http.Response, string) {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	require.NoError(t, err)
	client := &http.Client{Transport: &MiddlewareTransport{Base: http.DefaultTransport}}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	read, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return resp, string(read)
}

func TestRepositoryUpdateIsRetriedOnTransientStateError(t *testing.T) {
	noRetryDelays(t)

	for name, transient := range map[string]string{"STARTED": transientStarted, "STOPPED": transientStopped} {
		t.Run("update rejected with "+name+" then accepted", func(t *testing.T) {
			server := newScriptedServer(t, reply{http.StatusBadRequest, transient})

			resp, _ := do(t, context.Background(), http.MethodPut, server.URL+updatePath, `{"name":"my-repo"}`)

			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
			assert.EqualValues(t, 2, server.requests.Load())
			assert.Equal(t, []string{`{"name":"my-repo"}`, `{"name":"my-repo"}`}, server.bodies, "the retry must send the same body")
		})
	}

	t.Run("delete rejected with a transient state error is retried too", func(t *testing.T) {
		server := newScriptedServer(t, reply{http.StatusBadRequest, transientStarted}, reply{http.StatusBadRequest, transientStopped})

		resp, _ := do(t, context.Background(), http.MethodDelete, server.URL+deletePath, "")

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.EqualValues(t, 3, server.requests.Load())
	})

	t.Run("a persistent transient error gives up after four requests and the body is still readable", func(t *testing.T) {
		replies := make([]reply, 10)
		for i := range replies {
			replies[i] = reply{http.StatusBadRequest, transientStarted}
		}
		server := newScriptedServer(t, replies...)

		resp, body := do(t, context.Background(), http.MethodPut, server.URL+updatePath, `{"name":"my-repo"}`)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, body, "Invalid state: STARTED", "the caller must still be able to read Nexus's message")
		assert.EqualValues(t, 4, server.requests.Load())
	})
}

func TestRepositoryRequestsAreNotRetriedOtherwise(t *testing.T) {
	noRetryDelays(t)

	cases := map[string]struct {
		method string
		path   string
		body   string
		reply  reply
	}{
		"FAILED state means a broken repository":  {http.MethodPut, updatePath, `{}`, reply{http.StatusBadRequest, failedState}},
		"some other 400":                          {http.MethodPut, updatePath, `{}`, reply{http.StatusBadRequest, `{"id":"*","message":"\"Name already in use\""}`}},
		"a 500":                                   {http.MethodPut, updatePath, `{}`, reply{http.StatusInternalServerError, `{"id":"*"}`}},
		"a create (POST) is never repeated":       {http.MethodPost, "/service/rest/v1/repositories/npm/proxy", `{}`, reply{http.StatusBadRequest, transientStarted}},
		"a read (GET) is not a lifecycle request": {http.MethodGet, updatePath, "", reply{http.StatusBadRequest, transientStarted}},
		"an update to something that is not a repository": {http.MethodPut, "/service/rest/v1/security/realms/active", `{}`,
			reply{http.StatusBadRequest, transientStarted}},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			server := newScriptedServer(t, tc.reply)

			resp, _ := do(t, context.Background(), tc.method, server.URL+tc.path, tc.body)

			assert.Equal(t, tc.reply.status, resp.StatusCode)
			assert.EqualValues(t, 1, server.requests.Load())
		})
	}
}

func TestRepositoryRetryStopsWhenContextIsCancelled(t *testing.T) {
	original := repositoryRetryDelays
	repositoryRetryDelays = []time.Duration{time.Hour, time.Hour, time.Hour}
	t.Cleanup(func() { repositoryRetryDelays = original })
	server := newScriptedServer(t, reply{http.StatusBadRequest, transientStarted})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	started := time.Now()
	resp, body := do(t, ctx, http.MethodPut, server.URL+updatePath, `{}`)

	assert.Less(t, time.Since(started), 5*time.Second, "must not wait out the hour-long delay")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "the last response is returned")
	assert.Contains(t, body, "Invalid state")
	assert.EqualValues(t, 1, server.requests.Load())
}

func TestRepositoryRequestWithoutReplayableBodyIsNotRetried(t *testing.T) {
	noRetryDelays(t)
	server := newScriptedServer(t, reply{http.StatusBadRequest, transientStarted})

	req, err := http.NewRequest(http.MethodPut, server.URL+updatePath, io.NopCloser(strings.NewReader(`{}`)))
	require.NoError(t, err)
	require.Nil(t, req.GetBody, "a plain io.ReadCloser body cannot be recreated")

	resp, err := (&MiddlewareTransport{Base: http.DefaultTransport}).RoundTrip(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.EqualValues(t, 1, server.requests.Load(), "resending a consumed body would send an empty request")
}
