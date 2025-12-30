/*
Copyright © 2020 Yevhen Lebid ealebed@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dockerhub

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name        string
		authToken   string
		method      string
		url         string
		payload     io.Reader
		wantAuth    bool
		wantMethod  string
		validateReq func(*testing.T, *http.Request)
	}{
		{
			name:       "GET request with existing auth token",
			authToken:  "test-token-123",
			method:     http.MethodGet,
			url:        "https://example.com/api",
			payload:    nil,
			wantAuth:   true,
			wantMethod: http.MethodGet,
			validateReq: func(t *testing.T, req *http.Request) {
				authHeader := req.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "JWT ") {
					t.Errorf("Authorization header should start with 'JWT ', got %v", authHeader)
				}
				if !strings.Contains(authHeader, "test-token-123") {
					t.Errorf("Authorization header should contain token, got %v", authHeader)
				}
			},
		},
		{
			name:       "POST request with payload",
			authToken:  "token-456",
			method:     http.MethodPost,
			url:        "https://example.com/api",
			payload:    bytes.NewBufferString(`{"key":"value"}`),
			wantAuth:   true,
			wantMethod: http.MethodPost,
			validateReq: func(t *testing.T, req *http.Request) {
				if req.Method != http.MethodPost {
					t.Errorf("Request method = %v, want %v", req.Method, http.MethodPost)
				}
				if req.URL.String() != "https://example.com/api" {
					t.Errorf("Request URL = %v, want https://example.com/api", req.URL.String())
				}
			},
		},
		{
			name:       "DELETE request",
			authToken:  "delete-token",
			method:     http.MethodDelete,
			url:        "https://example.com/resource",
			payload:    nil,
			wantAuth:   true,
			wantMethod: http.MethodDelete,
			validateReq: func(t *testing.T, req *http.Request) {
				if req.Method != http.MethodDelete {
					t.Errorf("Request method = %v, want %v", req.Method, http.MethodDelete)
				}
			},
		},
		{
			name:       "PUT request with payload",
			authToken:  "put-token",
			method:     http.MethodPut,
			url:        "https://example.com/update",
			payload:    bytes.NewBufferString(`{"data":"updated"}`),
			wantAuth:   true,
			wantMethod: http.MethodPut,
			validateReq: func(t *testing.T, req *http.Request) {
				if req.Method != http.MethodPut {
					t.Errorf("Request method = %v, want %v", req.Method, http.MethodPut)
				}
			},
		},
		{
			name:       "empty auth token",
			authToken:  "",
			method:     http.MethodGet,
			url:        "https://example.com/api",
			payload:    nil,
			wantAuth:   false,
			wantMethod: http.MethodGet,
			validateReq: func(t *testing.T, req *http.Request) {
				// When auth token is empty, GetAuthToken will be called
				// which requires environment variables, so we can't fully test this
				// but we can verify the request structure is created
				if req.URL.String() != "https://example.com/api" {
					t.Errorf("Request URL = %v, want https://example.com/api", req.URL.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("testorg", "")
			client.AuthToken = tt.authToken

			// For empty token case, skip the actual request creation test
			// as it would require environment variables
			if tt.authToken == "" {
				t.Skip("Skipping test that requires GetAuthToken (needs env vars)")
				return
			}

			req, err := client.NewRequest(tt.method, tt.url, tt.payload)
			if err != nil {
				t.Fatalf("NewRequest() error = %v", err)
			}

			if req == nil {
				t.Fatal("NewRequest() returned nil request")
			}

			if req.Method != tt.wantMethod {
				t.Errorf("NewRequest() method = %v, want %v", req.Method, tt.wantMethod)
			}

			if tt.wantAuth {
				authHeader := req.Header.Get("Authorization")
				if authHeader == "" {
					t.Error("NewRequest() Authorization header should be set")
				}
				if !strings.HasPrefix(authHeader, "JWT ") {
					t.Errorf("NewRequest() Authorization header should start with 'JWT ', got %v", authHeader)
				}
			}

			if tt.validateReq != nil {
				tt.validateReq(t, req)
			}
		})
	}
}

func TestNewRequestInvalidURL(t *testing.T) {
	client := NewClient("testorg", "")
	client.AuthToken = "test-token"

	// Test with invalid URL
	_, err := client.NewRequest(http.MethodGet, "://invalid-url", nil)
	if err == nil {
		t.Error("NewRequest() with invalid URL should return error")
	}
}

func TestAuthResponseJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    *AuthResponse
		wantErr bool
	}{
		{
			name:    "valid auth response",
			json:    `{"token":"test-token-123"}`,
			want:    &AuthResponse{Token: "test-token-123"},
			wantErr: false,
		},
		{
			name:    "empty token",
			json:    `{"token":""}`,
			want:    &AuthResponse{Token: ""},
			wantErr: false,
		},
		{
			name:    "missing token field",
			json:    `{}`,
			want:    &AuthResponse{Token: ""},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			json:    `{"token":invalid}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got AuthResponse
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got.Token != tt.want.Token {
				t.Errorf("Unmarshal Token = %v, want %v", got.Token, tt.want.Token)
			}
		})
	}
}

func TestAuthResponseMarshal(t *testing.T) {
	resp := &AuthResponse{Token: "test-token-456"}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Verify it contains the token
	if !strings.Contains(string(data), "test-token-456") {
		t.Errorf("Marshaled JSON should contain token, got %v", string(data))
	}

	// Verify it can be unmarshaled back
	var got AuthResponse
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if got.Token != resp.Token {
		t.Errorf("Round-trip Token = %v, want %v", got.Token, resp.Token)
	}
}

func TestConstants(t *testing.T) {
	// Test that constants are properly defined
	if BaseURL == "" {
		t.Error("BaseURL should not be empty")
	}

	if RepositoriesURL == "" {
		t.Error("RepositoriesURL should not be empty")
	}

	// Verify BaseURL format
	if !strings.HasPrefix(BaseURL, "https://") {
		t.Errorf("BaseURL should start with https://, got %v", BaseURL)
	}

	// Verify RepositoriesURL contains BaseURL
	if !strings.Contains(RepositoriesURL, BaseURL) {
		t.Errorf("RepositoriesURL should contain BaseURL, got %v", RepositoriesURL)
	}
}
