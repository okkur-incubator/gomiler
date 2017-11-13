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
	"os"
	"strconv"
	"strings"
	"time"
)

// Struct to be used for milestone
// milestoneApi .....
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

//Struct to get ID from main API
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

// Function to get project ID from the gitLabAPI
func getProjectID(baseURL string, token string, projectname string, namespace string) (string, error) {
	projects := []gitLabAPI{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		logger.Println(err)
	}
	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)
	}

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("fail to read response data")

	}
	json.Unmarshal(respByte, &projects)
	defer resp.Body.Close()
	for _, p := range projects {
		// To check for message404 error and that project is not forund
		if p.Name == "message" {
			fmt.Println(p.Name)
		}

		if p.Name == projectname && p.NameSpace.Path == namespace {
			fmt.Println(strconv.Itoa(p.ID))
			return strconv.Itoa(p.ID), nil
		}
	}
	return "", fmt.Errorf("project %s not found", projectname)
}

// It is getting milestones data from the milestone API and returning list of closed milestones
func getMilestones(baseURL string, token string, projectID string) ([]string, error) {
	milestones := []milestoneAPI{}
	list := []string{}
	strurl := []string{baseURL, projectID, "/milestones"}
	url := strings.Join(strurl, "")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Println(err)
	}
	req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)

	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("fail to read response data")
	}

	json.Unmarshal(respByte, &milestones)
	defer resp.Body.Close()
	for _, m := range milestones {
		if m.State != "closed" {
			titleTuple := []string{"Title: ", m.Title, ", DueDate", m.DueDate}
			y := strings.Join(titleTuple, " ")
			list = append(list, y)
		}
	}
	fmt.Println(list)
	return list, nil
}

// CreateMilestoneData is used to check the due date using the time package of python
func createMilestoneData(advance int) []string {
	today := time.Now().Local()
	list := []string{}
	for i := 0; i < advance; i++ {
		date := today.AddDate(0, 0, i)                                                                      // To Format to ISOFormat and it converts to string so can be used in list directly
		datelist := []string{"Title: ", date.Format("2006-01-02"), "  due_date", date.Format("2006-01-02")} // Was unable to get a map in a list so made this
		y := strings.Join(datelist, ",")
		list = append(list, y)
	}

	fmt.Println(list)
	return list
}

func createMilestones(baseURL string, token string, projectID string, milestones []string) (string, error) {
	strurl := []string{baseURL, projectID, "/milestones"}
	url := strings.Join(strurl, "")
	client := &http.Client{}
	for _, m := range milestones {
		mbyte := bytes.NewReader([]byte(m))
		req, err := http.NewRequest("POST", url, mbyte)
		if err != nil {
			return "", err
		}
		req.Header.Add("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()
	}
	return ("Milestones Created: " + strings.Join(milestones, "")), nil
}

func main() {
	// Declaring variables for flags
	var Token, baseURL, Namespace, Project string
	var Advance int
	// Command Line Parsing Starts
	flag.StringVar(&Token, "Token", "bVYFTJaYtgAZesSofKbq", "Gitlab api key/token.")
	flag.StringVar(&baseURL, "baseURL", "dev.seetheprogress.eu", "Gitlab api base url")
	flag.StringVar(&Namespace, "Namespace", "okkur", "Namespace to use in Gitlab")
	flag.StringVar(&Project, "ProjectName", "dailymile_test", "Project to use in Gitlab")
	flag.IntVar(&Advance, "Advance", 30, "Define timeframe to generate milestones in advance.")
	flag.Parse() //Command Line Parsing Ends

	// Initializing logger
	LoggerSetup(os.Stdout)
	// Calling getProjectID
	baseurl := "https://" + baseURL + "/api/v4" + "/projects/"
	projectID, err := getProjectID(baseurl, Token, Project, Namespace)
	if err != nil {
		logger.Println(err)
	}
	oldMilestone, err := getMilestones(baseurl, Token, projectID)
	if err != nil {
		logger.Println(err)
	}
	newMilestone := createMilestoneData(Advance)

	for index, y := range newMilestone {
		for _, z := range oldMilestone {
			if z == y {
				newMilestone = append(newMilestone[:index], newMilestone[(index+1):]...)
			}
		}

	}
	message, err := createMilestones(baseurl, Token, projectID, newMilestone)
	if err != nil {
		logger.Println(err)
	}
	logger.Println(message)
}
