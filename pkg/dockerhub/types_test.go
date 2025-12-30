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
	"testing"
	"time"
)

func TestRepositoryJSON(t *testing.T) {
	tests := []struct {
		name string
		repo *Repository
	}{
		{
			name: "complete repository",
			repo: &Repository{
				User:              "testuser",
				Name:              "testrepo",
				Namespace:         "testnamespace",
				RepositoryType:    "image",
				Status:            1,
				Description:       "Test repository",
				IsPrivate:         false,
				IsAutomated:       true,
				CanEdit:           true,
				StarCount:         10,
				PullCount:         100,
				LastUpdated:       time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				IsMigrated:        false,
				CollaboratorCount: 5,
				Affiliation:       "owner",
				HubUser:           "testuser",
			},
		},
		{
			name: "minimal repository",
			repo: &Repository{
				Name:   "minimal",
				Status: 0,
			},
		},
		{
			name: "private repository",
			repo: &Repository{
				Name:      "private",
				IsPrivate: true,
				Status:    1,
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

			// Verify key fields
			if got.Name != tt.repo.Name {
				t.Errorf("Unmarshal Name = %v, want %v", got.Name, tt.repo.Name)
			}
			if got.Status != tt.repo.Status {
				t.Errorf("Unmarshal Status = %v, want %v", got.Status, tt.repo.Status)
			}
			if got.IsPrivate != tt.repo.IsPrivate {
				t.Errorf("Unmarshal IsPrivate = %v, want %v", got.IsPrivate, tt.repo.IsPrivate)
			}
		})
	}
}

func TestTagJSON(t *testing.T) {
	tests := []struct {
		name string
		tag  *Tag
	}{
		{
			name: "complete tag",
			tag: &Tag{
				Creator:         1,
				ID:              123,
				ImageID:         "sha256:abc123",
				Images:          []*Image{},
				LastUpdated:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				LastUpdater:     2,
				LastUpdaterUser: "updater",
				Name:            "v1.0.0",
				Repository:      456,
				FullSize:        1024000,
				V2:              true,
				TagStatus:       "active",
				TagLastPulled:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
				TagLastPushed:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "minimal tag",
			tag: &Tag{
				Name: "latest",
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

			// Verify key fields
			if got.Name != tt.tag.Name {
				t.Errorf("Unmarshal Name = %v, want %v", got.Name, tt.tag.Name)
			}
			if got.FullSize != tt.tag.FullSize {
				t.Errorf("Unmarshal FullSize = %v, want %v", got.FullSize, tt.tag.FullSize)
			}
		})
	}
}

func TestImageJSON(t *testing.T) {
	tests := []struct {
		name  string
		image *Image
	}{
		{
			name: "complete image",
			image: &Image{
				Architecture: "amd64",
				Features:     "features",
				Variant:      "variant",
				Digest:       "sha256:abc123",
				OS:           "linux",
				OSFeatures:   "osfeatures",
				OSVersion:    "20.04",
				Size:         1024000,
				Status:       "active",
				LastPulled:   time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
				LastPushed:   time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "minimal image",
			image: &Image{
				OS:     "linux",
				Size:   0,
				Status: "active",
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

			// Verify key fields
			if got.OS != tt.image.OS {
				t.Errorf("Unmarshal OS = %v, want %v", got.OS, tt.image.OS)
			}
			if got.Size != tt.image.Size {
				t.Errorf("Unmarshal Size = %v, want %v", got.Size, tt.image.Size)
			}
		})
	}
}

func TestRepositoryListJSON(t *testing.T) {
	repoList := &RepositoryList{
		Count:    2,
		Next:     "https://next.page",
		Previous: "https://prev.page",
		Results: []*Repository{
			{Name: "repo1", Status: 1},
			{Name: "repo2", Status: 1},
		},
	}

	// Test marshaling
	data, err := json.Marshal(repoList)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Test unmarshaling
	var got RepositoryList
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Verify fields
	if got.Count != repoList.Count {
		t.Errorf("Unmarshal Count = %v, want %v", got.Count, repoList.Count)
	}
	if len(got.Results) != len(repoList.Results) {
		t.Errorf("Unmarshal Results length = %v, want %v", len(got.Results), len(repoList.Results))
	}
}

func TestTagListJSON(t *testing.T) {
	tagList := &TagList{
		Count:    3,
		Next:     "https://next.page",
		Previous: "https://prev.page",
		Results: []*Tag{
			{Name: "v1.0.0", FullSize: 1000},
			{Name: "v1.1.0", FullSize: 2000},
			{Name: "latest", FullSize: 3000},
		},
	}

	// Test marshaling
	data, err := json.Marshal(tagList)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	// Test unmarshaling
	var got TagList
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Verify fields
	if got.Count != tagList.Count {
		t.Errorf("Unmarshal Count = %v, want %v", got.Count, tagList.Count)
	}
	if len(got.Results) != len(tagList.Results) {
		t.Errorf("Unmarshal Results length = %v, want %v", len(got.Results), len(tagList.Results))
	}
}
