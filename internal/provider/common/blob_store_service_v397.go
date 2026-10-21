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
	sonatyperepoV397 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v397"
)

// blobStoreServiceV397 implements BlobStoreService against NXRM API client V397 (targets NXRM 3.97.0+).
type blobStoreServiceV397 struct {
	client *sonatyperepoV397.APIClient
}

// V397 Implementation - Bridges to V397 API via jsonBridge
func (s *blobStoreServiceV397) CreateS3BlobStore(ctx context.Context, body sonatyperepoV382.S3BlobStoreApiModel) (*http.Response, error) {
	var v397Body sonatyperepoV397.S3BlobStoreApiModel
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateS3BlobStore(ctx).S3BlobStoreApiModel(v397Body).Execute()
}

func (s *blobStoreServiceV397) GetS3BlobStore(ctx context.Context, name string) (*sonatyperepoV382.S3BlobStoreApiModel, *http.Response, error) {
	v397Response, httpResponse, err := s.client.BlobStoreAPI.GetS3BlobStore(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.S3BlobStoreApiModel
	if err := jsonBridge(v397Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV397) UpdateS3BlobStore(ctx context.Context, name string, body sonatyperepoV382.S3BlobStoreApiModel) (*http.Response, error) {
	var v397Body sonatyperepoV397.S3BlobStoreApiModel
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateS3BlobStore(ctx, name).S3BlobStoreApiModel(v397Body).Execute()
}

func (s *blobStoreServiceV397) CreateFileBlobStore(ctx context.Context, body sonatyperepoV382.FileBlobStoreApiCreateRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.FileBlobStoreApiCreateRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateBlobstoresFile(ctx).FileBlobStoreApiCreateRequest(v397Body).Execute()
}

func (s *blobStoreServiceV397) GetFileBlobStoreConfiguration(ctx context.Context, name string) (*sonatyperepoV382.FileBlobStoreApiModel, *http.Response, error) {
	v397Response, httpResponse, err := s.client.BlobStoreAPI.GetBlobstoresFile(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.FileBlobStoreApiModel
	if err := jsonBridge(v397Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV397) UpdateFileBlobStore(ctx context.Context, name string, body sonatyperepoV382.FileBlobStoreApiUpdateRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.FileBlobStoreApiUpdateRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateBlobstoresFile(ctx, name).FileBlobStoreApiUpdateRequest(v397Body).Execute()
}

func (s *blobStoreServiceV397) CreateGroupBlobStore(ctx context.Context, body sonatyperepoV382.GroupBlobStoreApiCreateRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.GroupBlobStoreApiCreateRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateBlobstoresGroup(ctx).GroupBlobStoreApiCreateRequest(v397Body).Execute()
}

func (s *blobStoreServiceV397) GetGroupBlobStoreConfiguration(ctx context.Context, name string) (*sonatyperepoV382.GroupBlobStoreApiModel, *http.Response, error) {
	v397Response, httpResponse, err := s.client.BlobStoreAPI.GetBlobstoresGroup(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.GroupBlobStoreApiModel
	if err := jsonBridge(v397Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV397) UpdateGroupBlobStore(ctx context.Context, name string, body sonatyperepoV382.GroupBlobStoreApiUpdateRequest) (*http.Response, error) {
	var v397Body sonatyperepoV397.GroupBlobStoreApiUpdateRequest
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateBlobstoresGroup(ctx, name).GroupBlobStoreApiUpdateRequest(v397Body).Execute()
}

func (s *blobStoreServiceV397) CreateBlobStore1(ctx context.Context, body sonatyperepoV382.AzureBlobStoreApiModel) (*http.Response, error) {
	var v397Body sonatyperepoV397.AzureBlobStoreApiModel
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateBlobstoresAzure(ctx).AzureBlobStoreApiModel(v397Body).Execute()
}

func (s *blobStoreServiceV397) GetBlobStore1(ctx context.Context, name string) (*sonatyperepoV382.AzureBlobStoreApiModel, *http.Response, error) {
	v397Response, httpResponse, err := s.client.BlobStoreAPI.GetBlobstoresAzure(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.AzureBlobStoreApiModel
	if err := jsonBridge(v397Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV397) UpdateBlobStore1(ctx context.Context, name string, body sonatyperepoV382.AzureBlobStoreApiModel) (*http.Response, error) {
	var v397Body sonatyperepoV397.AzureBlobStoreApiModel
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateBlobstoresAzure(ctx, name).AzureBlobStoreApiModel(v397Body).Execute()
}

func (s *blobStoreServiceV397) CreateBlobStore2(ctx context.Context, body sonatyperepoV382.GoogleCloudBlobstoreApiModel) (*http.Response, error) {
	var v397Body sonatyperepoV397.GoogleCloudBlobstoreApiModel
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.CreateBlobstoresGoogle(ctx).GoogleCloudBlobstoreApiModel(v397Body).Execute()
}

func (s *blobStoreServiceV397) GetBlobStore2(ctx context.Context, name string) (*sonatyperepoV382.GoogleCloudBlobstoreApiModel, *http.Response, error) {
	v397Response, httpResponse, err := s.client.BlobStoreAPI.GetBlobstoresGoogle(ctx, name).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result sonatyperepoV382.GoogleCloudBlobstoreApiModel
	if err := jsonBridge(v397Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return &result, httpResponse, nil
}

func (s *blobStoreServiceV397) UpdateBlobStore2(ctx context.Context, name string, body sonatyperepoV382.GoogleCloudBlobstoreApiModel) (*http.Response, error) {
	var v397Body sonatyperepoV397.GoogleCloudBlobstoreApiModel
	if err := jsonBridge(body, &v397Body); err != nil {
		return nil, err
	}
	return s.client.BlobStoreAPI.UpdateBlobstoresGoogle(ctx, name).GoogleCloudBlobstoreApiModel(v397Body).Execute()
}

func (s *blobStoreServiceV397) ListBlobStores(ctx context.Context) ([]sonatyperepoV382.GenericBlobStoreApiResponse, *http.Response, error) {
	v397Response, httpResponse, err := s.client.BlobStoreAPI.ListBlobstores(ctx).Execute()
	if err != nil {
		return nil, httpResponse, err
	}
	var result []sonatyperepoV382.GenericBlobStoreApiResponse
	if err := jsonBridge(v397Response, &result); err != nil {
		return nil, httpResponse, err
	}
	return result, httpResponse, nil
}
