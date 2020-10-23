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
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/fatih/color"
)

// ListTags returns list of docker image tags for selected image from docker hub
func (c *Client) ListTags(image string) ([]*Tag, error) {
	var tags = []*Tag{}
	output, err := c.listTagsRequest(image, "")
	if err != nil {
		return nil, err
	}

	tags = append(tags, output.Results...)
	next := output.Next

	for {
		if next == "" {
			return tags, nil
		}
		output, err := c.listTagsRequest(image, next)
		if err != nil {
			return nil, err
		}

		tags = append(tags, output.Results...)
		next = output.Next
	}
}

func (c *Client) listTagsRequest(image, next string) (*TagList, error) {
	var url string
	if next != "" {
		url = fmt.Sprint(next)
	} else {
		url = fmt.Sprintf("%s/%s/%s/tags/?page_size=100", RepositoriesURL, c.ORG, image)
	}

	data, _, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	output := &TagList{}

	if err = json.NewDecoder(bytes.NewReader(data)).Decode(output); err != nil {
		return nil, err
	}

	return output, nil
}

// GetTagsCount returns count docker image tag from docker hub for selected repository
func (c *Client) GetTagsCount(image string) (int, error) {
	url := fmt.Sprintf("%s/%s/%s/tags/?page_size=100", RepositoriesURL, c.ORG, image)

	data, _, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return -1, err
	}

	output := &TagList{}

	if err = json.NewDecoder(bytes.NewReader(data)).Decode(output); err != nil {
		return -1, err
	}

	tagsCount := output.Count

	return tagsCount, nil
}

// deleteDockerImageTag delete docker image tag from docker hub
/* curl \
   -H "Authorization: JWT ${TOKEN}" \
   -X DELETE https://hub.docker.com/v2/repositories/${ORG}/${IMAGE}/tags/${TAG}/
*/
func (c *Client) deleteDockerImageTag(image string, tag string) error {
	if _, _, err := c.doRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s/tags/%s/", RepositoriesURL, c.ORG, image, tag), nil); err != nil {
		color.Red("Error while deleting docker image tag: %s", err)
	}

	return nil
}

// GetLatestTag returns latest (by LastUpdated field) docker image tag from docker hub
func (c *Client) GetLatestTag(image string) (string, error) {
	tags, err := c.ListTags(image)
	if err != nil {
		return "", err
	}

	return tags[0].Name, nil
}

// TruncateTags deletes docker image tags tags that match `regularExpression` OR are older than `expiredRange` except latest `leaveTagsCounter` ones
func (c *Client) TruncateTags(image string, truncateOld bool, regularExpression string) error {
	var tagsToRemove []string
	var leaveTagsCounter = 0

	tags, err := c.ListTags(image)
	if err != nil {
		return err
	}

	loc, _ := time.LoadLocation("UTC")
	currentTime := time.Now().In(loc)
	expiredRange := (time.Hour * 24 * 30)

	if regularExpression != "" {
		regexPattern := fmt.Sprintf(`(?i)%s`, regularExpression)
		for _, tag := range tags {
			matched, _ := regexp.MatchString(regexPattern, tag.Name)
			if matched {
				tagsToRemove = append(tagsToRemove, tag.Name)
			}
		}
	} else {
		for _, tag := range tags {
			lastUpdatedAt := tag.LastUpdated.In(loc)
			diff := currentTime.Sub(lastUpdatedAt)
			if diff.Hours() > expiredRange.Hours() {
				tagsToRemove = append(tagsToRemove, tag.Name)
				leaveTagsCounter = 25
			}
		}
	}

	for i := leaveTagsCounter; i < len(tagsToRemove); i++ {
		color.Green("\u2714  Delete tag %s", BW(tagsToRemove[i]))
		if err := c.deleteDockerImageTag(image, tagsToRemove[i]); err != nil {
			color.Red("Error while deleting image tag: %s", err)
		}
	}

	return nil
}
