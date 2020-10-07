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

// NewDockerhubListRepositoriesCmd returns new docker repositories list command
func NewDockerhubListRepositoriesCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "returns list all dockerhub repositories",
		Long:    "returns list all dockerhub organization repositories",
		Example: "dha list",
		Run: func(cmd *cobra.Command, args []string) {
			listDockerhubRepos(cmd.InheritedFlags())
		},
	}

	return cmd
}

// listDockerhubRepos returns list of all Dockerhub repositories
func listDockerhubRepos(flags *pflag.FlagSet) {
	org, _, err := dockerhub.GetFlags(flags)
	if err != nil {
		color.Red("Error: %s", err)
	}

	repositories, err := dockerhub.NewClient(org, "").ListRepositories()
	if err != nil {
		color.Red("Error: %s", err)
	}

	for repoCount, repo := range repositories {
		tagsCount, err := dockerhub.NewClient(org, "").GetTagsCount(repo.Name)
		if err != nil {
			color.Red("Error: %s", err)
		}

		if tagsCount == 0 {
			fmt.Printf("| Image %-5d | %-55s | [%s]\n", repoCount+1, dockerhub.BW(repo.Name), dockerhub.BR(tagsCount))
		} else if tagsCount >= 50 {
			fmt.Printf("| Image %-5d | %-55s | [%s]\n", repoCount+1, dockerhub.BW(repo.Name), dockerhub.BY(tagsCount))
		} else {
			fmt.Printf("| Image %-5d | %-55s | [%s]\n", repoCount+1, dockerhub.BW(repo.Name), dockerhub.BW(tagsCount))
		}
	}
}
