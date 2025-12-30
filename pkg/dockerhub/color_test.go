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
	"strings"
	"testing"
)

func TestColorFunctions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(...interface{}) string
		input    string
		validate func(*testing.T, string)
	}{
		{
			name:  "BG with valid input",
			fn:    BG,
			input: "test message",
			validate: func(t *testing.T, got string) {
				if got == "" {
					t.Error("BG() returned empty string")
				}
				if !strings.Contains(got, "test message") {
					t.Errorf("BG() output = %v, should contain 'test message'", got)
				}
			},
		},
		{
			name:  "BW with valid input",
			fn:    BW,
			input: "white text",
			validate: func(t *testing.T, got string) {
				if got == "" {
					t.Error("BW() returned empty string")
				}
				if !strings.Contains(got, "white text") {
					t.Errorf("BW() output = %v, should contain 'white text'", got)
				}
			},
		},
		{
			name:  "BY with valid input",
			fn:    BY,
			input: "yellow text",
			validate: func(t *testing.T, got string) {
				if got == "" {
					t.Error("BY() returned empty string")
				}
				if !strings.Contains(got, "yellow text") {
					t.Errorf("BY() output = %v, should contain 'yellow text'", got)
				}
			},
		},
		{
			name:  "BR with valid input",
			fn:    BR,
			input: "red text",
			validate: func(t *testing.T, got string) {
				if got == "" {
					t.Error("BR() returned empty string")
				}
				if !strings.Contains(got, "red text") {
					t.Errorf("BR() output = %v, should contain 'red text'", got)
				}
			},
		},
		{
			name:  "BG with empty string",
			fn:    BG,
			input: "",
			validate: func(t *testing.T, got string) {
				// Empty input returns empty string, which is expected behavior
				if got != "" {
					t.Errorf("BG() with empty input should return empty string, got %v", got)
				}
			},
		},
		{
			name:  "BW with special characters",
			fn:    BW,
			input: "test@123!@#$%",
			validate: func(t *testing.T, got string) {
				if got == "" {
					t.Error("BW() returned empty string")
				}
			},
		},
		{
			name:  "BY with multiline input",
			fn:    BY,
			input: "line1\nline2\nline3",
			validate: func(t *testing.T, got string) {
				if got == "" {
					t.Error("BY() returned empty string")
				}
			},
		},
		{
			name:  "BR with unicode characters",
			fn:    BR,
			input: "测试消息 🚀",
			validate: func(t *testing.T, got string) {
				if got == "" {
					t.Error("BR() returned empty string")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.fn(tt.input)
			tt.validate(t, got)
		})
	}
}

func TestColorFunctionsConsistency(t *testing.T) {
	input := "test input"

	bgResult := BG(input)
	bwResult := BW(input)
	byResult := BY(input)
	brResult := BR(input)

	// All functions should return non-empty strings
	if bgResult == "" || bwResult == "" || byResult == "" || brResult == "" {
		t.Error("All color functions should return non-empty strings")
	}

	// All results should contain the input text
	if !strings.Contains(bgResult, input) {
		t.Errorf("BG() result should contain input: %v", bgResult)
	}
	if !strings.Contains(bwResult, input) {
		t.Errorf("BW() result should contain input: %v", bwResult)
	}
	if !strings.Contains(byResult, input) {
		t.Errorf("BY() result should contain input: %v", byResult)
	}
	if !strings.Contains(brResult, input) {
		t.Errorf("BR() result should contain input: %v", brResult)
	}
}
