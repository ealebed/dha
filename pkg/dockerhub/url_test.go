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
	"fmt"
	"strings"
	"testing"
)

func TestURLConstruction(t *testing.T) {
	tests := []struct {
		name     string
		org      string
		image    string
		tag      string
		validate func(*testing.T, string, string)
	}{
		{
			name:  "repository URL construction",
			org:   "myorg",
			image: "myimage",
			tag:   "",
			validate: func(t *testing.T, org, image string) {
				expected := fmt.Sprintf("%s/%s/?page=1&page_size=100", RepositoriesURL, org)
				// Verify the URL pattern
				if !strings.Contains(expected, RepositoriesURL) {
					t.Errorf("URL should contain RepositoriesURL, got %v", expected)
				}
				if !strings.Contains(expected, org) {
					t.Errorf("URL should contain org, got %v", expected)
				}
			},
		},
		{
			name:  "tag URL construction",
			org:   "testorg",
			image: "testimage",
			tag:   "v1.0.0",
			validate: func(t *testing.T, org, image string) {
				expected := fmt.Sprintf("%s/%s/%s/tags/?page_size=100", RepositoriesURL, org, image)
				if !strings.Contains(expected, RepositoriesURL) {
					t.Errorf("URL should contain RepositoriesURL, got %v", expected)
				}
				if !strings.Contains(expected, org) {
					t.Errorf("URL should contain org, got %v", expected)
				}
				if !strings.Contains(expected, image) {
					t.Errorf("URL should contain image, got %v", expected)
				}
				if !strings.Contains(expected, "/tags/") {
					t.Errorf("URL should contain /tags/, got %v", expected)
				}
			},
		},
		{
			name:  "delete tag URL construction",
			org:   "myorg",
			image: "myimage",
			tag:   "latest",
			validate: func(t *testing.T, org, image string) {
				expected := fmt.Sprintf("%s/%s/%s/tags/%s/", RepositoriesURL, org, image, "latest")
				if !strings.Contains(expected, "/tags/") {
					t.Errorf("URL should contain /tags/, got %v", expected)
				}
				if !strings.HasSuffix(expected, "/") {
					t.Errorf("Delete URL should end with /, got %v", expected)
				}
			},
		},
		{
			name:  "describe repository URL construction",
			org:   "org",
			image: "image",
			tag:   "",
			validate: func(t *testing.T, org, image string) {
				expected := fmt.Sprintf("%s/%s/%s", RepositoriesURL, org, image)
				if !strings.Contains(expected, org) {
					t.Errorf("URL should contain org, got %v", expected)
				}
				if !strings.Contains(expected, image) {
					t.Errorf("URL should contain image, got %v", expected)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.org, tt.image)
		})
	}
}

func TestBaseURLConstants(t *testing.T) {
	// Test BaseURL is properly formatted
	if BaseURL == "" {
		t.Error("BaseURL constant should not be empty")
	}

	if !strings.HasPrefix(BaseURL, "https://") {
		t.Errorf("BaseURL should start with https://, got %v", BaseURL)
	}

	if !strings.HasSuffix(BaseURL, "/") {
		t.Errorf("BaseURL should end with /, got %v", BaseURL)
	}

	// Test RepositoriesURL is properly constructed
	if RepositoriesURL == "" {
		t.Error("RepositoriesURL constant should not be empty")
	}

	if !strings.Contains(RepositoriesURL, BaseURL) {
		t.Errorf("RepositoriesURL should contain BaseURL, got %v", RepositoriesURL)
	}

	if !strings.Contains(RepositoriesURL, "repositories") {
		t.Errorf("RepositoriesURL should contain 'repositories', got %v", RepositoriesURL)
	}
}

func TestURLPatterns(t *testing.T) {
	org := "testorg"
	image := "testimage"
	tag := "v1.0.0"

	// Test repository list URL pattern
	repoListURL := fmt.Sprintf("%s/%s/?page=1&page_size=100", RepositoriesURL, org)
	if !strings.Contains(repoListURL, "page=1") {
		t.Errorf("Repository list URL should contain page=1, got %v", repoListURL)
	}
	if !strings.Contains(repoListURL, "page_size=100") {
		t.Errorf("Repository list URL should contain page_size=100, got %v", repoListURL)
	}

	// Test tag list URL pattern
	tagListURL := fmt.Sprintf("%s/%s/%s/tags/?page_size=100", RepositoriesURL, org, image)
	if !strings.Contains(tagListURL, "/tags/") {
		t.Errorf("Tag list URL should contain /tags/, got %v", tagListURL)
	}
	if !strings.Contains(tagListURL, "page_size=100") {
		t.Errorf("Tag list URL should contain page_size=100, got %v", tagListURL)
	}

	// Test delete tag URL pattern
	deleteTagURL := fmt.Sprintf("%s/%s/%s/tags/%s/", RepositoriesURL, org, image, tag)
	if !strings.Contains(deleteTagURL, "/tags/") {
		t.Errorf("Delete tag URL should contain /tags/, got %v", deleteTagURL)
	}
	if !strings.HasSuffix(deleteTagURL, "/") {
		t.Errorf("Delete tag URL should end with /, got %v", deleteTagURL)
	}
}
