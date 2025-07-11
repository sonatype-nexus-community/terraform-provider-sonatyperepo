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

package utils_test

import (
	"terraform-provider-sonatyperepo/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// ProviderConfig is a shared configuration to combine with the actual test configuration.
	ProviderConfig = `
provider "sonatyperepo" {
  username = ""
  password = ""
  url     = ""
}
`
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sonatyperepo": providerserver.NewProtocol6WithError(provider.New("test")()),
}

// func testAccPreCheck(t *testing.T) {
// 	// You can add code here to run prior to any test case execution, for example assertions
// 	// about the appropriate environment variables being set are common to see in a pre-check
// 	// function.
// }
