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

	sonatyperepoV382 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sonatyperepoV395 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v395"
)

// BlobStoreService abstracts blob store CRUD across NXRM API client generations.
// Every method's request/response shape is expressed in terms of the V382 generated types.
type BlobStoreService interface {
	// S3 Blob Store operations
	CreateS3BlobStore(ctx context.Context, body sonatyperepoV382.S3BlobStoreApiModel) (*http.Response, error)
	GetS3BlobStore(ctx context.Context, name string) (*sonatyperepoV382.S3BlobStoreApiModel, *http.Response, error)
	UpdateS3BlobStore(ctx context.Context, name string, body sonatyperepoV382.S3BlobStoreApiModel) (*http.Response, error)

	// File Blob Store operations
	CreateFileBlobStore(ctx context.Context, body sonatyperepoV382.FileBlobStoreApiCreateRequest) (*http.Response, error)
	GetFileBlobStoreConfiguration(ctx context.Context, name string) (*sonatyperepoV382.FileBlobStoreApiModel, *http.Response, error)
	UpdateFileBlobStore(ctx context.Context, name string, body sonatyperepoV382.FileBlobStoreApiUpdateRequest) (*http.Response, error)

	// Group Blob Store operations
	CreateGroupBlobStore(ctx context.Context, body sonatyperepoV382.GroupBlobStoreApiCreateRequest) (*http.Response, error)
	GetGroupBlobStoreConfiguration(ctx context.Context, name string) (*sonatyperepoV382.GroupBlobStoreApiModel, *http.Response, error)
	UpdateGroupBlobStore(ctx context.Context, name string, body sonatyperepoV382.GroupBlobStoreApiUpdateRequest) (*http.Response, error)

	// Azure Cloud Storage (ACS) Blob Store operations (named CreateBlobStore1/GetBlobStore1/UpdateBlobStore1 in V382)
	CreateBlobStore1(ctx context.Context, body sonatyperepoV382.AzureBlobStoreApiModel) (*http.Response, error)
	GetBlobStore1(ctx context.Context, name string) (*sonatyperepoV382.AzureBlobStoreApiModel, *http.Response, error)
	UpdateBlobStore1(ctx context.Context, name string, body sonatyperepoV382.AzureBlobStoreApiModel) (*http.Response, error)

	// Google Cloud Blob Store operations (named CreateBlobStore2/GetBlobStore2/UpdateBlobStore2 in V382)
	CreateBlobStore2(ctx context.Context, body sonatyperepoV382.GoogleCloudBlobstoreApiModel) (*http.Response, error)
	GetBlobStore2(ctx context.Context, name string) (*sonatyperepoV382.GoogleCloudBlobstoreApiModel, *http.Response, error)
	UpdateBlobStore2(ctx context.Context, name string, body sonatyperepoV382.GoogleCloudBlobstoreApiModel) (*http.Response, error)

	// List all blob stores
	ListBlobStores(ctx context.Context) ([]sonatyperepoV382.GenericBlobStoreApiResponse, *http.Response, error)
}

// blobStoreServiceV382 implements BlobStoreService against NXRM API client V382 (targets NXRM < 3.94.0).
type blobStoreServiceV382 struct {
	client *sonatyperepoV382.APIClient
}

// blobStoreServiceV395 implements BlobStoreService against NXRM API client V395 (targets NXRM 3.94.0+).
type blobStoreServiceV395 struct {
	client *sonatyperepoV395.APIClient
}

// ============================================================================
// V382 Implementation - Pure passthrough
// ============================================================================

func (s *blobStoreServiceV382) CreateS3BlobStore(ctx context.Context, body sonatyperepoV382.S3BlobStoreApiModel) (*http.Response, error) {
	return s.client.BlobStoreAPI.CreateS3BlobStore(ctx).Body(body).Execute()
}

func (s *blobStoreServiceV382) GetS3BlobStore(ctx context.Context, name string) (*sonatyperepoV382.S3BlobStoreApiModel, *http.Response, error) {
	return s.client.BlobStoreAPI.GetS3BlobStore(ctx, name).Execute()
}

func (s *blobStoreServiceV382) UpdateS3BlobStore(ctx context.Context, name string, body sonatyperepoV382.S3BlobStoreApiModel) (*http.Response, error) {
	return s.client.BlobStoreAPI.UpdateS3BlobStore(ctx, name).Body(body).Execute()
}

func (s *blobStoreServiceV382) CreateFileBlobStore(ctx context.Context, body sonatyperepoV382.FileBlobStoreApiCreateRequest) (*http.Response, error) {
	return s.client.BlobStoreAPI.CreateFileBlobStore(ctx).Body(body).Execute()
}

