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

	"github.com/spf13/pflag"
)

func BenchmarkGetFlags(b *testing.B) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("org", "testorg", "organization")
	fs.Bool("dry-run", false, "dry run")
	fs.Set("org", "myorg")
	fs.Set("dry-run", "false")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = GetFlags(fs)
	}
}

func BenchmarkNewClient(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewClient("testorg", "https://test.com")
	}
}

func BenchmarkColorBG(b *testing.B) {
	input := "test message"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BG(input)
	}
}

func BenchmarkColorBW(b *testing.B) {
	input := "test message"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BW(input)
	}
}

func BenchmarkColorBY(b *testing.B) {
	input := "test message"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BY(input)
	}
}

func BenchmarkColorBR(b *testing.B) {
	input := "test message"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BR(input)
	}
}

func BenchmarkColorFunctionsLongInput(b *testing.B) {
	input := strings.Repeat("test message ", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BG(input)
		_ = BW(input)
		_ = BY(input)
		_ = BR(input)
	}
}

func BenchmarkNewRequest(b *testing.B) {
	client := NewClient("testorg", "https://test.com")
	client.AuthToken = "test-token-123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, _ := client.NewRequest("GET", "https://example.com/api", nil)
		if req != nil {
			_ = req.Method
		}
	}
}
