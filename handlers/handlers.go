package handlers

import (
	"donation/db"
	generatetoken "donation/generateToken"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

type NotAllowedHandler struct{}

var tmpl *template.Template

var u db.User

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
		http.Redirect(w, r, "/error", http.StatusMethodNotAllowed)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		fmt.Println(err)
	}

	r.ParseForm()

	u.Login = r.FormValue("userNickname")
	u.Mail = r.FormValue("userEmail")
	u.Password = r.FormValue("userPassword")
	u.Token = generatetoken.GenerateToken()

	db.InsertUser(db.User{
		Login:    u.Login,
		Mail:     u.Mail,
		Password: u.Password,
		Token:    u.Token,
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Donation Handler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	w.WriteHeader(http.StatusOK)

	nickname := u.Login
	password := u.Password
	user := db.FindUserNicknameAndPassword(nickname, password)

	err := tmpl.ExecuteTemplate(w, "register.html", user)
	if err != nil {
		fmt.Println(err)
	}
}

func DonationHanler(w http.ResponseWriter, r *http.Request) {
	log.Println("Donation Handler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)

	if r.URL.Path != "/donation" {
		http.Redirect(w, r, "/error", http.StatusMethodNotAllowed)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	err := tmpl.ExecuteTemplate(w, "donation.html", nil)
	if err != nil {
		fmt.Println(err)
	}

	d := db.Donate{}

	r.ParseForm()

	d.LoginStrimer = r.FormValue("loginStrimer")
	d.NameSub = r.FormValue("userNickname")
	d.Value, _ = strconv.Atoi(r.FormValue("Value"))
	d.Text = r.FormValue("Text")

	db.InsertDonate(db.Donate{
		LoginStrimer: d.LoginStrimer,
		NameSub:      d.NameSub,
		Value:        d.Value,
		Text:         d.Text,
	})
}

// SliceToJSON encodes a slice with JSON records
func SliceToJSON(slice interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(slice)
}
