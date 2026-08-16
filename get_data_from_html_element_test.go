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

func TestGetImagesFromHTMLRelative(t *testing.T) {
	tests := []struct {
		name      string
		inputURL  string
		inputBody string
		expected  []string
	}{
		{
			name:      "Get img src",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="https://crawler-test.com/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:      "Get img src relative",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			expected:  []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:     "Get img src multi",
			inputURL: "https://crawler-test.com",
			inputBody: `<html><body>
		<img src="https://crawler-test.com/logo.png" alt="Logo">
		<img src="/logo1.png" alt="Logo">
		<img src="https://crawler-test.com/logo2.png" alt="Logo">
		<img src="/logo3.png" alt="Logo">
			</body></html>`,
			expected: []string{
				"https://crawler-test.com/logo.png",
				"https://crawler-test.com/logo1.png",
				"https://crawler-test.com/logo2.png",
				"https://crawler-test.com/logo3.png",
			},
		},
		{
			name:      "No a href",
			inputURL:  "https://crawler-test.com",
			inputBody: `<html><body><img alt="Logo"></body></html>`,
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

			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: \n expected text: %v \n   actual text: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

// func TestExtractPageData(t *testing.T) {
// 	inputURL := "https://crawler-test.com"
// 	inputBody := `<html><body>
//         <h1>Test Title</h1>
//         <p>This is the first paragraph.</p>
//         <a href="/link1">Link 1</a>
//         <img src="/image1.jpg" alt="Image 1">
//     </body></html>`
//
// 	actual := extractPageData(inputBody, inputURL)
//
// 	expected := PageData{
// 		URL:            "https://crawler-test.com",
// 		Heading:        "Test Title",
// 		FirstParagraph: "This is the first paragraph.",
// 		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
// 		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
// 	}
//
// 	if !reflect.DeepEqual(actual, expected) {
// 		t.Errorf("\nexpected %+v\n     got %+v", expected, actual)
// 	}
// }

func TestExtractPageData(t *testing.T) {
	cases := []struct {
		name    string
		pageURL string
		html    string
		want    PageData
	}{
		{
			name:    "basic: h1, main paragraph, relative link and img",
			pageURL: "https://crawler-test.com",
			html: `
<html>
  <body>
    <h1>Hello World</h1>
    <main><p>First paragraph inside main.</p></main>
    <a href="/about">About</a>
    <img src="/logo.png" alt="Logo">
  </body>
</html>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Hello World",
				FirstParagraph: "First paragraph inside main.",
				OutgoingLinks:  []string{"https://crawler-test.com/about"},
				ImageURLs:      []string{"https://crawler-test.com/logo.png"},
			},
		},
		{
			name:    "fallback paragraph when no <main>",
			pageURL: "https://crawler-test.com",
			html: `
<html>
  <body>
    <h1>Title</h1>
    <p>Outside paragraph wins.</p>
    <a href="/x">x</a>
    <img src="/img.png">
  </body>
</html>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Title",
				FirstParagraph: "Outside paragraph wins.",
				OutgoingLinks:  []string{"https://crawler-test.com/x"},
				ImageURLs:      []string{"https://crawler-test.com/img.png"},
			},
		},
		{
			name:    "malformed HTML still parsed; absolute link and image",
			pageURL: "https://crawler-test.com",
			html: `
<html body>
  <h1>Messy</h1>
  <a href="https://other.com/path">Other</a>
  <img src="https://cdn.boot.dev/banner.jpg">
</html body>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "Messy",
				FirstParagraph: "", // no <p> present
				OutgoingLinks:  []string{"https://other.com/path"},
				ImageURLs:      []string{"https://cdn.boot.dev/banner.jpg"},
			},
		},
		{
			name:    "no h1 and no paragraph",
			pageURL: "https://crawler-test.com",
			html: `
<html>
  <body>
    <a href="/only-link">Only link</a>
    <img src="/only.png">
  </body>
</html>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "",
				FirstParagraph: "",
				OutgoingLinks:  []string{"https://crawler-test.com/only-link"},
				ImageURLs:      []string{"https://crawler-test.com/only.png"},
			},
		},
		{
			name:    "multiple links and images preserve order",
			pageURL: "https://crawler-test.com",
			html: `
<html><body>
  <h1>t</h1>
  <main><p>p</p></main>
  <a href="/a1">a1</a>
  <a href="https://x.dev/a2">a2</a>
  <img src="/i1.png">
  <img src="https://x.dev/i2.png">
</body></html>`,
			want: PageData{
				URL:            "https://crawler-test.com",
				Heading:        "t",
				FirstParagraph: "p",
				OutgoingLinks: []string{
					"https://crawler-test.com/a1",
					"https://x.dev/a2",
				},
				ImageURLs: []string{
					"https://crawler-test.com/i1.png",
					"https://x.dev/i2.png",
				},
			},
		},
		{
			name:    "invalid base URL → empty link/image slices",
			pageURL: `:\\invalidBaseURL`,
			html: `
<html>
  <body>
    <h1>Title</h1>
    <p>Paragraph</p>
    <a href="/path">path</a>
    <img src="/logo.png">
  </body>
</html>`,
			want: PageData{
				URL:            `:\\invalidBaseURL`,
				Heading:        "Title",
				FirstParagraph: "Paragraph",
				OutgoingLinks:  []string{},
				ImageURLs:      []string{},
			},
		},
	}

	for _, tc := range cases {
		// tc := tc // shadow the loop variable.
		t.Run(tc.name, func(t *testing.T) {
			got := extractPageData(tc.html, tc.pageURL)

			if got.URL != tc.want.URL {
				t.Errorf("URL: want %q, got %q", tc.want.URL, got.URL)
			}
			if got.Heading != tc.want.Heading {
				t.Errorf("Heading: want %q, got %q", tc.want.Heading, got.Heading)
			}
			if got.FirstParagraph != tc.want.FirstParagraph {
				t.Errorf("FirstParagraph: want %q, got %q", tc.want.FirstParagraph, got.FirstParagraph)
			}
			if !reflect.DeepEqual(got.OutgoingLinks, tc.want.OutgoingLinks) {
				t.Errorf("OutgoingLinks: want %v, got %v", tc.want.OutgoingLinks, got.OutgoingLinks)
			}
			if !reflect.DeepEqual(got.ImageURLs, tc.want.ImageURLs) {
				t.Errorf("ImageURLs: want %v, got %v", tc.want.ImageURLs, got.ImageURLs)
			}
		})
	}
}
