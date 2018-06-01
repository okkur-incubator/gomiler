/*
Copyright 2017 - The GoMiler Authors

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

package gitlab

import (
	"encoding/json"
	"fmt"
	"github.com/okkur/gomiler/utils"
	"github.com/peterhellberg/link"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
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

// GetProjectID function that gets a project's ID from the gitLabAPI
func GetProjectID(baseURL string, token string, projectname string, namespace string) (string, error) {
	strURL := []string{baseURL, "/projects/"}
	URL := strings.Join(strURL, "")
	u, _ := url.Parse(URL)
	q := u.Query()
	q.Set("search", projectname)
	u.RawQuery = q.Encode()
	apiData, err := paginate(u.String(), token)
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

func createGitlabMilestoneMap(gitlabAPI []gitlabAPI) map[string]utils.Milestone {
	milestones := map[string]utils.Milestone{}
	for _, v := range gitlabAPI {
		var m utils.Milestone
		m.DueDate = v.DueDate
		m.ID = strconv.Itoa(v.ID)
		m.Title = v.Title
		milestones[v.Title] = m
	}

	return milestones
}

func paginate(URL string, token string) ([][]byte, error) {
	apiData := make([][]byte, 1)
	client := &http.Client{}
	paginate := true
	for paginate == true {
		paginate = false
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		respByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		apiData = append(apiData, respByte)
		defer resp.Body.Close()

		// Retrieve next page header
		linkHeader := resp.Header.Get("Link")
		parsedHeader := link.Parse(linkHeader)
		for _, elem := range parsedHeader {
			if elem.Rel != "next" {
				continue
			}

			// Prevent break and modify URL for next iteration
			if elem.Rel == "next" {
				URL = elem.URI
				paginate = true
			}
		}
	}
	return apiData, nil
}

// Get and return currently active milestones
func getActiveMilestones(baseURL string, token string, projectID string) ([]gitlabAPI, error) {
	var state string
	state = "active"
	return getMilestones(baseURL, token, projectID, state)
}

// Get and return inactive milestones
func getInactiveMilestones(baseURL string, token string, project string) ([]gitlabAPI, error) {
	state := "closed"
	return getMilestones(baseURL, token, project, state)
}

// ReactivateClosedMilestones reactivates closed milestones that occur in the future
func ReactivateClosedMilestones(milestones map[string]utils.Milestone, baseURL string, token string,
	project string, logger *log.Logger) (map[string]utils.Milestone, error) {
	client := &http.Client{}
	var strURL []string
	for _, v := range milestones {
		milestoneID := v.ID
		strURL = []string{baseURL, "/projects/", project, "/milestones/", milestoneID}
		URL := strings.Join(strURL, "")
		var req *http.Request
		var err error
		// Overwrite state information in URL
		u, _ := url.Parse(URL)
		q := u.Query()
		q.Set("state_event", "activate")
		u.RawQuery = q.Encode()
		req, err = http.NewRequest("PUT", u.String(), nil)
		if err != nil {
			logger.Println(err)
		}
		req.Header.Add("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
	}
	// copy map of milestones with states changed to active for testing purposes
	reactivatedMilestones := map[string]utils.Milestone{}
	for k, v := range milestones {
		v.State = "active"
		reactivatedMilestones[k] = v
	}

	return reactivatedMilestones, nil
}

func getMilestones(baseURL string, token string, project string, state string) ([]gitlabAPI, error) {
	var strURL []string
	var URL, newURL string
	var apiData [][]byte
	strURL = []string{baseURL, "/projects/", project, "/milestones"}
	URL = strings.Join(strURL, "")
	u, _ := url.Parse(URL)
	q := u.Query()
	q.Set("state", state)
	u.RawQuery = q.Encode()
	newURL = u.String()
	apiData, err := paginate(newURL, token)
	if err != nil {
		return nil, err
	}
	milestones := []gitlabAPI{}
	tmpM := []gitlabAPI{}
	for _, v := range apiData {
		json.Unmarshal(v, &tmpM)
		milestones = append(milestones, tmpM...)
	}
	return milestones, nil
}

func createMilestones(baseURL string, token string, project string, milestones map[string]utils.Milestone) error {
	client := &http.Client{}
	var strURL []string
	strURL = []string{baseURL, "/projects/", project, "/milestones"}
	URL := strings.Join(strURL, "")
	params := url.Values{}
	for _, v := range milestones {
		var req *http.Request
		var err error

		params.Set("dueDate", v.DueDate)
		params.Set("title", v.Title)
		req, err = http.NewRequest("POST", URL, strings.NewReader((params.Encode())))
		if err != nil {
			return err
		}
		req.Header.Add("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	return nil
}

// CreateAndDisplayNewMilestones creates and displays new milestones
func CreateAndDisplayNewMilestones(baseURL string, token string,
	projectID string, milestoneData map[string]utils.Milestone, logger *log.Logger) error {
	activeMilestonesAPI, err := getActiveMilestones(baseURL, token, projectID)
	if err != nil {
		return err
	}
	activeMilestones := createGitlabMilestoneMap(activeMilestonesAPI)

	// copy map of active milestones
	newMilestones := map[string]utils.Milestone{}
	for k, v := range milestoneData {
		newMilestones[k] = v
	}
	for k := range milestoneData {
		for ok := range activeMilestones {
			if k == ok {
				delete(newMilestones, k)
			}
		}
	}
	if len(newMilestones) == 0 {
		logger.Println("No milestone creation needed")
	} else {
		logger.Println("New milestones:")
		var keys []string
		for k := range newMilestones {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			logger.Printf("Title: %s - Due Date: %s", newMilestones[key].Title, newMilestones[key].DueDate)
		}
		err = createMilestones(baseURL, token, projectID, newMilestones)
		if err != nil {
			return (err)
		}
	}
	return nil
}

// GetClosedMilestones gets closed milestones
func GetClosedMilestones(baseURL string, token string, projectID string, milestoneData map[string]utils.Milestone) (map[string]utils.Milestone, error) {
	closedMilestonesAPI, err := getInactiveMilestones(baseURL, token, projectID)
	if err != nil {
		return nil, err
	}
	closedGitlabMilestones := createGitlabMilestoneMap(closedMilestonesAPI)

	// copy map of closed milestones
	milestones := map[string]utils.Milestone{}
	for k := range milestoneData {
		for ek, ev := range closedGitlabMilestones {
			if k == ek {
				milestones[ek] = ev
			}
		}
	}
	return milestones, nil
}
