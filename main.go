package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Database []struct {
		Tag      string `json:"tag"`
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
}

type person struct {
	Email                string
	Username             string
	Password_reset_token string
	Link                 string
}

func pswreset(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("templates/pswreset.html"))
	tpl.Execute(w, nil)

}

func psw_reset(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var username = r.Form["username"]
	var email = r.Form["email"]
	if person, person_err := check_user(username[0], email[0]); person_err != nil {
		var tpl = template.Must(template.ParseFiles("templates/error.html"))
		tpl.Execute(w, person_err)

	} else {
		var tpl = template.Must(template.ParseFiles("templates/index.html"))
		fmt.Printf("person %v", person)
		tpl.Execute(w, person)

	}
}

func check_user(username string, email string) (person, error) {
	//opening the config files
	if file, err := os.ReadFile("configs.json"); err != nil {
		panic(err.Error())
	} else {
		config := &Config{}
		err := json.Unmarshal(file, config)
		if err != nil {
			log.Panic(err)
		}

		// Connect to the database
		dbconfig := config.Database[0]
		db, err := sql.Open("mysql", dbconfig.Username+":"+dbconfig.Password+"@tcp("+dbconfig.Host+")/"+dbconfig.Name)
		if err != nil {
			return person{}, err
		}

		// Test the connection
		err = db.Ping()
		if err != nil {
			panic(err.Error())
		}
		var query string
		query = fmt.Sprintf("SELECT username,email,password_reset_token FROM user where username='%s' AND email='%s'", (username), (email))
		rows, err := db.Query(query)
		if err != nil && err != sql.ErrNoRows {
			return person{}, err
		}

		p := person{}
		for rows.Next() {

			query_err := rows.Scan(&p.Username, &p.Email, &p.Password_reset_token)

			if query_err != nil {
				return person{}, query_err
			}
		}

		p.Link = fmt.Sprint("https://portal.dev01.int.betika.com/site/reset-password?token=", p.Password_reset_token)
		return p, nil
	}
}

func main() {

	http.HandleFunc("/", pswreset)

	http.HandleFunc("/psw_reset", psw_reset)

	http.ListenAndServe(":8000", nil)

}
