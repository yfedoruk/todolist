package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/yfedoruck/todolist/lang"
	"github.com/yfedoruck/todolist/pg"
	"github.com/yfedoruck/todolist/validate"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	DbUser     = "postgres"
	DbPassword = "1"
	DbName     = "todolist"
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
	TodoList []pg.Todo
	Error    string
}

type User struct {
	id int
}

func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func register(regData *RegisterData, db *sql.DB, user *User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			dir, err := os.Getwd()
			check(err)

			ViewPath := filepath.FromSlash("/src/github.com/yfedoruck/todolist/views/")
			t, err := template.ParseFiles(dir + ViewPath + "register.html")
			check(err)

			err = t.Execute(w, regData)
			check(err)
		} else {
			var validation = true

			err := r.ParseForm()
			check(err)
			username := r.PostFormValue("username")
			password := r.PostFormValue("password")
			email := r.PostFormValue("email")

			errUsername := validate.Username(r)
			if errUsername != nil {
				regData.Error.Username = errUsername.Error()
				validation = false
			} else {
				regData.Error.Username = ""
			}

			errEmail := validate.Email(r)
			if errEmail != nil {
				regData.Error.Email = errEmail.Error()
				validation = false
			} else {
				regData.Error.Email = ""
			}

			if !validation {
				regData.PreFill = RegisterField{
					email,
					username,
					password,
				}

				http.Redirect(w, r, "/register", http.StatusFound)
				return
			}

			isUniqueUsername := isUniqueUsername(r, regData, db)
			isUniqueEmail := isUniqueEmail(r, regData, db)

			if !isUniqueUsername || !isUniqueEmail {
				regData.PreFill = RegisterField{
					email,
					username,
					password,
				}
				http.Redirect(w, r, "/register", http.StatusFound)
				return
			} else {
				regData.PreFill = RegisterField{}
				regData.Error = RegisterErr{}
			}

			user.id = pg.RegisterUser(db, email, password, username)
			clearRegisterForm(regData)
			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
	})
}

func clearRegisterForm(data *RegisterData) {
	data.Error = RegisterErr{}
	data.PreFill = RegisterField{}
}

func clearLoginForm(ld *LoginData) {
	ld.Error = ""
	ld.PreFill = LoginField{}
}

func isUniqueUsername(r *http.Request, data *RegisterData, db *sql.DB) bool {

	username := r.PostFormValue("username")

	ok := pg.IsUniqueUsername(username, db)

	if ok {
		data.Error.Username = ""
	} else {
		data.Error.Username = lang.UsernameExists
	}

	return ok
}

func isUniqueEmail(r *http.Request, data *RegisterData, db *sql.DB) bool {

	ok := pg.IsUniqueEmail(r.PostFormValue("email"), db)

	if ok {
		data.Error.Email = ""
	} else {
		data.Error.Email = lang.EmailExists
	}

	return ok
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func todoListHandler(data *NotesListData, db *sql.DB, user *User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user.id == 0 {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			data.UserId = user.id
			data.TodoList = pg.TodoListData(user.id, db)
		}

		renderTemplate(w, "todolist", data)
	})
}

func renderTemplate(w http.ResponseWriter, tpl string, data interface{}) {
	dir, err := os.Getwd()
	check(err)

	ViewPath := filepath.FromSlash("/src/github.com/yfedoruck/todolist/views/")

	t, err := template.ParseFiles(dir + ViewPath + tpl + ".html")
	check(err)

	_ = t.Execute(w, data)
}

func loginHandler(db *sql.DB, loginData *LoginData, user *User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user.id != 0 {
			http.Redirect(w, r, "/todolist", http.StatusFound)
		}

		if r.Method == "GET" {
			renderTemplate(w, "login", loginData)
		} else {
			err := r.ParseForm()
			check(err)
			username := r.PostFormValue("username")
			password := r.PostFormValue("password")

			if username == "" {
				err = errors.New(lang.NameEmpty)
			}

			if password == "" {
				err = errors.New(lang.PassEmpty)
			}

			id, err := pg.LoginUser(db, username, password)
			if err != nil {
				loginData.Error = err.Error()
				loginData.PreFill = LoginField{
					Username: username,
					Password: password,
				}
				http.Redirect(w, r, "/login", http.StatusFound)
			} else {
				user.id = id
				clearLoginForm(loginData)
				http.Redirect(w, r, "/todolist", http.StatusFound)
			}
		}
	})
}

func addNoteHandler(notes *NotesListData, db *sql.DB, user *User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user.id == 0 {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			check(err)

			if len(r.Form["note"]) == 0 {
				panic("note not exists")
			}

			err = validate.Note(r)
			if err != nil {
				notes.Error = err.Error()
				http.Redirect(w, r, "/todolist", http.StatusFound)
				return
			} else {
				notes.Error = ""
			}

			pg.AddNote(db, user.id, r.PostFormValue("note"))

			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
	})
}

func closeDb(db *sql.DB) {
	err := db.Close()
	check(err)
}

func logoutHandler(ld *LoginData, listData *NotesListData, regData *RegisterData, user *User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user.id = 0
		ld.Error = ""
		listData.Error = ""
		regData.Error = RegisterErr{}
		regData.PreFill = RegisterField{}
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func removeTodoHandler(db *sql.DB, user *User) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			check(err)

			if len(r.Form["id"]) == 0 {
				panic("id not exists")
			}
			if user.id == 0 {
				panic("user_id = 0")
			}

			ok, err := strconv.Atoi(r.PostFormValue("id"))
			check(err)

			pg.RemoveNote(ok, db)
			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
	})
}

func main() {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DbUser, DbPassword, DbName)
	db, err := sql.Open("postgres", dbInfo)
	check(err)
	defer closeDb(db)

	tables(db)

	var loginData = &LoginData{
		"/css/signin.css",
		"Sign in",
		"",
		LoginField{},
	}
	var notesListData = &NotesListData{
		"/css/signin.css",
		"Todo list",
		0,
		nil,
		"",
	}
	var registerData = &RegisterData{
		"/css/signin.css",
		"Sign in",
		RegisterErr{},
		RegisterField{},
	}
	user := &User{}

	fs := http.FileServer(http.Dir("./src/github.com/yfedoruck/todolist/static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	http.HandleFunc("/", root)
	http.Handle("/register", register(registerData, db, user))
	http.Handle("/login", loginHandler(db, loginData, user))
	http.Handle("/todolist", todoListHandler(notesListData, db, user))
	http.Handle("/add", addNoteHandler(notesListData, db, user))
	http.Handle("/remove", removeTodoHandler(db, user))
	http.Handle("/logout", logoutHandler(loginData, notesListData, registerData, user))

	err = http.ListenAndServe(":"+port(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
