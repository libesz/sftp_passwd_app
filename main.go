package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
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
		Addr:    ":8000",
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
	password := r.PostForm.Get("password")
	pwValid, err := validatePassword(username, password)
	if err != nil {
		//TODO!!!
		return
	}
	if pwValid {
		session, _ := store.Get(r, "session")
		session.Values["username"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		data := PageData{Function: "login"}
		data.ErrorMsg = "Wrong username or password!"
		templates.ExecuteTemplate(w, "index.html", data)
	}
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
	username, _ := session.Values["username"].(string)
	data := PageData{Function: "changepw", Session: session.Values}
	if password != password2 {
		data.ErrorMsg = "Password fields does not match."
	} else if len(password) < 8 {
		data.ErrorMsg = "Password must be at least 8 characters."
	} else {
		_, err := changePassword(username, password)
		if err != nil {
			//TODO
			return
		}
		data.SuccessMsg = "Password has been successfully updated."
	}
	templates.ExecuteTemplate(w, "index.html", data)
}

func validatePassword(username, password string) (bool, error) {
	file, err := os.Open(os.Getenv("USERS_CONF"))
	if err != nil {
		return false, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		if s[0] == username && s[1] == password {
			return true, nil
		}
	}
	return false, nil
}

func changePassword(username, password string) (bool, error) {
	file, err := os.OpenFile(os.Getenv("USERS_CONF"), os.O_RDWR, 0755)
	if err != nil {
		return false, err
	}
	defer file.Close()
	var users []string
	changed := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		log.Println("read:", s)
		ss := strings.Split(scanner.Text(), ":")
		if ss[0] == username {
			ss[1] = password
			s = strings.Join(ss, ":")
			changed = true
		}
		users = append(users, s)
	}
	if changed == true {
		file.Truncate(0)
		file.Seek(0,0)
		for _, v := range users {
			log.Println("write:", v)
			_, err := file.WriteString(v+"\n")
			if err != nil {
				log.Println(err)
			}
		}
	}
	return true, nil
}

