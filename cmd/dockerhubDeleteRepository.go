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
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ealebed/dha/pkg/dockerhub"
)

// DeleteRepositoryOptions represents options for docker delete repository command
type DeleteRepositoryOptions struct {
	imageName string
}

// NewDockerhubDeleteRepositoryCmd returns new docker delete repository command
func NewDockerhubDeleteRepositoryCmd() *cobra.Command {
	options := DeleteRepositoryOptions{}

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "delete the specified docker repository",
		Long:    "delete the specified docker repository",
		Example: "dha delete [--image=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteRepository(cmd.InheritedFlags(), options.imageName)
		},
	}

	cmd.Flags().StringVarP(&options.imageName, "image", "i", "", "docker image name for delete")
	if err := cmd.MarkFlagRequired("image"); err != nil {
		// Flag marking should not fail in normal operation
		return nil
	}

	return cmd
}

// deleteRepository deletes docker repository
func deleteRepository(flags *pflag.FlagSet, image string) error {
	org, dryRun, err := dockerhub.GetFlags(flags)
	if err != nil {
		color.Red("Error: %s", err)
	}

	if dryRun {
		color.Yellow("[DRY-RUN] Delete docker image repository: %s/%s", dockerhub.BW(org), dockerhub.BW(image))
	} else {
		color.Blue("===> %s %s", dockerhub.BW("Deleting docker image repository"), dockerhub.BG(org+"/"+image))
		if err := dockerhub.NewClient(org, "").DeleteRepository(image); err != nil {
			return fmt.Errorf("failed to delete repository: %w", err)
		}
		color.Green("Done \u2714")
	}

	return nil
}
