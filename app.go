package main

import (
    "fmt"
	"net/http"
	"html/template"
	"encoding/json"
	"log"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
)


type Status struct {
    Status bool
}


type Todos struct{
	Id int
	Todo, Status string
}

func dbConn() (db *sql.DB) {
    db, err := sql.Open("mysql", "user:pass(127.0.0.1:3306)/todo_go")
    if err != nil {
        panic(err.Error())
    }
    return db
}



func index(w http.ResponseWriter, r *http.Request) {
    data := ""
	
	t, err := template.ParseFiles("templates/index.gohtml") 
    if err != nil { // if there is an error
  	  log.Print("template parsing error: ", err) // log it
  	}
    err = t.Execute(w, data) 
    if err != nil { // if there is an error
  	  log.Print("template executing error: ", err) //log it
  	}
}

func get(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	selDB, err := db.Query("SELECT id, todo, status FROM todos ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
	}
	todo := Todos{}
	todos:= []Todos{}
	var id int
	var  _todo, status string
	
    for selDB.Next() {
        err = selDB.Scan(&id, &_todo, &status)
        if err != nil {
            panic(err.Error())
		}
		todo.Id 	= id
        todo.Todo 	= _todo
        todo.Status = status
        todos = append(todos, todo)
	}
	
	result, err := json.Marshal(todos)
    w.Write(result)
}

func add(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w)
	fmt.Println("*************")
	fmt.Println(r)
	fmt.Println("*************")
	fmt.Println(r.Form)
	fmt.Println("*************")

	var value string = r.FormValue("todo")

	db := dbConn()
	insert, err := db.Query("INSERT INTO todos(todo) VALUES ('"+value+"')")
	var statu Status
    if err != nil {
		statu = Status{false}
    }else{
		statu = Status{true}
	}
	defer insert.Close()
	
	result, err := json.Marshal(statu)
    w.Write(result)
}

func delete(w http.ResponseWriter, r *http.Request) {
	var id string = r.FormValue("id")

	db := dbConn()
	delForm, err := db.Prepare("DELETE FROM todos WHERE id=?")
	var statu Status
    if err != nil {
		statu = Status{false}
    }else{
		delForm.Exec(id)
		statu = Status{true}
	}
	
	result, err := json.Marshal(statu)
    w.Write(result)
}

func change(w http.ResponseWriter, r *http.Request){

	db := dbConn()
	status 	:= r.FormValue("status")
	id 		:= r.FormValue("id")
	insForm, err := db.Prepare("UPDATE todos SET status=? WHERE id=?")
	var statu Status
	if err != nil {
		statu = Status{false}
	}else{
		statu = Status{true}
		insForm.Exec(status, id)
	}
	result, err := json.Marshal(statu)
    w.Write(result)
}

func main() {
	fmt.Println("Server starting....")
	fmt.Println("http://localhost:8080")

	fs := http.FileServer(http.Dir("./static"))
  	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/add-todo", add)
	http.HandleFunc("/get-todo", get)
	http.HandleFunc("/delete-todo", delete)
	http.HandleFunc("/change-status-todo", change)
	http.ListenAndServe(":8080", nil)
}