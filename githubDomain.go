package main

// GithubPushEvent represents a Github webhook push event
type GithubPushEvent struct {
	After          string           `json:"after"`
	Commits        []GithubCommit   `json:"commits"`
	HeadCommit     GithubCommit     `json:"head_commit"`
	Pusher         GithubPusher     `json:"pusher"`
	Repository     GithubRepository `json:"repository"`
	InstallationID int              `json:"installation"`
}

// GithubCommit represents a Github commit
type GithubCommit struct {
	Author  GithubAuthor `json:"author"`
	Message string       `json:"message"`
	ID      string       `json:"id"`
}

// GithubAuthor represents a Github author
type GithubAuthor struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	UserName string `json:"username"`
}

// GithubPusher represents a Github pusher
type GithubPusher struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GithubRepository represents a Github repository
type GithubRepository struct {
	GitURL   string `json:"git_url"`
	HTMLURL  string `json:"html_url"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
}