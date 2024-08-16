package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	ID    int
	Login string
	Mail  string
	Token string
}

// FromJSON decodes a serialized JSON record - User{}
func (m *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(m)
}

// ToJSON encodes a User JSON record
func (m *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

// PostgreSQL connection details
var (
	Hostname = "localhost"
	Port     = 5432
	Username = "postgres"
	Password = "postgres"
	Database = "donation_rest"
)

func ConnectPostgres() *sql.DB {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Println(err)
		return nil
	}

	return db
}

// InsertUser is for adding a new user to the database
func InsertUser(u User) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	if IsUserValid(u) {
		log.Println("User", u.Login, "already exist!")
		return false
	}

	stmt, err := db.Prepare("INSERT INTO user_registed(Nickname, Mail, Token) values($1, $2, $3)")
	if err != nil {
		log.Println("AddUser:", err)
		return false
	}

	stmt.Exec(u.Login, u.Mail, u.Token)
	return true
}

// ListAllMessages if for returning all messages from the database table
func ListAllMessages() []User {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return []User{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user_registed \n")
	if err != nil {
		log.Println(err)
		return []User{}
	}

	all := []User{}
	var c1 int
	var c2, c3, c4 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4)
		temp := User{c1, c2, c3, c4}
		all = append(all, temp)
	}

	return all
}

// Same as on top, returns user record by name
func FindUserNickname(nickname, password string) User {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return User{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user_registed WHERE Nickname, Password = $1, $2\n", nickname, password)
	if err != nil {
		log.Println("Query:", err)
		return User{}
	}
	defer rows.Close()

	u := User{}
	var c1 int
	var c2, c3, c4 string

	for rows.Next() {
		err := rows.Scan(&c1, &c2, &c3, &c4)
		if err != nil {
			log.Println(err)
			return User{}
		}
		u = User{c1, c2, c3, c4}
	}

	return u
}

func IsUserValid(u User) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user_registed WHERE Nickname = $1 \n", u.Login)
	if err != nil {
		log.Println(err)
		return false
	}

	temp := User{}
	var c1 int
	var c2, c3, c4 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4)
		if err != nil {
			log.Println(err)
			return false
		}
		temp = User{c1, c2, c3, c4}
	}
	if u.Login == temp.Login {
		return true
	}
	return false
}