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
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ealebed/dha/pkg/dockerhub"
)

// RenewTagsOptions represents options for list tags command
type RenewTagsOptions struct {
	imageName string
	allImages bool
}

// NewDockerhubRenewTagsCmd returns new docker list tags command
func NewDockerhubRenewTagsCmd() *cobra.Command {
	options := RenewTagsOptions{}

	cmd := &cobra.Command{
		Use:     "renew",
		Short:   "renew tags from the provided dockerhub repository (image)",
		Long:    "renew tags from the provided dockerhub repository (image) or all organization repositories",
		Example: "dha renew [--image=...] || [--all]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return renewImageTags(cmd.InheritedFlags(), options.imageName, options.allImages)
		},
	}

	cmd.Flags().StringVarP(&options.imageName, "image", "i", "", "docker image name for getting tags")
	cmd.Flags().BoolVar(&options.allImages, "all", false, "renew tags in all organization repositories")

	return cmd
}

// renewImageTags renew tags from the provided dockerhub repository (image)
func renewImageTags(flags *pflag.FlagSet, image string, allImages bool) error {
	boldGreen := color.New(color.FgGreen, color.Bold).SprintFunc()
	boldWhite := color.New(color.FgWhite, color.Bold).SprintFunc()

	org, err := flags.GetString("org")
	if err != nil {
		color.Red("Error: %s", err)
	}

	dryRun, err := flags.GetBool("dry-run")
	if err != nil {
		color.Red("Error: %s", err)
	}

	if dryRun {
		color.Yellow("[DRY-RUN] Renewing tags for docker image repository: %s/%s", boldWhite(org), boldWhite(image))
	} else {
		if !allImages && image == "" {
			color.Red("You should provide image or set flag --all")
			os.Exit(1)
		} else if allImages && image == "" {
			repositories, err := dockerhub.NewClient(org, "").ListRepositories()
			if err != nil {
				color.Red("Error: %s", err)
			}
			for repoCount, repo := range repositories {
				color.Blue("===> %s %s %s/%s ", boldWhite("Processing docker image repository"), boldGreen(org+"/"+repo.Name), boldWhite(repoCount+1), boldWhite(len(repositories)))
				dockerhub.NewClient(org, "").RenewDockerImage(repo.Name)
			}
		} else {
			dockerhub.NewClient(org, "").RenewDockerImage(image)
		}
	}

	return nil
}
