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

// Struct to be used for milestone
// Milestone ....
type gitLabAPI struct {
	ID          int               `json:"id"`
	Iid         int               `json:"iid"`
	ProjectID   int               `json:"project_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	StartDate   string            `json:"start_date"`
	DueDate     string            `json:"due_date"`
	State       string            `json:"state"`
	UpdatedAt   *time.Time        `json:"updated_at"`
	CreatedAt   *time.Time        `json:"created_at"`
	Name        string            `json:"name"`
	NameSpace   map[string]string `json:"namespace"`
}

// Initialization of logging variable
var logger *log.Logger

// LoggerSetup Initialization
func LoggerSetup(info io.Writer) {
	logger = log.New(info, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Function to get project ID from the gitLabAPI
func getProjectID(baseURL string, Token string, Projectname string, Namespace string) string {
	project := gitLabAPI{}
	urls := "https://" + baseURL + "/projects"
	page := 1

	for {

		strPage := strconv.Itoa(page)
		s := []string{urls, "?page=", strPage}
		completeURL := strings.Join(s, "")
		client := &http.Client{}
		req, err := http.NewRequest("GET", completeURL, nil)
		if err != nil {
			logger.Println(err)
			break
		}
		req.Header.Add("PRIVATE-TOKEN", Token)
		resp, err := client.Do(req)
		if err != nil {
			logger.Println(err)
			break
		}
		respByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Println("fail to read response data")
			break
		}
		json.Unmarshal(respByte, project)
		defer resp.Body.Close()
		fmt.Println(resp.Body)
		fmt.Println(project)
		if project.Name == "message" {
			logger.Println(project.Name)
		}
		if project.Name == Projectname && project.NameSpace["path"] == Namespace {
			return strconv.Itoa(project.ID)
		}
		if project.Title == "" {
			break
		}
		page++
	}
	return strconv.Itoa(project.ID)
}

func getMilestones(baseURL string, token string, projectID string) []string {
	project := gitLabAPI{}
	list := []string{}
	strurl := []string{"https://", baseURL, "/projects/", projectID, "/milestones"}
	urls := strings.Join(strurl, "")
	page := 1

	for {
		strPage := strconv.Itoa(page)
		s := []string{urls, "?page=", strPage}
		url := strings.Join(s, "")
		client := &http.Client{}
		re := regexp.MustCompile("^[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]$")
		req, err := http.NewRequest("GET", url, nil)
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
		respByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Println("fail to read response data")
			break
		}
		if project.State != "closed" && re.MatchString(project.Title) {

			titleList := []string{"Title: ", project.Title, "due_date", project.DueDate}
			y := strings.Join(titleList, ",")
			list = append(list, y)
		}

		json.Unmarshal(respByte, project)
		defer resp.Body.Close()

		if project.Name == "" {
			break
		}
		page++

	}
	return list
}

// CreateMilestoneData is used to check the due date using the time package of golang
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
	client := &http.Client{}
	for _, m := range milestones {
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			logger.Println(err)
			break
		}
		req.Header.Add("", m)
		req.Header.Add("PRIVATE-TOKEN", token)
		resp, err := client.Do(req)
		if err != nil {
			logger.Println(err)
			break
		}
		defer resp.Body.Close()
	}
	return ("Milestones Created" + strings.Join(milestones, ","))
}
func main() {
	// Declaring variables for flags
	var Token, APIBase, Namespace, Project string
	var Advance int
	// Command Line Parsing Starts
	flag.StringVar(&Token, "Token", " ", "Gitlab api key/token.")
	flag.StringVar(&APIBase, "baseURL", " ", "Gitlab api base url")
	flag.StringVar(&Namespace, "Namespace", " ", "Namespace to use in Gitlab")
	flag.StringVar(&Project, "ProjectName", " ", "Project to use in Gitlab")
	flag.IntVar(&Advance, "Advance", 30, "Define timeframe to generate milestones in advance.")
	flag.Parse() //Command Line Parsing Ends

	// Initializing logger
	LoggerSetup(os.Stdout)
	// Calling getProjectID

	projectID := getProjectID(APIBase, Token, Project, Namespace)
	newMilestone := getMilestones(APIBase, Token, projectID)
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
