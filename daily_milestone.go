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

// Struct to be used for project
// project ....
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
func getProjectID(baseURL string, Token string, Projectname string, Namespace string) (string, error) {
	project := []gitLabAPI{}
	urls := "https://" + baseURL + "/api/v4" + "/projects"
	client := &http.Client{}
	req, err := http.NewRequest("GET", urls, nil)
	if err != nil {
		logger.Println(err)

	}
	req.Header.Add("PRIVATE-TOKEN", Token)
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)
	}

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("fail to read response data")

	}
	json.Unmarshal(respByte, &project)
	defer resp.Body.Close()
	for _, p := range project {
		if p.Name == "message" {
			fmt.Println(p.Name)
		}

		if p.Name == Projectname && p.NameSpace.Path == Namespace {
			fmt.Println(strconv.Itoa(p.ID))
			return strconv.Itoa(p.ID), nil
		}
	}
	return "", fmt.Errorf("project %s not found", Projectname)
}

// It is getting milestones data from the milestone API
func getMilestones(baseURL string, token string, projectID string) ([]string, error) {
	project := []milestoneAPI{}
	list := []string{}
	strurl := []string{"https://", baseURL, "/api/v4", "/projects/", projectID, "/milestones"}
	urls := strings.Join(strurl, "")
	client := &http.Client{}
	req, err := http.NewRequest("GET", urls, nil)
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

	json.Unmarshal(respByte, &project)
	defer resp.Body.Close()
	for _, p := range project {
		if p.State != "closed" {

			titleList := []string{"Title: ", p.Title, ", DueDate", p.DueDate}
			y := strings.Join(titleList, " ")
			list = append(list, y)
		}
	}

	return list, nil
}

// CreateMilestoneData is used to check the due date using the time package of python
func createMilestoneData(advance int) []string {
	today := time.Now().Local()
	date := today.AddDate(0, 0, advance)
	list := []string{}
	/*for i := 0; i < advance; i++ {
		date := today.AddDate(0, 0, i)                          // To Format to ISOFormat and it converts to string so can be used in list directly
		datelist := []string{"Title: ", date, "due_date", date} // Was unable to get a map in a list so made this
		y := strings.Join(datelist, ",")
		list = append(list, y)
	}*/
	fmt.Println(date)
	fmt.Println(list)
	return list
}

func createMilestones(baseURL string, token string, projectID string, milestones []string) string {
	strurl := []string{"https://", baseURL, "/projects/", projectID, "/milestones"}
	url := strings.Join(strurl, "")
	client := &http.Client{}
	for _, m := range milestones {
		mbyte := bytes.NewReader([]byte(m))
		req, err := http.NewRequest("POST", url, mbyte)
		if err != nil {
			logger.Println(err)
			break
		}
		req.Header.Add("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			logger.Println(err)
			break
		}
		defer resp.Body.Close()
	}
	return ("Milestones Created: " + strings.Join(milestones, ""))
}

func main() {
	// Declaring variables for flags
	var Token, APIBase, Namespace, Project string
	var Advance int
	APIBase = ""
	// Command Line Parsing Starts
	flag.StringVar(&Token, "Token", "bVYFTJaYtgAZesSofKbq", "Gitlab api key/token.")
	flag.StringVar(&APIBase, "baseURL", "dev.seetheprogress.eu", "Gitlab api base url")
	flag.StringVar(&Namespace, "Namespace", "okkur", "Namespace to use in Gitlab")
	flag.StringVar(&Project, "ProjectName", "dailymile_test", "Project to use in Gitlab")
	flag.IntVar(&Advance, "Advance", 30, "Define timeframe to generate milestones in advance.")
	flag.Parse() //Command Line Parsing Ends

	// Initializing logger
	LoggerSetup(os.Stdout)
	// Calling getProjectID
	projectID, err := getProjectID(APIBase, Token, Project, Namespace)
	if err != nil {
		logger.Println(err)
	}
	newMilestone, err := getMilestones(APIBase, Token, projectID)
	if err != nil {
		logger.Println(err)
	}
	oldMilestone := createMilestoneData(Advance)

	for index, y := range newMilestone {
		for _, z := range oldMilestone {
			if z == y {
				newMilestone = append(newMilestone[:index], newMilestone[(index+1):]...)
			}
		}

	}
	message := createMilestones(APIBase, Token, projectID, newMilestone)
	logger.Println(message)
}
