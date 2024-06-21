package main

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
)

const emailBad = `
<input style="margin-bottom:5xp;box-shadow: 0 0 3px #CC0000;" type="email" name="email" hx-post="/validate/email" placeholder="Email" required="true" value={{ . }}>
<label style="position: relative;width: 60%;display: block;margin: auto;height: auto;color: red;" for="email">Email is already in use</label>
`
const emailGood = `
<input style="box-shadow: 0 0 3px #36cc00;" type="email" name="email" hx-post="/validate/email" placeholder="Email" required="true" value={{ . }}>
`
const usernameBad = `
<input style="margin-bottom:5xp;box-shadow: 0 0 3px #CC0000;" type="txt" name="username" hx-post="/validate/username" placeholder="Username" required="true" value={{ . }}>
<label style="position: relative;width: 60%;display: block;margin: auto;height: auto;color: red;" for="username">Username is already in use</label>
`
const usernameGood = `
<input style="box-shadow: 0 0 3px #36cc00;" type="txt" name="username" hx-post="/validate/username" placeholder="Username" required="true" value={{ . }}>
`

const loginBad = `
	<form action="/logmein">
		<label class="loginlabel" for="chk" aria-hidden="true">Login</label>
		<input type="txt" name="username" placeholder="Username" required="true" value="{{ . }}">
		<input type="password" name="password" placeholder="Password" required="true" value="">
		<label style="position: relative;width: 60%;display: block;margin: auto;height: auto;color: red;" for="password">Username or password incorrect</label>
		<button style="box-shadow: 0 0 3px #CC0000;" type="submit" hx-post="/logmein" hx-swap="none">Login</button>
	</form>
`

func emailValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Check if the username is already in use
		var email string
		err := db.QueryRow("SELECT email FROM users WHERE email = $1", r.FormValue("email")).Scan(&email)
		if !errors.Is(err, sql.ErrNoRows) {
			t := template.Must(template.New("status").Parse(emailBad))
			t.Execute(w, r.FormValue("email"))
			return
		}
		t := template.Must(template.New("status").Parse(emailGood))
		t.Execute(w, r.FormValue("email"))
	}
}

func usernameValidator(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Check if the username is already in use
		var username string
		err := db.QueryRow("SELECT username FROM users WHERE username = $1", r.FormValue("username")).Scan(&username)
		if !errors.Is(err, sql.ErrNoRows) {
			t := template.Must(template.New("status").Parse(usernameBad))
			t.Execute(w, r.FormValue("username"))
			return
		}
		t := template.Must(template.New("status").Parse(usernameGood))
		t.Execute(w, r.FormValue("username"))
	}
}
