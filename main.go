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

func main() {
	var err error
	// Read the config file
	cfg, err := structs.GetConfig("config.yaml")

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

	// Get a user from the database
	user1, err := authenticateUser(db, "pero", "123456")
	if err != nil {
		log.Println(err)
	}
	fmt.Print(user1)

	fileServer := http.FileServer(http.Dir("templates/"))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/style.css", styleHandler(fileServer))
	http.HandleFunc("/login", loginPageHandler)
	http.HandleFunc("/logmein", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	// Create a custom file server
	// Serve up the index page

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

