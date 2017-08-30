package main

import (
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
)

var logger *log.Logger

//LoggerSetup Initialization
func LoggerSetup(info io.Writer) {
	logger = log.New(info, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func getProjectID(BaseURL string, Token string, Projectname string, Namespace string) string {
	urls := "https://" + BaseURL + "/projects"
	page := 1
	headers := url.Values{}
	headers.Set("PrivateToken", Token)
	strPage := strconv.Itoa(page)
	s := []string{urls, "?page=", strPage}
	completeURL := strings.Join(s, "")
	for {
		resp, err := http.PostForm(completeURL, headers)
		//fmt.Println(resp) just to check what output resp is giving
		fmt.Println(resp)
		if err != nil {
			fmt.Println(err)
		}

		// ioutil.ReadAll is being used
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			var projectID map[string]interface{}
			errs := json.Unmarshal([]byte(responseData), &projectID)
			if errs != nil {
				panic(err)
			}
			name := projectID["Title"]

			nameSp := projectID["NameSpace"]["Path"] //not sure how to get this
			if name == "message" {
				logger.Println(info)
			}
			if name == Projectname && nameSp == Namespace {
				return projectID["ID"]
			}
			if projectID == []{
				break
			}
		}
		//fmt.Println(responseData) just to check output of ioutil
	}
	page = page+1
	return
}

//func getMilestones(BaseUrl string, Token string, ProjectId string) []string {

//}
//func createMilestonesData(Advance int) []string {
//today := time.Now()
//var list []string
//for i := 0; i < Advance; i++ {

//}

//}

//func createMilestones(BaseUrl string, Token string, Project_Id string, Milestones string) string {

//}

func main() {
	// Declaring variables for flags
	var APIKey, APIBase, Namespace, Project string
	var Advance int

	// Command Line Parsing Starts
	flag.StringVar(&APIKey, "Token", " ", "Gitlab api key/token.")
	flag.StringVar(&APIBase, "BaseUrl", " ", "Gitlab api base url")
	flag.StringVar(&Namespace, "Namespace", " ", "Namespace to use in Gitlab")
	flag.StringVar(&Project, "ProjectName", " ", "Project to use in Gitlab")
	flag.IntVar(&Advance, "Advance", 30, "Define timeframe to generate milestones in advance.")
	flag.Parse() //Command Line Parsing Ends

	// Initializing logger
	LoggerSetup(os.Stdout)

}
