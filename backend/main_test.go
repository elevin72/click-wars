package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	w := httptest.NewRecorder()
	rootHandler(w, httptest.NewRequest("GET", "/", nil))

	expected := "<!DOCTYPE html>"
	if !strings.HasPrefix(w.Body.String(), expected) {
		t.Fatalf("Expected response to be start with %s", expected)
	}
}

func TestStaticFiles(t *testing.T) {
	tt := []struct {
		name               string
		method             string
		want               string
		expectedStatusCode int
	}{
		{
			name:               "/static/js/script.js",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "/static/css/style.css",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "/static/css/not-found",
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			staticHandler(w, httptest.NewRequest("GET", tc.name, nil))
			if statusCode := w.Result().StatusCode; statusCode != tc.expectedStatusCode {
				t.Fatalf("Expected %d, got %d", tc.expectedStatusCode, statusCode)
			}
		})
	}
}
