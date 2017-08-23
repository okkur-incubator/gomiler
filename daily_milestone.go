package main

import "flag"

func get_project_id(base_url string, token string, projectname string, namespace string) string {
	headers := map[string]string{"PrivateToken": token}
	url := "https://" + base_url + "/projects"
	page := 1
	for {

	}

}

//func get_milestones(base_url string, token string, project_id string)[]string{}
//func create_milestone_data(advance int) []string{}
//func create_milestones(base_url string, token string, project_id string, milestones string) string{}

func main() {
	var apiKey, apiBase, namespace, project string
	var advance int
	flag.StringVar(&apiKey, "token", " ", "Gitlab api key/token.")
	flag.StringVar(&apiBase, "base_url", " ", "Gitlab api base url")
	flag.StringVar(&namespace, "namespace", " ", "Namespace to use in Gitlab")
	flag.StringVar(&project, "projectname", " ", "Project to use in Gitlab")
	flag.IntVar(&advance, "Advance", 30, "Define timeframe to generate milestones in advance.")
	flag.Parse()

}
