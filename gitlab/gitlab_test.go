/*
Copyright 2017 - The GoMiler Authors

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

package gitlab

import (
	"encoding/json"
	"github.com/okkur/gomiler/utils"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
	"log"
	"os"
	"testing"
)

var (
	logger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func TestGetProjectID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "gitlab.com" + "/api/v4"
	projects := [1]gitlabAPI{}

	projects[0].Name = "test"
	projects[0].NameSpace.Path = "test"
	projects[0].ID = 1

	jsonStruct, err := json.Marshal(projects)
	if err != nil {
		t.Error(err)
	}

	httpmock.RegisterResponder("GET", "https://gitlab.com/api/v4/projects/",
		httpmock.NewStringResponder(200, string(jsonStruct)))
	res, err := GetProjectID(mockURL, "213123", "test", "test")

	if res != "1" && err != nil {
		t.Errorf("Expected %s, got %s", "1", res)
	}
}

func TestGitlabCreateAndDisplayNewMilestones(t *testing.T) {
	milestoneData := utils.CreateMilestoneData(10, "daily", nil, "gitlab")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "gitlab.com" + "/api/v4"
	MockGitlabAPIGetRequest(mockURL)
	MockGitlabAPIPostRequest(mockURL)
	err := CreateAndDisplayNewMilestones(mockURL, "213123", "1", milestoneData, logger)
	if err != nil {
		t.Error(err)
	}
}

func TestPaginate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	MockPaginate("https://example.com", []byte("testing"))
	apiData, err := paginate("https://example.com", "token")
	if err != nil {
		t.Errorf("Expected %s, got error %s: ", "testing", err)
	}
	if string(apiData[1]) != "testing" {
		t.Errorf("Expected %s, got %s", "testing", string(apiData[1]))
	}
}

func TestPaginateFailWhenURLisWrong(t *testing.T) {
	_, err := paginate("https://example.c_m", "token")
	if err == nil {
		t.Errorf("Expected to get an error when url is wrong")
	}
}

func TestGetActiveMilestones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "gitlab.com" + "/api/v4"
	MockGitlabAPIGetRequest(mockURL)
	activeMilestonesAPI, err := getActiveMilestones(mockURL, "token", "1")
	if err != nil {
		t.Error(err)
	}
	activeMilestones := createGitlabMilestoneMap(activeMilestonesAPI)
	for _, v := range activeMilestones {
		if v.State != "active" {
			t.Errorf("Expected %s, got %s", "active", v.State)
		}
	}
}

func TestGetInactiveMilestones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "gitlab.com" + "/api/v4"
	MockGitlabAPIGetRequest(mockURL)
	inactiveMilestonesAPI, err := getInactiveMilestones(mockURL, "token", "1")
	if err != nil {
		t.Error(err)
	}
	inactiveMilestones := createGitlabMilestoneMap(inactiveMilestonesAPI)
	for _, v := range inactiveMilestones {
		if v.State != "closed" {
			t.Errorf("Expected %s, got %s", "closed", v.State)
		}
	}
}
