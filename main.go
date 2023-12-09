package main

import (
	// "fmt"
	"html/template"
	"net/http"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type ToDo struct {
	Title string
	Description string
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))

	todoArr := []ToDo{}


	rows, err := db.Query("select id, title, description, completed from todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var id int
		var title, description, completed string
		err = rows.Scan(&id, &title, &description, &completed)
		if err != nil {
			log.Fatal(err)
		}

		todoArr = append(todoArr, ToDo{Title: title, Description: description})
	}


	todos := map[string][]ToDo {
		"ToDos": todoArr,
	}
	
	tmpl.Execute(w, todos)
}

func addToDo(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	description := r.PostFormValue("description")

	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("insert into todos(title, description, completed) values(?, ?, 'N')")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(title, description)
	checkErr(err)

	err = tx.Commit()
	checkErr(err)
	
	tmpl := template.Must(template.ParseFiles(("index.html")))
	tmpl.ExecuteTemplate(w, "todos-item", ToDo{Title: title, Description: description})
}

func main() {
	var err error // error is scoped locally so that "=" can be used in the following line instead of ":=" which would overrride global "db"
	db, err = sql.Open("sqlite3", "./db.sqlite")
	checkErr(err)
	defer db.Close()

	dbInit := `
		create table if not exists todos (
			id integer not null primary key,
			title text not null,
			description text,
			completed text not null
		);
	`
	_, err = db.Exec(dbInit)
	checkErr(err)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/add", addToDo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkErr(err error) {
	if (err != nil) {
		panic(err)
	}
}