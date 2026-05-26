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

package provider_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

func TestMain(m *testing.M) {
	if os.Getenv("TF_ACC_HA_MODE") == "1" && os.Getenv("TF_ACC_HA_BLOB_STORE_PATH") != "" {
		log.Println("Setting up resources for Sonatype Nexus Repository in HA Mode...")

		clientConfiguration := v3.NewConfiguration()
		clientConfiguration.Servers = []v3.ServerConfiguration{
			{
				URL:         fmt.Sprintf("%s%s", strings.TrimRight(os.Getenv("NXRM_SERVER_URL"), "/"), "/service/rest"),
				Description: "Sonatype Nexus Repository Server",
			},
		}
		nxrmClient := v3.NewAPIClient(clientConfiguration)
		ctx := context.WithValue(
			context.Background(),
			v3.ContextBasicAuth,
			v3.BasicAuth{UserName: os.Getenv("NXRM_SERVER_USERNAME"), Password: os.Getenv("NXRM_SERVER_PASSWORD")},
		)

		// Create Default Blobstore
		createDefaultBlobStore(nxrmClient, &ctx)

		// Create Maven Central Proxy Repository, then wait until all cluster
		// nodes have replicated it before running tests.
		createMavenCentralProxy(nxrmClient, &ctx)
		waitForRepositoryReplication(os.Getenv("NXRM_SERVER_URL"), os.Getenv("NXRM_SERVER_USERNAME"), os.Getenv("NXRM_SERVER_PASSWORD"), "maven-central", 3)

	} else {
		log.Println("Continuing in non-HA Mode...")
	}

	// Run Tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

func createDefaultBlobStore(nxrmClient *v3.APIClient, ctx *context.Context) {
	httpResponse, err := nxrmClient.BlobStoreAPI.CreateFileBlobStore(*ctx).Body(
		v3.FileBlobStoreApiCreateRequest{
			Name: v3.PtrString("default"),
			Path: v3.PtrString(os.Getenv("TF_ACC_HA_BLOB_STORE_PATH")),
		},
	).Execute()

	if err != nil || (httpResponse != nil && httpResponse.StatusCode != http.StatusNoContent) {
		log.Printf("Failed to create default Blob Store: %v", err)
		if httpResponse != nil {
			log.Printf("API Response: %d", httpResponse.StatusCode)
		}
	}
}

// waitForRepositoryReplication polls the load-balancer until repoName is visible
// on nodeCount consecutive round-robin responses, confirming all cluster nodes
// have replicated the repository before tests begin.
func waitForRepositoryReplication(serverURL, username, password, repoName string, nodeCount int) {
	if serverURL == "" {
		return
	}
	listURL := fmt.Sprintf("%s/service/rest/v1/repositories", strings.TrimRight(serverURL, "/"))
	httpClient := &http.Client{Timeout: 15 * time.Second}

	consecutive := 0
	for attempt := 0; attempt < 40 && consecutive < nodeCount; attempt++ {
		req, err := http.NewRequest(http.MethodGet, listURL, nil)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		req.SetBasicAuth(username, password)

		resp, err := httpClient.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			consecutive = 0
			log.Printf("waitForRepositoryReplication: request failed (attempt %d): %v", attempt+1, err)
			time.Sleep(3 * time.Second)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var repos []struct {
			Name string `json:"name"`
		}
		found := false
		if json.Unmarshal(body, &repos) == nil {
			for _, r := range repos {
				if r.Name == repoName {
					found = true
					break
				}
			}
		}

		if found {
			consecutive++
			log.Printf("waitForRepositoryReplication: '%s' confirmed (%d/%d consecutive)", repoName, consecutive, nodeCount)
		} else {
			consecutive = 0
			log.Printf("waitForRepositoryReplication: '%s' not yet visible (attempt %d), retrying...", repoName, attempt+1)
			time.Sleep(3 * time.Second)
		}
	}

	if consecutive < nodeCount {
		log.Printf("waitForRepositoryReplication: WARNING — '%s' may not be visible on all nodes after polling", repoName)
	}
}

func createMavenCentralProxy(nxrmClient *v3.APIClient, ctx *context.Context) {
	httpClient := v3.NewHttpClientAttributesWithPreemptiveAuth()
	httpClient.AutoBlock = v3.PtrBool(true)
	httpClient.Blocked = v3.PtrBool(false)
	httpResponse, err := nxrmClient.RepositoryManagementAPI.CreateMavenProxyRepository(*ctx).Body(
		v3.MavenProxyRepositoryApiRequest{
			Name:       "maven-central",
			Online:     true,
			HttpClient: *httpClient,
			NegativeCache: v3.NegativeCacheAttributes{
				Enabled:    true,
				TimeToLive: 1440,
			},
			Proxy: v3.ProxyAttributes{
				RemoteUrl: v3.PtrString("https://repo1.maven.org/maven2/"),
			},
			Storage: v3.StorageAttributes{
				BlobStoreName:               "default",
				StrictContentTypeValidation: true,
			},
			Maven: v3.MavenAttributes{
				ContentDisposition: v3.PtrString("INLINE"),
				LayoutPolicy:       v3.PtrString("STRICT"),
				VersionPolicy:      v3.PtrString("RELEASE"),
			},
		},
	).Execute()

	if err != nil || (httpResponse != nil && httpResponse.StatusCode != http.StatusCreated) {
		log.Printf("Failed to create maven-central Proxy Repository: %v", err)
		if httpResponse != nil {
			log.Printf("API Response: %d", httpResponse.StatusCode)
		}
	}
}
