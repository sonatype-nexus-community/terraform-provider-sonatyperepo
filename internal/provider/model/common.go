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

package model

import (
	"fmt"
	"reflect"
	"strconv"
)

// StructToMap converts any struct to a map[string]string using reflection
// It uses the "tfsdk" tag to determine key names and handles type conversions
func StructToMap(v interface{}) *map[string]string {
	result := make(map[string]string)

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return &result
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Get the tfsdk tag, fallback to field name
		key := field.Tag.Get("tfsdk")
		if key == "" {
			key = field.Name
		}

		// Skip empty tags marked with "-"
		if key == "-" {
			continue
		}

		// Convert field value to string
		result[key] = fieldValueToString(fieldVal)
	}

	return &result
}

// fieldValueToString converts a reflect.Value to its string representation
func fieldValueToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Ptr:
		if v.IsNil() {
			return ""
		}
		return fieldValueToString(v.Elem())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
