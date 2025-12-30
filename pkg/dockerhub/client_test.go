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
	"testing"
	"time"

	"github.com/spf13/pflag"
)

func TestGetFlags(t *testing.T) {
	tests := []struct {
		name       string
		setupFlags func(*pflag.FlagSet)
		wantOrg    string
		wantDryRun bool
		wantErr    bool
	}{
		{
			name: "valid flags with org and dry-run false",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "testorg", "organization")
				fs.Bool("dry-run", false, "dry run")
				fs.Set("org", "myorg")
				fs.Set("dry-run", "false")
			},
			wantOrg:    "myorg",
			wantDryRun: false,
			wantErr:    false,
		},
		{
			name: "valid flags with org and dry-run true",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "testorg", "organization")
				fs.Bool("dry-run", true, "dry run")
				fs.Set("org", "myorg")
				fs.Set("dry-run", "true")
			},
			wantOrg:    "myorg",
			wantDryRun: true,
			wantErr:    false,
		},
		{
			name: "missing org flag",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.Bool("dry-run", false, "dry run")
			},
			wantOrg:    "",
			wantDryRun: true,
			wantErr:    true,
		},
		{
			name: "missing dry-run flag",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "testorg", "organization")
				fs.Set("org", "myorg")
			},
			wantOrg:    "myorg",
			wantDryRun: true,
			wantErr:    true,
		},
		{
			name: "empty org value",
			setupFlags: func(fs *pflag.FlagSet) {
				fs.String("org", "", "organization")
				fs.Bool("dry-run", false, "dry run")
				fs.Set("org", "")
				fs.Set("dry-run", "false")
			},
			wantOrg:    "",
			wantDryRun: false,
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

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		org       string
		url       string
		wantOrg   string
		wantURL   string
		checkFunc func(*testing.T, *Client)
	}{
		{
			name:    "with custom URL",
			org:     "testorg",
			url:     "https://custom.docker.com",
			wantOrg: "testorg",
			wantURL: "https://custom.docker.com",
			checkFunc: func(t *testing.T, c *Client) {
				if c.Client == nil {
					t.Error("NewClient() Client should not be nil")
				}
				if c.Client.Timeout != 30*time.Second {
					t.Errorf("NewClient() Client.Timeout = %v, want %v", c.Client.Timeout, 30*time.Second)
				}
				if c.Header == nil {
					t.Error("NewClient() Header should not be nil")
				}
				if c.Header.Get("Content-Type") != "application/json" {
					t.Errorf("NewClient() Header Content-Type = %v, want application/json", c.Header.Get("Content-Type"))
				}
			},
		},
		{
			name:    "with empty URL (should use default)",
			org:     "myorg",
			url:     "",
			wantOrg: "myorg",
			wantURL: "https://hub.docker.com",
			checkFunc: func(t *testing.T, c *Client) {
				if c.Client == nil {
					t.Error("NewClient() Client should not be nil")
				}
				if c.Header == nil {
					t.Error("NewClient() Header should not be nil")
				}
			},
		},
		{
			name:    "with empty org",
			org:     "",
			url:     "https://test.com",
			wantOrg: "",
			wantURL: "https://test.com",
			checkFunc: func(t *testing.T, c *Client) {
				if c.Client == nil {
					t.Error("NewClient() Client should not be nil")
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
