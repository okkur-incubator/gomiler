package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GitlabAPI struct
type gitlabAPI struct {
	ID          int        `json:"id"`
	Iid         int        `json:"iid"`
	ProjectID   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   string     `json:"start_date"`
	DueDate     string     `json:"due_date"`
	State       string     `json:"state"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   *time.Time `json:"created_at"`
	Name        string     `json:"name"`
	NameSpace   struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		Kind     string `json:"kind"`
		FullPath string `json:"full_path"`
	} `json:"namespace"`
}

// Function to get project ID from the gitLabAPI
func (g *GoMiler) getProjectID(baseURL string, token string, projectname string, namespace string, api string) (string, error) {
	strURL := []string{baseURL, "/projects/"}
	URL := strings.Join(strURL, "")
	u, _ := url.Parse(URL)
	q := u.Query()
	q.Set("search", projectname)
	u.RawQuery = q.Encode()
	apiData, err := paginate(u.String(), token, api)
	if err != nil {
		return "", err
	}
	projects := []gitlabAPI{}
	tmpM := []gitlabAPI{}
	for _, v := range apiData {
		json.Unmarshal(v, &tmpM)
		projects = append(projects, tmpM...)
	}
	for _, p := range projects {
		// Check for returned error messages
		if p.Name == "message" {
			return "", fmt.Errorf("api returned error %s", "error")
		}
		if p.Name == projectname && p.NameSpace.Path == namespace {
			return strconv.Itoa(p.ID), nil
		}
	}

	return "", fmt.Errorf("project %s not found", projectname)
}

func (g *GoMiler) getGitlabMilestones(gomiler []GoMiler) []gitlabAPI {
	milestones := []gitlabAPI{}
	tmpM := []gitlabAPI{}
	for range gomiler {
		json.Unmarshal(g.JSON, &tmpM)
		milestones = append(milestones, tmpM...)
	}
	return milestones
}

func (g *GoMiler) createGitlabMilestoneMap(gitlabAPI []gitlabAPI, api string) map[string]milestone {
	milestones := map[string]milestone{}
	for _, v := range gitlabAPI {
		var m milestone
		m.DueDate = v.DueDate
		m.ID = strconv.Itoa(v.ID)
		m.Title = v.Title
		milestones[v.Title] = m
	}

	return milestones
}
