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
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// jsonBridge re-shapes a struct from one NXRM API client generation into the equivalent
// struct of another generation via a JSON round-trip. This is safe because every field
// that exists in both generations shares the same JSON tag; fields that only exist on one
// side are silently dropped (write) or left at their zero value (read). Known same-named-
// field renames (e.g. RoutingRule -> RoutingRuleName) must be patched explicitly by the
// caller after this bridge runs. Used by every domain service adapter that canonicalizes
// V395 responses/requests into the V382 generated shape the model-mapping layer expects.
//
// Generated structs commonly have a custom UnmarshalJSON that calls decoder.DisallowUnknownFields(),
// so before the final unmarshal into dst, every JSON object in the marshaled src is pruned
// (recursively, at every nesting level) down to the field names dst's type actually declares --
// otherwise a field that only exists on the src side (e.g. a newer generation added a field)
// would fail the whole decode with "json: unknown field ...".
func jsonBridge(src, dst any) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	pruned, err := pruneUnknownJSONFields(b, reflect.TypeOf(dst))
	if err != nil {
		return err
	}
	return json.Unmarshal(pruned, dst)
}

// bridgeFromResponse populates dst (a pointer) from apiSrc, the already-decoded response of the
// source generation's own generated client. If the source generation's client itself failed to
// decode the HTTP response (err != nil) -- which happens when the live server returns a field
// newer than the vendored client's generated struct declares, a version-skew gap in the client
// module rather than a real request failure -- and the HTTP call itself succeeded (a response
// object exists with a non-error status), this falls back to bridging directly from the raw
// response body instead of giving up. Any other error (network failure, 4xx/5xx status) is
// returned unchanged.
func bridgeFromResponse(apiSrc any, httpResponse *http.Response, err error, dst any) error {
	if err == nil {
		return jsonBridge(apiSrc, dst)
	}
	if httpResponse == nil || httpResponse.StatusCode >= 300 || httpResponse.Body == nil {
		return err
	}

	body, readErr := io.ReadAll(httpResponse.Body)
	_ = httpResponse.Body.Close()
	if readErr != nil {
		return err
	}
	httpResponse.Body = io.NopCloser(bytes.NewReader(body))

	pruned, pruneErr := pruneUnknownJSONFields(body, reflect.TypeOf(dst))
	if pruneErr != nil {
		return err
	}
	if json.Unmarshal(pruned, dst) != nil {
		return err
	}
	return nil
}

// pruneUnknownJSONFields recursively drops object keys that have no corresponding field
// (by JSON tag) on t, descending into nested objects/arrays along the way. Types it can't
// reason about (maps, interfaces, scalars) are passed through unchanged.
func pruneUnknownJSONFields(raw json.RawMessage, t reflect.Type) (json.RawMessage, error) {
	t = indirect(t)
	if t == nil {
		return raw, nil
	}

	// A JSON null must stay null -- unmarshaling it into a map (the struct case below)
	// succeeds silently with a nil map, which would otherwise fall through to backfilling
	// every known key with its zero value and fabricate a non-null object in its place.
	if trimmed := bytes.TrimSpace(raw); string(trimmed) == "null" {
		return raw, nil
	}

	if elemType, ok := nullableElemType(t); ok {
		return pruneUnknownJSONFields(raw, elemType)
	}

	switch t.Kind() {
	case reflect.Struct:
		return pruneStructFields(raw, t)
	case reflect.Slice, reflect.Array:
		return pruneSliceElements(raw, t.Elem())
	default:
		return raw, nil
	}
}

// indirect strips pointer indirection, following it all the way down to the underlying type.
func indirect(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// pruneStructFields drops object keys that have no corresponding field (by JSON tag) on t,
// recursing into each retained field's value, and backfills declared fields the object is
// missing. Not a JSON object (e.g. null, or a mismatched shape) is left for the real unmarshal
// to accept or reject.
func pruneStructFields(raw json.RawMessage, t reflect.Type) (json.RawMessage, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return raw, nil
	}
	known := jsonFieldTypes(t)
	for key, val := range obj {
		fieldType, ok := known[key]
		if !ok {
			delete(obj, key)
			continue
		}
		prunedVal, err := pruneUnknownJSONFields(val, fieldType)
		if err != nil {
			return nil, err
		}
		obj[key] = prunedVal
	}
	// Some generated structs manually check that every declared field is present as a
	// JSON key before decoding, even when the field's Go type has a perfectly usable
	// zero value (e.g. a bool). If the source generation dropped a field the destination
	// still declares -- a version-skew removal, not just an encoding difference -- that
	// check fails the whole decode unless we backfill the zero value ourselves.
	for key, fieldType := range known {
		if _, exists := obj[key]; exists {
			continue
		}
		if zero, ok := zeroJSONValue(fieldType); ok {
			obj[key] = zero
		}
	}
	return json.Marshal(obj)
}

