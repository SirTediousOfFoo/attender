// Description: This file contains the User struct.
package main

import (
	"database/sql"
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
}

func getUserFromDB(db *sql.DB, id int) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, name, surname, username, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Surname, &user.Username, &user.Email)
	if err != nil {
		user.Authenticated = false
		return user, err
	}
	user.Authenticated = true
	return user, nil
}

func authenticateUser(db *sql.DB, username, password string) (User, error) {
	var user User

	err := db.QueryRow("SELECT id, username FROM users WHERE username = $1 AND password = $2", username, password).Scan(&user.ID, &user.Username)

	if err != nil {
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
		log.Print("error in checkAuthenticated", err)
		return 0, err
	}

	return id, err
}
