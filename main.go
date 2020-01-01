package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "1"
	DB_NAME     = "todolist"
)

var (
	db       *sql.DB
	err      error
	UserData User
)

func root(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "root")
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		dir, err := os.Getwd()
		check(err)

		ViewPath := filepath.FromSlash("/src/github.com/yfedoruck/todolist/views/")
		t, err := template.ParseFiles(dir + ViewPath + "register.html")

		fmt.Println(dir + ViewPath + "register.html")
		check(err)

		_ = t.Execute(w, nil)
	} else {
		err = r.ParseForm()
		check(err)

		if len(r.Form["email"]) == 0 {
			panic("email not exists")
		}

		var lastInsertId int
		err = db.QueryRow("INSERT into account (email,password,username) VALUES ($1,$2,$3) returning id;", r.Form["email"][0], r.Form["password"][0], r.Form["username"][0]).Scan(&lastInsertId)
		check(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func todoList(w http.ResponseWriter, r *http.Request) {
	list := todoListData(UserData.id)
	fmt.Println(list)

	data := TodoListData{
		"/css/signin.css",
		"Todo list",
		UserData.id,
		list,
	}

	renderTemplate(w, "todolist", data)
}

func renderTemplate(w http.ResponseWriter, tpl string, data interface{}) {
	dir, err := os.Getwd()
	check(err)

	ViewPath := filepath.FromSlash("/src/github.com/yfedoruck/todolist/views/")

	t, err := template.ParseFiles(dir + ViewPath + tpl + ".html")
	check(err)

	_ = t.Execute(w, data)
}

type LoginData struct {
	Css   string
	Title string
}

type TodoListData struct {
	Css      string
	Title    string
	UserId   int
	TodoList []Todo
}
type Todo struct {
	Id     int
	Todo   string
	Status bool
}
type User struct {
	id int
}

func login(w http.ResponseWriter, r *http.Request) {

	data := LoginData{
		"/css/signin.css",
		"Sign in",
	}

	if r.Method == "GET" {
		renderTemplate(w, "login", data)
	} else {

		err = r.ParseForm()
		check(err)

		fmt.Println("username = " + r.PostFormValue("username"))

		if len(r.Form["username"]) == 0 {
			panic("username not exists")
		}
		if len(r.Form["password"]) == 0 {
			panic("password not exists")
		}

		rows, err := db.Query("SELECT id, email FROM account WHERE username = $1 and password=$2 limit 1;", r.Form["username"][0], r.Form["password"][0])
		check(err)

		if rows.Next() == false {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		} else {
			var id int
			var email string
			err = rows.Scan(&id, &email)
			check(err)
			fmt.Println("id = " + string(id))
			fmt.Println("email = " + email)
			UserData = User{
				id,
			}

			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
		//var email string
		//err = rows.Scan(&email)
		//check(err)
		//fmt.Println("email = " + email)

	}
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err = r.ParseForm()
		check(err)

		if len(r.Form["note"]) == 0 {
			panic("note not exists")
		}

		if UserData.id == 0 {
			panic("user_id = 0")
		}

		var lastInsertId int
		err = db.QueryRow("INSERT into public.todo_list (user_id,todo,status) VALUES ($1,$2,$3) returning id;", UserData.id, r.PostFormValue("note"), true).Scan(&lastInsertId)
		check(err)

		http.Redirect(w, r, "/todolist", http.StatusFound)
	}
}

func todoListData(userId int) []Todo {
	rows, err := db.Query("SELECT id, todo, status FROM  public.todo_list where user_id = $1", userId)
	check(err)

	var id int
	var todo string
	var status bool
	var list []Todo
	for rows.Next() {
		err = rows.Scan(&id, &todo, &status)
		check(err)
		td := Todo{id, todo, status}
		fmt.Println(userId, id, todo, status)
		list = append(list, td)
	}

	return list
}

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err = sql.Open("postgres", dbinfo)
	check(err)
	defer db.Close()

	fs := http.FileServer(http.Dir("./src/github.com/yfedoruck/todolist/static/css"))
	//http.Handle("/src/todolist/static/css/signin.css", fs)
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	//var h Hello
	//http.HandleFunc("/", root)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/todolist", todoList)
	http.HandleFunc("/add", addTodo)
	http.HandleFunc("/remove", removeTodo)

	err := http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func removeTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err = r.ParseForm()
		check(err)

		if len(r.Form["id"]) == 0 {
			panic("id not exists")
		}
		if UserData.id == 0 {
			panic("user_id = 0")
		}

		stmt, err := db.Prepare("Delete from todo_list where id=$1")
		check(err)
		_, err = stmt.Exec(r.PostFormValue("id"))
		check(err)
		http.Redirect(w, r, "/todolist", http.StatusFound)
	}
}
