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
	"time"
)

// Repository represents docker repository information returned from hub.docker.com
type Repository struct {
	User              string    `json:"user"`
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	RepositoryType    string    `json:"repository_type"`
	Status            int       `json:"status"`
	Description       string    `json:"description"`
	IsPrivate         bool      `json:"is_private"`
	IsAutomated       bool      `json:"is_automated"`
	CanEdit           bool      `json:"can_edit"`
	StarCount         int       `json:"star_count"`
	PullCount         int       `json:"pull_count"`
	LastUpdated       time.Time `json:"last_updated"`
	IsMigrated        bool      `json:"is_migrated"`
	CollaboratorCount int       `json:"collaborator_count"`
	Affiliation       string    `json:"affiliation"`
	HubUser           string    `json:"hub_user"`
}

// RepositoryList represents the search repositories results from hub.docker.com
type RepositoryList struct {
	Count    int           `json:"count"`
	Next     string        `json:"next"`
	Previous string        `json:"previous"`
	Results  []*Repository `json:"results"`
}

// Image represents docker image information returned from hub.docker.com
type Image struct {
	Architecture string    `json:"architecture"`
	Features     string    `json:"features"`
	Variant      string    `json:"variant"`
	Digest       string    `json:"digest"`
	OS           string    `json:"os"`
	OSFeatures   string    `json:"os_features"`
	OSVersion    string    `json:"os_version"`
	Size         int       `json:"size"`
	Status       string    `json:"status"`
	LastPulled   time.Time `json:"last_pulled"`
	PastPushed   time.Time `json:"last_pushed"`
}

// Tag represents docker tag information returned from hub.docker.com
type Tag struct {
	Creator         int64     `json:"creator"`
	ID              int64     `json:"id"`
	ImageID         string    `json:"image_id"`
	Images          []*Image  `json:"images"`
	LastUpdated     time.Time `json:"last_updated"`
	LastUpdater     int64     `json:"last_updater"`
	LastUpdaterUser string    `json:"last_updater_username"`
	Name            string    `json:"name"`
	Repository      int64     `json:"repository"`
	FullSize        int       `json:"full_size"`
	V2              bool      `json:"v2"`
	TagStatus       string    `json:"tag_status"`
	TagLastPulled   time.Time `json:"tag_last_pulled"`
	TagLstaPushed   time.Time `json:"tag_last_pushed"`
}

// TagList represents the search tags results from hub.docker.com
type TagList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []*Tag `json:"results"`
}
