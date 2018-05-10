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
	httpmock "gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"strconv"
	"time"
)

// MockGitlabAPI populates a []gitlabAPI with mock API data
func MockGitlabAPI() []gitlabAPI {
	currentTime := time.Now()
	gitlabAPImock := []gitlabAPI{}
	mock := gitlabAPI{}
	for i := 0; i < 10; i++ {
		mock.ID = i
		mock.Iid = i
		mock.ProjectID = 1
		mock.Title = "test" + strconv.Itoa(i)
		mock.Description = "test" + strconv.Itoa(i)
		mock.StartDate = "test" + strconv.Itoa(i)
		mock.DueDate = "test" + strconv.Itoa(i)
		if i%2 == 0 {
			mock.State = "closed"
		} else {
			mock.State = "active"
		}
		mock.UpdatedAt = &currentTime
		mock.CreatedAt = &currentTime
		mock.Name = "test" + strconv.Itoa(i)
		mock.NameSpace.ID = i
		mock.NameSpace.Name = "test" + strconv.Itoa(i)
		mock.NameSpace.Path = "test" + strconv.Itoa(i)
		mock.NameSpace.Kind = "test" + strconv.Itoa(i)
		mock.NameSpace.FullPath = "test" + strconv.Itoa(i)
	}

	return gitlabAPImock
}

// MockGitlabAPIRequest creates a mock responder and sends back mock JSON data
func MockGitlabAPIRequest(URL string) {
	json := MockGitlabAPI()
	httpmock.Activate()
	httpmock.RegisterResponder("GET", URL,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, json)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}
