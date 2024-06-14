package main


import (

	"fmt"
	"net/http"
	"html/template"
)

type User struct {
	FullName     string
	Password     string
	UserID       int
	Authenticated bool
	Email        string
}

func main() {
	// Create a new user
	user := User{
		FullName: "John Doe",
		Password: "password",
		UserID:   1,
		Email:    "john@doe.net",
		Authenticated: false,
	}
	// Create a custom file server
	fileServer := http.FileServer(http.Dir("templates/"))
	// Serve up the index page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the user has a session cookie
		sessionCookie, err := r.Cookie("sessionId")
		if err != nil || sessionCookie.Value == "" {
			user.Authenticated = false
		} else {
			if sessionCookie.Value == "123" {
				user.Authenticated = true
			} else {
				user.Authenticated = false
			}
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
	})
	// Set the correct MIME type for CSS files
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		fileServer.ServeHTTP(w, r)
	})

	//set cookie with session id when accessing /login
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:  "sessionId",
			Value: "123",
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	//remove cookie when accessing /logout
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:   "sessionId",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	
	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
