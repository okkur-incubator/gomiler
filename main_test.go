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
	httpmock "gopkg.in/jarcoal/httpmock.v1"
	"testing"
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

func TestValidateGithubBaseURLScheme(t *testing.T) {
	URL := "api.github.com"
	baseURL, err := validateBaseURLScheme(URL)
	if baseURL != "https://api.github.com" && err != nil {
		t.Errorf("Expected %s, got %s", "https://api.github.com", baseURL)
	}
}

func TestValidateGitlabBaseURLScheme(t *testing.T) {
	URL := "gitlab.com"
	baseURL, err := validateBaseURLScheme(URL)
	if baseURL != "https://gitlab.com" && err != nil {
		t.Errorf("Expected %s, got %s", "https://gitlab.com", baseURL)
	}
}
