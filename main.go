package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"sirtediousoffoo/attender/structs"
)

var db *sql.DB
var cfg *structs.Config

func main() {
	var err error
	// Read the config file
	cfg, err = structs.GetConfig("config.yaml")

	// Connect to a postgres database
	var sqlInfo = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)

	db, err = sql.Open("postgres", sqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Check if the connection is working
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	fileServer := http.FileServer(http.Dir("templates/"))

	http.HandleFunc("/style.css", styleHandler(fileServer))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/aglogo.jpg")
	})
	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/logmein", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/termsandconditions", termsAndConditionsHandler)
	http.HandleFunc("/gdpr", gdprHandler)
	http.HandleFunc("/signmeup", signupHandler)
	http.HandleFunc("/validate/email", emailValidator)
	http.HandleFunc("/validate/username", usernameValidator)
	http.HandleFunc("/{$}", indexHandler)
	http.HandleFunc("/attend", attendHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/admin", adminViewHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/profile", profileHandler)

	// Serve up the index page
	connString := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Println("Server started on http://" + connString)
	http.ListenAndServe(connString, nil)
}
