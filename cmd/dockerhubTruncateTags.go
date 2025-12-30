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
	"os"
	"regexp"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/ealebed/dha/pkg/dockerhub"
)

// TruncateTagsOptions represents options for truncate command
type TruncateTagsOptions struct {
	imageName            string
	imageNameRegex       string
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
		Example: "dha truncate [--image=...] || [--imageRegEx=...] || [--all] [--inactive=...] || [--tagRegEx=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return truncateTags(cmd.InheritedFlags(), options.imageName, options.imageNameRegex, options.allImages, options.truncateInactiveTags, options.imageTagRegex)
		},
	}

	cmd.Flags().StringVarP(&options.imageName, "image", "i", "", "docker image name for truncating tags")
	cmd.Flags().StringVar(&options.imageNameRegex, "imageRegEx", "", "docker image name, matching specified regular expression string")
	cmd.Flags().BoolVar(&options.allImages, "all", false, "truncate tags in all organization repositories")
	cmd.Flags().BoolVar(&options.truncateInactiveTags, "inactive", false, "truncate inactive image tags (tags that haven't been pushed or pulled in over a month)")
	cmd.Flags().StringVar(&options.imageTagRegex, "tagRegEx", "", "truncate image tags, matching specified regular expression string")

	return cmd
}

// truncateTags truncate tags in docker repository except latest 30 ones
func truncateTags(flags *pflag.FlagSet, image, imageRegex string, allImages, truncateInactive bool, tagRegex string) error {
	org, dryRun, err := dockerhub.GetFlags(flags)
	if err != nil {
		color.Red("Error: %s", err)
	}

	if dryRun {
		color.Yellow("[DRY-RUN] Truncating tags for docker image repository: %s/%s", dockerhub.BW(org), dockerhub.BW(image))
		return nil
	}

	if err := validateTruncateFlags(truncateInactive, tagRegex, allImages, image, imageRegex); err != nil {
		return err
	}

	if allImages && (image == "" || imageRegex == "") {
		return truncateAllRepositories(org, tagRegex, truncateInactive)
	}

	if !allImages && image == "" && imageRegex != "" {
		return truncateRepositoriesByRegex(org, imageRegex, truncateInactive, tagRegex)
	}

	return truncateSingleRepository(org, image, truncateInactive, tagRegex)
}

func validateTruncateFlags(truncateInactive bool, tagRegex string, allImages bool, image, imageRegex string) error {
	if !truncateInactive && tagRegex == "" {
		color.Red("You should provide RegExp for image tag or set flag '--inactive'")
		os.Exit(1)
	}
	if !allImages && image == "" && imageRegex == "" {
		color.Red("You should provide image (fixed name or RegExp) or set flag '--all'")
		os.Exit(1)
	}
	return nil
}

func truncateAllRepositories(org, tagRegex string, truncateInactive bool) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	availableRoutines := runtime.NumCPU()
	routineReady := make(chan bool)

	repositories, err := dockerhub.NewClient(org, "").ListRepositories()
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}

	limiter := time.Tick(300 * time.Millisecond)

	for repoCount, repo := range repositories {
		<-limiter
		if availableRoutines == 0 {
			<-routineReady
			availableRoutines++
		}
		availableRoutines--

		go truncater(repoCount, len(repositories), org, tagRegex, truncateInactive, repo, routineReady)
	}

	for availableRoutines < runtime.NumCPU() {
		<-routineReady
		availableRoutines++
	}

	return nil
}

func truncateRepositoriesByRegex(org, imageRegex string, truncateInactive bool, tagRegex string) error {
	repositories, err := dockerhub.NewClient(org, "").ListRepositories()
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}

	regexPattern := fmt.Sprintf(`(?i)%s`, imageRegex)
	var repositoriesToTruncate []string
	for _, repo := range repositories {
		matched, _ := regexp.MatchString(regexPattern, repo.Name)
		if matched {
			repositoriesToTruncate = append(repositoriesToTruncate, repo.Name)
		}
	}

	for _, image := range repositoriesToTruncate {
		color.Blue("===> %s %s ", dockerhub.BW("Processing docker image repository"), dockerhub.BG(org+"/"+image))
		if err := dockerhub.NewClient(org, "").TruncateTags(image, truncateInactive, tagRegex); err != nil {
			color.Red("Error truncating tags for %s: %v", image, err)
		}
		dockerhub.BG("Done \u2714")
	}

	return nil
}

func truncateSingleRepository(org, image string, truncateInactive bool, tagRegex string) error {
	color.Blue("===> %s %s ", dockerhub.BW("Processing docker image repository"), dockerhub.BG(org+"/"+image))
	if err := dockerhub.NewClient(org, "").TruncateTags(image, truncateInactive, tagRegex); err != nil {
		return fmt.Errorf("failed to truncate tags: %w", err)
	}
	dockerhub.BG("Done \u2714")
	return nil
}

func truncater(repoCount, repositories int, org, tagRegex string, truncateInactive bool, repo *dockerhub.Repository, routineReady chan bool) {
	msg := "Processing docker image repository"
	repoName := org + "/" + repo.Name
	color.Blue("===> %s %s %s/%s ", dockerhub.BW(msg), dockerhub.BG(repoName), dockerhub.BW(repoCount+1), dockerhub.BW(repositories))
	if err := dockerhub.NewClient(org, "").TruncateTags(repo.Name, truncateInactive, tagRegex); err != nil {
		color.Red("Error truncating tags for %s: %v", repo.Name, err)
	}
	dockerhub.BG("Done \u2714")

	routineReady <- true
}
