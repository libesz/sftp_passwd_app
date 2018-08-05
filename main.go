package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var templates *template.Template
var store = sessions.NewCookieStore([]byte("temp"))

func main() {
	r := mux.NewRouter()
	templates = template.Must(template.ParseGlob("templates/*.html"))

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/logout", logoutGetHandler).Methods("GET")
	r.HandleFunc("/changepw", changePwGetHandler).Methods("GET")
	r.HandleFunc("/changepw", changePwPostHandler).Methods("POST")

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

type PageData struct {
	Function   string
	Session    map[interface{}]interface{}
	ErrorMsg   string
	SuccessMsg string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	session, _ := store.Get(r, "session")
	data := PageData{Function: "welcome", Session: session.Values}
	templates.ExecuteTemplate(w, "index.html", data)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	session, _ := store.Get(r, "session")
	data := PageData{Function: "login", Session: session.Values}
	templates.ExecuteTemplate(w, "index.html", data)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	//password := r.PostForm.Get("password")
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["username"] = nil
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

func changePwGetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	session, _ := store.Get(r, "session")
	data := PageData{Function: "changepw", Session: session.Values}
	templates.ExecuteTemplate(w, "index.html", data)
}

func changePwPostHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	r.ParseForm()
	password := r.PostForm.Get("password")
	password2 := r.PostForm.Get("password2")
	session, _ := store.Get(r, "session")
	data := PageData{Function: "changepw", Session: session.Values}
	if password != password2 {
		data.ErrorMsg = "Password fields does not match."
	} else if len(password) < 8 {
		data.ErrorMsg = "Password must be at least 8 characters."
	} else {
		data.SuccessMsg = "Password has been successfully updated."
	}
	templates.ExecuteTemplate(w, "index.html", data)
}
