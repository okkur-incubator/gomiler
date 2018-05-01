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
	"strconv"
	"testing"
	"time"
)

func TestCreateMilestoneDataDaily(t *testing.T) {
	milestones := utils.CreateMilestoneData(30, "daily", nil, "github")
	today := time.Now().Local().Format("2006-01-02")
	todayFormatted := time.Now().Local().Format(time.RFC3339)
	if milestones[today].DueDate != todayFormatted {
		t.Errorf("Expected %s, got %s", today, milestones[today].DueDate)
	}
}

func TestCreateMilestoneDataWeekly(t *testing.T) {
	milestones := utils.CreateMilestoneData(20, "weekly", nil, "github")
	today := time.Now().Local()
	lastDay := utils.LastDayWeek(today)
	year, week := lastDay.ISOWeek()
	title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
	expected := lastDay.Format(time.RFC3339)
	if milestones[title].DueDate != expected {
		t.Errorf("Expected %s, got %s", expected, milestones[title].DueDate)
	}
}

func TestCreateMilestoneDataMonthly(t *testing.T) {
	milestones := utils.CreateMilestoneData(2, "monthly", nil, "github")
	currentMonth := time.Now().Local().Format("2006-01")
	expected := utils.LastDayMonth(time.Now().Local().Year(), int(time.Now().Local().Month()), time.UTC).Format(time.RFC3339)
	if milestones[currentMonth].DueDate != expected {
		t.Errorf("Expected %s, got %s", expected, milestones[currentMonth].DueDate)
	}
}
