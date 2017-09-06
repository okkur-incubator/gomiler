package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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
	strPage := strconv.Itoa(page)
	s := []string{urls, "?page=", strPage}
	completeURL := strings.Join(s, "")
	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", completeURL, nil)
		if err != nil {
			fmt.Println(err)
			break
		}
		req.Header.Add("Token", "text/plain")
		req.Header.Add("User-Agent", "SampleClient/1.0")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Print(err)
			break
		}
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)
		return "xjson" // JSON part to be added
	}
	return "json" // JSON part to be added
}

func main() {
	// Declaring variables for flags
	var APIKey, APIBase, Namespace, Project string
	var Advance int
	APIBase = "lol"
	// Command Line Parsing Starts
	flag.StringVar(&APIKey, "Token", "lol", "Gitlab api key/token.")
	flag.StringVar(&APIBase, "BaseURL", " ", "Gitlab api base url")
	flag.StringVar(&Namespace, "Namespace", " ", "Namespace to use in Gitlab")
	flag.StringVar(&Project, "ProjectName", " ", "Project to use in Gitlab")
	flag.IntVar(&Advance, "Advance", 30, "Define timeframe to generate milestones in advance.")
	flag.Parse() //Command Line Parsing Ends

	// Initializing logger
	LoggerSetup(os.Stdout)
	//logger.Println("info")
	getProjectID(APIBase, APIKey, Project, Namespace)
}
