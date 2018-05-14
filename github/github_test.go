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
	"github.com/okkur/gomiler/utils"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
	"log"
	"os"
	"testing"
)

var (
	logger = log.New(os.Stderr, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func TestGithubCreateAndDisplayNewMilestones(t *testing.T) {
	milestoneData := utils.CreateMilestoneData(10, "daily", nil, "github")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "api.github.com"
	MockGithubAPIGetRequest(mockURL)
	MockGithubAPIPostRequest(mockURL)
	err := CreateAndDisplayNewMilestones(mockURL, "213123", "1", milestoneData, logger)
	if err != nil {
		t.Error(err)
	}
}
