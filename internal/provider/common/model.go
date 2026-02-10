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
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sonatyperepo "github.com/sonatype-nexus-community/nexus-repo-api-client-go/v3"
	sharederr "github.com/sonatype-nexus-community/terraform-provider-shared/errors"
)

type SonatypeDataSourceData struct {
	Auth         sonatyperepo.BasicAuth
	BaseUrl      string
	Client       *sonatyperepo.APIClient
	NxrmVersion  SystemVersion
	NxrmWritable bool
}

func (p *SonatypeDataSourceData) CheckWritableAndGetVersion(ctx context.Context, respDiags *diag.Diagnostics, versionHint *string) {
	httpResponse, err := p.Client.StatusAPI.IsWritable(ctx).Execute()
	if err != nil {
		sharederr.HandleAPIError(
			"Sonatype Nexus Repository is not writable or contactable",
			&err,
			httpResponse,
			respDiags,
		)
		return
	}

	if httpResponse.StatusCode == http.StatusOK {
		p.NxrmWritable = true
		var nxrmVersion SystemVersion

		if versionHint != nil {
			nxrmVersion = ParseServerHeaderToVersion(*versionHint)
			tflog.Debug(ctx, fmt.Sprintf("Version Hint: %s", nxrmVersion.String()))
		} else {
			nxrmVersion = ParseServerHeaderToVersion(httpResponse.Header.Get("server"))
			tflog.Debug(ctx, fmt.Sprintf("Server Header: %s", nxrmVersion.String()))
		}
		p.NxrmVersion = nxrmVersion
	}

	tflog.Info(ctx, fmt.Sprintf("Determined Sonatype Nexus Repository to be version %s", p.NxrmVersion.String()))
}

func GetStringAsInt8(s string) *int8 {
	i64, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		return nil
	}
	i8 := int8(i64)
	return &i8
}

func FindAllGroups(re *regexp.Regexp, s string) map[string]string {
	matches := re.FindStringSubmatch(s)
	subnames := re.SubexpNames()
	if matches == nil || len(matches) != len(subnames) {
		return nil
	}

	matchMap := map[string]string{}
	for i := 1; i < len(matches); i++ {
		matchMap[subnames[i]] = matches[i]
	}
	return matchMap
}

var nxrmServerVersionExp = regexp.MustCompile(`^NEXUS\/(?P<MAJOR>\d+)\.(?P<MINOR>\d+)\.(?P<PATCH>\d+)\-(?P<BUILD>\d+)\s+\((?P<EDITION>\w+)\)$`)

type SystemVersion struct {
	Major      int8
	Minor      int8
	Patch      int8
	Build      int8
	ProVersion bool
}

func (s *SystemVersion) NewerThan(major, minor, patch, build int) bool {
	if s.Major > int8(major) {
		return true
	} else if s.Major == int8(major) {
		if s.Minor > int8(minor) {
			return true
		} else if s.Minor == int8(minor) {
			if s.Patch > int8(patch) {
				return true
			} else if s.Patch == int8(patch) {
				return s.Build > int8(build)
			}
		}
	}
	return false
}

func (s *SystemVersion) OlderThan(major, minor, patch, build int) bool {
	return !s.NewerThan(major, minor, patch, build)
}

func (s *SystemVersion) SemVerString() string {
	return fmt.Sprintf("%d.%d.%d-%d", s.Major, s.Minor, s.Patch, s.Build)
}

func (s *SystemVersion) String() string {
	return fmt.Sprintf("%d.%d.%d-%d (PRO=%t)", s.Major, s.Minor, s.Patch, s.Build, s.ProVersion)
}

func (s *SystemVersion) RequiresLowerCaseRepostioryNameDocker() bool {
	return s.NewerThan(3, 89, 0, 0)
}

func (s *SystemVersion) SupportsCapabilities() bool {
	return s.NewerThan(3, 84, 0, 0)
}

func ParseServerHeaderToVersion(headerStr string) SystemVersion {
	match := FindAllGroups(nxrmServerVersionExp, strings.ToUpper(headerStr))
	sysVersion := SystemVersion{}
	for k, v := range match {
		switch k {
		case "MAJOR":
			sysVersion.Major = *GetStringAsInt8(v)
		case "MINOR":
			sysVersion.Minor = *GetStringAsInt8(v)
		case "PATCH":
			sysVersion.Patch = *GetStringAsInt8(v)
		case "BUILD":
			sysVersion.Build = *GetStringAsInt8(v)
		case "EDITION":
			if strings.TrimSpace(v) == "PRO" {
				sysVersion.ProVersion = true
			} else {
				sysVersion.ProVersion = false
			}
		}
	}
	return sysVersion
}
