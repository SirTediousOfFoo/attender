// Description: This file contains the User struct.
package main

import (
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

// User holds data for our current user
type User struct {
	ID            int
	Name          string
	Surname       string
	Username      string
	Email         string
	Authenticated bool
	Stats         Stats
	Admin         bool
}

// Stats holds attendance data for our current user, probably will be nil for most cases except for the stats page
type Stats struct {
	AttendedTotal   int
	AttendedMonthly int
	AttendedYearly  int
}

func getUserFromDB(db *sql.DB, id int) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, name, surname, username, email, admin FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Surname, &user.Username, &user.Email, &user.Admin)
	if err != nil {
		user.Authenticated = false
		return user, err
	}
	user.Authenticated = true
	return user, nil
}

func authenticateUser(db *sql.DB, username, password string) (User, error) {

	var user User
	// Hash the password
	password = password + cfg.PasswordSalt
	password = hex.EncodeToString(sha512.New512_256().Sum([]byte(password))[:])
	err := db.QueryRow("SELECT id, username FROM users WHERE username = $1 AND password = $2", username, password).Scan(&user.ID, &user.Username)

	if err != nil {
		log.Println(err, password)
		return User{Username: username, Authenticated: false}, err
	}

	user.Authenticated = true
	return user, nil
}

func checkAuthenticated(db *sql.DB, r *http.Request) (int, error) {
	var id int
	sessionCookie, err := r.Cookie("sessionId")
	if err != nil {
		return 0, err
	}
	// Check if current userID is in the session table and not expired
	err = db.QueryRow("SELECT userid FROM sessions WHERE sessionid = $1 AND expirydate > $2", sessionCookie.Value, time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, err
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		user, err := authenticateUser(db, r.FormValue("username"), r.FormValue("password"))
		log.Println(user, err)
		if err != nil {
			t := template.Must(template.New("status").Parse(loginBad))
			t.Execute(w, r.FormValue("username"))
			return
		}

		if user.Authenticated {
			cookie := createSessionCookie(db, user.ID)
			http.SetCookie(w, cookie)
		} else {
			cookie := &http.Cookie{
				Name:   "sessionId",
				Value:  "",
				MaxAge: -1,
			}
			http.SetCookie(w, cookie)
		}
		w.Header().Add("HX-Location", "/")
		w.Header().Add("HX-Refresh", "true")
		w.Header().Add("hx-push-url", "/")
		http.Redirect(w, r, "/", http.StatusOK)
		return
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// First we make sure the user doesn't already exist
		var out string
		err := db.QueryRow("SELECT email FROM users WHERE email = $1", r.FormValue("email")).Scan(&out)
		if !errors.Is(err, sql.ErrNoRows) {
			t := template.Must(template.New("status").Parse(emailBad))
			t.Execute(w, out)
			return
		}

		if r.FormValue("name") == "" || r.FormValue("surname") == "" || r.FormValue("username") == "" || r.FormValue("email") == "" || r.FormValue("password") == "" || r.FormValue("gdpr") != "on" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Hash the password
		password := r.FormValue("password") + cfg.PasswordSalt
		password = hex.EncodeToString(sha512.New512_256().Sum([]byte(password))[:])

		// Insert the user into the database
		_, err = db.Exec("INSERT INTO users (name, surname, username, email, password) VALUES ($1, $2, $3, $4, $5)", r.FormValue("name"), r.FormValue("surname"), r.FormValue("username"), r.FormValue("email"), password)
		log.Println(r.FormValue("name"), r.FormValue("surname"), r.FormValue("username"), r.FormValue("email"), password)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Authenticate the user so we can log him in right away
		user, err := authenticateUser(db, r.FormValue("username"), r.FormValue("password"))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Create a session cookie
		if user.Authenticated {
			cookie := createSessionCookie(db, user.ID)
			http.SetCookie(w, cookie)
		} else {
			cookie := &http.Cookie{
				Name:   "sessionId",
				Value:  "",
				MaxAge: -1,
			}
			http.SetCookie(w, cookie)
			return
		}
		//send an htmx hx-redirect header with the response
		w.Header().Add("HX-Location", "/")
		w.Header().Add("HX-Refresh", "true")
		w.Header().Add("hx-push-url", "/")
		http.Redirect(w, r, "/", http.StatusOK)
		return
	}
}
