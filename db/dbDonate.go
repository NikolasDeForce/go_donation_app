package db

import (
	"log"

	_ "github.com/lib/pq"
)

type Donate struct {
	ID           int
	LoginStrimer string
	NameSub      string
	Value        int
	Text         string
}

// InsertDonate is for adding a new donate to the database
func InsertDonate(d Donate) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	if IsDonateValid(d) {
		log.Println("Donate", d.Text, "already exist!")
		return false
	}

	stmt, err := db.Prepare("INSERT INTO user_registed(LoginStrimer, NameSub, Value, Text) values($1, $2, $3, $4)")
	if err != nil {
		log.Println("AddDonate:", err)
		return false
	}

	stmt.Exec(d.LoginStrimer, d.NameSub, d.Value, d.Text)
	return true
}

// ListAllMessages if for returning all messages from the database table
func ListAllDonates() []Donate {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return []Donate{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM donation_list \n")
	if err != nil {
		log.Println(err)
		return []Donate{}
	}

	all := []Donate{}
	var c1, c4 int
	var c2, c3, c5 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5)
		temp := Donate{c1, c2, c3, c4, c5}
		all = append(all, temp)
	}

	return all
}

func IsDonateValid(d Donate) bool {
	db := ConnectPostgres()
	if db == nil {
		log.Println("Cannot connect to PostreSQL!")
		db.Close()
		return false
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM donation_list WHERE LoginStrimer, NameSub, Value, Text = $1, $2, $3, $4 \n", d.LoginStrimer, d.NameSub, d.Value, d.Text)
	if err != nil {
		log.Println(err)
		return false
	}

	temp := Donate{}
	var c1, c4 int
	var c2, c3, c5 string

	for rows.Next() {
		err = rows.Scan(&c1, &c2, &c3, &c4, &c5)
		if err != nil {
			log.Println(err)
			return false
		}
		temp = Donate{c1, c2, c3, c4, c5}
	}
	if d.LoginStrimer == temp.LoginStrimer && d.NameSub == temp.NameSub && d.Value == temp.Value && d.Text == temp.Text {
		return true
	}
	return false
}
