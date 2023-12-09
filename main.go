package main

import (
	// "fmt"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type ToDo struct {
	Id int64
	Title string
	Description string
	Completed bool
}

func getToDos() map[string][]ToDo{
	todoArr := []ToDo{}


	rows, err := db.Query("select id, title, description, completed from todos order by completed, dtm, id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var id int64
		var title, description string
		var completed bool
		err = rows.Scan(&id, &title, &description, &completed)
		if err != nil {
			log.Fatal(err)
		}

		todoArr = append(todoArr, ToDo{Id: id, Title: title, Description: description, Completed: completed})
	}


	todos := map[string][]ToDo {
		"ToDos": todoArr,
	}
	return todos
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	todos := getToDos()
	tmpl.Execute(w, todos)
}

func toggleToDo(w http.ResponseWriter, r *http.Request) {
	strId := strings.TrimPrefix(r.URL.Path, "/toggle/")
	id, err := strconv.Atoi(strId)
	checkErr(err);

	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("update todos set completed=(case when completed = true then false else true end), dtm=(case when completed = true then null else datetime() end) where id=?")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkErr(err)

	err = tx.Commit()
	checkErr(err)

	todos := getToDos()
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.ExecuteTemplate(w, "todos-list", todos)
}

func addToDo(w http.ResponseWriter, r *http.Request) {
	title := r.PostFormValue("title")
	description := r.PostFormValue("description")

	tx, err := db.Begin()
	checkErr(err)
	stmt, err := tx.Prepare("insert into todos(title, description, completed, dtm) values(?, ?, FALSE, NULL)")
	checkErr(err)
	defer stmt.Close()

	result, err := stmt.Exec(title, description)
	checkErr(err)
	
	id, err := result.LastInsertId()
	checkErr(err)

	err = tx.Commit()
	checkErr(err)
	
	tmpl := template.Must(template.ParseFiles(("index.html")))
	tmpl.ExecuteTemplate(w, "todos-item", ToDo{Id: id, Title: title, Description: description, Completed: false})
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
			completed boolean not null,
			dtm date
		);
	`
	_, err = db.Exec(dbInit)
	checkErr(err)

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/add", addToDo)
	http.HandleFunc("/toggle/", toggleToDo)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func checkErr(err error) {
	if (err != nil) {
		panic(err)
	}
}