func (s *blobStoreServiceV382) GetFileBlobStoreConfiguration(ctx context.Context, name string) (*sonatyperepoV382.FileBlobStoreApiModel, *http.Response, error) {
	return s.client.BlobStoreAPI.GetFileBlobStoreConfiguration(ctx, name).Execute()
}

func (s *blobStoreServiceV382) UpdateFileBlobStore(ctx context.Context, name string, body sonatyperepoV382.FileBlobStoreApiUpdateRequest) (*http.Response, error) {
	return s.client.BlobStoreAPI.UpdateFileBlobStore(ctx, name).Body(body).Execute()
}

func (s *blobStoreServiceV382) CreateGroupBlobStore(ctx context.Context, body sonatyperepoV382.GroupBlobStoreApiCreateRequest) (*http.Response, error) {
	return s.client.BlobStoreAPI.CreateGroupBlobStore(ctx).Body(body).Execute()
}

func (s *blobStoreServiceV382) GetGroupBlobStoreConfiguration(ctx context.Context, name string) (*sonatyperepoV382.GroupBlobStoreApiModel, *http.Response, error) {
	return s.client.BlobStoreAPI.GetGroupBlobStoreConfiguration(ctx, name).Execute()
}

func (s *blobStoreServiceV382) UpdateGroupBlobStore(ctx context.Context, name string, body sonatyperepoV382.GroupBlobStoreApiUpdateRequest) (*http.Response, error) {
	return s.client.BlobStoreAPI.UpdateGroupBlobStore(ctx, name).Body(body).Execute()
}

func (s *blobStoreServiceV382) CreateBlobStore1(ctx context.Context, body sonatyperepoV382.AzureBlobStoreApiModel) (*http.Response, error) {
	return s.client.BlobStoreAPI.CreateBlobStore1(ctx).Body(body).Execute()
}

func (s *blobStoreServiceV382) GetBlobStore1(ctx context.Context, name string) (*sonatyperepoV382.AzureBlobStoreApiModel, *http.Response, error) {
	return s.client.BlobStoreAPI.GetBlobStore1(ctx, name).Execute()
}

func (s *blobStoreServiceV382) UpdateBlobStore1(ctx context.Context, name string, body sonatyperepoV382.AzureBlobStoreApiModel) (*http.Response, error) {
	return s.client.BlobStoreAPI.UpdateBlobStore1(ctx, name).Body(body).Execute()
}

func (s *blobStoreServiceV382) CreateBlobStore2(ctx context.Context, body sonatyperepoV382.GoogleCloudBlobstoreApiModel) (*http.Response, error) {
	return s.client.BlobStoreAPI.CreateBlobStore2(ctx).Body(body).Execute()
}

func (s *blobStoreServiceV382) GetBlobStore2(ctx context.Context, name string) (*sonatyperepoV382.GoogleCloudBlobstoreApiModel, *http.Response, error) {
	return s.client.BlobStoreAPI.GetBlobStore2(ctx, name).Execute()
}

func (s *blobStoreServiceV382) UpdateBlobStore2(ctx context.Context, name string, body sonatyperepoV382.GoogleCloudBlobstoreApiModel) (*http.Response, error) {
	return s.client.BlobStoreAPI.UpdateBlobStore2(ctx, name).Body(body).Execute()
}

func (s *blobStoreServiceV382) ListBlobStores(ctx context.Context) ([]sonatyperepoV382.GenericBlobStoreApiResponse, *http.Response, error) {
	return s.client.BlobStoreAPI.ListBlobStores(ctx).Execute()
}

// ============================================================================
// V395 Implementation - Bridges to V395 API via jsonBridge
// ============================================================================

func (s *blobStoreServiceV395) CreateS3BlobStore(ctx context.Context, body sonatyperepoV382.S3BlobStoreApiModel) (*http.Response, error) {
	var v395Body sonatyperepoV395.S3BlobStoreApiModel
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateS3BlobStore(ctx).S3BlobStoreApiModel(v395Body).Execute()
}

