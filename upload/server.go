
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
	INTERNAL = "/internal"
	UPLOAD = "/upload"
	LOGOUT = "/logout"
	SLASH = "/"
	BLANK = ""

	LOGIN_HTML = "login.html"
	INTERNAL_HTML = "internal.html"
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
		Path: SLASH,
	}
	http.SetCookie(writer, cookie)
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
		fmt.Println("Invalid session")
		return false
	}
	fmt.Println(cookie.Value)
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
		t, _ := template.ParseFiles(LOGIN_HTML)
		t.Execute(writer, nil)
	} else {
		request.ParseForm()
		name := request.FormValue("username")
		pass := request.FormValue("password")

		if valid := isValid(name, pass);(valid) {
			startSession(writer, name) // start a new session
			http.Redirect(writer, request, UPLOAD, 301)
		} else {
			http.Redirect(writer, request, LOGIN, 301) // invalid login, try again
		}
	}
}

func internal (writer http.ResponseWriter, request *http.Request) {
	if (!validSession(request)) {
		http.Redirect(writer, request, LOGIN, 301)
	}

	if (request.Method == "GET") {
		t, _ := template.ParseFiles(INTERNAL_HTML)
		t.Execute(writer, nil)
	}
}

// logs user out of site
func logout (writer http.ResponseWriter, request *http.Request) {
	endSession(writer)
	http.Redirect(writer, request, LOGIN, 301)
}

 func upload (writer http.ResponseWriter, request *http.Request) {
 	if (!validSession(request)) {
 		http.Redirect(writer, request, LOGIN, 301)
 	}

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

	 	out, err := os.Create(destName)
	 	checkError(err)
	 	defer out.Close()

	 	// write the content from POST to the file
	 	_, err = io.Copy(out, file)
	 	checkError(err)

	 	fmt.Fprintf(writer, "File uploaded successfully : ")
	 	fmt.Fprintf(writer, header.Filename)
	 }
 }

func main () {
	//http.HandleFunc("/", login)
	http.HandleFunc(LOGIN, login)
	http.HandleFunc(INTERNAL, internal)
	http.HandleFunc(LOGOUT, logout)
	http.HandleFunc(UPLOAD, upload)

	err := http.ListenAndServe(":9090", nil) // begin listening on port 9090
	checkError(err)
}

func checkError (err error) {
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}