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
	"strconv"
	"strings"
	"time"
)

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

// Struct to get ID from main API
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

// Initialization of logging variable
var logger *log.Logger

// LoggerSetup Initialization
func LoggerSetup(info io.Writer) {
	logger = log.New(info, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
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

// Function to get project ID from the gitLabAPI
func getProjectID(baseURL string, token string, projectname string, namespace string) (string, error) {
	projects := []gitLabAPI{}

	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("PRIVATE-TOKEN", token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(respByte, &projects)
	defer resp.Body.Close()
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
func getMilestones(baseURL string, token string, projectID string) (map[string]string, error) {
	milestones := []milestoneAPI{}
	m := map[string]string{}
	strURL := []string{baseURL, projectID, "/milestones?state=active&per_page=100"}
	URL := strings.Join(strURL, "")
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return m, err
	}

	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		return m, err
	}

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return m, err
	}

	json.Unmarshal(respByte, &milestones)
	defer resp.Body.Close()
	m = map[string]string{}

	for _, milestone := range milestones {
		m[milestone.Title] = milestone.DueDate
	}

	return m, nil
}

// CreateMilestoneData creates new milestones with title and due date
func createMilestoneData(advance int, timeInterval string) map[string]string {
	today := time.Now().Local()
	m := map[string]string{}

	switch {
	case timeInterval == "daily":
		for i := 0; i < advance; i++ {
			date := today.AddDate(0, 0, i).Format("2006-01-02")
			m[date] = date
		}

	case timeInterval == "weekly":
		for i := 0; i < advance; i++ {
			lastDay := lastDayWeek(today)
			year, week := lastDay.ISOWeek()
			title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
			dueDate := lastDay.Format("2006-01-02")
			m[title] = dueDate
			lastDay = lastDay.AddDate(0, 0, 7)
		}

	case timeInterval == "monthly":
		for i := 0; i < advance; i++ {
			date := today.AddDate(0, i, 0)
			lastday := lastDayMonth(date.Year(), int(date.Month()), time.UTC)
			title := date.Format("2006-01")
			dueDate := lastday.Format("2006-01-02")
			m[title] = dueDate
		}

	default:
		logger.Println("Error: Not Correct TimeInterval")
		return m
	}

	return m
}

func createMilestones(baseURL string, token string, projectID string, milestones map[string]string) error {
	strURL := []string{baseURL, projectID, "/milestones"}
	URL := strings.Join(strURL, "")
	client := &http.Client{}

	for k, v := range milestones {
		urlV := url.Values{}
		urlV.Set("title", k)
		urlV.Set("due_date", v)
		mbyte := bytes.NewReader([]byte(urlV.Encode()))
		req, err := http.NewRequest("POST", URL, mbyte)
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

func main() {
	// Declaring variables for flags
	var token, baseURL, namespace, project, timeInterval string
	var advance int
	// Command Line Parsing Starts
	flag.StringVar(&token, "Token", "jGWPwqQUuf37b", "Gitlab api key/token")
	flag.StringVar(&timeInterval, "TimeInterval", "daily", "Set milestone to daily, weekly or monthly")
	flag.StringVar(&baseURL, "BaseURL", "dev.example.com", "Gitlab api base url")
	flag.StringVar(&namespace, "Namespace", "someNamespace", "Namespace to use in Gitlab")
	flag.StringVar(&project, "ProjectName", "someProject", "Project to use in Gitlab")
	flag.IntVar(&advance, "Advance", 30, "Define timeframe to generate milestones in advance")
	flag.Parse() //Command Line Parsing Ends

	// Initializing logger
	LoggerSetup(os.Stdout)

	// Calling getProjectID
	baseURL = "https://" + baseURL + "/api/v4" + "/projects/"
	projectID, err := getProjectID(baseURL, token, project, namespace)
	if err != nil {
		logger.Println(err)
		// TODO: check for authentication error (currently it only says project not found)
	}

	oldMilestones, err := getMilestones(baseURL, token, projectID)
	if err != nil {
		logger.Println(err)
	}

	milestoneData := createMilestoneData(advance, strings.ToLower(timeInterval))

	// copy map
	newMilestones := map[string]string{}
	for k, v := range milestoneData {
		newMilestones[k] = v
	}

	for k, _ := range milestoneData {
		for ok, _ := range oldMilestones {
			if k == ok {
				delete(newMilestones, k)
			}
		}
	}

	if len(newMilestones) == 0 {
		logger.Println("No milestone creation needed")
	} else {
		logger.Println("New milestones:")
		for _, milestone := range newMilestones {
			logger.Printf("Title: %s - Due Date: %s", milestone.Title, milestone.DueDate)
		}
	}

	err = createMilestones(baseURL, token, projectID, newMilestones)
	if err != nil {
		logger.Println(err)
	}

	logger.Println(newMilestones) // TODO: Add final logging message with milestones created
}
