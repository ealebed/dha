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

// ListTagsOptions represents options for list tags command
type ListTagsOptions struct {
	imageName string
}

// NewDockerhubListTagsCmd returns new docker list tags command
func NewDockerhubListTagsCmd() *cobra.Command {
	options := ListTagsOptions{}

	cmd := &cobra.Command{
		Use:     "get",
		Short:   "returns list tags from the provided dockerhub repository (image)",
		Long:    "returns list tags from the provided dockerhub repository (image)",
		Example: "dha get [--image=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listImageTags(cmd.InheritedFlags(), options.imageName)
		},
	}

	cmd.Flags().StringVarP(&options.imageName, "image", "i", "", "docker image name for getting tags")
	cmd.MarkFlagRequired("image")

	return cmd
}

// listImageTags returns list tags from the provided dockerhub repository (image)
func listImageTags(flags *pflag.FlagSet, image string) error {
	org, _, err := dockerhub.GetFlags(flags)
	if err != nil {
		color.Red("Error: %s", err)
	}

	tags, err := dockerhub.NewClient(org, "").ListTags(image)
	if err != nil {
		color.Red("Error: %s", err)
	}

	for count, tag := range tags {
		fmt.Printf("| Tag %-3d | %-30s | %s\n", count+1, dockerhub.BW(tag.Name), tag.LastUpdated)
	}

	return nil
}
