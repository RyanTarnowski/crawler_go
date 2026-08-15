package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTMLBasic(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name:      "Get H1 text",
			inputBody: "<html><body><h1>Test H1 Title</h1></body></html>",
			expected:  "Test H1 Title",
		},
		{
			name:      "Get H2 text",
			inputBody: "<html><body><h2>Test H2 Title</h2></body></html>",
			expected:  "Test H2 Title",
		},
		{
			name:      "Missing H1 and H2",
			inputBody: "<html><body></body></html>",
			expected:  "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: \n expected text: %v \n   actual text: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		expected  string
	}{
		{
			name: "Get Main P text",
			inputBody: `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "Get P text",
			inputBody: `<html><body>
		<p>paragraph.</p>
	</body></html>`,
			expected: "paragraph.",
		},
		{
			name:      "Missing P",
			inputBody: "<html><body></body></html>",
			expected:  "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.inputBody)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: \n expected text: %v \n   actual text: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "Get a href",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com"},
		},
		{
			name:      "Get a href relative",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a href="/test"><span>Boot.dev</span></a></body></html>`,
			expected:  []string{"https://crawler-test.com/test"},
		},
		{
			name:     "Get a href multi",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
			<a href="https://crawler-test.com/test"><span>Boot.dev</span></a>
			<a href="/test1"><span>Boot.dev</span></a>
			<a href="/test2"><span>Boot.dev</span></a>
			<a href="/test3"><span>Boot.dev</span></a>
			</body></html>`,
			expected: []string{
				"https://crawler-test.com/test",
				"https://crawler-test.com/test1",
				"https://crawler-test.com/test2",
				"https://crawler-test.com/test3",
			},
		},
		{
			name:      "No a href",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><a <span>Boot.dev</span></a></body></html>`,
			expected:  []string{},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}

			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: \n expected text: %v \n   actual text: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}
