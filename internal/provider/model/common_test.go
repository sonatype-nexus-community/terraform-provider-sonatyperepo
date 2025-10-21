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

package model_test

import (
	"terraform-provider-sonatyperepo/internal/provider/model"
	"testing"
)

// Test structs
type SimpleStruct struct {
	Name  string `nxrm:"name"`
	Age   int    `nxrm:"age"`
	Email string `nxrm:"email"`
}

type StructWithNoTags struct {
	Field1 string
	Field2 int
}

type StructWithMixedTags struct {
	Tagged   string `nxrm:"custom_name"`
	Untagged string
	Skipped  string `nxrm:"-"`
}

type AllTypesStruct struct {
	StringField  string  `nxrm:"str"`
	IntField     int     `nxrm:"int"`
	Int8Field    int8    `nxrm:"int8"`
	Int16Field   int16   `nxrm:"int16"`
	Int32Field   int32   `nxrm:"int32"`
	Int64Field   int64   `nxrm:"int64"`
	UintField    uint    `nxrm:"uint"`
	Uint8Field   uint8   `nxrm:"uint8"`
	Uint16Field  uint16  `nxrm:"uint16"`
	Uint32Field  uint32  `nxrm:"uint32"`
	Uint64Field  uint64  `nxrm:"uint64"`
	Float32Field float32 `nxrm:"float32"`
	Float64Field float64 `nxrm:"float64"`
	BoolField    bool    `nxrm:"bool"`
}

type StructWithPointers struct {
	StringPtr *string `nxrm:"string_ptr"`
	IntPtr    *int    `nxrm:"int_ptr"`
	BoolPtr   *bool   `nxrm:"bool_ptr"`
}

type NestedStruct struct {
	Simple SimpleStruct `nxrm:"simple"`
	Name   string       `nxrm:"name"`
}

func TestStructToMap_SimpleStruct(t *testing.T) {
	s := SimpleStruct{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	result := model.StructToMap(s)

	expected := map[string]string{
		"name":  "John Doe",
		"age":   "30",
		"email": "john@example.com",
	}

	if !mapsEqual(*result, expected) {
		t.Errorf("Expected %v, got %v", expected, *result)
	}
}

func TestStructToMap_PointerToStruct(t *testing.T) {
	s := &SimpleStruct{
		Name:  "Jane Doe",
		Age:   25,
		Email: "jane@example.com",
	}

	result := model.StructToMap(s)

	expected := map[string]string{
		"name":  "Jane Doe",
		"age":   "25",
		"email": "jane@example.com",
	}

	if !mapsEqual(*result, expected) {
		t.Errorf("Expected %v, got %v", expected, *result)
	}
}

func TestStructToMap_NoTags(t *testing.T) {
	s := StructWithNoTags{
		Field1: "value1",
		Field2: 42,
	}

	result := model.StructToMap(s)

	expected := map[string]string{
		"Field1": "value1",
		"Field2": "42",
	}

	if !mapsEqual(*result, expected) {
		t.Errorf("Expected %v, got %v", expected, *result)
	}
}

func TestStructToMap_MixedTags(t *testing.T) {
	s := StructWithMixedTags{
		Tagged:   "tagged_value",
		Untagged: "untagged_value",
		Skipped:  "should_be_skipped",
	}

	result := model.StructToMap(s)

	if len(*result) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(*result))
	}

	if (*result)["custom_name"] != "tagged_value" {
		t.Errorf("Expected custom_name='tagged_value', got '%s'", (*result)["custom_name"])
	}

	if (*result)["Untagged"] != "untagged_value" {
		t.Errorf("Expected Untagged='untagged_value', got '%s'", (*result)["Untagged"])
	}

	if _, exists := (*result)["Skipped"]; exists {
		t.Error("Skipped field should not be in result")
	}

	if _, exists := (*result)["-"]; exists {
		t.Error("Field with '-' tag should not be in result")
	}
}

func TestStructToMap_AllTypes(t *testing.T) {
	s := AllTypesStruct{
		StringField:  "test",
		IntField:     -42,
		Int8Field:    -8,
		Int16Field:   -16,
		Int32Field:   -32,
		Int64Field:   -64,
		UintField:    42,
		Uint8Field:   8,
		Uint16Field:  16,
		Uint32Field:  32,
		Uint64Field:  64,
		Float32Field: 3.14,
		Float64Field: 2.718281828,
		BoolField:    true,
	}

	result := model.StructToMap(s)

	tests := []struct {
		key      string
		expected string
	}{
		{"str", "test"},
		{"int", "-42"},
		{"int8", "-8"},
		{"int16", "-16"},
		{"int32", "-32"},
		{"int64", "-64"},
		{"uint", "42"},
		{"uint8", "8"},
		{"uint16", "16"},
		{"uint32", "32"},
		{"uint64", "64"},
		{"float32", "3.140000104904175"}, // float32 precision when converted to float64
		{"float64", "2.718281828"},
		{"bool", "true"},
	}

	for _, tt := range tests {
		if (*result)[tt.key] != tt.expected {
			t.Errorf("Expected %s='%s', got '%s'", tt.key, tt.expected, (*result)[tt.key])
		}
	}
}

