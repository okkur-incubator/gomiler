package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger *log.Logger

//LoggerSetup Initialization
func LoggerSetup(info io.Writer) {
	logger = log.New(info, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func getProjectId(BaseUrl string, Token string, Projectname string, Namespace string) string {
	headers := map[string]string{"PrivateToken": Token}
	url := "https://" + BaseUrl + "/projects"

	page := 1
	strPage := strconv.Itoa(page)
	s := []string{url, "?page=", strPage}
	completeUrl := strings.Join(s, "")
	for {
		r, err := http.PostForm(completeUrl, url.Values())

	}

}

//func getMilestones(BaseUrl string, Token string, ProjectId string)[]string{}
func createMilestonesData(Advance int) []string {
	today := time.Now()
	var list []int
	for i := 0; i < Advance; i++ {
		date := today + time.After(i*time.Minute)

	}

}

///func createMilestones(BaseUrl string, Token string, Project_Id string, Milestones string) string{}

func main() {
	// Declaring variables for flags
	var APIKey, APIBase, Namespace, Project string
	var Advance int

	//Command Line Parsing Starts
	flag.StringVar(&APIKey, "Token", " ", "Gitlab api key/token.")
	flag.StringVar(&APIBase, "BaseUrl", " ", "Gitlab api base url")
	flag.StringVar(&Namespace, "Namespace", " ", "Namespace to use in Gitlab")
	flag.StringVar(&Project, "ProjectName", " ", "Project to use in Gitlab")
	flag.IntVar(&Advance, "Advance", 30, "Define timeframe to generate milestones in advance.")
	flag.Parse() //Command Line Parsing Ends

	//initializing logger
	LoggerSetup(os.Stdout)

}
