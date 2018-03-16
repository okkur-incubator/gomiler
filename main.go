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

	"github.com/peterhellberg/link"
)

// Struct to be used for milestone queries
type milestone struct {
	DueDate string
	ID      string
	Title   string
}

// Struct to be used for milestone
type milestoneAPI struct {
	ID          int        `json:"id"`
	Iid         int        `json:"iid"`
	ProjectID   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	StartDate   string     `json:"start_date"`
	DueDate     string     `json:"due_date"`
}

// Struct for GitLab API
type gitLabAPI struct {
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

// Struct for GitHub API
type gitHubAPI struct {
	URL         string     `json:"url"`
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Number      int        `json:"number"`
	State       string     `json:"state"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DueDate     string     `json:"due_on"`
}

// Initialization of logging variable
var logger *log.Logger

// LoggerSetup Initialization
func LoggerSetup(info io.Writer) {
	logger = log.New(info, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func checkAPI(baseURL string, token string, namespace string, project string) (string, error) {
	gitlab := "gitlab"
	github := "github"
	gitlabURL := baseURL + "/api/v4/version"
	githubURL := baseURL + "/repos/" + namespace + "/" + project
	m := map[string]string{
		gitlab: gitlabURL,
		github: githubURL,
	}
	client := &http.Client{}
	for k, v := range m {
		req, err := http.NewRequest("GET", v, nil)
		if err != nil {
			return "", err
		}
		if k == gitlab {
			req.Header.Add("PRIVATE-TOKEN", token)
		}
		if k == github {
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

// Function to get last day of the month
func lastDayMonth(year int, month int, timezone *time.Location) time.Time {
	t := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	return t
}

// last day of week
func lastDayWeek(lastDay time.Time) time.Time {
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
		switch {
		case api == "gitlab":
			req.Header.Add("PRIVATE-TOKEN", token)
		case api == "github":
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

// Function to get project ID from the gitLabAPI
func getProjectID(baseURL string, token string, projectname string, namespace string, api string) (string, error) {
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
	projects := []gitLabAPI{}
	tmpM := []gitLabAPI{}
	for _, v := range apiData {
		json.Unmarshal(v, &tmpM)
		projects = append(projects, tmpM...)
	}
	for _, p := range projects {
		// Check for returned error messages
		if p.Name == "message" {
			return "", fmt.Errorf("api returned error %s", "error")
			// TODO: give back error/message returned by api
		}
		if p.Name == projectname && p.NameSpace.Path == namespace {
			return strconv.Itoa(p.ID), nil
		}
	}

	return "", fmt.Errorf("project %s not found", projectname)
}

// Get and return currently active milestones
func getActiveMilestones(baseURL string, token string, projectID string, api string) ([]milestoneAPI, error) {
	state := "active"
	return getMilestones(baseURL, token, projectID, state, api)
}

// Get and return inactive milestones
func getInactiveMilestones(baseURL string, token string, projectID string, api string) ([]milestoneAPI, error) {
	state := "closed"
	return getMilestones(baseURL, token, projectID, state, api)
}

func reactivateClosedMilestones(milestones map[string]milestone, baseURL string, token string, projectID string, api string) error {
	client := &http.Client{}
	for _, v := range milestones {
		milestoneID := v.ID
		strURL := []string{baseURL, "/projects/", projectID, "/milestones/", milestoneID}
		URL := strings.Join(strURL, "")

		// Overwrite state information in URL
		u, _ := url.Parse(URL)
		q := u.Query()
		q.Set("state_event", "activate")
		u.RawQuery = q.Encode()
		req, err := http.NewRequest("PUT", u.String(), nil)
		if err != nil {
			logger.Println(err)
		}
		switch {
		case api == "gitlab":
			req.Header.Add("PRIVATE-TOKEN", token)
		case api == "github":
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

func getMilestones(baseURL string, token string, projectID string, state string, api string) ([]milestoneAPI, error) {
	var strURL []string
	var URL, newURL string
	var apiData [][]byte
	switch {
	case api == "gitlab":
		strURL = []string{baseURL, "/projects/", projectID, "/milestones"}
		URL = strings.Join(strURL, "")
		u, _ := url.Parse(URL)
		q := u.Query()
		q.Set("state", state)
		u.RawQuery = q.Encode()
		newURL = u.String()
	case api == "github":
		strURL = []string{baseURL, "/milestones"}
		newURL = strings.Join(strURL, "")
	}
	apiData, err := paginate(newURL, token, api)
	if err != nil {
		return nil, err
	}
	milestones := []milestoneAPI{}
	tmpM := []milestoneAPI{}
	for _, v := range apiData {
		json.Unmarshal(v, &tmpM)
		milestones = append(milestones, tmpM...)
	}
	return milestones, nil
}

// CreateMilestoneData creates new milestones with title and due date
func createMilestoneData(advance int, timeInterval string, api string) map[string]milestone {
	today := time.Now().Local()
	milestones := map[string]milestone{}
	switch {
	case timeInterval == "daily":
		for i := 0; i < advance; i++ {
			var m milestone
			var dueDate string
			switch {
			case api == "gitlab":
				dueDate = today.AddDate(0, 0, i).Format("2006-01-02")
			case api == "github":
				dueDate = today.AddDate(0, 0, i).Format(time.RFC3339)
			}
			title := dueDate
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
		}
	case timeInterval == "weekly":
		for i := 0; i < advance; i++ {
			var m milestone
			var dueDate string
			lastDay := lastDayWeek(today)
			year, week := lastDay.ISOWeek()
			title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
			switch {
			case api == "gitlab":
				dueDate = today.AddDate(0, 0, i).Format("2006-01-02")
				dueDate = lastDay.Format("2006-01-02")
			case api == "github":
				dueDate = today.AddDate(0, 0, i).Format(time.RFC3339)
				dueDate = lastDay.Format(time.RFC3339)
			}
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
			today = lastDay.AddDate(0, 0, 7)
		}
	case timeInterval == "monthly":
		for i := 0; i < advance; i++ {
			var m milestone
			var dueDate string
			date := today.AddDate(0, i, 0)
			lastDay := lastDayMonth(date.Year(), int(date.Month()), time.UTC)
			title := date.Format("2006-01")
			switch {
			case api == "gitlab":
				dueDate = today.AddDate(0, 0, i).Format("2006-01-02")
				dueDate = lastDay.Format("2006-01-02")
			case api == "github":
				dueDate = today.AddDate(0, 0, i).Format(time.RFC3339)
				dueDate = lastDay.Format(time.RFC3339)
			}
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
		}
	default:
		logger.Println("Error: Not Correct TimeInterval")
		return milestones
	}

	return milestones
}

func createMilestones(baseURL string, token string, projectID string, milestones map[string]milestone, api string) error {
	client := &http.Client{}
	var strURL []string
	switch {
	case api == "gitlab":
		strURL = []string{baseURL, "/projects/", projectID, "/milestones"}
	case api == "github":
		strURL = []string{baseURL, "/milestones"}
	}
	URL := strings.Join(strURL, "")
	params := url.Values{}
	for _, v := range milestones {
		var req *http.Request
		switch {
		case api == "gitlab":
			params.Set("dueDate", v.DueDate)
			params.Set("title", v.Title)
			req, err := http.NewRequest("POST", URL, strings.NewReader((params.Encode())))
			if err != nil {
				return err
			}
			req.Header.Add("PRIVATE-TOKEN", token)
		case api == "github":
			params.Set("due_on", v.DueDate)
			create := struct {
				Title string `json:"title"`
				DueDate string `json:"due_on"`
			}{
				Title: v.Title,
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

func createMilestoneMap(milestoneAPI []milestoneAPI) map[string]milestone {
	milestones := map[string]milestone{}
	for _, v := range milestoneAPI {
		var m milestone
		m.DueDate = v.DueDate
		m.ID = strconv.Itoa(v.ID)
		m.Title = v.Title
		milestones[v.Title] = m
	}

	return milestones
}

func createAndDisplayNewMilestones(baseURL string, token string,
	projectID string, milestoneData map[string]milestone, api string) error {
	activeMilestonesAPI, err := getActiveMilestones(baseURL, token, projectID, api)
	if err != nil {
		return err
	}
	activeMilestones := createMilestoneMap(activeMilestonesAPI)
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
	closedMilestones := createMilestoneMap(closedMilestonesAPI)
	// copy map of closed milestones
	editMilestones := map[string]milestone{}
	for k := range milestoneData {
		for ek, ev := range closedMilestones {
			if k == ek {
				editMilestones[ek] = ev
			}
		}
	}
	return editMilestones, nil
}

func validateBaseURLScheme(baseURL string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	scheme := u.Scheme
	if scheme != "" {
		return baseURL, nil
	}
	URL := "https://" + baseURL
	return URL, nil
}

func main() {
	// Declaring variables for flags
	var token, baseURL, namespace, project, timeInterval string
	var advance int
	// Command Line Parsing Starts
	flag.StringVar(&token, "token", "jGWPwqQUuf37b", "Gitlab api key/token")
	flag.StringVar(&timeInterval, "time-interval", "daily", "Set milestone to daily, weekly or monthly")
	flag.StringVar(&baseURL, "base-url", "dev.example.com", "Gitlab api base url")
	flag.StringVar(&namespace, "namespace", "someNamespace", "Namespace to use in Gitlab")
	flag.StringVar(&project, "project", "someProject", "Project to use in Gitlab")
	flag.IntVar(&advance, "advance", 30, "Define timeframe to generate milestones in advance")
	flag.Parse() //Command Line Parsing Ends

	// Initializing logger
	LoggerSetup(os.Stdout)

	// Check which API to use
	api, err := checkAPI(baseURL, token, namespace, project)
	if err != nil {
		logger.Fatal(err)
	}
	milestoneData := createMilestoneData(advance, strings.ToLower(timeInterval), api)

	// Validate baseURL scheme
	URL, err := validateBaseURLScheme(baseURL)
	if err != nil {
		logger.Println(err)
	}
	// Calling getProjectID
	var newBaseURL, projectID string
	switch {
	case api == "gitlab":
		newBaseURL = URL + "/api/v4"
		projectID, err = getProjectID(newBaseURL, token, project, namespace, api)
		if err != nil {
			logger.Fatal(err)
			// TODO: check for authentication error (currently it only says project not found)
		}
		err = createAndDisplayNewMilestones(newBaseURL, token, projectID, milestoneData, api)
		if err != nil {
			logger.Println(err)
		}
		editMilestones, err := getClosedMilestones(newBaseURL, token, projectID, milestoneData, api)
		if err != nil {
			logger.Println(err)
		}
		err = reactivateClosedMilestones(editMilestones, newBaseURL, token, projectID, api)
		if err != nil {
			logger.Println(err)
		}
	case api == "github":
		newBaseURL = URL + "/repos/" + namespace + "/" + project
		err = createAndDisplayNewMilestones(newBaseURL, token, "", milestoneData, api)
		if err != nil {
			logger.Println(err)
		}
	}
}
