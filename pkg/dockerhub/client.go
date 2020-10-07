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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

// BaseURL represents Docker Hub endpoint
const BaseURL = "https://hub.docker.com/v2/"

// RepositoriesURL represents Docker Hub repositories endpoint
var RepositoriesURL = BaseURL + "repositories"

// AuthResponse represents auth response
type AuthResponse struct {
	Token string `json:"token"`
}

// Client represents new HTTP client
type Client struct {
	*http.Client
	Header    http.Header
	AuthToken string
	URL       string
	ORG       string
}

func GetFlags(flags *pflag.FlagSet) (string, bool, error) {
	org, err := flags.GetString("org")
	if err != nil {
		return "", true, err
	}

	dryRun, err := flags.GetBool("dry-run")
	if err != nil {
		return org, true, err
	}

	return org, dryRun, nil
}

func NewClient(org, url string) *Client {
	c := &http.Client{
		Timeout: time.Second * 30,
	}
	if url == "" {
		url = "https://hub.docker.com"
	}

	h := http.Header{}
	h.Set("Content-Type", "application/json")

	return &Client{
		Client: c,
		Header: h,
		URL:    url,
		ORG:    org,
	}
}

// GetAuthToken returns JWT Token from docker hub login page
/* curl --silent \
   -H "Content-Type: application/json" \
   -X POST \
   -d '{"username": "'${DOCKERHUB_USERNAME}'", "password": "'${DOCKERHUB_PASSWORD}'"}' \
   https://hub.docker.com/v2/users/login/ | jq -r .token
*/
func (c *Client) GetAuthToken() (string, error) {
	payload := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, os.Getenv("DOCKERHUB_USERNAME"), os.Getenv("DOCKERHUB_PASSWORD"))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users/login", BaseURL), bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	accessToken := &AuthResponse{}
	if err = json.NewDecoder(resp.Body).Decode(accessToken); err != nil {
		return "", err
	}

	c.AuthToken = accessToken.Token
	if accessToken.Token == "" {
		color.Red("failed to log into the registry")
		return "", err
	}

	return accessToken.Token, nil
}

func (c *Client) NewRequest(method, url string, payload io.Reader) (*http.Request, error) {
	if c.AuthToken == "" {
		token, err := c.GetAuthToken()
		if err != nil {
			return nil, err
		}
		c.AuthToken = token
	}

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("JWT %s", c.AuthToken))

	return req, nil
}

func (c *Client) doRequest(method, url string, payload io.Reader) (data []byte, status int, err error) {
	request, err := c.NewRequest(method, url, payload)
	if err != nil {
		return nil, 0, err
	}

	response, err := c.Client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, 0, err
	}

	if (method == http.MethodGet) && (response.StatusCode != http.StatusOK) {
		color.Red("HTTP error!\nURL: %s\nstatus code: %d\nbody:\n%s\n", url, response.StatusCode, response.Body)
		return nil, response.StatusCode, err
	}

	return body, response.StatusCode, nil
}
