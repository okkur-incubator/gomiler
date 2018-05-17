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
	httpmock "gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// MockGithubAPI populates a []githubAPI with mock API data
func MockGithubAPI(state string) []githubAPI {
	currentTime := time.Now()
	githubAPImock := []githubAPI{}
	mock := githubAPI{}
	for i := 0; i < 10; i++ {
		mock.ID = i
		mock.Title = "test" + strconv.Itoa(i)
		mock.Description = "test" + strconv.Itoa(i)
		if state == "open" {
			mock.State = "open"
		} else {
			mock.State = "closed"				
		}
		mock.CreatedAt = &currentTime
		mock.UpdatedAt = &currentTime
		mock.StartDate = "test" + strconv.Itoa(i)
		mock.DueDate = "test" + strconv.Itoa(i)
		mock.Number = i

		githubAPImock = append(githubAPImock, mock)
	}

	return githubAPImock
}

// MockGithubAPIGetRequest creates a mock responder for GET requests and sends back mock JSON data
func MockGithubAPIGetRequest(URL string, state string) {
	json := MockGithubAPI(state)
	httpmock.Activate()
	var strURL []string
	strURL = []string{URL, "1", "/milestones"}
	newURL := strings.Join(strURL, "")
	httpmock.RegisterResponder("GET", newURL,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, json)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}

// MockGithubAPIPostRequest creates a mock responder for POST requests and sends back mock JSON data
func MockGithubAPIPostRequest(URL string, state string) {
	json := MockGithubAPI(state)
	var strURL []string
	strURL = []string{URL, "1", "/milestones"}
	newURL := strings.Join(strURL, "")
	httpmock.RegisterResponder("POST", newURL,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, json)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}

// MockGithubAPIPatchRequest creates a mock responder for PUT requests and sends back mock JSON data
func MockGithubAPIPatchRequest(URL string, state string, id string) {
	json := MockGithubAPI(state)
	var strURL []string
	strURL = []string{URL, "1", "/milestones/", id}
	newURL := strings.Join(strURL, "")
	httpmock.RegisterResponder("PATCH", newURL,
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, json)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}

// MockPaginate creates a mock responder to return a byte slice
func MockPaginate(url string, data []byte) {
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(200, data)
			return resp, nil
		},
	)
}
