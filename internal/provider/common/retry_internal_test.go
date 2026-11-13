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
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// noDelays makes retries immediate in tests; the list has the same length as the production one,
// so the number of attempts is unchanged.
func noDelays(t *testing.T) {
	t.Helper()
	original := capabilityListRetryDelays
	capabilityListRetryDelays = make([]time.Duration, len(original))
	t.Cleanup(func() { capabilityListRetryDelays = original })
}

// scripted returns a call function that answers with the given HTTP statuses in turn (the last
// one repeats), and a counter of how many calls were made.
func scripted(statuses ...int) (func() (string, *http.Response, error), *atomic.Int32) {
	var calls atomic.Int32
	return func() (string, *http.Response, error) {
		n := int(calls.Add(1))
		status := statuses[len(statuses)-1]
		if n <= len(statuses) {
			status = statuses[n-1]
		}
		response := &http.Response{StatusCode: status}
		if status >= 300 {
			return "", response, errors.New(http.StatusText(status))
		}
		return "ok", response, nil
	}, &calls
}

func TestRetryOnServerError(t *testing.T) {
	delays := []time.Duration{0, 0}

	t.Run("success on the first call is not retried", func(t *testing.T) {
		fn, calls := scripted(200)
		result, response, err := retryOnServerError(context.Background(), "test", delays, fn)
		require.NoError(t, err)
		assert.Equal(t, "ok", result)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.EqualValues(t, 1, calls.Load())
	})

	t.Run("a 500 then success is retried and succeeds", func(t *testing.T) {
		fn, calls := scripted(500, 200)
		result, response, err := retryOnServerError(context.Background(), "test", delays, fn)
		require.NoError(t, err)
		assert.Equal(t, "ok", result)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.EqualValues(t, 2, calls.Load())
	})

	t.Run("a persistent 5xx stops after len(delays)+1 calls and returns the last error", func(t *testing.T) {
		fn, calls := scripted(503)
		_, response, err := retryOnServerError(context.Background(), "test", delays, fn)
		require.Error(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, response.StatusCode)
		assert.EqualValues(t, 3, calls.Load())
	})

	t.Run("a 4xx is not retried", func(t *testing.T) {
		for _, status := range []int{400, 401, 403, 404} {
			fn, calls := scripted(status)
			_, response, err := retryOnServerError(context.Background(), "test", delays, fn)
			require.Error(t, err)
			assert.Equal(t, status, response.StatusCode)
			assert.EqualValues(t, 1, calls.Load(), "status %d", status)
		}
	})

	t.Run("an error with no HTTP response is not retried", func(t *testing.T) {
		var calls atomic.Int32
		_, response, err := retryOnServerError(context.Background(), "test", delays, func() (string, *http.Response, error) {
			calls.Add(1)
			return "", nil, errors.New("connection refused")
		})
		require.Error(t, err)
		assert.Nil(t, response)
		assert.EqualValues(t, 1, calls.Load())
	})

	t.Run("a cancelled context stops the wait between attempts", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		fn, calls := scripted(500)

		started := time.Now()
		_, response, err := retryOnServerError(ctx, "test", []time.Duration{time.Hour, time.Hour}, fn)

		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
		assert.EqualValues(t, 1, calls.Load(), "must not retry once the context is cancelled")
		assert.Less(t, time.Since(started), 5*time.Second, "must not wait out the delay")
	})
}

// capabilityListServer answers every request with a 500 for the first `failures` requests and
// then with a one-capability list, and counts the requests it received.
func capabilityListServer(t *testing.T, failures int32, failureStatus int) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if requests.Add(1) <= failures {
			w.WriteHeader(failureStatus)
			_, _ = w.Write([]byte(`{"id":"*","message":"transient"}`))
			return
		}
		_, _ = w.Write([]byte(`[{"id":"cap-1","type":"rapture.baseurl","enabled":true,"notes":"n","properties":{"url":"http://x"}}]`))
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

// capabilityServices builds each client generation's CapabilityService against a server URL.
func capabilityServices(url string) map[string]CapabilityService {
	c382 := sonatyperepoV382.NewConfiguration()
	c382.Servers = sonatyperepoV382.ServerConfigurations{{URL: url}}
	c395 := sonatyperepoV395.NewConfiguration()
	c395.Servers = sonatyperepoV395.ServerConfigurations{{URL: url}}
	c397 := sonatyperepoV397.NewConfiguration()
	c397.Servers = sonatyperepoV397.ServerConfigurations{{URL: url}}
	return map[string]CapabilityService{
		"V382": &capabilityServiceV382{client: sonatyperepoV382.NewAPIClient(c382)},
		"V395": &capabilityServiceV395{client: sonatyperepoV395.NewAPIClient(c395)},
		"V397": &capabilityServiceV397{client: sonatyperepoV397.NewAPIClient(c397)},
	}
}

// TestCapabilityListRetriesServerErrors covers #537 for every client generation: a 500 from the
// capability list (what Nexus returns, intermittently, while repositories are being deleted) is
// retried and the capability is found, instead of surfacing as an error.
func TestCapabilityListRetriesServerErrors(t *testing.T) {
	noDelays(t)

	t.Run("a transient 500 is retried", func(t *testing.T) {
		server, requests := capabilityListServer(t, 1, http.StatusInternalServerError)
		for name, service := range capabilityServices(server.URL) {
			requests.Store(0)

			capabilities, response, err := service.List(context.Background())

			require.NoError(t, err, name)
			require.Len(t, capabilities, 1, name)
			assert.Equal(t, "cap-1", *capabilities[0].Id, name)
			assert.Equal(t, http.StatusOK, response.StatusCode, name)
			assert.EqualValues(t, 2, requests.Load(), "%s: one failed request and one retry", name)
		}
	})

	t.Run("a persistent 500 gives up after three requests", func(t *testing.T) {
		server, requests := capabilityListServer(t, 1000, http.StatusInternalServerError)
		for name, service := range capabilityServices(server.URL) {
			requests.Store(0)

			_, response, err := service.List(context.Background())

			require.Error(t, err, name)
			require.NotNil(t, response, name)
			assert.Equal(t, http.StatusInternalServerError, response.StatusCode, name)
			assert.EqualValues(t, 3, requests.Load(), name)
		}
	})

	t.Run("a 403 is not retried", func(t *testing.T) {
		server, requests := capabilityListServer(t, 1000, http.StatusForbidden)
		for name, service := range capabilityServices(server.URL) {
			requests.Store(0)

			_, response, err := service.List(context.Background())

			require.Error(t, err, name)
			assert.Equal(t, http.StatusForbidden, response.StatusCode, name)
			assert.EqualValues(t, 1, requests.Load(), name)
		}
	})
}
