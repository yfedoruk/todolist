package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/yfedoruck/todolist/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	DbUser     = "postgres"
	DbPassword = "1"
	DbName     = "todolist"

	UsernameExists = "username already exists"
	EmailExists    = "email already exists"
)

var (
	db       *sql.DB
	err      error
	UserData User
)

type LoginData struct {
	Css     string
	Title   string
	Error   string
	PreFill LoginField
}

type RegisterData struct {
	Css     string
	Title   string
	Error   RegisterErr
	PreFill RegisterField
}
type RegisterErr struct {
	Email    string
	Username string
	Password string
}
type RegisterField struct {
	Email    string
	Username string
	Password string
}
type LoginField struct {
	Username string
	Password string
}

type NotesListData struct {
	Css      string
	Title    string
	UserId   int
	TodoList []Todo
	Error    string
}
type Todo struct {
	Id     int
	Todo   string
	Status bool
}
type User struct {
	id int
}

var loginData = LoginData{
	"/css/signin.css",
	"Sign in",
	"",
	LoginField{},
}

var registerData = RegisterData{
	"/css/signin.css",
	"Sign in",
	RegisterErr{},
	RegisterField{},
}
var notesListData = NotesListData{
	"/css/signin.css",
	"Todo list",
	0,
	nil,
	"",
}

func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		dir, err := os.Getwd()
		check(err)

		ViewPath := filepath.FromSlash("/src/github.com/yfedoruck/todolist/views/")
		t, err := template.ParseFiles(dir + ViewPath + "register.html")
		check(err)

		_ = t.Execute(w, registerData)
	} else {
		var validation = true

		errUsername := middleware.Username(r)
		if errUsername != nil {
			registerData.Error.Username = errUsername.Error()
			validation = false
		} else {
			registerData.Error.Username = ""
		}

		errEmail := middleware.Email(r)
		if errEmail != nil {
			registerData.Error.Email = errEmail.Error()
			validation = false
		} else {
			registerData.Error.Email = ""
		}

		if !validation {
			registerData.PreFill = RegisterField{
				r.PostFormValue("email"),
				r.PostFormValue("username"),
				r.PostFormValue("password"),
			}

			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		isUniqueUsername := isUniqueUsername(r)
		isUniqueEmail := isUniqueEmail(r)

		if !isUniqueUsername || !isUniqueEmail {
			registerData.PreFill = RegisterField{
				r.PostFormValue("email"),
				r.PostFormValue("username"),
				r.PostFormValue("password"),
			}
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		} else {
			registerData.PreFill = RegisterField{}
			registerData.Error = RegisterErr{}
		}

		registerUser(r)
		http.Redirect(w, r, "/todolist", http.StatusFound)
	}
}

func registerUser(r *http.Request) {
	err = r.ParseForm()
	check(err)

	if len(r.Form["email"]) == 0 {
		panic("email not exists")
	}

	var lastInsertId int
	dbErr := db.QueryRow("INSERT into account (email,password,username) VALUES ($1,$2,$3) returning id;", r.PostFormValue("email"), r.PostFormValue("password"), r.PostFormValue("username")).Scan(&lastInsertId)
	check(dbErr)

	UserData = User{
		lastInsertId,
	}
	registerData.Error = RegisterErr{}
}

func isUniqueUsername(r *http.Request) bool {
	var result bool

	rows, err := db.Query("SELECT id FROM account WHERE username = $1 limit 1;", r.PostFormValue("username"))
	check(err)

	if rows.Next() {
		registerData.Error.Username = UsernameExists
		result = false
	} else {
		registerData.Error.Username = ""
		result = true
	}

	return result
}

func isUniqueEmail(r *http.Request) bool {
	var result bool

	//strings.Contains(err, "account_username_key")
	rows, err := db.Query("SELECT id FROM account WHERE email = $1 limit 1;", r.PostFormValue("email"))
	check(err)

	if rows.Next() {
		registerData.Error.Email = EmailExists
		result = false
	} else {
		registerData.Error.Email = ""
		result = true
	}

	return result
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func todoList(w http.ResponseWriter, r *http.Request) {
	if UserData.id == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		notesListData.UserId = UserData.id
		notesListData.TodoList = todoListData(UserData.id)
	}

	renderTemplate(w, "todolist", notesListData)
}

func renderTemplate(w http.ResponseWriter, tpl string, data interface{}) {
	dir, err := os.Getwd()
	check(err)

	ViewPath := filepath.FromSlash("/src/github.com/yfedoruck/todolist/views/")

	t, err := template.ParseFiles(dir + ViewPath + tpl + ".html")
	check(err)

	_ = t.Execute(w, data)
}

func login(w http.ResponseWriter, r *http.Request) {

	if UserData.id != 0 {
		http.Redirect(w, r, "/todolist", http.StatusFound)
	}

	if r.Method == "GET" {
		renderTemplate(w, "login", loginData)
	} else {

		err = r.ParseForm()
		check(err)

		fmt.Println("username = " + r.FormValue("username"))

		if len(r.Form["username"]) == 0 {
			panic("username not exists")
		}
		if len(r.Form["password"]) == 0 {
			panic("password not exists")
		}

		rows, err := db.Query("SELECT id, email FROM account WHERE username = $1 and password=$2 limit 1;", r.PostFormValue("username"), r.PostFormValue("password"))
		check(err)

		if rows.Next() == false {
			loginData.Error = "wrong username or password"
			loginData.PreFill = LoginField{
				Username: r.PostFormValue("username"),
				Password: r.PostFormValue("password"),
			}
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

func addNote(w http.ResponseWriter, r *http.Request) {
	if UserData.id == 0 {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	if r.Method == "POST" {
		err = r.ParseForm()
		check(err)

		if len(r.Form["note"]) == 0 {
			panic("note not exists")
		}

		err = middleware.Note(r)
		if err != nil {
			notesListData.Error = err.Error()
			http.Redirect(w, r, "/todolist", http.StatusFound)
			return
		} else {
			notesListData.Error = ""
		}

		var lastInsertId int
		err = db.QueryRow("INSERT into public.todo_list (user_id,todo,status) VALUES ($1,$2,$3) returning id;", UserData.id, r.PostFormValue("note"), true).Scan(&lastInsertId)
		check(err)

		http.Redirect(w, r, "/todolist", http.StatusFound)
	}
}

func todoListData(userId int) []Todo {
	rows, err := db.Query("SELECT id, todo, status FROM  public.todo_list where user_id = $1 ORDER BY id DESC", userId)
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

func closeDb() {
	err = db.Close()
	check(err)
}

func main() {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DbUser, DbPassword, DbName)
	db, err = sql.Open("postgres", dbInfo)
	check(err)
	defer closeDb()

	fs := http.FileServer(http.Dir("./src/github.com/yfedoruck/todolist/static/css"))
	//http.Handle("/src/todolist/static/css/signin.css", fs)
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	//var h Hello
	http.HandleFunc("/", root)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/todolist", todoList)
	http.HandleFunc("/add", addNote)
	http.HandleFunc("/remove", removeTodo)
	http.HandleFunc("/logout", logout)

	err := http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	UserData.id = 0
	loginData.Error = ""
	notesListData.Error = ""
	registerData.Error = RegisterErr{}
	registerData.PreFill = RegisterField{}
	http.Redirect(w, r, "/login", http.StatusFound)
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
