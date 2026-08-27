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

package testutil

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"

	semver "github.com/hashicorp/go-version"
	v3 "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
)

var CurrenTestNxrmVersion = common.ParseServerHeaderToVersion(fmt.Sprintf("Nexus/%s (PRO)", os.Getenv("NXRM_VERSION")))

// SkipIfNoIqServer skips a test unless TF_ACC_IQ_SERVER=1 is set. This mirrors the
// TF_ACC_S3_BLOB_STORE / TF_ACC_GCS_BLOB_STORE opt-in pattern (see
// internal/provider/blob_store/blob_store_s3_resource_test.go) for tests that need a real,
// licensed Sonatype IQ Server connected to the NXRM instance under test - infrastructure that
// isn't present by default, so a plain local `TF_ACC=1 go test ./...` run skips these instead
// of failing. See https://github.com/sonatype-nexus-community/terraform-provider-sonatyperepo/issues/285.
func SkipIfNoIqServer(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC_IQ_SERVER") != "1" {
		t.Skip("Sonatype IQ Server tests require a live IQ Server connected to Sonatype Nexus Repository Manager - set TF_ACC_IQ_SERVER=1 to enable")
	}
}

// IqServerUrl returns the URL of the Sonatype IQ Server under test, defaulting to the
// well-known local port used by both this provider's CI workflows and the sibling
// terraform-provider-sonatypeiq project.
func IqServerUrl() string {
	if v := os.Getenv("IQ_SERVER_URL"); v != "" {
		return v
	}
	return "http://localhost:8070"
}

// IqServerUsername returns the username to authenticate NXRM's connection to IQ Server with.
func IqServerUsername() string {
	if v := os.Getenv("IQ_SERVER_USERNAME"); v != "" {
		return v
	}
	return "admin"
}

// IqServerPassword returns the password to authenticate NXRM's connection to IQ Server with.
func IqServerPassword() string {
	if v := os.Getenv("IQ_SERVER_PASSWORD"); v != "" {
		return v
	}
	return "admin123"
}

// nxrmV3Client builds a raw NXRM API client from NXRM_SERVER_URL/USERNAME/PASSWORD, bypassing
// the provider entirely - used for imperative test setup, mirroring the pattern already used in
// internal/provider/provider_test.go's TestMain for HA fixture bootstrapping.
func nxrmV3Client() (*v3.APIClient, context.Context) {
	clientConfiguration := v3.NewConfiguration()
	clientConfiguration.Servers = []v3.ServerConfiguration{
		{
			URL:         fmt.Sprintf("%s%s", strings.TrimRight(os.Getenv("NXRM_SERVER_URL"), "/"), "/service/rest"),
			Description: "Sonatype Nexus Repository Server",
		},
	}
	client := v3.NewAPIClient(clientConfiguration)
	ctx := context.WithValue(
		context.Background(),
		v3.ContextBasicAuth,
		v3.BasicAuth{UserName: os.Getenv("NXRM_SERVER_USERNAME"), Password: os.Getenv("NXRM_SERVER_PASSWORD")},
	)
	return client, ctx
}

// ConfigureIqConnection points the NXRM instance under test (NXRM_SERVER_URL) at the Sonatype IQ
// Server described by IqServerUrl/IqServerUsername/IqServerPassword. Used both to bootstrap the
// connection once before any acceptance test runs (TestMain, gated on TF_ACC_IQ_SERVER=1) and to
// restore it after tests that intentionally disable it (e.g. TestAccSystemIqConnectionResource's
// implicit end-of-test destroy), so a disable in one package can't leave it down for others
// running concurrently.
func ConfigureIqConnection() error {
	client, ctx := nxrmV3Client()

	_, httpResponse, err := client.ManageSonatypeRepositoryFirewallConfigurationAPI.UpdateConfiguration1(ctx).Body(v3.IqConnectionXo{
		AuthenticationType:  "USER",
		Enabled:             v3.PtrBool(true),
		FailOpenModeEnabled: v3.PtrBool(false),
		ShowLink:            v3.PtrBool(true),
		Url:                 v3.PtrString(IqServerUrl()),
		Username:            v3.PtrString(IqServerUsername()),
		Password:            v3.PtrString(IqServerPassword()),
	}).Execute()
	if err != nil {
		return err
	}
	if httpResponse != nil && httpResponse.StatusCode != http.StatusNoContent && httpResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response configuring Sonatype IQ Server connection: %s", httpResponse.Status)
	}
	return nil
}

// VerifyIqConnection asks NXRM to verify its currently configured Sonatype IQ Server connection,
// returning whether it succeeded and, on failure, the reason NXRM gave.
func VerifyIqConnection() (success bool, reason string, err error) {
	client, ctx := nxrmV3Client()

	verification, _, err := client.ManageSonatypeRepositoryFirewallConfigurationAPI.VerifyIqConnection(ctx).Execute()
	if err != nil {
		return false, "", err
	}
	if verification == nil || verification.Success == nil {
		return false, "", nil
	}
	if verification.Reason != nil {
		reason = *verification.Reason
	}
	return *verification.Success, reason, nil
}

func SkipIfNxrmVersionEq(t *testing.T, v *common.SystemVersion) {
	t.Helper()

	if v.Major == CurrenTestNxrmVersion.Major && v.Minor == CurrenTestNxrmVersion.Minor && v.Patch == CurrenTestNxrmVersion.Patch {
		t.Skipf("NXRM Version is == %s - skipping", v.String())
	}
}

func SkipIfNxrmVersionInRange(t *testing.T, low *common.SystemVersion, high *common.SystemVersion) {
	t.Helper()

	inRange, err := VersionInRange(&CurrenTestNxrmVersion, low, high)

	if err != nil {
		t.Errorf("Error comparing versions: %v", err)
		t.FailNow()
	}

	if inRange {
		t.Skipf("NXRM Version within range %s and %s - skipping", low.String(), high.String())
	}
}

func VersionInRange(ver *common.SystemVersion, low *common.SystemVersion, high *common.SystemVersion) (bool, error) {
	thisVersion, err := semver.NewVersion(ver.SemVerString())
	if err != nil {
		return false, err
	}

	lowVersion, err := semver.NewVersion(low.SemVerString())
	if err != nil {
		return false, err
	}

	highVersion, err := semver.NewVersion(high.SemVerString())
	if err != nil {
		return false, err
	}

	if lowVersion.LessThanOrEqual(thisVersion) && highVersion.GreaterThanOrEqual(thisVersion) {
		return true, nil
	}

	return false, nil
}
