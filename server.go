package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func checkAuthenticated(r *http.Request) User {
	sessionCookie, err := r.Cookie("sessionId")
	user := User{}
	if err != nil || sessionCookie.Value == "" {
			user.Authenticated = false
	} else {
			if sessionCookie.Value == "123" {
					user.Authenticated = true
			} else {
					user.Authenticated = false
			}
	}
	return user
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user has a session cookie
	user := checkAuthenticated(r)

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
	user := checkAuthenticated(r)

	// Render the index.gohtml template
	tmpl, err := template.ParseFiles("templates/login.gohtml", "templates/userMenu.gohtml")
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	user, err := authenticateUser(db, r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		log.Println(err)
	}
	fmt.Print(user)

	if user.Authenticated {
		cookie := &http.Cookie{
				Name:  "sessionId",
				Value: "123",
		}
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
