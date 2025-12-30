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
	"time"

	"github.com/spf13/pflag"
)

func TestGetFlagsEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		setupFlags func(*pflag.FlagSet)
		wantOrg    string
		wantDryRun bool
		wantErr    bool
	}{
		{
			name: "very long org name",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "", "organization")
				fs.Bool("dry-run", false, "dry run")
				longOrg := strings.Repeat("a", 1000)
				fs.Set("org", longOrg)
				fs.Set("dry-run", "false")
			},
			wantOrg:    strings.Repeat("a", 1000),
			wantDryRun: false,
			wantErr:    false,
		},
		{
			name: "org with special characters",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "", "organization")
				fs.Bool("dry-run", false, "dry run")
				fs.Set("org", "org-with-dashes_and_underscores")
				fs.Set("dry-run", "false")
			},
			wantOrg:    "org-with-dashes_and_underscores",
			wantDryRun: false,
			wantErr:    false,
		},
		{
			name: "org with unicode characters",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "", "organization")
				fs.Bool("dry-run", false, "dry run")
				fs.Set("org", "组织名称")
				fs.Set("dry-run", "false")
			},
			wantOrg:    "组织名称",
			wantDryRun: false,
			wantErr:    false,
		},
		{
			name: "whitespace only org",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "", "organization")
				fs.Bool("dry-run", false, "dry run")
				fs.Set("org", "   ")
				fs.Set("dry-run", "false")
			},
			wantOrg:    "   ",
			wantDryRun: false,
			wantErr:    false,
		},
		{
			name: "dry-run with various string representations",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "", "organization")
				fs.Bool("dry-run", false, "dry run")
				fs.Set("org", "testorg")
				fs.Set("dry-run", "1") // Alternative true representation
			},
			wantOrg:    "testorg",
			wantDryRun: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			tt.setupFlags(fs)

			gotOrg, gotDryRun, err := GetFlags(fs)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotOrg != tt.wantOrg {
				t.Errorf("GetFlags() org = %v, want %v", gotOrg, tt.wantOrg)
			}

			if gotDryRun != tt.wantDryRun {
				t.Errorf("GetFlags() dryRun = %v, want %v", gotDryRun, tt.wantDryRun)
			}
		})
	}
}

func TestNewClientEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		org       string
		url       string
		wantOrg   string
		wantURL   string
		checkFunc func(*testing.T, *Client)
	}{
		{
			name:    "very long org name",
			org:     strings.Repeat("a", 500),
			url:     "https://test.com",
			wantOrg: strings.Repeat("a", 500),
			wantURL: "https://test.com",
			checkFunc: func(t *testing.T, c *Client) {
				if len(c.ORG) != 500 {
					t.Errorf("NewClient() ORG length = %v, want 500", len(c.ORG))
				}
			},
		},
		{
			name:    "URL with port",
			org:     "testorg",
			url:     "https://test.com:8080",
			wantOrg: "testorg",
			wantURL: "https://test.com:8080",
			checkFunc: func(t *testing.T, c *Client) {
				if c.URL != "https://test.com:8080" {
					t.Errorf("NewClient() URL = %v, want https://test.com:8080", c.URL)
				}
			},
		},
		{
			name:    "URL with path",
			org:     "testorg",
			url:     "https://test.com/api/v1",
			wantOrg: "testorg",
			wantURL: "https://test.com/api/v1",
			checkFunc: func(t *testing.T, c *Client) {
				if c.URL != "https://test.com/api/v1" {
					t.Errorf("NewClient() URL = %v, want https://test.com/api/v1", c.URL)
				}
			},
		},
		{
			name:    "URL with query parameters",
			org:     "testorg",
			url:     "https://test.com?param=value",
			wantOrg: "testorg",
			wantURL: "https://test.com?param=value",
			checkFunc: func(t *testing.T, c *Client) {
				if c.URL != "https://test.com?param=value" {
					t.Errorf("NewClient() URL = %v, want https://test.com?param=value", c.URL)
				}
			},
		},
		{
			name:    "org with newlines",
			org:     "org\nwith\nnewlines",
			url:     "https://test.com",
			wantOrg: "org\nwith\nnewlines",
			wantURL: "https://test.com",
			checkFunc: func(t *testing.T, c *Client) {
				if c.ORG != "org\nwith\nnewlines" {
					t.Errorf("NewClient() ORG = %v, want org\\nwith\\nnewlines", c.ORG)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewClient(tt.org, tt.url)

			if got.ORG != tt.wantOrg {
				t.Errorf("NewClient() ORG = %v, want %v", got.ORG, tt.wantOrg)
			}

			if got.URL != tt.wantURL {
				t.Errorf("NewClient() URL = %v, want %v", got.URL, tt.wantURL)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, got)
			}
		})
	}
}

func TestClientInitialization(t *testing.T) {
	client := NewClient("testorg", "https://test.com")

	// Verify all fields are initialized
	if client.Client == nil {
		t.Error("Client.Client should not be nil")
	}

	if client.Header == nil {
		t.Error("Client.Header should not be nil")
	}

	if client.AuthToken != "" {
		t.Errorf("Client.AuthToken should be empty initially, got %v", client.AuthToken)
	}

	// Verify timeout is set correctly
	if client.Client.Timeout != 30*time.Second {
		t.Errorf("Client.Client.Timeout = %v, want %v", client.Client.Timeout, 30*time.Second)
	}

	// Verify header is set correctly
	if client.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Client.Header Content-Type = %v, want application/json", client.Header.Get("Content-Type"))
	}
}

func TestClientMultipleInstances(t *testing.T) {
	// Test that multiple clients are independent
	client1 := NewClient("org1", "https://test1.com")
	client2 := NewClient("org2", "https://test2.com")

	if client1.ORG == client2.ORG {
		t.Error("Client instances should have independent ORG values")
	}

	if client1.URL == client2.URL {
		t.Error("Client instances should have independent URL values")
	}

	// Verify they have separate HTTP clients
	if client1.Client == client2.Client {
		t.Error("Client instances should have separate HTTP clients")
	}

	// Verify they have separate headers (check by comparing Content-Type values)
	// Since maps can't be compared directly, we verify they're separate by checking
	// that modifying one doesn't affect the other
	originalCT := client2.Header.Get("Content-Type")
	client1.Header.Set("Content-Type", "modified")
	if client2.Header.Get("Content-Type") != originalCT {
		t.Error("Client instances should have separate headers")
	}
	// Restore original
	client1.Header.Set("Content-Type", "application/json")
}
