package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Assuming fetchVersion is defined in the same package

func TestFetchVersion(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   interface{}
		mockStatusCode int
		jsonKey        string
		want           string
		wantErr        bool
	}{
		{
			name:           "valid response",
			mockResponse:   map[string]string{"version": "1.0.0"},
			mockStatusCode: http.StatusOK,
			jsonKey:        "version",
			want:           "1.0.0",
			wantErr:        false,
		},
		{
			name:           "key not found",
			mockResponse:   map[string]string{"invalidKey": "1.0.0"},
			mockStatusCode: http.StatusOK,
			jsonKey:        "version",
			want:           "",
			wantErr:        true,
		},
		{
			name:           "invalid json",
			mockResponse:   "invalid json",
			mockStatusCode: http.StatusOK,
			jsonKey:        "version",
			want:           "",
			wantErr:        true,
		},
		{
			name:           "bad request",
			mockResponse:   map[string]string{"version": "1.0.0"},
			mockStatusCode: http.StatusBadRequest,
			jsonKey:        "version",
			want:           "",
			wantErr:        true,
		},
		{
			name: "nested valid response",
			mockResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"version": "2.48.1",
				},
			},
			mockStatusCode: http.StatusOK,
			jsonKey:        "data.version",
			want:           "2.48.1",
			wantErr:        false,
		},
		{
			name: "nested key not found",
			mockResponse: map[string]interface{}{
				"data": map[string]interface{}{
					"invalidKey": "2.48.1",
				},
			},
			mockStatusCode: http.StatusOK,
			jsonKey:        "data.version",
			want:           "",
			wantErr:        true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.mockStatusCode)
				json.NewEncoder(w).Encode(tc.mockResponse)
			}))
			defer server.Close()

			// Call the function under test
			got, err := fetchVersion(server.URL, tc.jsonKey)

			// Check for error
			if (err != nil) != tc.wantErr {
				t.Errorf("fetchVersion() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// Check the function response
			if got != tc.want {
				t.Errorf("fetchVersion() = %v, want %v", got, tc.want)
			}
		})
	}
}
