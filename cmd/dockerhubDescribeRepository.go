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

// DescribeRepositoryOptions represents options for get command
type DescribeRepositoryOptions struct {
	imageName string
}

// NewDockerhubDescribeRepositoryCmd returns new docker get repository command
func NewDockerhubDescribeRepositoryCmd() *cobra.Command {
	options := DescribeRepositoryOptions{}

	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "returns info about provided dockerhub repository (image)",
		Long:    "returns detailed information about provided dockerhub repository (image)",
		Example: "dha describe [--image=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return describeRepository(cmd.InheritedFlags(), options.imageName)
		},
	}

	cmd.Flags().StringVarP(&options.imageName, "image", "i", "", "docker image name for getting information")
	if err := cmd.MarkFlagRequired("image"); err != nil {
		// Flag marking should not fail in normal operation
		return nil
	}

	return cmd
}

// describeRepository returns information about the provided dockerhub repository (image)
func describeRepository(flags *pflag.FlagSet, image string) error {
	org, _, err := dockerhub.GetFlags(flags)
	if err != nil {
		color.Red("Error: %s", err)
	}

	repoInfo, err := dockerhub.NewClient(org, "").DescribeRepository(image)
	if err != nil {
		color.Red("Error: %s", err)
	}

	color.Blue("User: " + repoInfo.User +
		"\nName: " + repoInfo.Name +
		"\nNamespace: " + repoInfo.Namespace +
		"\nRepositoryType: " + repoInfo.RepositoryType +
		"\nStatus: " + fmt.Sprintf("%d", repoInfo.Status) +
		"\nDescription: " + repoInfo.Description +
		"\nIsPrivate: " + fmt.Sprintf("%t", repoInfo.IsPrivate) +
		"\nIsAutomated: " + fmt.Sprintf("%t", repoInfo.IsAutomated) +
		"\nCanEdit: " + fmt.Sprintf("%t", repoInfo.CanEdit) +
		"\nStarCount: " + fmt.Sprintf("%d", repoInfo.StarCount) +
		"\nPullCount: " + fmt.Sprintf("%d", repoInfo.PullCount) +
		"\nLastUpdated: " + repoInfo.LastUpdated.String() +
		"\nIsMigrated: " + fmt.Sprintf("%t", repoInfo.IsMigrated) +
		"\nCollaboratorCount: " + fmt.Sprintf("%d", repoInfo.CollaboratorCount) +
		"\nAffiliation: " + repoInfo.Affiliation +
		"\nHubUser: " + repoInfo.HubUser)

	return nil
}
