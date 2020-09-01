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

package dockerhub

import (
	"os/exec"
	"regexp"
	"time"

	"github.com/fatih/color"
)

// RenewDockerImage renew docker image tags older than 30 days from docker hub
func (c *Client) RenewDockerImage(image string) error {
	tags, err := c.ListTags(image)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("UTC")
	currentTime := time.Now().In(loc)
	expiredRange := (time.Hour * 24 * 20)
	validTag := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{2}\-\d{2}\.\d{2}$`)

	for _, tag := range tags {
		imageReference := c.ORG + "/" + image + ":" + tag.Name
		if !validTag.MatchString(tag.Name) {
			color.Yellow("	Skip %s - invalid tag", imageReference)
		} else {
			lastUpdatedAt := tag.LastUpdated.In(loc)

			diff := currentTime.Sub(lastUpdatedAt)
			if diff.Hours() > expiredRange.Hours() {
				commandPull(imageReference)
				commandPush(imageReference)
				commandRmi(imageReference)
			} else {
				color.Yellow("	Skip %s - tag newer than %v hours", imageReference, expiredRange.Hours())
			}
		}
	}

	return nil
}

// helper function to create the `docker pull` command.
func commandPull(imageReference string) {
	color.Green("	<== Pulling from dockerHub %s", imageReference)
	exec.Command("docker", "pull", imageReference).Run()

	// debug
	// cmd := exec.Command("docker", "pull", imageReference)
	// stdout, err := cmd.Output()
	// if err != nil {
	// 	color.Red("Error while pulling docker image: %s", err)
	// }

	// fmt.Println(string(stdout))
}

// helper function to create the `docker push` command.
func commandPush(imageReference string) {
	color.Green("	==> Pushing to dockerHub %s", imageReference)
	exec.Command("docker", "push", imageReference).Run()
}

// helper function to create the `docker image rm` command.
func commandRmi(imageReference string) {
	color.Green("	Removing from localhost %s", imageReference)
	exec.Command("docker", "image", "rm", imageReference).Run()
}
