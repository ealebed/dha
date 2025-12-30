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

func TestRepositoryEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		repo *Repository
	}{
		{
			name: "repository with zero values",
			repo: &Repository{},
		},
		{
			name: "repository with negative values",
			repo: &Repository{
				Status:            -1,
				StarCount:         -10,
				PullCount:         -100,
				CollaboratorCount: -5,
			},
		},
		{
			name: "repository with very large values",
			repo: &Repository{
				StarCount:         999999999,
				PullCount:         999999999,
				CollaboratorCount: 999999999,
			},
		},
		{
			name: "repository with future date",
			repo: &Repository{
				LastUpdated: time.Now().Add(24 * time.Hour),
			},
		},
		{
			name: "repository with past date",
			repo: &Repository{
				LastUpdated: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.repo)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Test unmarshaling
			var got Repository
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Verify key fields match
			if got.Status != tt.repo.Status {
				t.Errorf("Unmarshal Status = %v, want %v", got.Status, tt.repo.Status)
			}
		})
	}
}

func TestTagEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		tag  *Tag
	}{
		{
			name: "tag with zero values",
			tag:  &Tag{},
		},
		{
			name: "tag with negative ID",
			tag: &Tag{
				ID: -1,
			},
		},
		{
			name: "tag with very large size",
			tag: &Tag{
				FullSize: 999999999999,
			},
		},
		{
			name: "tag with empty name",
			tag: &Tag{
				Name: "",
			},
		},
		{
			name: "tag with very long name",
			tag: &Tag{
				Name: strings.Repeat("a", 1000),
			},
		},
		{
			name: "tag with multiple images",
			tag: &Tag{
				Images: []*Image{
					{OS: "linux", Size: 1000},
					{OS: "windows", Size: 2000},
					{OS: "darwin", Size: 3000},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.tag)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Test unmarshaling
			var got Tag
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Verify key fields match
			if got.Name != tt.tag.Name {
				t.Errorf("Unmarshal Name = %v, want %v", got.Name, tt.tag.Name)
			}
		})
	}
}

func TestImageEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		image *Image
	}{
		{
			name:  "image with zero values",
			image: &Image{},
		},
		{
			name: "image with negative size",
			image: &Image{
				Size: -1000,
			},
		},
		{
			name: "image with very large size",
			image: &Image{
				Size: 999999999999,
			},
		},
		{
			name: "image with empty OS",
			image: &Image{
				OS: "",
			},
		},
		{
			name: "image with all architectures",
			image: &Image{
				Architecture: "amd64",
				OS:           "linux",
				Variant:      "variant",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.image)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Test unmarshaling
			var got Image
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			// Verify key fields match
			if got.Size != tt.image.Size {
				t.Errorf("Unmarshal Size = %v, want %v", got.Size, tt.image.Size)
			}
		})
	}
}

func TestRepositoryListEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		repoList *RepositoryList
		validate func(*testing.T, *RepositoryList)
	}{
		{
			name: "empty repository list",
			repoList: &RepositoryList{
				Count:   0,
				Results: []*Repository{},
			},
			validate: func(t *testing.T, got *RepositoryList) {
				if got.Count != 0 {
					t.Errorf("Count = %v, want 0", got.Count)
				}
				if len(got.Results) != 0 {
					t.Errorf("Results length = %v, want 0", len(got.Results))
				}
			},
		},
		{
			name: "repository list with nil results",
			repoList: &RepositoryList{
				Count:   0,
				Results: nil,
			},
			validate: func(t *testing.T, got *RepositoryList) {
				if got.Count != 0 {
					t.Errorf("Count = %v, want 0", got.Count)
				}
			},
		},
		{
			name: "repository list with large count",
			repoList: &RepositoryList{
				Count:   999999,
				Results: []*Repository{},
			},
			validate: func(t *testing.T, got *RepositoryList) {
				if got.Count != 999999 {
					t.Errorf("Count = %v, want 999999", got.Count)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.repoList)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Test unmarshaling
			var got RepositoryList
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, &got)
			}
		})
	}
}

func TestTagListEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		tagList  *TagList
		validate func(*testing.T, *TagList)
	}{
		{
			name: "empty tag list",
			tagList: &TagList{
				Count:   0,
				Results: []*Tag{},
			},
			validate: func(t *testing.T, got *TagList) {
				if got.Count != 0 {
					t.Errorf("Count = %v, want 0", got.Count)
				}
				if len(got.Results) != 0 {
					t.Errorf("Results length = %v, want 0", len(got.Results))
				}
			},
		},
		{
			name: "tag list with nil results",
			tagList: &TagList{
				Count:   0,
				Results: nil,
			},
			validate: func(t *testing.T, got *TagList) {
				if got.Count != 0 {
					t.Errorf("Count = %v, want 0", got.Count)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.tagList)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}

			// Test unmarshaling
			var got TagList
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, &got)
			}
		})
	}
}
