package main

import (
	"strconv"
	"time"
)

func generateCalendar(month time.Month, year int) []string {
	var calendar []string
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	calendar = append(calendar, "<tr>\n")

	if firstOfMonth.Weekday() == time.Sunday {
		for i := 0; i < 6; i++ {
			calendar = append(calendar, "<td ></td>\n")
		}
	} else if firstOfMonth.Weekday() != time.Monday {
		for i := 0; i < int(firstOfMonth.Weekday())-1; i++ {
			calendar = append(calendar, "<td ></td>\n")
		}
	}

	for i := 1; i <= lastOfMonth.Day(); i++ {
		currentDate := time.Date(year, month, i, 0, 0, 0, 0, time.UTC)
		var dayAttend int
		db.QueryRow("SELECT COUNT(*) from  attendance WHERE date = $1", currentDate).Scan(&dayAttend)
		if currentDate.Weekday() == time.Monday {
			calendar = append(calendar, "<tr>\n")
		}
		var dayEntry string
		if dayAttend != 0 {
			dayEntry = "<td id='trainingday'>" + strconv.Itoa(i) + "<div style='position:relative;left:50%;color:black;'>" + strconv.Itoa(dayAttend) + "</div></td>\n"
		} else {
			dayEntry = "<td>" + strconv.Itoa(i) + "<div style='position:relative;left:50%;color:black;'></div></td>\n"
		}
		calendar = append(calendar, dayEntry)

		if currentDate.Weekday() == time.Sunday || i == lastOfMonth.Day() {
			calendar = append(calendar, "</tr>\n")
		}
	}

	return calendar
}
