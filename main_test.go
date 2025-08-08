package main

import (
	"reflect"
	"testing"
)

func TestParseExportString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Env
	}{
		{
			input: `export TESTING=$'projects/test/secrets/testing';export TEST=$'projects/test/secrets/test';`,
			expected: Env{
				"projects/test/secrets/testing": "projects/test/secrets/testing",
				"projects/test/secrets/test":    "projects/test/secrets/test",
			},
		},
		{
			name:  "Multiple exports",
			input: `export TESTING=$'testing';export TEST=$'test';`,
			expected: Env{
				"testing": "testing",
				"test":    "test",
			},
		},
		{
			name:  "Single export",
			input: `export FOO=$'bar';`,
			expected: Env{
				"bar": "bar",
			},
		},
		{
			name:     "Empty string",
			input:    ``,
			expected: Env{},
		},
		{
			name:     "No matches",
			input:    `echo 'not an export'`,
			expected: Env{},
		},
		{
			name:  "Handles underscores and numbers",
			input: `export VAR_1=$'value1';export VAR_2=$'value2';`,
			expected: Env{
				"value1": "value1",
				"value2": "value2",
			},
		},
		{
			name:  "Handles empty values",
			input: `export EMPTY=$'';`,
			expected: map[string]string{
				"": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseExportString(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got: %#v, expected: %#v", result, tt.expected)
			}
		})
	}
}
