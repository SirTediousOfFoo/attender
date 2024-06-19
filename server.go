package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user has a session cookie
	userID, err := checkAuthenticated(db, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	user, err := getUserFromDB(db, userID)
	if err != nil {
		log.Println(err)
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
	// Check if the user has a session cookie
	id, err := checkAuthenticated(db, r)
	if err == nil && id != 0 {
		// Redirect to the index page
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(db, r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		//TODO rediredct to error page
		log.Println("am i here", err)
	}
	fmt.Print(user)

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
	http.Redirect(w, r, "/", http.StatusSeeOther)
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

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Hash the password
		password := r.FormValue("password") + cfg.PasswordSalt
		password = hex.EncodeToString(sha512.New512_256().Sum([]byte(password))[:])
		// Insert the user into the database
		_, err := db.Exec("INSERT INTO users (name, surname, username, email, password) VALUES ($1, $2, $3, $4, $5)", r.FormValue("name"), r.FormValue("surname"), r.FormValue("username"), r.FormValue("email"), password)
		log.Println(r.FormValue("name"), r.FormValue("surname"), r.FormValue("username"), r.FormValue("email"), password)
		if err != nil {
			log.Println("aaaaa", err)
		}
		// Authenticate the user so we can log him in right away
		user, err := authenticateUser(db, r.FormValue("username"), password)
		if err != nil {
			log.Println("ovdje", err)
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
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
