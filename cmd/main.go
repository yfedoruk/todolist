package main

import (
	"errors"
	_ "github.com/lib/pq"
	"github.com/yfedoruck/todolist/pkg/lang"
	"github.com/yfedoruck/todolist/pkg/validate"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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

func root(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func register(regData *RegisterData, db postgres) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {
			viewDir := BasePath() + filepath.FromSlash("/views/")
			t, err := template.ParseFiles(viewDir + "register.html")
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

			id := db.RegisterUser(password, username, email)

			cookie{
				Name: username,
				Id:   id,
			}.set(w)

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

func isUniqueUsername(r *http.Request, data *RegisterData, db postgres) bool {

	username := r.PostFormValue("username")

	ok := db.IsUniqueUsername(username)

	if ok {
		data.Error.Username = ""
	} else {
		data.Error.Username = lang.UsernameExists
	}

	return ok
}

func isUniqueEmail(r *http.Request, data *RegisterData, db postgres) bool {

	ok := db.IsUniqueEmail(r.PostFormValue("email"))

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

func todoListHandler(data *NotesListData, db postgres) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, err := r.Cookie("auth")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user := cookie{}
		user.decode(auth.Value)
		data.UserId = user.Id
		data.TodoList = db.TodoListData(user.Id)
		renderTemplate(w, "todolist", data)
	})
}

func renderTemplate(w http.ResponseWriter, tpl string, data interface{}) {
	viewDir := BasePath() + filepath.FromSlash("/views/")

	t, err := template.ParseFiles(viewDir + tpl + ".html")
	check(err)

	_ = t.Execute(w, data)
}

func loginHandler(db postgres, loginData *LoginData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("auth"); err == nil {
			http.Redirect(w, r, "/todolist", http.StatusFound)
			return
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

			id, err := db.LoginUser(username, password)
			if err != nil {
				loginData.Error = err.Error()
				loginData.PreFill = LoginField{
					Username: username,
					Password: password,
				}
				http.Redirect(w, r, "/login", http.StatusFound)
			} else {
				cookie{
					Name: username,
					Id:   id,
				}.set(w)
				clearLoginForm(loginData)
				http.Redirect(w, r, "/todolist", http.StatusFound)
			}
		}
	})
}

func addNoteHandler(notes *NotesListData, db postgres) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, err := r.Cookie("auth")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
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

			user := cookie{}
			user.decode(auth.Value)
			db.AddNote(user.Id, r.PostFormValue("note"))

			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
	})
}

func logoutHandler(ld *LoginData, listData *NotesListData, regData *RegisterData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ld.Error = ""
		listData.Error = ""
		regData.Error = RegisterErr{}
		regData.PreFill = RegisterField{}
		removeCookie(w)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}

func removeTodoHandler(db postgres) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			check(err)

			if len(r.Form["id"]) == 0 {
				panic("id not exists")
			}

			if _, err := r.Cookie("auth"); err != nil {
				log.Println("user id not found")
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			ok, err := strconv.Atoi(r.PostFormValue("id"))
			check(err)

			db.RemoveNote(ok)
			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
	})
}

func main() {
	log.Println("45454545545vvvvvv33333dcdcdcdcdc")
	db := postgres{}
	db.Connect()
	defer db.Close()
	db.Tables()

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

	cssDir := BasePath() + filepath.FromSlash("/static/css")
	fs := http.FileServer(http.Dir(cssDir))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	http.HandleFunc("/", root)
	http.Handle("/register", register(registerData, db))
	http.Handle("/login", loginHandler(db, loginData))
	http.Handle("/todolist", todoListHandler(notesListData, db))
	http.Handle("/add", addNoteHandler(notesListData, db))
	http.Handle("/remove", removeTodoHandler(db))
	http.Handle("/logout", logoutHandler(loginData, notesListData, registerData))

	err := http.ListenAndServe(":"+port(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
