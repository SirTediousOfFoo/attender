package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func createSessionCookie(db *sql.DB, userID int) *http.Cookie {
	sessionID := uuid.New().String()
	expirydate := time.Now().Add(30 * 24 * time.Hour) // 30 days cookie duration
	insertStatement := `INSERT INTO sessions (sessionid, userid, expirydate) VALUES ($1, $2, $3)`
	_, err := db.Exec(insertStatement, sessionID, userID, expirydate)
	if err != nil {
		panic(err)
	}

	cookie := &http.Cookie{
		Name:    "sessionId",
		Value:   sessionID,
		Expires: expirydate,
	}
	return cookie
}
