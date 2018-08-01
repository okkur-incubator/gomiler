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

package main

import (
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestGithubCheckAPI(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "api.github.com"
	httpmock.RegisterResponder("GET", "https://api.github.com/repos/test/test",
		httpmock.NewStringResponder(200, ""))
	res, err := checkAPI(mockURL, "token", "test", "test")
	if res != "github" && res != "" && err != nil {
		t.Errorf("Expected %s, got %s", "github", res)
	}
}

func TestGitlabCheckAPI(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "gitlab.com"
	httpmock.RegisterResponder("GET", "https://gitlab.com/api/v4/version",
		httpmock.NewStringResponder(200, ""))
	res, err := checkAPI(mockURL, "token", "test", "test")
	if res != "gitlab" && res != "" && err != nil {
		t.Errorf("Expected %s, got %s", "gitlab", res)
	}
}

func TestGithubCheckAPIwithInvalidToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "api.github.com"
	httpmock.RegisterResponder("GET", "https://api.github.com/repos/test/test",
		httpmock.NewStringResponder(403, ""))
	_, err := checkAPI(mockURL, "token", "test", "test")
	if err == nil {
		t.Errorf("Expected to get an error when token is invalid")
	}
}

func TestGitlabCheckAPIwithInvalidToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockURL := "https://" + "gitlab.com"
	httpmock.RegisterResponder("GET", "https://gitlab.com/api/v4/version",
		httpmock.NewStringResponder(403, ""))
	_, err := checkAPI(mockURL, "token", "test", "test")
	if err == nil {
		t.Errorf("Expected to get an error when token is invalid")
	}
}

func TestValidateBaseURLScheme(t *testing.T) {
	URL := "example.com"
	baseURL, err := validateBaseURLScheme(URL)
	if baseURL != "https:/example.com" && err != nil {
		t.Errorf("Expected %s, got %s", "https://example.com", baseURL)
	}
}

func TestValidateBaseURLSchemeWhenSchemeAlreadyExists(t *testing.T) {
	URL := "https://example.com"
	baseURL, err := validateBaseURLScheme(URL)
	if baseURL != "https://example.com" && err != nil {
		t.Errorf("Expected %s, got %s", "https://example.com", baseURL)
	}
}

func TestPaginate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	pages := mockPaginate("https://example.com")
	apiData, err := paginate("https://example.com", "token")
	if err != nil {
		t.Error(err)
	}
	if len(apiData) != pages {
		t.Errorf("Expected %d, got %d", pages, len(apiData))
	}
}

func TestPaginateFailWhenURLisWrong(t *testing.T) {
	_, err := paginate("https://example.c_m", "token")
	if err == nil {
		t.Errorf("Expected to get an error when url is wrong")
	}
}

// MockPaginate creates a mock responder to return a byte slice
func mockPaginate(url string) int {
	linkHeader := []string{
		"<http://example.com/page=1>; rel=\"next\", <http://example.com/page=3>; rel=\"last\"",
		"<http://example.com/page=3>; rel=\"next\", <http://example.com/page=3>; rel=\"last\"",
		"<http://example.com/page=2>; rel=\"first\", <http://example.com/page=3>; rel=\"last\"",
	}
	links := []string{
		"http://example.com/page=1",
		"http://example.com/page=3",
		"http://example.com/page=2",
	}
	for i, link := range linkHeader {
		httpmock.RegisterResponder("GET", links[i],
			func(req *http.Request) (*http.Response, error) {
				resp := httpmock.NewStringResponse(200, "testing")
				resp.Header.Set("Link", link)
				return resp, nil
			},
		)
	}
	httpmock.RegisterResponder("GET", url,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(200, "testing")
			resp.Header.Set("Link", linkHeader[0])
			return resp, nil
		},
	)
	return len(linkHeader)
}
