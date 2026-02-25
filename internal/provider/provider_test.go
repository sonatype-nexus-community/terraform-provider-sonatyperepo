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
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

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
		nxrmClient.BlobStoreAPI.CreateFileBlobStore(ctx).Body(
			v3.FileBlobStoreApiCreateRequest{
				Name: v3.PtrString("default"),
				Path: v3.PtrString(os.Getenv("TF_ACC_HA_BLOB_STORE_PATH")),
			},
		).Execute()
	} else {
		log.Println("Continuing in non-HA Mode...")
	}

	// Run Tests
	exitCode := m.Run()
	os.Exit(exitCode)
}
