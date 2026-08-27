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
	"path/filepath"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"
	"time"

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

// iqConnectionLockPath returns a fixed path shared by every package's test binary within a
// single `go test ./...` invocation (each package compiles to a separate binary and can run
// concurrently, so an in-process mutex can't coordinate across them). It deliberately avoids
// os.TempDir() (e.g. /tmp): that directory is world-writable, and SonarQube rightly flags any
// predictable path built from it as CWE-377/379 - another local user could plant or pre-create
// that exact path. os.UserCacheDir() is per-user by definition (~/.cache, ~/Library/Caches,
// %LocalAppData%, ...), so a subdirectory under it stays deterministic and shared for every test
// binary run by the same user - which the lock depends on - without that exposure. The 0700
// subdirectory this process creates is still applied on top as a second layer of defense.
func iqConnectionLockPath() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve user cache directory: %w", err)
	}
	dir := filepath.Join(base, "terraform-provider-sonatyperepo")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "iq-connection.lock"), nil
}

// iqConnectionLockStaleAfter bounds how long a lock file is honored before it's treated as
// abandoned (e.g. left behind by a crashed local test run) and removed so a new lock can be
// acquired. CI runners are ephemeral and never see a stale lock; this only protects local runs.
const iqConnectionLockStaleAfter = 5 * time.Minute

// LockIqConnection acquires a simple, cross-process file lock guarding writes to the shared
// sonatyperepo_system_iq_connection singleton, so TestMain's bootstrap
// (internal/provider/provider_test.go) and TestAccSystemIqConnectionResource
// (internal/provider/system) - separate `go test` binaries that `go test ./...` can run
// concurrently - can't interleave their writes to it and produce spurious drift. Blocks
// (polling) until acquired or timeout elapses; the returned func must be called to release.
func LockIqConnection(timeout time.Duration) (unlock func(), err error) {
	path, err := iqConnectionLockPath()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare IQ connection lock directory: %w", err)
	}
	deadline := time.Now().Add(timeout)
	for {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			_ = f.Close()
			return func() { _ = os.Remove(path) }, nil
		}
		if !os.IsExist(err) {
			return nil, err
		}
		if info, statErr := os.Stat(path); statErr == nil && time.Since(info.ModTime()) > iqConnectionLockStaleAfter {
			_ = os.Remove(path)
			continue
		}
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("timed out after %s waiting for IQ connection lock at %s", timeout, path)
		}
		time.Sleep(100 * time.Millisecond)
	}
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
