package main

import "testing"

type normalizeTestStruct struct {
	name     string
	inputURL string
	expected string
}

func TestNormalizeURL(t *testing.T) {
	tests := []normalizeTestStruct{
		{
			name:     "remove schems",
			inputURL: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove schems",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove schems",
			inputURL: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove schems",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if tc.expected != actual {
				t.Errorf("Test %v - '%s FAIL: expected '%s', got '%s'", i, tc.name, tc.expected, actual)
			}

		})
	}

}
