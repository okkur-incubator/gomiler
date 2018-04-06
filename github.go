package main

import (
	"encoding/json"
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

func (g *GoMiler) getGithubMilestones(gomiler []GoMiler) []githubAPI {
	milestones := []githubAPI{}
	tmpM := []githubAPI{}
	for range gomiler {
		json.Unmarshal(g.JSON, &tmpM)
		milestones = append(milestones, tmpM...)
	}
	return milestones
}

func (g *GoMiler) createGithubMilestoneMap(githubAPI []githubAPI, api string) map[string]milestone {
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
