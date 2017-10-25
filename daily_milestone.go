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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Struct to be used for project
// project ....
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

// Struct to be used for milestone
// milestone .....

type milestone struct {
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
	//page := 1
	//strPage := strconv.Itoa(page)
	//s := []string{urls, "?page=", strPage}
	//completeURL := strings.Join(s, "")
	//fmt.Println(completeURL)
	client := &http.Client{}
	req, err := http.NewRequest("GET", urls, nil)
	fmt.Println(urls)
	if err != nil {
		logger.Println(err)

	}
	req.Header.Add("Token", Token)
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
	fmt.Println(resp.Body)
	fmt.Println(len(project))
	for _, p := range project {
		if p.Name == "message" {
			fmt.Println(p.Name)
		}
		fmt.Println(p.Name)
		fmt.Println(p.NameSpace.Path)
		if p.Name == Projectname /*&& p.NameSpace.Path == Namespace*/ {
			return strconv.Itoa(p.ID), nil
		}
		if p.Title == "" {
			break
		}

	}
	return "", fmt.Errorf("project %s not found", Projectname)
}

func getMilestones(baseURL string, token string, projectID string) ([]string, error) {
	project := []gitLabAPI{}
	list := []string{}
	strurl := []string{"https://", baseURL, "/projects/", projectID, "/milestones"}
	urls := strings.Join(strurl, "")

	page := 1
	strPage := strconv.Itoa(page)
	s := []string{urls, "/", strPage, "/milestones"}
	url := strings.Join(s, "")
	client := &http.Client{}
	re := regexp.MustCompile("^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]$")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Println(err)

	}
	//req.Header.Add("PRIVATE-TOKEN", token)
	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)

	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("fail to read response data")

	}

	for _, p := range project {

		if p.State != "closed" && re.MatchString(p.Title) {

			titleList := []string{"Title: ", p.Title, "due_date", p.DueDate}
			y := strings.Join(titleList, ",")
			list = append(list, y)
		}
		errjson := json.Unmarshal(respByte, &project)
		if errjson != nil {
			logger.Println("Error")
			logger.Println(errjson)
		}
		defer resp.Body.Close()

		if p.Name == "" {
			break
		}
		page++

	}
	return list, nil
}

// CreateMilestoneData is used to check the due date using the time package of python
func createMilestoneData(advance int) []string {
	today := time.Now().Local()
	list := []string{}
	for i := 0; i < advance; i++ {
		date := today.AddDate(0, 0, i)
		dateiso := date.Format("2009-01-02")                            // To Format to ISOFormat and it converts to string so can be used in list directly
		datelist := []string{"Title: ", "dateiso", "due_date", dateiso} // Was unable to get a map in a list so made this
		y := strings.Join(datelist, ",")
		list = append(list, y)

	}
	return list
}

func createMilestones(baseURL string, token string, projectID string, milestones []string) string {
	strurl := []string{"https://", baseURL, "/projects/", projectID, "/milestones"}
	url := strings.Join(strurl, ",")
	fmt.Println(url)
	client := &http.Client{}
	for _, m := range milestones {
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			logger.Println(err)
			break
		}
		req.Header.Add("", m)
		//req.Header.Add("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			logger.Println(err)
			break
		}
		defer resp.Body.Close()
	}
	return ("Milestones Created" + strings.Join(milestones, ""))
}

func main() {
	// Declaring variables for flags
	var Token, APIBase, Namespace, Project string
	var Advance int
	APIBase = "lol"
	// Command Line Parsing Starts
	flag.StringVar(&Token, "Token", "lol", "Gitlab api key/token.")
	flag.StringVar(&APIBase, "baseURL", "dev.seetheprogress.eu", "Gitlab api base url")
	flag.StringVar(&Namespace, "Namespace", "okkur", "Namespace to use in Gitlab")
	flag.StringVar(&Project, "ProjectName", "syna", "Project to use in Gitlab")
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
