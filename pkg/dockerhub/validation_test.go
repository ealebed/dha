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
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestClientFieldValidation(t *testing.T) {
	client := NewClient("testorg", "https://test.com")

	// Validate Client field
	if client.Client == nil {
		t.Fatal("Client.Client should not be nil")
	}

	// Validate Header field
	if client.Header == nil {
		t.Fatal("Client.Header should not be nil")
	}

	// Validate URL field
	if client.URL == "" {
		t.Error("Client.URL should not be empty")
	}

	// Validate ORG field (can be empty, but should be set)
	if client.ORG != "testorg" {
		t.Errorf("Client.ORG = %v, want testorg", client.ORG)
	}

	// Validate AuthToken starts empty
	if client.AuthToken != "" {
		t.Errorf("Client.AuthToken should be empty initially, got %v", client.AuthToken)
	}
}

func TestJSONRoundTrip(t *testing.T) {
	// Test Repository round-trip
	repo := &Repository{
		Name:   "test",
		Status: 1,
	}
	data, err := json.Marshal(repo)
	if err != nil {
		t.Fatalf("json.Marshal() Repository error = %v", err)
	}
	var repoUnmarshaled Repository
	if err := json.Unmarshal(data, &repoUnmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() Repository error = %v", err)
	}
	if repoUnmarshaled.Name != repo.Name {
		t.Errorf("Repository round-trip Name = %v, want %v", repoUnmarshaled.Name, repo.Name)
	}

	// Test Tag round-trip
	tag := &Tag{
		Name:     "v1.0.0",
		FullSize: 1000,
	}
	data, err = json.Marshal(tag)
	if err != nil {
		t.Fatalf("json.Marshal() Tag error = %v", err)
	}
	var tagUnmarshaled Tag
	if err := json.Unmarshal(data, &tagUnmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() Tag error = %v", err)
	}
	if tagUnmarshaled.Name != tag.Name {
		t.Errorf("Tag round-trip Name = %v, want %v", tagUnmarshaled.Name, tag.Name)
	}

	// Test Image round-trip
	image := &Image{
		OS:   "linux",
		Size: 1000,
	}
	data, err = json.Marshal(image)
	if err != nil {
		t.Fatalf("json.Marshal() Image error = %v", err)
	}
	var imageUnmarshaled Image
	if err := json.Unmarshal(data, &imageUnmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() Image error = %v", err)
	}
	if imageUnmarshaled.OS != image.OS {
		t.Errorf("Image round-trip OS = %v, want %v", imageUnmarshaled.OS, image.OS)
	}

	// Test AuthResponse round-trip
	authResp := &AuthResponse{
		Token: "test-token",
	}
	data, err = json.Marshal(authResp)
	if err != nil {
		t.Fatalf("json.Marshal() AuthResponse error = %v", err)
	}
	var authRespUnmarshaled AuthResponse
	if err := json.Unmarshal(data, &authRespUnmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() AuthResponse error = %v", err)
	}
	if authRespUnmarshaled.Token != authResp.Token {
		t.Errorf("AuthResponse round-trip Token = %v, want %v", authRespUnmarshaled.Token, authResp.Token)
	}
}

func TestConstantsAreSet(t *testing.T) {
	// Test BaseURL
	if BaseURL == "" {
		t.Error("BaseURL constant should be set")
	}

	// Test RepositoriesURL
	if RepositoriesURL == "" {
		t.Error("RepositoriesURL constant should be set")
	}

	// Verify they're related
	if !strings.Contains(RepositoriesURL, BaseURL) {
		t.Errorf("RepositoriesURL should contain BaseURL: %v vs %v", RepositoriesURL, BaseURL)
	}
}

func TestClientDefaults(t *testing.T) {
	// Test default URL
	client := NewClient("testorg", "")
	if client.URL != "https://hub.docker.com" {
		t.Errorf("Default URL = %v, want https://hub.docker.com", client.URL)
	}

	// Test custom URL
	client2 := NewClient("testorg", "https://custom.com")
	if client2.URL != "https://custom.com" {
		t.Errorf("Custom URL = %v, want https://custom.com", client2.URL)
	}
}

func TestClientHeaderDefaults(t *testing.T) {
	client := NewClient("testorg", "https://test.com")

	// Verify Content-Type is set
	contentType := client.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Default Content-Type = %v, want application/json", contentType)
	}

	// Verify header is not nil
	if client.Header == nil {
		t.Fatal("Header should not be nil")
	}
}

func TestClientTimeout(t *testing.T) {
	client := NewClient("testorg", "https://test.com")

	// Verify timeout is set to 30 seconds
	if client.Client.Timeout != 30*time.Second {
		t.Errorf("Client timeout = %v, want 30s", client.Client.Timeout)
	}
}
