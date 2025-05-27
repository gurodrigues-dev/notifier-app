package stringcommon

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{" ", false},
		{"abc", false},
		{"0", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Empty(%q)", tt.input), func(t *testing.T) {
			result := Empty(tt.input)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSerializeToJSON(t *testing.T) {
	type sampleStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name        string
		input       any
		expected    []byte
		expectError bool
	}{
		{
			name:     "Sample Struct",
			input:    sampleStruct{Name: "Alice", Age: 30},
			expected: []byte(`{"name":"Alice","age":30}`),
		},
		{
			name:     "Map",
			input:    map[string]int{"one": 1, "two": 2},
			expected: []byte(`{"one":1,"two":2}`),
		},
		{
			name:        "Channel",
			input:       make(chan int),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SerializeToJSON(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("error was expected but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("error not expected: %v", err)
				return
			}

			var expectedMap, resultMap map[string]any
			if err := json.Unmarshal(tt.expected, &expectedMap); err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			if err := json.Unmarshal(result, &resultMap); err != nil {
				t.Fatalf("unmarshal result error: %v", err)
			}
			if !reflect.DeepEqual(resultMap, expectedMap) {
				t.Errorf("result not expected: got %s, want %s", result, tt.expected)
			}
		})
	}
}
