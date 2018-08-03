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

package utils

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestLastDayMonth(t *testing.T) {
	date := time.Now().Local()
	lastDay := LastDayMonth(date.Year(), int(date.Month()), time.UTC)
	expectedDay := time.Date(date.Year(), time.Month(date.Month())+1, 0, 0, 0, 0, 0, time.UTC)
	if lastDay != expectedDay {
		t.Errorf("Expected %v, got %v", expectedDay, lastDay)
	}
}

func TestLastDayWeek(t *testing.T) {
	date := time.Now().Local()
	lastDay := LastDayWeek(date)
	if lastDay.Weekday() != time.Sunday {
		t.Errorf("Expected %s, got %s", time.Sunday, lastDay.Weekday())
	}
}

func TestGithubCreateMilestoneDataDaily(t *testing.T) {
	milestones, err := CreateMilestoneData(30, "daily", nil, "github")
	if err != nil {
		t.Error(err)
	}
	today := time.Now().Local().Format("2006-01-02")
	todayFormatted := time.Now().Local().Format(time.RFC3339)
	if milestones[today].DueDate != todayFormatted {
		t.Errorf("Expected %s, got %s", today, milestones[today].DueDate)
	}
}

func TestGithubCreateMilestoneDataWeekly(t *testing.T) {
	milestones, err := CreateMilestoneData(20, "weekly", nil, "github")
	if err != nil {
		t.Error(err)
	}
	today := time.Now().Local()
	lastDay := LastDayWeek(today)
	year, week := lastDay.ISOWeek()
	title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
	expected := lastDay.Format(time.RFC3339)
	if milestones[title].DueDate != expected {
		t.Errorf("Expected %s, got %s", expected, milestones[title].DueDate)
	}
}

func TestGithubCreateMilestoneDataMonthly(t *testing.T) {
	milestones, err := CreateMilestoneData(2, "monthly", nil, "github")
	if err != nil {
		t.Error(err)
	}
	currentMonth := time.Now().Local().Format("2006-01")
	expected := LastDayMonth(time.Now().Local().Year(), int(time.Now().Local().Month()), time.UTC).Format(time.RFC3339)
	if milestones[currentMonth].DueDate != expected {
		t.Errorf("Expected %s, got %s", expected, milestones[currentMonth].DueDate)
	}
}

func TestGitlabCreateMilestoneDataDaily(t *testing.T) {
	milestones, err := CreateMilestoneData(30, "daily", nil, "gitlab")
	if err != nil {
		t.Error(err)
	}
	today := time.Now().Local().Format("2006-01-02")
	if milestones[today].DueDate != today {
		t.Errorf("Expected %s, got %s", today, milestones[today].DueDate)
	}
}

func TestGitlabCreateMilestoneDataWeekly(t *testing.T) {
	milestones, err := CreateMilestoneData(20, "weekly", nil, "gitlab")
	if err != nil {
		t.Error(err)
	}
	today := time.Now().Local()
	lastDay := LastDayWeek(today)
	year, week := lastDay.ISOWeek()
	title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
	expected := lastDay.Format("2006-01-02")
	if milestones[title].DueDate != expected {
		t.Errorf("Expected %s, got %s", expected, milestones[title].DueDate)
	}
}

func TestGitlabCreateMilestoneDataMonthly(t *testing.T) {
	milestones, err := CreateMilestoneData(2, "monthly", nil, "gitlab")
	if err != nil {
		t.Error(err)
	}
	currentMonth := time.Now().Local().Format("2006-01")
	expected := LastDayMonth(time.Now().Local().Year(), int(time.Now().Local().Month()), time.UTC).Format("2006-01-02")
	if milestones[currentMonth].DueDate != expected {
		t.Errorf("Expected %s, got %s", expected, milestones[currentMonth].DueDate)
	}
}

func TestGithubCreateMilestoneDataDailyWrongDueDate(t *testing.T) {
	milestones, err := CreateMilestoneData(30, "daily", nil, "github")
	if err != nil {
		t.Error(err)
	}
	today := time.Now().Local().Format("2006-01-02")
	expected := time.Now().Local().Format(time.RFC3339)
	if milestones[today].DueDate == today {
		t.Errorf("Expected %s, got %s", expected, milestones[today].DueDate)
	}
}

func TestGithubCreateMilestoneDataWeeklyWrongDueDate(t *testing.T) {
	milestones, err := CreateMilestoneData(20, "weekly", nil, "github")
	if err != nil {
		t.Error(err)
	}
	today := time.Now().Local()
	lastDay := LastDayWeek(today)
	year, week := lastDay.ISOWeek()
	title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
	expected := lastDay.Format(time.RFC3339)
	if milestones[title].DueDate == title {
		t.Errorf("Expected %s, got %s", expected, milestones[title].DueDate)
	}
}

func TestGithubCreateMilestoneDataMonthlyWrongDueDate(t *testing.T) {
	milestones, err := CreateMilestoneData(2, "monthly", nil, "github")
	if err != nil {
		t.Error(err)
	}
	currentMonth := time.Now().Local().Format("2006-01")
	expected := LastDayMonth(time.Now().Local().Year(), int(time.Now().Local().Month()), time.UTC).Format(time.RFC3339)
	if milestones[currentMonth].DueDate == currentMonth {
		t.Errorf("Expected %s, got %s", expected, milestones[currentMonth].DueDate)
	}
}

func TestGitlabCreateMilestoneDataWeeklyWrongDueDate(t *testing.T) {
	milestones, err := CreateMilestoneData(20, "weekly", nil, "gitlab")
	if err != nil {
		t.Error(err)
	}
	today := time.Now().Local()
	lastDay := LastDayWeek(today)
	year, week := lastDay.ISOWeek()
	title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
	expected := lastDay.Format("2006-01-02")
	if milestones[title].DueDate == title {
		t.Errorf("Expected %s, got %s", expected, milestones[title].DueDate)
	}
}

func TestGitlabCreateMilestoneDataMonthlyWrongDueDate(t *testing.T) {
	milestones, err := CreateMilestoneData(2, "monthly", nil, "gitlab")
	if err != nil {
		t.Error(err)
	}
	currentMonth := time.Now().Local().Format("2006-01")
	expected := LastDayMonth(time.Now().Local().Year(), int(time.Now().Local().Month()), time.UTC).Format("2006-01-02")
	if milestones[currentMonth].DueDate == currentMonth {
		t.Errorf("Expected %s, got %s", expected, milestones[currentMonth].DueDate)
	}
}

func TestGithubCreateMilestoneDataWrongInterval(t *testing.T) {
	_, err := CreateMilestoneData(30, "2", nil, "github")
	if err == nil {
		t.Errorf("Expected to get an error when interval invalid")
	}
}

func TestGitlabCreateMilestoneDataWrongInterval(t *testing.T) {
	_, err := CreateMilestoneData(30, "2", nil, "gitlab")
	if err == nil {
		t.Errorf("Expected to get an error when interval invalid")
	}
}

func TestPaginate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	pages := MockPaginate("https://example.com")
	apiData, err := Paginate("https://example.com", "github", "token")
	if err != nil {
		t.Error(err)
	}
	if len(apiData) != pages {
		t.Errorf("Expected %d, got %d", pages, len(apiData))
	}
}

func TestPaginateFailWhenURLisWrong(t *testing.T) {
	_, err := Paginate("https://example.c_m", "github", "token")
	if err == nil {
		t.Errorf("Expected to get an error when url is wrong")
	}
}

// MockPaginate creates a mock responder to return a byte slice
func MockPaginate(url string) int {
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
