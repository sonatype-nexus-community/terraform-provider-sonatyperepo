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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// listBodyAsReturnedByNxrm is the shape NXRM 3.98.1 returns for the list endpoint: a JSON array,
// although the published API specification (and so the generated client) says it is one object.
const listBodyAsReturnedByNxrm = `[ {
  "id" : "f67d69bd-3113-42e3-b4ff-11f5d64d72d0",
  "name" : "tf-probe-1",
  "description" : "terraform provider probe",
  "roleId" : "nx-anonymous",
  "createdBy" : "admin",
  "createdAt" : "2026-11-08T00:15:20.794601Z"
}, {
  "id" : "316a8b11-ebe0-4e7b-a29a-6c36d85b569a",
  "name" : "tf-probe-2",
  "roleId" : "nx-anonymous",
  "createdBy" : "admin",
  "createdAt" : "2026-11-08T00:15:20.838723Z",
  "expiresAt" : "2026-12-08T00:15:20.838723Z"
} ]`

func TestDecodeServiceAccountTokenList(t *testing.T) {
	t.Run("array, as returned by NXRM", func(t *testing.T) {
		tokens, err := decodeServiceAccountTokenList([]byte(listBodyAsReturnedByNxrm))
		require.NoError(t, err)
		require.Len(t, tokens, 2)
		assert.Equal(t, "tf-probe-1", *tokens[0].Name)
		assert.Equal(t, "terraform provider probe", *tokens[0].Description)
		assert.Nil(t, tokens[0].ExpiresAt, "a token that never expires has no expiresAt")
		assert.Equal(t, "tf-probe-2", *tokens[1].Name)
		assert.Nil(t, tokens[1].Description)
		require.NotNil(t, tokens[1].ExpiresAt)
		assert.Equal(t, 2026, tokens[1].ExpiresAt.Year())
	})

	t.Run("empty array", func(t *testing.T) {
		tokens, err := decodeServiceAccountTokenList([]byte("[ ]"))
		require.NoError(t, err)
		assert.Empty(t, tokens)
	})

	t.Run("single object, as the API specification describes", func(t *testing.T) {
		tokens, err := decodeServiceAccountTokenList([]byte(`{"id":"abc","name":"one"}`))
		require.NoError(t, err)
		require.Len(t, tokens, 1)
		assert.Equal(t, "abc", *tokens[0].Id)
	})

	t.Run("empty body", func(t *testing.T) {
		tokens, err := decodeServiceAccountTokenList(nil)
		require.NoError(t, err)
		assert.Empty(t, tokens)
	})

	t.Run("not JSON", func(t *testing.T) {
		_, err := decodeServiceAccountTokenList([]byte("<html>nope</html>"))
		assert.Error(t, err)
	})

	t.Run("malformed array", func(t *testing.T) {
		_, err := decodeServiceAccountTokenList([]byte(`[{"id":`))
		assert.Error(t, err)
	})
}

// newServiceAccountTokenTestService returns a V397 service pointed at a fake NXRM that answers
// every request with the given status and body.
func newServiceAccountTokenTestService(t *testing.T, status int, body string) ServiceAccountTokenService {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	configuration := sonatyperepoV397.NewConfiguration()
	configuration.Servers = sonatyperepoV397.ServerConfigurations{{URL: server.URL}}
	return NewServiceAccountTokenServiceV397(sonatyperepoV397.NewAPIClient(configuration))
}

// TestListServiceAccountTokensV397DecodesArray runs the real generated client against a fake
// server returning the array NXRM actually sends. The generated client cannot decode this into
// its declared single-object type, so this guards the workaround in ListServiceAccountTokens.
func TestListServiceAccountTokensV397DecodesArray(t *testing.T) {
	service := newServiceAccountTokenTestService(t, http.StatusOK, listBodyAsReturnedByNxrm)

	tokens, httpResponse, err := service.ListServiceAccountTokens(context.Background())

	require.NoError(t, err)
	require.NotNil(t, httpResponse)
	assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
	require.Len(t, tokens, 2)
	assert.Equal(t, "316a8b11-ebe0-4e7b-a29a-6c36d85b569a", *tokens[1].Id)
}

func TestListServiceAccountTokensV397EmptyList(t *testing.T) {
	service := newServiceAccountTokenTestService(t, http.StatusOK, "[ ]")

	tokens, _, err := service.ListServiceAccountTokens(context.Background())

	require.NoError(t, err)
	assert.Empty(t, tokens)
}

func TestListServiceAccountTokensV397ReturnsHttpErrors(t *testing.T) {
	service := newServiceAccountTokenTestService(t, http.StatusForbidden, `{"id":"*","message":"forbidden"}`)

	tokens, httpResponse, err := service.ListServiceAccountTokens(context.Background())

	require.Error(t, err)
	require.NotNil(t, httpResponse)
	assert.Equal(t, http.StatusForbidden, httpResponse.StatusCode)
	assert.Nil(t, tokens)
}

func TestServiceAccountTokenServiceUnsupported(t *testing.T) {
	service := NewServiceAccountTokenServiceUnsupported()
	ctx := context.Background()

	_, _, createErr := service.CreateServiceAccountToken(ctx, sonatyperepoV397.ServiceAccountCreateXO{})
	_, _, listErr := service.ListServiceAccountTokens(ctx)
	_, deleteErr := service.DeleteServiceAccountToken(ctx, "id")

	for _, err := range []error{createErr, listErr, deleteErr} {
		require.Error(t, err)
		assert.Contains(t, err.Error(), "3.97.0")
	}
}

// TestNewServicesServiceAccountTokensByVersion checks which servers get the real API and which
// get the adapter that returns a clear error. Only servers before 3.97.0 can be exercised
// without a network call.
func TestNewServicesServiceAccountTokensByVersion(t *testing.T) {
	for _, header := range []string{"Nexus/3.82.0-08 (PRO)", "Nexus/3.93.0-06 (PRO)", "Nexus/3.94.1-06 (PRO)", "Nexus/3.96.0-09 (PRO)"} {
		t.Run(header, func(t *testing.T) {
			services := NewServices(
				ParseServerHeaderToVersion(header),
				sonatyperepoV382.NewAPIClient(sonatyperepoV382.NewConfiguration()),
				sonatyperepoV395.NewAPIClient(sonatyperepoV395.NewConfiguration()),
				sonatyperepoV397.NewAPIClient(sonatyperepoV397.NewConfiguration()),
			)

			_, _, err := services.ServiceAccountTokens.ListServiceAccountTokens(context.Background())

			require.Error(t, err)
			assert.Contains(t, err.Error(), "require Nexus Repository Manager 3.97.0 or later")
		})
	}

	t.Run("3.97.0 and later use the real API", func(t *testing.T) {
		for _, header := range []string{"Nexus/3.97.0-13 (PRO)", "Nexus/3.98.1-01 (PRO)"} {
			services := NewServices(
				ParseServerHeaderToVersion(header),
				sonatyperepoV382.NewAPIClient(sonatyperepoV382.NewConfiguration()),
				sonatyperepoV395.NewAPIClient(sonatyperepoV395.NewConfiguration()),
				sonatyperepoV397.NewAPIClient(sonatyperepoV397.NewConfiguration()),
			)
			assert.IsType(t, &serviceAccountTokenServiceV397{}, services.ServiceAccountTokens, header)
		}
	})
}
