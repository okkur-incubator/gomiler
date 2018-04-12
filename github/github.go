package github

import (
	"strconv"
	"time"
)

// GithubAPI struct
type githubAPI struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	StartDate   string     `json:"start_date"`
	DueDate     string     `json:"due_on"`
	Number      int        `json:"number"`
}

// Struct to be used for milestone queries
type milestone struct {
	DueDate string
	ID      string
	Title   string
	State   string
	Number  int
}

func createGithubMilestoneMap(githubAPI []githubAPI, api string) map[string]milestone {
	milestones := map[string]milestone{}
	for _, v := range githubAPI {
		var m milestone
		m.DueDate = v.DueDate
		m.ID = strconv.Itoa(v.ID)
		m.Title = v.Title
		if api == "github" {
			m.State = v.State
			m.Number = v.Number
		}
		milestones[v.Title] = m
	}

	return milestones
}
