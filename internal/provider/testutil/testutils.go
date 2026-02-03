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
	"fmt"
	"os"
	"strings"
	"terraform-provider-sonatyperepo/internal/provider/common"
	"testing"

	semver "github.com/hashicorp/go-version"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var CurrenTestNxrmVersion = common.ParseServerHeaderToVersion(fmt.Sprintf("Nexus/%s (PRO)", os.Getenv("NXRM_VERSION")))

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

// PatchHCL parses HCL, finds the block/attribute at the path, and updates it.
func PatchHCL(hclStr string, path string, newValue string) (string, error) {
	f, diags := hclwrite.ParseConfig([]byte(hclStr), "main.tf", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return "", fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	parts := strings.Split(path, ".")
	currentBody := f.Body()

	// 1. Walk down the Blocks
	i := 0
	for i < len(parts) {
		foundBlock := false
		for n := 0; n <= 2 && i+n < len(parts); n++ {
			blockType := parts[i]
			labels := parts[i+1 : i+1+n]

			block := currentBody.FirstMatchingBlock(blockType, labels)
			if block != nil {
				currentBody = block.Body()
				i += 1 + n
				foundBlock = true
				break
			}
		}
		if !foundBlock {
			break
		}
	}

	// 2. Find the Attribute
	if i >= len(parts) {
		return "", fmt.Errorf("path resolved to a Block, but expected an Attribute")
	}

	attrName := parts[i]
	attr := currentBody.GetAttribute(attrName)
	if attr == nil {
		return "", fmt.Errorf("attribute '%s' not found", attrName)
	}

	// 3. Patch logic
	remainingPath := parts[i+1:]

	if len(remainingPath) == 0 {
		// Top level attribute patch
		currentBody.SetAttributeValue(attrName, cty.StringVal(newValue))
	} else {
		// Nested Map patch (Token Surgery)
		newTokens, err := patchMapTokens(attr.Expr().BuildTokens(nil), remainingPath, newValue)
		if err != nil {
			return "", err
		}
		currentBody.SetAttributeRaw(attrName, newTokens)
	}

	return string(f.Bytes()), nil
}

// RemoveHCL removes a block or attribute at the specified path.
func RemoveHCL(hclStr string, path string) (string, error) {
	f, diags := hclwrite.ParseConfig([]byte(hclStr), "main.tf", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return "", fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	parts := strings.Split(path, ".")
	currentBody := f.Body()

	// Navigate to the PARENT of the item we want to remove
	// We iterate up to len(parts)-1 because the last part is the item to delete
	i := 0
	for i < len(parts)-1 {
		foundBlock := false
		// Look ahead 0, 1, or 2 steps to match block labels
		for n := 0; n <= 2 && i+n < len(parts)-1; n++ {
			blockType := parts[i]
			labels := parts[i+1 : i+1+n]

			block := currentBody.FirstMatchingBlock(blockType, labels)
			if block != nil {
				currentBody = block.Body()
				i += 1 + n
				foundBlock = true
				break
			}
		}
		if !foundBlock {
			// If we can't find a block, the path is invalid or refers to a nested map attribute
			// which this function does not support deep deletion for (e.g. proxy.remote_url).
			return "", fmt.Errorf("could not resolve parent block at '%s'", parts[i])
		}
	}

	targetName := parts[len(parts)-1]

	// 1. Try to remove it as an Attribute
	// (e.g. `online = true` or `storage = { ... }`)
	attr := currentBody.GetAttribute(targetName)
	if attr != nil {
		currentBody.RemoveAttribute(targetName)
		return string(f.Bytes()), nil
	}

	// 2. Try to remove it as a Block
	// (e.g. `storage { ... }`)
	// blocks := currentBody.Blocks() // This returns all blocks
	// We need to match specifically. Since we don't have the labels for the target
	// (the path just ends in "storage"), we assume the user targets a block type with 0 labels
	// or specific name logic.

	// Scan for blocks with this type (e.g. "storage")
	foundBlock := false
	for _, b := range currentBody.Blocks() {
		if b.Type() == targetName {
			currentBody.RemoveBlock(b)
			foundBlock = true
			// If you want to remove ALL blocks of this type, don't break.
			// If you only want the first, break. Usually, HCL removal implies specific targeting,
			// but for named blocks (like `lifecycle`), removing all matching types is standard.
		}
	}

	if foundBlock {
		return string(f.Bytes()), nil
	}

	return "", fmt.Errorf("item '%s' not found (neither attribute nor block)", targetName)
}

func patchMapTokens(tokens hclwrite.Tokens, path []string, newValue string) (hclwrite.Tokens, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path in token patcher")
	}
	targetKey := path[0]

	var newTokens hclwrite.Tokens
	foundKey := false
	replaced := false

	for j := 0; j < len(tokens); j++ {
		tok := tokens[j]
		newTokens = append(newTokens, tok)

		// 1. Find the Key
		if !foundKey && tok.Type == hclsyntax.TokenIdent && string(tok.Bytes) == targetKey {
			foundKey = true
			continue
		}

		// 2. Find and Replace the Value
		if foundKey && !replaced {
			// Skip equals, colons, and whitespace to find the start of the value
			if tok.Type == hclsyntax.TokenEqual ||
				tok.Type == hclsyntax.TokenColon ||
				isWhitespace(tok.Bytes) {
				continue
			}

			// We found the value start!
			// Remove the token we just added (the old value start)
			newTokens = newTokens[:len(newTokens)-1]

			// If the old value was a string, it started with OQuote and we need to skip until CQuote.
			// If it was a number/bool (Ident), it's just one token.
			// For simplicity in this specific "string-to-string" patch scenario:
			if tok.Type == hclsyntax.TokenOQuote {
				// We need to skip tokens until we find the closing quote in the ORIGINAL stream
				// to ensure we remove the entire old string.
				for k := j + 1; k < len(tokens); k++ {
					if tokens[k].Type == hclsyntax.TokenCQuote {
						j = k // Fast forward the main loop
						break
					}
				}
			}

			// Inject new string value
			newTokens = append(newTokens, &hclwrite.Token{
				Type:  hclsyntax.TokenOQuote,
				Bytes: []byte(`"`),
			})
			newTokens = append(newTokens, &hclwrite.Token{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte(newValue),
			})
			newTokens = append(newTokens, &hclwrite.Token{
				Type:  hclsyntax.TokenCQuote,
				Bytes: []byte(`"`),
			})

			replaced = true
		}
	}

	if !foundKey {
		return nil, fmt.Errorf("key '%s' not found in attribute map", targetKey)
	}

	return newTokens, nil
}

// Helper to identify whitespace tokens blindly without relying on token types
func isWhitespace(b []byte) bool {
	s := string(b)
	return strings.TrimSpace(s) == ""
}
