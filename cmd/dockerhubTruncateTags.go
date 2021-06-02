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
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ealebed/dha/pkg/dockerhub"
)

// TruncateTagsOptions represents options for truncate command
type TruncateTagsOptions struct {
	imageName            string
	allImages            bool
	truncateInactiveTags bool
	imageTagRegex        string
}

// NewDockerhubTruncateTagsCmd returns new docker truncate tags command
func NewDockerhubTruncateTagsCmd() *cobra.Command {
	options := TruncateTagsOptions{}

	cmd := &cobra.Command{
		Use:     "truncate",
		Short:   "truncate tags in the specified docker repository",
		Long:    "truncate tags in the specified docker image repository (by default, except latest 30 ones)",
		Example: "dha truncate [--image=...] || [--all] [--inactive=...] || [--regEx=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return truncateTags(cmd.InheritedFlags(), options.imageName, options.allImages, options.truncateInactiveTags, options.imageTagRegex)
		},
	}

	cmd.Flags().StringVarP(&options.imageName, "image", "i", "", "docker image name for truncating tags")
	cmd.Flags().BoolVar(&options.allImages, "all", false, "truncate tags in all organization repositories")
	cmd.Flags().BoolVar(&options.truncateInactiveTags, "inactive", false, "truncate inactive image tags (tags that haven't been pushed or pulled in over a month)")
	cmd.Flags().StringVar(&options.imageTagRegex, "regEx", "", "truncate image tags, matching specified regular expression string")

	return cmd
}

// truncateTags truncate tags in docker repository except latest 30 ones
func truncateTags(flags *pflag.FlagSet, image string, allImages, truncateInactive bool, regEx string) error {
	org, dryRun, err := dockerhub.GetFlags(flags)
	if err != nil {
		color.Red("Error: %s", err)
	}

	if dryRun {
		color.Yellow("[DRY-RUN] Truncating tags for docker image repository: %s/%s", dockerhub.BW(org), dockerhub.BW(image))
	} else {
		if !truncateInactive && regEx == "" {
			color.Red("You should provide RegExp for image tag or set flag '--inactive'")
			os.Exit(1)
		}
		if !allImages && image == "" {
			color.Red("You should provide image or set flag '--all'")
			os.Exit(1)
		} else if allImages && image == "" {
			runtime.GOMAXPROCS(runtime.NumCPU())
			availableRoutines := runtime.NumCPU()
			routineReady := make(chan bool)

			repositories, err := dockerhub.NewClient(org, "").ListRepositories()
			if err != nil {
				color.Red("Error: %s", err)
			}

			for repoCount, repo := range repositories {
				if availableRoutines == 0 {
					<-routineReady
					availableRoutines = availableRoutines + 1
				}
				availableRoutines = availableRoutines - 1

				go truncater(repoCount, len(repositories), org, regEx, truncateInactive, repo, routineReady)
			}

			for availableRoutines < runtime.NumCPU() {
				<-routineReady
				availableRoutines = availableRoutines + 1
			}
		} else {
			color.Blue("===> %s %s ", dockerhub.BW("Processing docker image repository"), dockerhub.BG(org+"/"+image))
			dockerhub.NewClient(org, "").TruncateTags(image, truncateInactive, regEx)
			dockerhub.BG("Done \u2714")
		}
	}

	return nil
}

func truncater(repoCount, repositories int, org, regEx string, truncateInactive bool, repo *dockerhub.Repository, routineReady chan bool) {
	color.Blue("===> %s %s %s/%s ", dockerhub.BW("Processing docker image repository"), dockerhub.BG(org+"/"+repo.Name), dockerhub.BW(repoCount+1), dockerhub.BW(repositories))
	dockerhub.NewClient(org, "").TruncateTags(repo.Name, truncateInactive, regEx)
	dockerhub.BG("Done \u2714")

	routineReady <- true
}
