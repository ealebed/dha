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

package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewCmdRoot(t *testing.T) {
	tests := []struct {
		name     string
		out      io.Writer
		validate func(*testing.T, *cobra.Command)
	}{
		{
			name: "with bytes.Buffer output",
			out:  &bytes.Buffer{},
			validate: func(t *testing.T, cmd *cobra.Command) {
				if cmd == nil {
					t.Fatal("NewCmdRoot() returned nil command")
				}
				if !cmd.SilenceUsage {
					t.Error("NewCmdRoot() SilenceUsage should be true")
				}
				if !cmd.SilenceErrors {
					t.Error("NewCmdRoot() SilenceErrors should be true")
				}
				if cmd.Version == "" {
					t.Error("NewCmdRoot() Version should be set")
				}
			},
		},
		{
			name: "with nil output",
			out:  nil,
			validate: func(t *testing.T, cmd *cobra.Command) {
				if cmd == nil {
					t.Fatal("NewCmdRoot() returned nil command")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCmdRoot(tt.out)
			if tt.validate != nil {
				tt.validate(t, cmd)
			}
		})
	}
}

func TestNewCmdRootSubcommands(t *testing.T) {
	cmd := NewCmdRoot(&bytes.Buffer{})

	// Verify all subcommands are added
	expectedCommands := []string{
		"delete", "del",
		"describe",
		"list", "ls",
		"get",
		"renew",
		"truncate",
	}

	commands := cmd.Commands()
	if len(commands) < 6 {
		t.Errorf("NewCmdRoot() should have at least 6 subcommands, got %d", len(commands))
	}

	// Check that expected commands exist
	commandMap := make(map[string]bool)
	for _, c := range commands {
		commandMap[c.Use] = true
		for _, alias := range c.Aliases {
			commandMap[alias] = true
		}
	}

	for _, expected := range expectedCommands {
		if !commandMap[expected] {
			t.Errorf("Expected command/alias '%s' not found", expected)
		}
	}
}

func TestNewCmdRootFlags(t *testing.T) {
	cmd := NewCmdRoot(&bytes.Buffer{})

	// Verify persistent flags are set
	orgFlag := cmd.PersistentFlags().Lookup("org")
	if orgFlag == nil {
		t.Error("NewCmdRoot() should have 'org' persistent flag")
	}

	dryRunFlag := cmd.PersistentFlags().Lookup("dry-run")
	if dryRunFlag == nil {
		t.Error("NewCmdRoot() should have 'dry-run' persistent flag")
	}

	// Verify flag types
	if orgFlag.Value.Type() != "string" {
		t.Errorf("org flag type = %v, want string", orgFlag.Value.Type())
	}

	if dryRunFlag.Value.Type() != "bool" {
		t.Errorf("dry-run flag type = %v, want bool", dryRunFlag.Value.Type())
	}
}

func TestNewDockerhubDeleteRepositoryCmd(t *testing.T) {
	cmd := NewDockerhubDeleteRepositoryCmd()

	if cmd == nil {
		t.Fatal("NewDockerhubDeleteRepositoryCmd() returned nil")
	}

	if cmd.Use != "delete" {
		t.Errorf("Command Use = %v, want delete", cmd.Use)
	}

	// Check aliases
	if len(cmd.Aliases) == 0 || cmd.Aliases[0] != "del" {
		t.Errorf("Command Aliases = %v, want [del]", cmd.Aliases)
	}

	// Check required flag
	imageFlag := cmd.Flags().Lookup("image")
	if imageFlag == nil {
		t.Error("Command should have 'image' flag")
	}
}

func TestNewDockerhubDescribeRepositoryCmd(t *testing.T) {
	cmd := NewDockerhubDescribeRepositoryCmd()

	if cmd == nil {
		t.Fatal("NewDockerhubDescribeRepositoryCmd() returned nil")
	}

	if cmd.Use != "describe" {
		t.Errorf("Command Use = %v, want describe", cmd.Use)
	}

	// Check required flag
	imageFlag := cmd.Flags().Lookup("image")
	if imageFlag == nil {
		t.Error("Command should have 'image' flag")
	}
}

func TestNewDockerhubListRepositoriesCmd(t *testing.T) {
	cmd := NewDockerhubListRepositoriesCmd()

	if cmd == nil {
		t.Fatal("NewDockerhubListRepositoriesCmd() returned nil")
	}

	if cmd.Use != "list" {
		t.Errorf("Command Use = %v, want list", cmd.Use)
	}

	// Check aliases
	if len(cmd.Aliases) == 0 || cmd.Aliases[0] != "ls" {
		t.Errorf("Command Aliases = %v, want [ls]", cmd.Aliases)
	}

	// Check expand flag
	expandFlag := cmd.PersistentFlags().Lookup("expand")
	if expandFlag == nil {
		t.Error("Command should have 'expand' flag")
	}
}

func TestNewDockerhubListTagsCmd(t *testing.T) {
	cmd := NewDockerhubListTagsCmd()

	if cmd == nil {
		t.Fatal("NewDockerhubListTagsCmd() returned nil")
	}

	if cmd.Use != "get" {
		t.Errorf("Command Use = %v, want get", cmd.Use)
	}

	// Check required flag
	imageFlag := cmd.Flags().Lookup("image")
	if imageFlag == nil {
		t.Error("Command should have 'image' flag")
	}
}

func TestNewDockerhubRenewTagsCmd(t *testing.T) {
	cmd := NewDockerhubRenewTagsCmd()

	if cmd == nil {
		t.Fatal("NewDockerhubRenewTagsCmd() returned nil")
	}

	if cmd.Use != "renew" {
		t.Errorf("Command Use = %v, want renew", cmd.Use)
	}

	// Check flags
	imageFlag := cmd.Flags().Lookup("image")
	if imageFlag == nil {
		t.Error("Command should have 'image' flag")
	}

	allFlag := cmd.Flags().Lookup("all")
	if allFlag == nil {
		t.Error("Command should have 'all' flag")
	}
}

func TestNewDockerhubTruncateTagsCmd(t *testing.T) {
	cmd := NewDockerhubTruncateTagsCmd()

	if cmd == nil {
		t.Fatal("NewDockerhubTruncateTagsCmd() returned nil")
	}

	if cmd.Use != "truncate" {
		t.Errorf("Command Use = %v, want truncate", cmd.Use)
	}

	// Check flags
	imageFlag := cmd.Flags().Lookup("image")
	if imageFlag == nil {
		t.Error("Command should have 'image' flag")
	}

	imageRegexFlag := cmd.Flags().Lookup("imageRegEx")
	if imageRegexFlag == nil {
		t.Error("Command should have 'imageRegEx' flag")
	}

	allFlag := cmd.Flags().Lookup("all")
	if allFlag == nil {
		t.Error("Command should have 'all' flag")
	}

	inactiveFlag := cmd.Flags().Lookup("inactive")
	if inactiveFlag == nil {
		t.Error("Command should have 'inactive' flag")
	}

	tagRegexFlag := cmd.Flags().Lookup("tagRegEx")
	if tagRegexFlag == nil {
		t.Error("Command should have 'tagRegEx' flag")
	}
}
