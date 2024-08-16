package handlers

import (
	generatetoken "donation/generateToken"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

type NotAllowedHandler struct{}

var tmpl *template.Template

// ServeHTTP implements http.Handler.
func (h NotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	MethodNotAllowedHandler(w, r)
}

// MethodNotAllowedHandler is executed when the HTTP method is incorrect
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host, "with method")
	Body := "Method not allowed!\n"
	fmt.Fprintf(w, "%s", Body)
}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Main Handler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	if r.URL.Path != "/" {
		http.Error(w, "Error: NOT FOUND", http.StatusNotFound)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		fmt.Println(err)
	}

	// u := db.User{}

	r.ParseForm()

	nickname := r.FormValue("userNickname")
	mail := r.FormValue("userEmail")
	token := generatetoken.GenerateToken()

	if nickname != "" && mail != "" {
		http.Redirect(w, r, "/register", http.StatusPermanentRedirect)
	}

	log.Println(nickname)
	log.Println(mail)
	log.Println(token)

	// db.InsertUser(db.User{
	// 	Nickname: u.Nickname,
	// 	Mail:     u.Mail,
	// 	Token:    u.Token,
	// })
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Register Handler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	w.WriteHeader(http.StatusOK)

	err := tmpl.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		fmt.Println(err)
	}

}

// SliceToJSON encodes a slice with JSON records
func SliceToJSON(slice interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(slice)
}
