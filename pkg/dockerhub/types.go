package dockerhub

import (
	"time"
)

// Repository represents docker repository information returned from hub.docker.com
type Repository struct {
	User           string    `json:"user"`
	Name           string    `json:"name"`
	Namespace      string    `json:"namespace"`
	RepositoryType string    `json:"repository_type"`
	Status         int       `json:"status"`
	Description    string    `json:"description"`
	IsPrivate      bool      `json:"is_private"`
	IsAutomated    bool      `json:"is_automated"`
	CanEdit        bool      `json:"can_edit"`
	StarCount      int       `json:"star_count"`
	PullCount      int       `json:"pull_count"`
	LastUpdated    time.Time `json:"last_updated"`
	IsMigrated     bool      `json:"is_migrated"`
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
	Size         int    `json:"size"`
	Digest       string `json:"digest"`
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
	OSVersion    string `json:"os_version"`
	OSFeatures   string `json:"os_features"`
	Variant      string `json:"variant"`
	Features     string `json:"features"`
}

// Tag represents docker tag information returned from hub.docker.com
type Tag struct {
	Name            string    `json:"name"`
	FullSize        int       `json:"full_size"`
	Images          []*Image  `json:"images"`
	ID              int64     `json:"id"`
	Repository      int64     `json:"repository"`
	Creator         int64     `json:"creator"`
	LastUpdater     int64     `json:"last_updater"`
	LastUpdaterUser string    `json:"last_updater_username"`
	ImageID         string    `json:"image_id"`
	V2              bool      `json:"v2"`
	LastUpdated     time.Time `json:"last_updated"`
}

// TagList represents the search tags results from hub.docker.com
type TagList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []*Tag `json:"results"`
}
