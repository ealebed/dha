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

// listRepoOptions represents options for list command
type listRepoOptions struct {
	expand bool
}

// NewDockerhubListRepositoriesCmd returns new docker repositories list command
func NewDockerhubListRepositoriesCmd() *cobra.Command {
	options := &listRepoOptions{}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "returns list all dockerhub repositories",
		Long:    "returns list all dockerhub organization repositories",
		Example: "dha list",
		Run: func(cmd *cobra.Command, args []string) {
			listDockerhubRepos(cmd.InheritedFlags(), options)
		},
	}

	// Note that false here means defaults to false, and flips to true if the flag is present.
	cmd.PersistentFlags().BoolVarP(&options.expand, "expand", "x", false, "expand docker repositories list payload to include size and pull count")

	return cmd
}

// listDockerhubRepos returns list of all Dockerhub repositories
func listDockerhubRepos(flags *pflag.FlagSet, options *listRepoOptions) {
	org, _, err := dockerhub.GetFlags(flags)
	if err != nil {
		color.Red("Error: %s", err)
	}

	repositories, err := dockerhub.NewClient(org, "").ListRepositories()
	if err != nil {
		color.Red("Error: %s", err)
	}

	if options.expand {
		fmt.Printf("| Image Num   | %-44s | %-7s | %-7s | %s\n", "Name", "Pulls ", "AvgSize (MB)", "Tags Count")
	} else {
		fmt.Printf("| Image Num   | %-55s | %s\n", "Name", "Tags Count")
	}

	for repoCount, repo := range repositories {
		tagsCount, err := dockerhub.NewClient(org, "").GetTagsCount(repo.Name)
		if err != nil {
			color.Red("Error: %s", err)
		}

		if options.expand {
			avgSize, err := dockerhub.NewClient(org, "").GetAvgTagsSize(repo.Name)
			if err != nil {
				color.Red("Error: %s", err)
			}
			if tagsCount == 0 {
				fmt.Printf("| Image %-5d | %-55s | %-7d | %-12.2f | [%s]\n", repoCount+1, dockerhub.BW(repo.Name), repo.PullCount, avgSize, dockerhub.BR(tagsCount))
			} else if tagsCount >= 50 {
				fmt.Printf("| Image %-5d | %-55s | %-7d | %-12.2f | [%s]\n", repoCount+1, dockerhub.BW(repo.Name), repo.PullCount, avgSize, dockerhub.BY(tagsCount))
			} else {
				fmt.Printf("| Image %-5d | %-55s | %-7d | %-12.2f | [%s]\n", repoCount+1, dockerhub.BW(repo.Name), repo.PullCount, avgSize, dockerhub.BW(tagsCount))
			}
		}
		if !options.expand {
			fmt.Printf("| Image %-5d | %-55s | [%d]\n", repoCount+1, repo.Name, tagsCount)
		}
	}
}
