// Description: This file contains the User struct.
package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type User struct {
	Id            int
	Name          string
	Surname       string
	Username      string
	Email         string
	Authenticated bool
}

func getUserFromDB(db *sql.DB, id int) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, name, surname, username, email FROM user WHERE id = $1", id).Scan(&user.Id, &user.Name, &user.Surname, &user.Username, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func authenticateUser(db *sql.DB, username, password string) (User, error) {
	var user User

	err := db.QueryRow("SELECT id, username FROM users WHERE username = $1 AND password = $2", username, password).Scan(&user.Id, &user.Username)

	if err != nil {
		return User{Username: username, Authenticated: false}, err
	}

	user.Authenticated = true
	return user, nil
}