func TestStructToMap_WithPointers(t *testing.T) {
	str := "pointer_string"
	num := 100
	b := false

	s := StructWithPointers{
		StringPtr: &str,
		IntPtr:    &num,
		BoolPtr:   &b,
	}

	result := model.StructToMap(s)

	expected := map[string]string{
		"string_ptr": "pointer_string",
		"int_ptr":    "100",
		"bool_ptr":   "false",
	}

	if !mapsEqual(*result, expected) {
		t.Errorf("Expected %v, got %v", expected, *result)
	}
}

func TestStructToMap_WithNilPointers(t *testing.T) {
	s := StructWithPointers{
		StringPtr: nil,
		IntPtr:    nil,
		BoolPtr:   nil,
	}

	result := model.StructToMap(s)

	expected := map[string]string{
		"string_ptr": "",
		"int_ptr":    "",
		"bool_ptr":   "",
	}

	if !mapsEqual(*result, expected) {
		t.Errorf("Expected %v, got %v", expected, *result)
	}
}

func TestStructToMap_EmptyStruct(t *testing.T) {
	type EmptyStruct struct{}

	s := EmptyStruct{}
	result := model.StructToMap(s)

	if len(*result) != 0 {
		t.Errorf("Expected empty map, got %v", *result)
	}
}

func TestStructToMap_ZeroValues(t *testing.T) {
	s := SimpleStruct{
		Name:  "",
		Age:   0,
		Email: "",
	}

	result := model.StructToMap(s)

	expected := map[string]string{
		"name":  "",
		"age":   "0",
		"email": "",
	}

	if !mapsEqual(*result, expected) {
		t.Errorf("Expected %v, got %v", expected, *result)
	}
}

func TestStructToMap_NonStructInput(t *testing.T) {
	// Test with string
	result := model.StructToMap("not a struct")
	if len(*result) != 0 {
		t.Errorf("Expected empty map for string input, got %v", *result)
	}

	// Test with int
	result = model.StructToMap(42)
	if len(*result) != 0 {
		t.Errorf("Expected empty map for int input, got %v", *result)
	}

	// Test with nil
	result = model.StructToMap(nil)
	if len(*result) != 0 {
		t.Errorf("Expected empty map for nil input, got %v", *result)
	}
}

func TestStructToMap_NestedStruct(t *testing.T) {
	s := NestedStruct{
		Simple: SimpleStruct{
			Name:  "Nested",
			Age:   20,
			Email: "nested@example.com",
		},
		Name: "Parent",
	}

	result := model.StructToMap(s)

	// The nested struct should be formatted as a string representation
	if len(*result) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(*result))
	}

	if (*result)["name"] != "Parent" {
		t.Errorf("Expected name='Parent', got '%s'", (*result)["name"])
	}

	// The Simple field will be converted to string using fmt.Sprintf
	if (*result)["simple"] == "" {
		t.Error("Expected non-empty string for nested struct")
	}
}

func TestStructToMap_BooleanValues(t *testing.T) {
	type BoolStruct struct {
		TrueField  bool `nxrm:"true_field"`
		FalseField bool `nxrm:"false_field"`
	}

	s := BoolStruct{
		TrueField:  true,
		FalseField: false,
	}

	result := model.StructToMap(s)

	expected := map[string]string{
		"true_field":  "true",
		"false_field": "false",
	}

	if !mapsEqual(*result, expected) {
		t.Errorf("Expected %v, got %v", expected, *result)
	}
}

func TestStructToMap_NegativeNumbers(t *testing.T) {
	type NegativeStruct struct {
		NegInt   int     `nxrm:"neg_int"`
		NegFloat float64 `nxrm:"neg_float"`
	}

	s := NegativeStruct{
		NegInt:   -999,
		NegFloat: -123.456,
	}

	result := model.StructToMap(s)

	if (*result)["neg_int"] != "-999" {
		t.Errorf("Expected neg_int='-999', got '%s'", (*result)["neg_int"])
	}

	if (*result)["neg_float"] != "-123.456" {
		t.Errorf("Expected neg_float='-123.456', got '%s'", (*result)["neg_float"])
	}
}

// Helper function to compare maps
func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
