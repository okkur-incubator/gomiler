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
	"fmt"
	"log"
	"strconv"
	"time"
)

// Milestone struct to be used for milestone queries
type Milestone struct {
	DueDate string
	ID      string
	Title   string
	State   string
	Number  int
}

// LastDayMonth function to get last day of the month
func LastDayMonth(year int, month int, timezone *time.Location) time.Time {
	t := time.Date(year, time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC)
	return t
}

// LastDayWeek function to get last day of the week
func LastDayWeek(lastDay time.Time) time.Time {
	if lastDay.Weekday() != time.Sunday {
		for lastDay.Weekday() != time.Sunday {
			lastDay = lastDay.AddDate(0, 0, +1)
		}
		return lastDay
	}
	return lastDay
}

// CreateMilestoneData creates new milestones with title and due date
func CreateMilestoneData(advance int, interval string, logger *log.Logger, api string) (map[string]Milestone, error) {
	today := time.Now().Local()
	milestones := map[string]Milestone{}
	switch interval {
	case "daily":
		for i := 0; i < advance; i++ {
			var m Milestone
			var dueDate string
			title := today.AddDate(0, 0, i).Format("2006-01-02")
			switch api {
			case "gitlab":
				dueDate = today.AddDate(0, 0, i).Format("2006-01-02")
			case "github":
				dueDate = today.AddDate(0, 0, i).Format(time.RFC3339)
			}
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
		}
	case "weekly":
		for i := 0; i < advance; i++ {
			var m Milestone
			var dueDate string
			lastDay := LastDayWeek(today)
			year, week := lastDay.ISOWeek()
			title := strconv.Itoa(year) + "-w" + strconv.Itoa(week)
			switch api {
			case "gitlab":
				dueDate = lastDay.Format("2006-01-02")
			case "github":
				dueDate = lastDay.Format(time.RFC3339)
			}
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
			today = lastDay.AddDate(0, 0, 7)
		}
	case "monthly":
		for i := 0; i < advance; i++ {
			var m Milestone
			var dueDate string
			date := today.AddDate(0, i, 0)
			lastDay := LastDayMonth(date.Year(), int(date.Month()), time.UTC)
			title := date.Format("2006-01")
			switch api {
			case "gitlab":
				dueDate = lastDay.Format("2006-01-02")
			case "github":
				dueDate = lastDay.Format(time.RFC3339)
			}
			m.Title = title
			m.DueDate = dueDate
			milestones[title] = m
		}
	default:
		err := fmt.Errorf("Error: Invalid interval")
		return nil, err
	}

	return milestones, nil
}
