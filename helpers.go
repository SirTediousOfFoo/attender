package main

import (
	"database/sql"
	"log"
	"net/http"
	"sirtediousoffoo/attender/structs"
	"sort"
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

func getMembers(db *sql.DB, from, to time.Time) []structs.Member {
	rows, err := db.Query("SELECT u.id, u.name, u.surname, a.date from users u right JOIN  attendance a on u.id = a.userid WHERE a.date BETWEEN $1 AND $2", from, to)
	if err != nil {
		log.Println("Error getting members", err)
	}
	defer rows.Close()
	memberList := []structs.Member{}
	for rows.Next() {
		var m structs.Member
		var date time.Time
		err := rows.Scan(&m.ID, &m.Name, &m.Surname, &date)
		if err != nil {
			log.Println("Error scanning rows", err)
		}
		m.Dates = append(m.Dates, date)
		memberList = append(memberList, m)
	}
	sort.Slice(memberList, func(i, j int) bool {
		return memberList[i].Dates[0].Before(memberList[j].Dates[0])
	})
	return memberList
}
