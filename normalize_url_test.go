package main

import (
	"strings"
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		expected    string
		errorPrefix string
	}{
		{
			name:        "remove scheme",
			inputURL:    "https://www.boot.dev/blog/path",
			expected:    "www.boot.dev/blog/path",
			errorPrefix: "NormalizeURL error",
		},
		{
			name:        "remove scheme and trailing backslash",
			inputURL:    "https://www.boot.dev/blog/path/",
			expected:    "www.boot.dev/blog/path",
			errorPrefix: "NormalizeURL error",
		},
		{
			name:        "remove scheme",
			inputURL:    "http://www.boot.dev/blog/path",
			expected:    "www.boot.dev/blog/path",
			errorPrefix: "NormalizeURL error",
		},
		{
			name:        "remove scheme and trailing backslash",
			inputURL:    "http://www.boot.dev/blog/path/",
			expected:    "www.boot.dev/blog/path",
			errorPrefix: "NormalizeURL error",
		},
		{
			name:        "remove scheme, captital letters and trailing backslash",
			inputURL:    "HTTP://WWW.BOOT.DEV/BLOG/PATH/",
			expected:    "www.boot.dev/blog/path",
			errorPrefix: "NormalizeURL error",
		},
		{
			name:        "handle invalid URL",
			inputURL:    `:\\invalidURL`,
			expected:    "",
			errorPrefix: "NormalizeURL error",
		},
		// add more test cases here
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorPrefix) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: \n expected URL: %v \n   actual URL: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
