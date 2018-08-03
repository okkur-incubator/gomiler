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

package github

import (
	"log"
	"os"
	"testing"

	"github.com/okkur/gomiler/utils"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var logger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

func TestGithubCreateAndDisplayNewMilestones(t *testing.T) {
	milestoneData, err := utils.CreateMilestoneData(10, "daily", nil, "github")
	if err != nil {
		t.Error(err)
	}
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "api.github.com"
	MockGithubAPIGetRequest(mockURL, "open")
	MockGithubAPIPostRequest(mockURL, "open")
	err = CreateAndDisplayNewMilestones(mockURL, "213123", "1", milestoneData, logger)
	if err != nil {
		t.Error(err)
	}
}

func TestGetActiveMilestones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "api.github.com"
	MockGithubAPIGetRequest(mockURL, "open")
	activeMilestonesAPI, err := getActiveMilestones(mockURL, "token", "1")
	if err != nil {
		t.Error(err)
	}
	for _, v := range activeMilestonesAPI {
		if v.State != "open" {
			t.Errorf("Expected %s, got %s", "open", v.State)
		}
	}
}

func TestGetInactiveMilestones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "api.github.com"
	MockGithubAPIGetRequest(mockURL, "closed")
	inactiveMilestonesAPI, err := getInactiveMilestones(mockURL, "token", "1")
	if err != nil {
		t.Error(err)
	}
	for _, v := range inactiveMilestonesAPI {
		if v.State != "closed" {
			t.Errorf("Expected %s, got %s", "closed", v.State)
		}
	}
}

func TestReactivateClosedMilestones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "api.github.com"
	MockGithubAPIGetRequest(mockURL, "closed")
	inactiveMilestonesAPI, err := getInactiveMilestones(mockURL, "token", "1")
	if err != nil {
		t.Error(err)
	}
	inactiveMilestones := CreateGithubMilestoneMap(inactiveMilestonesAPI)
	for _, v := range inactiveMilestones {
		MockGithubAPIPatchRequest(mockURL, "open", v.ID)
	}
	reactivatedMilestones, err := ReactivateClosedMilestones(inactiveMilestones, mockURL, "token", "1")
	if err != nil {
		t.Error(err)
	}
	for _, v := range reactivatedMilestones {
		if v.State != "open" {
			t.Errorf("Expected %s, got %s", "open", v.State)
		}
	}
}
