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
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ealebed/dha/pkg/dockerhub"
)

// listRepoOptions represents options for list command
type listRepoOptions struct {
	expand bool
}

// listResult represents data for use in channel from goroutines
type listResult struct {
	repoName      string
	repoPullCount int
	avgSize       float64
	tagsCount     int
	lastUpdated   string
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
	var ret = []listResult{}
	var wg sync.WaitGroup
	chanRes := make(chan listResult)

	org, _, err := dockerhub.GetFlags(flags)
	if err != nil {
		color.Red("Error: %s", err)
	}

	repositories, err := dockerhub.NewClient(org, "").ListRepositories()
	if err != nil {
		color.Red("Error: %s", err)
	}

	go func() {
		for {
			r, opened := <-chanRes
			if !opened {
				break
			}
			ret = append(ret, r)
		}
	}()

	limiter := time.Tick(300 * time.Millisecond)

	for _, repo := range repositories {
		<-limiter
		wg.Add(1)
		go lister(org, repo, chanRes, &wg)
	}

	wg.Wait()
	close(chanRes)

	if options.expand {
		fmt.Printf("| Image Num   | %-44s | %-7s | %-7s | %-7s | %s\n", "Name", "Pulls Count", "AvgSize (MB)", "Tags Count", "Last Updated")
	}
	if !options.expand {
		fmt.Printf("| Image Num   | %-55s | %s\n", "Name", "Tags Count")
	}

	for repoCount, info := range ret {
		if options.expand {
			if info.tagsCount == 0 {
				fmt.Printf("| Image %-5d | %-55s | %-11d | %-12.2f | %-21s | %s\n", repoCount+1, dockerhub.BW(info.repoName), info.repoPullCount, info.avgSize, dockerhub.BR(info.tagsCount), info.lastUpdated)
			} else if info.tagsCount >= 50 {
				fmt.Printf("| Image %-5d | %-55s | %-11d | %-12.2f | %-21s | %s\n", repoCount+1, dockerhub.BW(info.repoName), info.repoPullCount, info.avgSize, dockerhub.BY(info.tagsCount), info.lastUpdated)
			} else {
				fmt.Printf("| Image %-5d | %-55s | %-11d | %-12.2f | %-21s | %s\n", repoCount+1, dockerhub.BW(info.repoName), info.repoPullCount, info.avgSize, dockerhub.BW(info.tagsCount), info.lastUpdated)
			}
		}
		if !options.expand {
			fmt.Printf("| Image %-5d | %-55s | %d\n", repoCount+1, info.repoName, info.tagsCount)
		}
	}

}

func lister(org string, repo *dockerhub.Repository, chanRes chan listResult, wg *sync.WaitGroup) {
	defer wg.Done()

	r := listResult{}

	tagsCount, err := dockerhub.NewClient(org, "").GetTagsCount(repo.Name)
	if err != nil {
		color.Red("Error: %s", err)
	}

	avgSize, err := dockerhub.NewClient(org, "").GetAvgTagsSize(repo.Name)
	if err != nil {
		color.Red("Error: %s", err)
	}

	r.repoName = repo.Name
	r.repoPullCount = repo.PullCount
	r.avgSize = avgSize
	r.tagsCount = tagsCount
	r.lastUpdated = repo.LastUpdated.String()

	chanRes <- r
}
