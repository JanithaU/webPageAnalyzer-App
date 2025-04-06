package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		method         string
		url            string
		expectedStatus int
		expectedBody   string
		formData       string
	}{
		{
			method:         "GET",
			url:            "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "<form", // Ensure the form is part of the response body
		},
		{
			method:         "POST",
			url:            "/",
			expectedStatus: http.StatusOK,
			expectedBody:   "Example Domain",
			formData:       "http://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			var req *http.Request
			var err error

			// Create request with form data if POST method
			if tt.method == "POST" {
				form := url.Values{}
				form.Add("url", tt.formData) // Add the URL to the form
				req, err = http.NewRequest(tt.method, tt.url, strings.NewReader(form.Encode()))
				req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			} else {
				req, err = http.NewRequest(tt.method, tt.url, nil)
			}

			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler(rr, req)

			resp := rr.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %v, but got %v", tt.expectedStatus, resp.StatusCode)
			}

			// Read response body and check expected content
			body := rr.Body.String()
			if !strings.Contains(body, tt.expectedBody) {
				t.Errorf("Expected body to contain %v, but got %v", tt.expectedBody, body)
			}
		})
	}
}
