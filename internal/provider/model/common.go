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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ParseBool(value string, defaultValue bool) bool {
	val, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return val
}

func ParseInt32(value string, defaultValue int32) int32 {
	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return defaultValue
	}

	return int32(val)
}

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

		// Get the nxrm tag, fallback to field name
		key := field.Tag.Get("nxrm")
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
	case reflect.Struct:
		// Handle types.Set
		if v.CanInterface() {
			if set, ok := v.Interface().(types.Set); ok {
				elements := set.Elements()
				var strs []string
				for _, elem := range elements {
					if strVal, ok := elem.(types.String); ok {
						strs = append(strs, strVal.ValueString())
					}
				}
				return strings.Join(strs, ",")
			}
		}

		// Handle Terraform types (types.String, types.Int32, etc.)
		if v.CanInterface() {
			// Check for ValueString() method (types.String)
			if method := v.MethodByName("ValueString"); method.IsValid() {
				result := method.Call(nil)
				if len(result) > 0 {
					return result[0].String()
				}
			}

			// Check for ValueInt32() method (types.Int32)
			if method := v.MethodByName("ValueInt32"); method.IsValid() {
				result := method.Call(nil)
				if len(result) > 0 {
					return strconv.FormatInt(result[0].Int(), 10)
				}
			}

			// Check for ValueInt64() method (types.Int64)
			if method := v.MethodByName("ValueInt64"); method.IsValid() {
				result := method.Call(nil)
				if len(result) > 0 {
					return strconv.FormatInt(result[0].Int(), 10)
				}
			}

			// Check for ValueBool() method (types.Bool)
			if method := v.MethodByName("ValueBool"); method.IsValid() {
				result := method.Call(nil)
				if len(result) > 0 {
					return strconv.FormatBool(result[0].Bool())
				}
			}

			// Check for ValueFloat64() method (types.Float64)
			if method := v.MethodByName("ValueFloat64"); method.IsValid() {
				result := method.Call(nil)
				if len(result) > 0 {
					return strconv.FormatFloat(result[0].Float(), 'f', -1, 64)
				}
			}
		}
		return fmt.Sprintf("%v", v.Interface())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