// pruneSliceElements applies pruneUnknownJSONFields to every element of a JSON array.
func pruneSliceElements(raw json.RawMessage, elemType reflect.Type) (json.RawMessage, error) {
	var arr []json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		return raw, nil
	}
	for i, v := range arr {
		prunedVal, err := pruneUnknownJSONFields(v, elemType)
		if err != nil {
			return nil, err
		}
		arr[i] = prunedVal
	}
	return json.Marshal(arr)
}

// nullableElemType detects the generated client's "NullableXxx" wrapper shape --
// an unexported `value *Xxx` field alongside an unexported `isSet bool` field -- used
// throughout both generations for optional/nullable properties. Its custom MarshalJSON/
// UnmarshalJSON serialize as Xxx directly (or null), not as the wrapper's own two-field
// Go layout, so pruning must recurse on Xxx rather than reflecting into "value"/"isSet"
// as if they were the real JSON keys -- otherwise a struct-valued Xxx (e.g. ProxySettingsXo)
// gets pruned down to nothing, dropping every real field it has.
func nullableElemType(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() != reflect.Struct || t.NumField() != 2 {
		return nil, false
	}
	valueField, ok := t.FieldByName("value")
	if !ok || valueField.Type.Kind() != reflect.Ptr {
		return nil, false
	}
	isSetField, ok := t.FieldByName("isSet")
	if !ok || isSetField.Type.Kind() != reflect.Bool {
		return nil, false
	}
	return valueField.Type.Elem(), true
}

// zeroJSONValue returns the JSON encoding of t's zero value, for the scalar kinds where
// backfilling a missing field is unambiguous and safe. Pointer, struct, slice, and map fields
// are deliberately left alone -- their "missing" state is either already valid (nil/omitted)
// or would require recursively synthesizing a nested zero object, which risks masking a real
// missing-data problem instead of a version-skew one.
func zeroJSONValue(t reflect.Type) (json.RawMessage, bool) {
	switch t.Kind() {
	case reflect.Bool:
		return json.RawMessage("false"), true
	case reflect.String:
		return json.RawMessage(`""`), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return json.RawMessage("0"), true
	default:
		return nil, false
	}
}

// jsonFieldTypes maps every JSON key t's fields would decode from to that field's type,
// merging in the fields of any anonymous (embedded) struct fields.
func jsonFieldTypes(t reflect.Type) map[string]reflect.Type {
	fields := make(map[string]reflect.Type)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			mergeEmbeddedFieldTypes(fields, f.Type)
			continue
		}
		if name, ok := jsonFieldName(f); ok {
			fields[name] = f.Type
		}
	}
	return fields
}

// mergeEmbeddedFieldTypes merges the JSON field types of an anonymous (embedded) struct field
// into fields, so promoted fields are addressable by their own JSON key.
func mergeEmbeddedFieldTypes(fields map[string]reflect.Type, embedded reflect.Type) {
	embedded = indirect(embedded)
	if embedded == nil || embedded.Kind() != reflect.Struct {
		return
	}
	for k, v := range jsonFieldTypes(embedded) {
		fields[k] = v
	}
}

// jsonFieldName returns the JSON key f would decode from, and false if it's excluded via
// a `json:"-"` tag.
func jsonFieldName(f reflect.StructField) (string, bool) {
	tag := f.Tag.Get("json")
	if tag == "-" {
		return "", false
	}
	if tag == "" {
		return f.Name, true
	}
	if idx := strings.Index(tag, ","); idx >= 0 {
		if tag[:idx] != "" {
			return tag[:idx], true
		}
		return f.Name, true
	}
	return tag, true
}
