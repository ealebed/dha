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
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/ealebed/dha/cmd/version"
)

// RootOptions implements global flags for all commands
type RootOptions struct {
	organization string
	dryRun       bool
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute(out io.Writer) error {
	cmd := NewCmdRoot(out)
	return cmd.Execute()
}

// NewCmdRoot returns new root command
func NewCmdRoot(out io.Writer) *cobra.Command {
	options := RootOptions{}

	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.String(),
	}

	cmd.PersistentFlags().StringVar(&options.organization, "org", os.Getenv("DOCKERHUB_USERNAME"), "repository source owner (user/organization)")
	cmd.PersistentFlags().BoolVar(&options.dryRun, "dry-run", true, "print output only")

	// create subcommands
	cmd.AddCommand(NewDockerhubDeleteRepositoryCmd())
	cmd.AddCommand(NewDockerhubDescribeRepositoryCmd())
	cmd.AddCommand(NewDockerhubListRepositoriesCmd())
	cmd.AddCommand(NewDockerhubListTagsCmd())
	cmd.AddCommand(NewDockerhubRenewTagsCmd())
	cmd.AddCommand(NewDockerhubTruncateTagsCmd())

	return cmd
}
