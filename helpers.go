package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

func getMonthYear(r *http.Request) (time.Month, int) {
	month := time.Now().Month()
	if r.URL.Query().Get("month") != "" {
		monthNum, err := strconv.Atoi(r.URL.Query().Get("month"))
		if err != nil {
			log.Println("Error converting month to int", err)
		}
		month = time.Month(monthNum)
	}

	year := time.Now().Year()
	if r.URL.Query().Get("year") != "" {
		var err error
		year, err = strconv.Atoi(r.URL.Query().Get("year"))
		if err != nil {
			log.Println("Error converting year to int", err)
		}
	}

	return month, year
}
