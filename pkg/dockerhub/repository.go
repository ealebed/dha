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

	"github.com/fatih/color"
)

// ListRepositories returns list of docker images from docker hub
func (c *Client) ListRepositories() (repos []*Repository, err error) {
	repos = []*Repository{}
	output, err := c.listRepositoriesRequest("")
	if err != nil {
		return nil, err
	}

	repos = append(repos, output.Results...)

	next := output.Next
	for {
		if next == "" {
			return repos, nil
		}
		output, err := c.listRepositoriesRequest(next)
		if err != nil {
			return nil, err
		}
		next = output.Next
		repos = append(repos, output.Results...)
	}
}

func (c *Client) listRepositoriesRequest(next string) (*RepositoryList, error) {
	var url string
	if next != "" {
		url = next
	} else {
		url = fmt.Sprintf("%s/%s/?page=1&page_size=100", RepositoriesURL, c.ORG)
	}

	data, _, err := c.doRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	output := &RepositoryList{}

	if err = json.NewDecoder(bytes.NewReader(data)).Decode(output); err != nil {
		return nil, err
	}

	return output, nil
}

// DescribeRepository print details about docker repository from docker hub
func (c *Client) DescribeRepository(image string) (*Repository, error) {
	data, _, err := c.doRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", RepositoriesURL, c.ORG, image), nil)
	if err != nil {
		return nil, err
	}

	repo := &Repository{}

	if err = json.NewDecoder(bytes.NewReader(data)).Decode(repo); err != nil {
		return nil, err
	}

	return repo, nil
}

// DeleteRepository delete docker repository from docker hub
/* curl \
   -H "Authorization: JWT ${TOKEN}" \
   -X DELETE \
   https://hub.docker.com/v2/repositories/${ORG}/${IMAGE}/
*/
func (c *Client) DeleteRepository(image string) error {
	if _, _, err := c.doRequest(http.MethodDelete, fmt.Sprintf("%s/%s/%s/", RepositoriesURL, c.ORG, image), nil); err != nil {
		color.Red("Error while deleting docker image: %s", err)
	}

	return nil
}
