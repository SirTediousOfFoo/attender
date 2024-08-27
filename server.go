package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"time"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(
			w,
			"That method is not allowed.",
			http.StatusMethodNotAllowed,
		)
		return
	}
	// Check if the user has a session cookie
	userID, err := checkAuthenticated(db, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := getUserFromDB(db, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Render the index.gohtml template
	tmpl, err := template.ParseFiles("templates/index.gohtml", "templates/userMenu.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func styleHandler(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		fileServer.ServeHTTP(w, r)
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(
			w,
			"That method is not allowed.",
			http.StatusMethodNotAllowed,
		)
		return
	}
	// Check if the user has a session cookie
	id, err := checkAuthenticated(db, r)
	if err == nil && id != 0 {
		// Redirect to the index page
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render the index.gohtml template
	tmpl, err := template.ParseFiles("templates/login.gohtml", "templates/userMenu.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "sessionId",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func termsAndConditionsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/tnc.html")
}

func attendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(
			w,
			"That method is not allowed.",
			http.StatusMethodNotAllowed,
		)
		return
	}

	// Check if the user has a session cookie
	userID, err := checkAuthenticated(db, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := getUserFromDB(db, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Render the index.gohtml template
	tmpl, err := template.ParseFiles("templates/attended.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Log the attendance in the db
	_, err = db.Exec("INSERT INTO attendance (userid, date) VALUES ($1, $2)", user.ID, time.Now())
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"attendance_unique\"" {
			tmpl, err = template.ParseFiles("templates/alreadyAttended.gohtml")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.Execute(w, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Return a random happy message
	sucessString := []string{"Yay", ":)", "Woo", "Huzzah", "Awesome", "Fantastic", "Yippee", "Alright", "^_^", "Cool", "Cool beans"}
	err = tmpl.Execute(w, sucessString[rand.Intn(len(sucessString))])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(
			w,
			"That method is not allowed.",
			http.StatusMethodNotAllowed,
		)
		return
	}
	// Check if the user has a session cookie
	userID, err := checkAuthenticated(db, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := getUserFromDB(db, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Render the stats.gohtml template
	month, year := getMonthYear(r)
	fmt.Println(month, year)

	tmpl, err := template.New("stats.gohtml").Funcs(template.FuncMap{
		"currMonth": func() string {
			return month.String()
		},
		"makeYearSelector": func() []int {
			var out []int
			minYear := 2024
			db.QueryRow("SELECT date_part('year', MIN(date)) FROM attendance WHERE userid = $1", user.ID).Scan(&minYear)
			minYear = minYear - 4
			for i := minYear; i <= time.Now().Year(); i++ {
				out = append(out, i)
			}
			return out
		},
		"currYear": func() int {
			return year
		},
		"youVsOthers": func() int {
			var out string
			err := db.QueryRow("SELECT coalesce(DIV(sum(totalForUser), count(totalForUser)), 0) from (SELECT count(date) as totalForUser FROM attendance WHERE date_part('year',date) = $1 AND date_part('month',date) = $2 AND userid != $3 GROUP BY userid) as totalForUser", year, int(month), user.ID).Scan(&out)
			fmt.Println("average", out, month, int(month))
			if err != nil {
				log.Println("Error scanning database", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return 0
			}
			outint, _ := strconv.Atoi(out)
			return outint
		},
	}).ParseFiles("templates/stats.gohtml", "templates/userMenu.gohtml")
	if err != nil {
		log.Println("Error parsing template", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendedTotal := db.QueryRow("SELECT COUNT(date) FROM attendance WHERE userid = $1", user.ID)
	err = attendedTotal.Scan(&user.Stats.AttendedTotal)
	if err != nil {
		log.Println("Error scanning database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendedYearly := db.QueryRow("SELECT COUNT(date) FROM attendance WHERE userid = $1 AND date_part('year', date) = $2", user.ID, year)
	err = attendedYearly.Scan(&user.Stats.AttendedYearly)
	if err != nil {
		log.Println("Error scanning database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attendedMonthly := db.QueryRow("SELECT COUNT(date) FROM attendance WHERE userid = $1 AND date_part('month', date) = $2 and date_part('year', date) = $3", user.ID, month, year)
	err = attendedMonthly.Scan(&user.Stats.AttendedMonthly)
	log.Println(user.Stats.AttendedMonthly)
	if err != nil {
		log.Println("Error scanning database", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, user)
	if err != nil {
		log.Println("Error executing template", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func adminViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(
			w,
			"That method is not allowed.",
			http.StatusMethodNotAllowed,
		)
		return
	}
	// Check if the user has a session cookie
	userID, err := checkAuthenticated(db, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user, err := getUserFromDB(db, userID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// only admins can view this page
	if !user.Admin {
		http.Error(w, "You are not an admin", http.StatusForbidden)
		return
	}
	month, year := getMonthYear(r)
	// Render the aminview.gohtml template
	tmpl, err := template.New("adminView.gohtml").Funcs(template.FuncMap{
		"currMonth": func() string {
			return month.String()
		},
		"makeYearSelector": func() []int {
			var out []int
			minYear := 2024
			db.QueryRow("SELECT date_part('year', MIN(date)) FROM attendance WHERE userid = $1", user.ID).Scan(&minYear)
			for i := minYear; i <= time.Now().Year(); i++ {
				out = append(out, i)
			}
			return out
		},
		"currYear": func() int {
			return year
		},
		"monthAverage": func() int {
			var out string
			err := db.QueryRow("SELECT DIV(attendances.Total, attendances.Cnt) from (SELECT count(*) as Total, greatest(count(distinct date),1) as Cnt from attendance where date_part('year', date) = $1 and date_part('month', date) = $2) as attendances", year, int(month)).Scan(&out)
			if err != nil {
				log.Println("Error scanning database", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return 0
			}
			outint, _ := strconv.Atoi(out)
			return outint
		},
		"yearAverage": func() int {
			var out string
			err := db.QueryRow("SELECT DIV(attendances.Total, attendances.Cnt) from (SELECT count(*) as Total, greatest(count(distinct date),1) as Cnt from attendance where date_part('year', date) = $1) as attendances", year).Scan(&out)
			if err != nil {
				log.Println("Error scanning database", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return 0
			}
			outint, _ := strconv.Atoi(out)
			return outint
		},
	}).ParseFiles("templates/adminView.gohtml", "templates/userMenu.gohtml")

	calTpl, err := template.ParseFiles("templates/calendar.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	stringFlat := strings.Join(generateCalendar(time.Month(month), year), " ")
	calTpl.Execute(w, template.HTML(stringFlat))

	memberList := getMembers(db, time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC), time.Now())
	membersTpl, err := template.ParseFiles("templates/memberAttendanceList.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	membersTpl.Execute(w, memberList)
}
