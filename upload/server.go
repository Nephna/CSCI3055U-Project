
package main

import (
	"fmt"
	"os"
	"io"
	"time"
	"net/http"
	"html/template"
)

const (
	LOGIN = "/login"
	MAIN = "/main"
	UPLOAD = "/upload"
	LOGOUT = "/logout"
	SLASH = "/"
	BLANK = ""

	LOGIN_HTML = "login.html"
	MAIN_HTML = "main.html"
	UPLOAD_HTML = "upload.html"
	LOGOUT_HTML = "logout.html"

	FILE_LOCATION = "." + string(os.PathSeparator) + "users" + string(os.PathSeparator)
)

// validates name and password
func isValid (name, pass string) (bool) {
	if (name != "" && pass != "") {
		return true
	}
	return false
}

func isFile (name string) (bool) {
	if (name != "") {
		return true
	}
	return false
}

// create a new session for specified user
func startSession (writer http.ResponseWriter, username string) {
	cookie := &http.Cookie {
		Name: "username",
		Value: username,
		Expires: time.Now().Add(3 * time.Minute),
		Path: SLASH,
	}
	http.SetCookie(writer, cookie)
}

// refreshes a currently active session or redirects to login page if the
// session has expired
func refreshSession (writer http.ResponseWriter, request *http.Request) {
	if (validSession(request)) {
		name := getName(request)
		startSession(writer, name)
		return
	}
	http.Redirect(writer, request, LOGIN, 301)
}

// ends the session by clearing the cookie
// TODO: fix to get working
func endSession (writer http.ResponseWriter) {
	cookie := &http.Cookie {
		Name: "username",
		Value: "",
		Path: BLANK,
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

// parses textbox input, validates it, and if valid, starts a new session for the user
func login (writer http.ResponseWriter, request *http.Request) {
	// already a valid session present, redirect to main page
	if (validSession(request)) {
		http.Redirect(writer, request, MAIN, 301)
	}

	if (request.Method == "GET") {
		t, _ := template.ParseFiles(LOGIN_HTML)
		t.Execute(writer, nil)
	} else {
		request.ParseForm()
		name := request.FormValue("username")
		pass := request.FormValue("password")

		if valid := isValid(name, pass);(valid) {
			startSession(writer, name) // start a new session
			http.Redirect(writer, request, MAIN, 301)
		} else {
			http.Redirect(writer, request, LOGIN, 301) // invalid login, try again
		}
	}
}

func internal (writer http.ResponseWriter, request *http.Request) {
	refreshSession(writer, request)

	if (request.Method == "GET") {
		t, _ := template.ParseFiles(MAIN_HTML)
		t.Execute(writer, nil)
	}
}

// logs user out of site
func logout (writer http.ResponseWriter, request *http.Request) {
	endSession(writer) // remove cookie
	http.Redirect(writer, request, LOGIN, 301)
}


func upload (writer http.ResponseWriter, request *http.Request) {
 	refreshSession(writer, request)

 	if (request.Method == "GET") {
 		t, _ := template.ParseFiles(UPLOAD_HTML)
 		t.Execute(writer, nil)
 	} else {
	 	// get file name from data in form
	 	file, header, err := request.FormFile("file")
	 	if (err != nil) {
	 		// file not found, try again
	 		http.Redirect(writer, request, UPLOAD, 301)
	 	}
	 	defer file.Close()

	 	// get account name
	 	name := getName(request)

	 	// TODO: error checking and handling
	 	dir := FILE_LOCATION + name + string(os.PathSeparator)
	 	err = os.MkdirAll(dir, 0766) // makes a new directory for this user

	 	destName := dir + header.Filename

	 	output, err := os.Create(destName)
	 	checkError(err)
	 	defer output.Close()

	 	// write the content from POST to the file
	 	_, err = io.Copy(output, file)
	 	checkError(err)

	 	// file uploaded successfully, redirect to main page
	 	http.Redirect(writer, request, MAIN, 301)
	 }
 }

func main () {
	if (len(os.Args) != 2) {
		fmt.Println("Error: incorrect amount of arguments supplied\nUsage: ", os.Args[0], "<port>")
		os.Exit(0)
	}

	// set url function handlers
	http.HandleFunc(LOGIN, login)
	http.HandleFunc(MAIN, internal)
	http.HandleFunc(LOGOUT, logout)
	http.HandleFunc(UPLOAD, upload)

	port := ":" + os.Args[1] // get port from user

	err := http.ListenAndServe(port, nil) // begin listening on port
	checkError(err)
}

func checkError (err error) {
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}