package main

import (
	"html/template"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
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
