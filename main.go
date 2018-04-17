/*
Copyright 2017 - The Dailymile Authors

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

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	github "github.com/okkur/gomiler/github"
	"github.com/peterhellberg/link"
	gitlab "github.com/okkur/gomiler/gitlab"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// GoMiler struct to be used as a generic struct for use with multiple APIs
type GoMiler struct {
	JSON []byte
	API  string
}

// Initialization of logging variable
var logger *log.Logger

// LoggerSetup Initialization
func LoggerSetup(info io.Writer) {
	logger = log.New(info, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func checkAPI(baseURL string, token string, namespace string, project string) (string, error) {
	gitlabURL := baseURL + "/api/v4/version"
	githubURL := baseURL + "/repos/" + namespace + "/" + project
	m := map[string]string{
		"gitlab": gitlabURL,
		"github": githubURL,
	}
	client := &http.Client{}
	for k, v := range m {
		req, err := http.NewRequest("GET", v, nil)
		if err != nil {
			return "", err
		}
		switch k {
		case "gitlab":
			req.Header.Add("PRIVATE-TOKEN", token)
		case "github":
			req.Header.Add("Accept", "application/vnd.github.inertia-preview+json")
			req.Header.Add("Authorization", "token "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			return k, nil
		}
	}
	return "", fmt.Errorf("Error: could not access GitLab or GitHub APIs")
}

// LastDayMonth function gets last day of the month
func LastDayMonth(year int, month int, timezone *time.Location) time.Time {
	t := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	return t
}

// LastDayWeek function gets last day of the week
func LastDayWeek(lastDay time.Time) time.Time {
	if lastDay.Weekday() != time.Sunday {
		for lastDay.Weekday() != time.Sunday {
			lastDay = lastDay.AddDate(0, 0, +1)
		}
		return lastDay
	}
	return lastDay
}

func paginate(URL string, token string, api string) ([][]byte, error) {
	apiData := make([][]byte, 1)
	client := &http.Client{}
	paginate := true
	for paginate == true {
		paginate = false
		req, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			return nil, err
		}
		switch api {
		case "gitlab":
			req.Header.Add("PRIVATE-TOKEN", token)
		case "github":
			req.Header.Add("Accept", "application/vnd.github.inertia-preview+json")
			req.Header.Add("Authorization", "token "+token)
		}
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
func getActiveMilestones(baseURL string, token string, projectID string, api string) ([]GoMiler, error) {
	var state string
	switch api {
	case "gitlab":
		state = "active"
	case "github":
		state = "open"
	}
	return getMilestones(baseURL, token, projectID, state, api)
}

// Get and return inactive milestones
func getInactiveMilestones(baseURL string, token string, project string, api string) ([]GoMiler, error) {
	state := "closed"
	return getMilestones(baseURL, token, project, state, api)
}

func reactivateClosedMilestones(milestones map[string]milestone, baseURL string, token string, project string, api string) error {
	client := &http.Client{}
	var strURL []string
	for _, v := range milestones {
		switch api {
		case "gitlab":
			milestoneID := v.ID
			strURL = []string{baseURL, "/projects/", project, "/milestones/", milestoneID}
		case "github":
			milestoneID := strconv.Itoa(v.Number)
			strURL = []string{baseURL, project, "/milestones/", milestoneID}
		}
		URL := strings.Join(strURL, "")
		var req *http.Request
		var err error
		switch api {
		case "gitlab":
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
		case "github":
			updatePatch := struct {
				State string `json:"state"`
			}{
				State: "open",
			}
			updatePatchBytes, err := json.Marshal(updatePatch)
			if err != nil {
				return err
			}
			req, err = http.NewRequest("PATCH", URL, bytes.NewReader(updatePatchBytes))
			if err != nil {
				return err
			}
			req.Header.Add("Accept", "application/vnd.github.inertia-preview+json")
			req.Header.Add("Authorization", "token "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	return nil
}

func getMilestones(baseURL string, token string, project string, state string, api string) ([]GoMiler, error) {
	var strURL []string
	var URL, newURL string
	var apiData [][]byte
	switch api {
	case "gitlab":
		strURL = []string{baseURL, "/projects/", project, "/milestones"}
	case "github":
		strURL = []string{baseURL, project, "/milestones"}
	}
	URL = strings.Join(strURL, "")
	u, _ := url.Parse(URL)
	q := u.Query()
	q.Set("state", state)
	u.RawQuery = q.Encode()
	newURL = u.String()
	apiData, err := paginate(newURL, token, api)
	if err != nil {
		return nil, err
	}
	milestones := []GoMiler{}
	tmpM := []GoMiler{}
	for _, v := range apiData {
		json.Unmarshal(v, &tmpM)
		milestones = append(milestones, tmpM...)
	}
	return milestones, nil
}

// CreateMilestoneData creates new milestones with title and due date
func createMilestoneData(advance int, interval string, api string) map[string]milestone {
	today := time.Now().Local()
	milestones := map[string]milestone{}
	switch interval {
	case "daily":
		for i := 0; i < advance; i++ {
			var m milestone
			var dueDate string
			title := today.AddDate(0, 0, i).Format("2006-01-02")
			switch api {
			case "gitlab":
				dueDate = today.AddDate(0, 0, i).Format("2006-01-02")
			case "github":
				dueDate = today.AddDate(0, 0, i).Format(time.RFC3339)
			}
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
		}
	case "weekly":
		for i := 0; i < advance; i++ {
			var m milestone
			var dueDate string
			lastDay := lastDayWeek(today)
			year, week := lastDay.ISOWeek()
			title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
			switch api {
			case "gitlab":
				dueDate = lastDay.Format("2006-01-02")
			case "github":
				dueDate = lastDay.Format(time.RFC3339)
			}
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
			today = lastDay.AddDate(0, 0, 7)
		}
	case "monthly":
		for i := 0; i < advance; i++ {
			var m milestone
			var dueDate string
			date := today.AddDate(0, i, 0)
			lastDay := lastDayMonth(date.Year(), int(date.Month()), time.UTC)
			title := date.Format("2006-01")
			switch api {
			case "gitlab":
				dueDate = lastDay.Format("2006-01-02")
			case "github":
				dueDate = lastDay.Format(time.RFC3339)
			}
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
		}
	default:
		logger.Println("Error: Incorrect interval")
		return milestones
	}

	return milestones
}

func createMilestones(baseURL string, token string, project string, milestones map[string]milestone, api string) error {
	client := &http.Client{}
	var strURL []string
	switch api {
	case "gitlab":
		strURL = []string{baseURL, "/projects/", project, "/milestones"}
	case "github":
		strURL = []string{baseURL, project, "/milestones"}
	}
	URL := strings.Join(strURL, "")
	params := url.Values{}
	for _, v := range milestones {
		var req *http.Request
		var err error
		switch api {
		case "gitlab":
			params.Set("dueDate", v.DueDate)
			params.Set("title", v.Title)
			req, err = http.NewRequest("POST", URL, strings.NewReader((params.Encode())))
			if err != nil {
				return err
			}
			req.Header.Add("PRIVATE-TOKEN", token)
		case "github":
			create := struct {
				Title   string `json:"title"`
				DueDate string `json:"due_on"`
			}{
				Title:   v.Title,
				DueDate: v.DueDate,
			}
			createBytes, err := json.Marshal(create)
			if err != nil {
				return err
			}
			req, err = http.NewRequest("POST", URL, bytes.NewReader(createBytes))
			if err != nil {
				return err
			}
			req.Header.Add("Accept", "application/vnd.github.inertia-preview+json")
			req.Header.Add("Authorization", "token "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}

	return nil
}

func createAndDisplayNewMilestones(baseURL string, token string,
	projectID string, milestoneData map[string]milestone, api string) error {
	activeMilestonesAPI, err := getActiveMilestones(baseURL, token, projectID, api)
	if err != nil {
		return err
	}
	gitlabMilestones := []gitlabAPI{}
	githubMilestones := []githubAPI{}
	activeMilestones := map[string]milestone{}
	var g GoMiler
	switch api {
	case "gitlab":
		gitlabMilestones = (*GoMiler).getGitlabMilestones(&g, activeMilestonesAPI)
		activeMilestones = (*GoMiler).createGitlabMilestoneMap(&g, gitlabMilestones, api)
	case "github":
		githubMilestones = (*GoMiler).getGithubMilestones(&g, activeMilestonesAPI)
		activeMilestones = (*GoMiler).createGithubMilestoneMap(&g, githubMilestones, api)
	}
	// copy map of active milestones
	newMilestones := map[string]milestone{}
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
		err = createMilestones(baseURL, token, projectID, newMilestones, api)
		if err != nil {
			return (err)
		}
	}
	return nil
}

func getClosedMilestones(baseURL string, token string, projectID string, milestoneData map[string]milestone, api string) (map[string]milestone, error) {
	closedMilestonesAPI, err := getInactiveMilestones(baseURL, token, projectID, api)
	if err != nil {
		return nil, err
	}
	closedGitlabMilestones := map[string]milestone{}
	closedGithubMilestones := map[string]milestone{}
	var g GoMiler
	switch api {
	case "gitlab":
		gitlabMilestones := (*GoMiler).getGitlabMilestones(&g, closedMilestonesAPI)
		closedGitlabMilestones = (*GoMiler).createGitlabMilestoneMap(&g, gitlabMilestones, api)
	case "github":
		githubMilestones := (*GoMiler).getGithubMilestones(&g, closedMilestonesAPI)
		closedGithubMilestones = (*GoMiler).createGithubMilestoneMap(&g, githubMilestones, api)
	}
	// copy map of closed milestones
	milestones := map[string]milestone{}
	for k := range milestoneData {
		switch api {
		case "gitlab":
			for ek, ev := range closedGitlabMilestones {
				if k == ek {
					milestones[ek] = ev
				}
			}
		case "github":
			for ek, ev := range closedGithubMilestones {
				if k == ek {
					milestones[ek] = ev
				}
			}
		}
	}
	return milestones, nil
}

func validateBaseURLScheme(baseURL string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u.Scheme = "https"
	q := u.Query()
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func main() {
	// Declaring variables for flags
	var token, baseURL, namespace, project, interval string
	var advance int
	// Command Line Parsing Starts
	flag.StringVar(&token, "token", "jGWPwqQUuf37b", "GitLab or GitHub API key/token")
	flag.StringVar(&interval, "interval", "daily", "Set milestone to daily, weekly or monthly")
	flag.StringVar(&baseURL, "url", "dev.example.com", "GitLab or GitHub API base URL")
	flag.StringVar(&namespace, "namespace", "someNamespace", "Namespace to use in GitLab or GitHub")
	flag.StringVar(&project, "project", "someProject", "Project to use in GitLab or GitHub")
	flag.IntVar(&advance, "advance", 30, "Define timeframe to generate milestones in advance")
	flag.Parse() //Command Line Parsing Ends

	// Initializing logger
	LoggerSetup(os.Stdout)

	// Validate baseURL scheme
	URL, err := validateBaseURLScheme(baseURL)
	if err != nil {
		logger.Println(err)
	}

	// Check which API to use
	api, err := checkAPI(URL, token, namespace, project)
	if err != nil {
		logger.Fatal(err)
	}
	milestoneData := createMilestoneData(advance, strings.ToLower(interval), api)
	
	// Calling getProjectID
	var newBaseURL, projectID string
	switch api {
	case "gitlab":
		newBaseURL = URL + "/api/v4"
		var g GoMiler
		projectID, err = (*GoMiler).getProjectID(&g, newBaseURL, token, project, namespace, api)
		if err != nil {
			logger.Fatal(err)
		}
		err = createAndDisplayNewMilestones(newBaseURL, token, projectID, milestoneData, api)
		if err != nil {
			logger.Println(err)
		}
		closedMilestones, err := getClosedMilestones(newBaseURL, token, projectID, milestoneData, api)
		if err != nil {
			logger.Println(err)
		}
		err = reactivateClosedMilestones(closedMilestones, newBaseURL, token, projectID, api)
		if err != nil {
			logger.Println(err)
		}
	case "github":
		newBaseURL = URL + "/repos/" + namespace + "/"
		err = createAndDisplayNewMilestones(newBaseURL, token, project, milestoneData, api)
		if err != nil {
			logger.Println(err)
		}
		closedMilestones, err := getClosedMilestones(newBaseURL, token, project, milestoneData, api)
		if err != nil {
			logger.Println(err)
		}
		err = reactivateClosedMilestones(closedMilestones, newBaseURL, token, project, api)
		if err != nil {
			logger.Println(err)
		}
	}
}
