package analysis_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JanithaU/webPageAnalyzer-App/internal/analysis"
)

func TestAnalyzePage(t *testing.T) {
	// Define the HTML content to be served by the test server
	const htmlContent = `
    <!DOCTYPE html>
    <html>
    <head>
        <title>Test Page</title>
    </head>
    <body>
        <h1>Main Heading</h1>
        <h2>Subheading</h2>
        <a href="http://external.com">External Link</a>
        <a href="/internal">Internal Link</a>
        <form action="/login">
            <input type="text" name="username">
            <input type="password" name="password">
        </form>
    </body>
    </html>`

	// Create a test server that serves the HTML content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlContent))
	}))
	defer server.Close()

	results, err := analysis.AnalyzePage(server.URL)
	if err != nil {
		t.Fatalf("AnalyzePage returned an error: %v", err)
	}

	if results.Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got '%s'", results.Title)
	}
	if results.HeadingCounts["h1"] != 1 {
		t.Errorf("Expected 1 <h1> tag, got %d", results.HeadingCounts["h1"])
	}
	if results.HeadingCounts["h2"] != 1 {
		t.Errorf("Expected 1 <h2> tag, got %d", results.HeadingCounts["h2"])
	}
	if results.InternalLinks != 1 {
		t.Errorf("Expected 1 internal link, got %d", results.InternalLinks)
	}
	if results.ExternalLinks != 1 {
		t.Errorf("Expected 1 external link, got %d", results.ExternalLinks)
	}
	if !results.LoginFormPresent {
		t.Errorf("Expected login form to be present")
	}
}
