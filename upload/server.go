
package main

import (
	"fmt"
	"os"
	"time"
	"net/http"
	"html/template"
)

// validates name and password
func isValid (name, pass string) (bool) {
	if (name != "" && pass != "") {
		return true
	}
	return false
}

// create a new session for specified user
func startSession (writer http.ResponseWriter, username string) {
	cookie := &http.Cookie{
		Name: "username",
		Value: username,
		Path: "/internal",
	}
	http.SetCookie(writer, cookie)
}

// ends the session by clearing the cookie
// TODO: fix to get working
func endSession (writer http.ResponseWriter, username string) {
	cookie := &http.Cookie{
		Name: "",
		Value: "",
		Path: "/",
		Expires: time.Now().Add(-24 * time.Hour),
		MaxAge: -1,
	}
	http.SetCookie(writer, cookie)
}

// returns whether the session is valid or not
func validSession (request *http.Request) (bool) {
	cookie, _ := request.Cookie("username")
	if (cookie == nil) {
		return false
	}
	return true
}

// returns a name from the session cookie
func getName (request *http.Request) (username string) {
	cookie, err := request.Cookie("username")
	checkError(err)
	return cookie.Value
}

func login (writer http.ResponseWriter, request *http.Request) {
	if (request.Method == "GET") {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(writer, nil)
	} else {
		request.ParseForm()
		name := request.FormValue("username")
		pass := request.FormValue("password")

		if valid := isValid(name, pass);(valid) {
			startSession(writer, name) // start a new session
			http.Redirect(writer, request, "/internal", 301)
		} else {
			http.Redirect(writer, request, "/login", 301) // invalid login, try again
		}
	}
}

func internal (writer http.ResponseWriter, request *http.Request) {
	if (!validSession(request)) {
		http.Redirect(writer, request, "/login", 301)
		return
	}

	if (request.Method == "GET") {
		t, _ := template.ParseFiles("internal.gtpl")
		t.Execute(writer, nil)
	}
}

// logs user out of site
func logout (writer http.ResponseWriter, request *http.Request) {
	if (!validSession(request)) {
		http.Redirect(writer, request, "/login", 301)
		return
	}

	request.ParseForm()
	name := getName(request)
	endSession(writer, name)
	http.Redirect(writer, request, "/login", 301)
}

func main () {
	http.HandleFunc("/", login)
	http.HandleFunc("/login", login)
	http.HandleFunc("/internal", internal)
	http.HandleFunc("/logout", logout)

	err := http.ListenAndServe(":9090", nil) // begin listening on port 9090
	checkError(err)
}

func checkError (err error) {
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}