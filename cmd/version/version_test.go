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

package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name         string
		version      string
		releasePhase string
		want         string
	}{
		{
			name:         "final release without phase",
			version:      "v1.0.0",
			releasePhase: "",
			want:         "v1.0.0",
		},
		{
			name:         "pre-release with dev phase",
			version:      "v1.0.0",
			releasePhase: "dev",
			want:         "v1.0.0-dev",
		},
		{
			name:         "pre-release with alpha phase",
			version:      "v0.5.0",
			releasePhase: "alpha",
			want:         "v0.5.0-alpha",
		},
		{
			name:         "pre-release with beta phase",
			version:      "v2.1.0",
			releasePhase: "beta",
			want:         "v2.1.0-beta",
		},
		{
			name:         "pre-release with rc phase",
			version:      "v1.2.3",
			releasePhase: "rc",
			want:         "v1.2.3-rc",
		},
		{
			name:         "empty version with phase",
			version:      "",
			releasePhase: "dev",
			want:         "-dev",
		},
		{
			name:         "empty version without phase",
			version:      "",
			releasePhase: "",
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			origVersion := Version
			origPhase := ReleasePhase

			// Set test values
			Version = tt.version
			ReleasePhase = tt.releasePhase

			// Run test
			got := String()

			// Restore original values
			Version = origVersion
			ReleasePhase = origPhase

			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringCurrentVersion(t *testing.T) {
	// Test that the current version function works
	got := String()

	// Should not be empty
	if got == "" {
		t.Error("String() should not return empty string for current version")
	}

	// Should contain the version
	if Version != "" && !strings.Contains(got, Version) {
		t.Errorf("String() = %v, should contain Version %v", got, Version)
	}

	// If ReleasePhase is set, should contain it
	if ReleasePhase != "" && !strings.Contains(got, ReleasePhase) {
		t.Errorf("String() = %v, should contain ReleasePhase %v", got, ReleasePhase)
	}
}