func (s *blobStoreServiceV395) GetS3BlobStore(ctx context.Context, name string) (*sonatyperepoV382.S3BlobStoreApiModel, *http.Response, error) {
	v395Response, httpResponse, err := s.client.BlobStoreAPI.GetS3BlobStore(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.S3BlobStoreApiModel
	if err := jsonBridge(v395Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV395) UpdateS3BlobStore(ctx context.Context, name string, body sonatyperepoV382.S3BlobStoreApiModel) (*http.Response, error) {
	var v395Body sonatyperepoV395.S3BlobStoreApiModel
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateS3BlobStore(ctx, name).S3BlobStoreApiModel(v395Body).Execute()
}

func (s *blobStoreServiceV395) CreateFileBlobStore(ctx context.Context, body sonatyperepoV382.FileBlobStoreApiCreateRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.FileBlobStoreApiCreateRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateBlobstoresFile(ctx).FileBlobStoreApiCreateRequest(v395Body).Execute()
}

func (s *blobStoreServiceV395) GetFileBlobStoreConfiguration(ctx context.Context, name string) (*sonatyperepoV382.FileBlobStoreApiModel, *http.Response, error) {
	v395Response, httpResponse, err := s.client.BlobStoreAPI.GetBlobstoresFile(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.FileBlobStoreApiModel
	if err := jsonBridge(v395Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV395) UpdateFileBlobStore(ctx context.Context, name string, body sonatyperepoV382.FileBlobStoreApiUpdateRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.FileBlobStoreApiUpdateRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateBlobstoresFile(ctx, name).FileBlobStoreApiUpdateRequest(v395Body).Execute()
}

func (s *blobStoreServiceV395) CreateGroupBlobStore(ctx context.Context, body sonatyperepoV382.GroupBlobStoreApiCreateRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.GroupBlobStoreApiCreateRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateBlobstoresGroup(ctx).GroupBlobStoreApiCreateRequest(v395Body).Execute()
}

func (s *blobStoreServiceV395) GetGroupBlobStoreConfiguration(ctx context.Context, name string) (*sonatyperepoV382.GroupBlobStoreApiModel, *http.Response, error) {
	v395Response, httpResponse, err := s.client.BlobStoreAPI.GetBlobstoresGroup(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.GroupBlobStoreApiModel
	if err := jsonBridge(v395Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV395) UpdateGroupBlobStore(ctx context.Context, name string, body sonatyperepoV382.GroupBlobStoreApiUpdateRequest) (*http.Response, error) {
	var v395Body sonatyperepoV395.GroupBlobStoreApiUpdateRequest
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateBlobstoresGroup(ctx, name).GroupBlobStoreApiUpdateRequest(v395Body).Execute()
}

func (s *blobStoreServiceV395) CreateBlobStore1(ctx context.Context, body sonatyperepoV382.AzureBlobStoreApiModel) (*http.Response, error) {
	var v395Body sonatyperepoV395.AzureBlobStoreApiModel
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateBlobstoresAzure(ctx).AzureBlobStoreApiModel(v395Body).Execute()
}

func (s *blobStoreServiceV395) GetBlobStore1(ctx context.Context, name string) (*sonatyperepoV382.AzureBlobStoreApiModel, *http.Response, error) {
	v395Response, httpResponse, err := s.client.BlobStoreAPI.GetBlobstoresAzure(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AzureBlobStoreApiModel
	if err := jsonBridge(v395Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV395) UpdateBlobStore1(ctx context.Context, name string, body sonatyperepoV382.AzureBlobStoreApiModel) (*http.Response, error) {
	var v395Body sonatyperepoV395.AzureBlobStoreApiModel
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateBlobstoresAzure(ctx, name).AzureBlobStoreApiModel(v395Body).Execute()
}

func (s *blobStoreServiceV395) CreateBlobStore2(ctx context.Context, body sonatyperepoV382.GoogleCloudBlobstoreApiModel) (*http.Response, error) {
	var v395Body sonatyperepoV395.GoogleCloudBlobstoreApiModel
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateBlobstoresGoogle(ctx).GoogleCloudBlobstoreApiModel(v395Body).Execute()
}

func (s *blobStoreServiceV395) GetBlobStore2(ctx context.Context, name string) (*sonatyperepoV382.GoogleCloudBlobstoreApiModel, *http.Response, error) {
	v395Response, httpResponse, err := s.client.BlobStoreAPI.GetBlobstoresGoogle(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.GoogleCloudBlobstoreApiModel
	if err := jsonBridge(v395Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV395) UpdateBlobStore2(ctx context.Context, name string, body sonatyperepoV382.GoogleCloudBlobstoreApiModel) (*http.Response, error) {
	var v395Body sonatyperepoV395.GoogleCloudBlobstoreApiModel
	if err := jsonBridge(body, &v395Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateBlobstoresGoogle(ctx, name).GoogleCloudBlobstoreApiModel(v395Body).Execute()
}

func (s *blobStoreServiceV395) ListBlobStores(ctx context.Context) ([]sonatyperepoV382.GenericBlobStoreApiResponse, *http.Response, error) {
	v395Response, httpResponse, err := s.client.BlobStoreAPI.ListBlobstores(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.GenericBlobStoreApiResponse
	if err := jsonBridge(v395Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
