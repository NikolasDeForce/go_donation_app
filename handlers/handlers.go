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

	"github.com/gorilla/mux"
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
	Body := "ERROR: METHOD NOT ALLOWED!\n"
	fmt.Fprintf(w, "%s", Body)
}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Main Handler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	w.WriteHeader(http.StatusOK)

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

func DonationHanler(w http.ResponseWriter, r *http.Request) {
	log.Println("Donation Handler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)
	w.WriteHeader(http.StatusOK)

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

	login := db.FindUserNickname(d.LoginStrimer)

	if login.Login == d.LoginStrimer {
		db.InsertDonate(db.Donate{
			LoginStrimer: d.LoginStrimer,
			NameSub:      d.NameSub,
			Value:        d.Value,
			Text:         d.Text,
		})
	}
}

// API handlers
func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Donation Handler Serving:", r.URL.Path, "from", r.Host, "with method", r.Method)

	login, ok := mux.Vars(r)["login"]
	if !ok {
		log.Println("login value not set!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	password, ok := mux.Vars(r)["password"]
	if !ok {
		log.Println("login value not set!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	l := db.FindUserNickname(login)

	if login == l.Login && password == l.Password {
		fmt.Fprintf(w, "Привет %v\nТвой токен - %v\nСохраните его в безопасное место\nВ случае утери обращайтесь в тех.поддержку", l.Login, l.Token)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "ERROR STATUS NOT FOUND", http.StatusNotFound)
	}
}

func GetDonatesHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDonatesHandler Serving:", r.URL.Path, "from", r.Host)

	token, ok := mux.Vars(r)["token"]
	if !ok {
		log.Println("token value not set!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t := db.FindUserToken(token)
	if t.Token == token {
		err := db.ListAllDonates(t.Login)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}
		for _, v := range err {
			fmt.Fprintf(w, "ID Доната: %v\nЛогин зрителя: %v\nЛогин стримера: %v\nСумма доната: %v\nСообщение: %v\n\n", v.ID, v.NameSub, v.LoginStrimer, v.Value, v.Text)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "ERROR STATUS NOT FOUND", http.StatusNotFound)
		log.Println("Token not found:", token)
	}
}

// SliceToJSON encodes a slice with JSON records
func SliceToJSON(slice interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(slice)
}
