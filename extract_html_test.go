package main

import (
	"net/url"
	"reflect"
	"testing"
)

type extractH1TestStruct struct {
	name      string
	inputHtml string
	expected  string
}

func TestExtractH1FromHtml(t *testing.T) {
	tests := []extractH1TestStruct{
		{
			name:      "Test case 1",
			inputHtml: "<html><head><title>Test</title></head><body><h1>Hello</h1></body></html>",
			expected:  "Hello",
		},
		{
			name: "Test case 2",
			inputHtml: `<html><body><h1>Welcome to Boot.dev</h1><main>
      				<p>Learn to code by building real projects.</p>
          			<p>This is the second paragraph.</p>
             		</main>
               		</body>
                </html>`,
			expected: "Welcome to Boot.dev",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getH1FromHTML(tc.inputHtml)
			if err != nil {
				t.Errorf("Test %v - '%s FAIL: unexpected error: %v", i, tc.name, err)
			}
			if actual != tc.expected {
				t.Errorf("Test %v - '%s FAIL: expected '%s', got '%s'", i, tc.name, tc.expected, actual)
			}
		})
	}

}

type extractParagraphTestStruct struct {
	name      string
	inputHtml string
	expected  string
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []extractParagraphTestStruct{
		{
			name: "main tag exists with paragraph",
			inputHtml: `<html><body>
        <p>Outside paragraph.</p>
        <main>
            <p>Main paragraph.</p>
        </main>
    </body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "no main tag fallback to first p",
			inputHtml: `<html><body>
        <p>First paragraph.</p>
        <p>Second paragraph.</p>
    </body></html>`,
			expected: "First paragraph.",
		},
		{
			name: "main tag with multiple paragraphs",
			inputHtml: `<html><body>
        <p>Outside paragraph.</p>
        <main>
            <p>First main paragraph.</p>
            <p>Second main paragraph.</p>
        </main>
    </body></html>`,
			expected: "First main paragraph.",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getFirstParagraphFromHTML(tc.inputHtml)
			if err != nil {
				t.Errorf("Test %v - '%s FAIL: unexpected error: %v", i, tc.name, err)
			}
			if actual != tc.expected {
				t.Errorf("Test %v - '%s FAIL: expected '%s', got '%s'", i, tc.name, tc.expected, actual)
			}
		})
	}

}

func TestGetURLsFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><a href="/path/to/page">Page</a></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/path/to/page"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetURLsFromHTMLMultiple(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
		<a href="https://blog.boot.dev">Boot.dev</a>
		<a href="/path1">Path 1</a>
		<a href="/path2">Path 2</a>
	</body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getURLsFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://blog.boot.dev",
		"https://blog.boot.dev/path1",
		"https://blog.boot.dev/path2",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLRelative(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><img src="/logo.png" alt="Logo"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://blog.boot.dev/logo.png"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLAbsolute(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body><img src="https://cdn.example.com/image.jpg" alt="Image"></body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"https://cdn.example.com/image.jpg"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetImagesFromHTMLMultiple(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
		<img src="/logo.png" alt="Logo">
		<img src="https://cdn.example.com/banner.jpg" alt="Banner">
		<img src="/favicon.ico">
	</body></html>`

	baseURL, err := url.Parse(inputURL)
	if err != nil {
		t.Errorf("couldn't parse input URL: %v", err)
		return
	}

	actual, err := getImagesFromHTML(inputBody, baseURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"https://blog.boot.dev/logo.png",
		"https://cdn.example.com/banner.jpg",
		"https://blog.boot.dev/favicon.ico",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestExtractPageData(t *testing.T) {
	inputURL := "https://blog.boot.dev"
	inputBody := `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://blog.boot.dev",
		H1:             "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://blog.boot.dev/link1"},
		ImageURLs:      []string{"https://blog.boot.dev/image1.jpg"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataWithMain(t *testing.T) {
	inputURL := "https://example.com"
	inputBody := `<html><body>
        <h1>Example Title</h1>
        <p>Outside paragraph</p>
        <main>
            <p>Main paragraph content</p>
        </main>
        <a href="/page1">Page 1</a>
        <a href="https://external.com">External</a>
        <img src="/logo.png" alt="Logo">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://example.com",
		H1:             "Example Title",
		FirstParagraph: "Main paragraph content",
		OutgoingLinks: []string{
			"https://example.com/page1",
			"https://external.com",
		},
		ImageURLs: []string{"https://example.com/logo.png"},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}

func TestExtractPageDataMultipleResources(t *testing.T) {
	inputURL := "https://blog.boot.dev/path"
	inputBody := `<html><body>
        <h1>Multi Resource Page</h1>
        <p>First paragraph</p>
        <a href="/link1">Link 1</a>
        <a href="/link2">Link 2</a>
        <a href="https://external.com/page">External Link</a>
        <img src="/img1.jpg" alt="Image 1">
        <img src="/img2.png" alt="Image 2">
        <img src="https://cdn.example.com/img3.jpg" alt="Image 3">
    </body></html>`

	actual := extractPageData(inputBody, inputURL)

	expected := PageData{
		URL:            "https://blog.boot.dev/path",
		H1:             "Multi Resource Page",
		FirstParagraph: "First paragraph",
		OutgoingLinks: []string{
			"https://blog.boot.dev/link1",
			"https://blog.boot.dev/link2",
			"https://external.com/page",
		},
		ImageURLs: []string{
			"https://blog.boot.dev/img1.jpg",
			"https://blog.boot.dev/img2.png",
			"https://cdn.example.com/img3.jpg",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, actual)
	}
}
