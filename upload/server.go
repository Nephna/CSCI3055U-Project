
package main

import (
	"fmt"
	"os"
	"io"
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
	cookie := &http.Cookie {
		Name: "username",
		Value: username,
		Path: "/internal",
	}
	http.SetCookie(writer, cookie)
}

// ends the session by clearing the cookie
// TODO: fix to get working
func endSession (writer http.ResponseWriter) {
	cookie := &http.Cookie {
		Name: "username",
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
	fmt.Println("Login:")
	if (request.Method == "GET") {
		t, _ := template.ParseFiles("login.html")
		fmt.Println("GET")
		t.Execute(writer, nil)
	} else {
		fmt.Println("POST")
		request.ParseForm()
		name := request.FormValue("username")
		pass := request.FormValue("password")

		if valid := isValid(name, pass);(valid) {
			startSession(writer, name) // start a new session
			http.Redirect(writer, request, "/upload", 301)
		} else {
			http.Redirect(writer, request, "/login", 301) // invalid login, try again
		}
	}
}

func internal (writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Internal:")
	if (!validSession(request)) {
		http.Redirect(writer, request, "/login", 301)
		return
	}

	if (request.Method == "GET") {
		fmt.Println("GET")
		t, _ := template.ParseFiles("internal.html")
		t.Execute(writer, nil)
	}
}

// logs user out of site
func logout (writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Logout:")
	validSession(request)
	endSession(writer)
	http.Redirect(writer, request, "/login", 301)
}

 func upload (writer http.ResponseWriter, request *http.Request) {
 	if (request.Method == "GET") {
 		t, _ := template.ParseFiles("upload.html")
 		t.Execute(writer, nil)
 	} else {
	 	// the FormFile function takes in the POST input id file
	 	file, header, err := request.FormFile("file")
	 	checkError(err)

	 	defer file.Close()

	 	out, err := os.Create("/tmp/uploadedfile")
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
	http.HandleFunc("/login", login)
	http.HandleFunc("/internal", internal)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/upload", upload)

	err := http.ListenAndServe(":9090", nil) // begin listening on port 9090
	checkError(err)
}

func checkError (err error) {
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}