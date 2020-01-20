package pg

import (
	"database/sql"
	"errors"
	"net/http"
)

type Todo struct {
	Id     int
	Todo   string
	Status bool
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func RegisterUser(r *http.Request, db *sql.DB) int {
	err := r.ParseForm()
	check(err)

	var lastInsertId int
	dbErr := db.QueryRow("INSERT into account (email,password,username) VALUES ($1,$2,$3) returning id;", r.PostFormValue("email"), r.PostFormValue("password"), r.PostFormValue("username")).Scan(&lastInsertId)
	check(dbErr)

	return lastInsertId
}

func LoginUser(db *sql.DB, r *http.Request) (int, error) {
	const (
		loginErr = "wrong username or password"
		nameErr  = "username not exists"
		passErr  = "password not exists"
	)
	err := r.ParseForm()
	check(err)

	if r.PostFormValue("username") == "" {
		return 0, errors.New(nameErr)
	}

	if r.PostFormValue("password") == "" {
		return 0, errors.New(passErr)
	}

	rows, err := db.Query("SELECT id, email FROM account WHERE username = $1 and password=$2 limit 1;", r.PostFormValue("username"), r.PostFormValue("password"))
	check(err)

	if rows.Next() == false {
		return 0, errors.New(loginErr)
	} else {
		var id int
		var email string
		err = rows.Scan(&id, &email)
		check(err)
		return id, nil
	}
}

func TodoListData(userId int, db *sql.DB) []Todo {
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
		list = append(list, td)
	}

	return list
}

func AddNote(r *http.Request, db *sql.DB, userId int) {
	err := r.ParseForm()
	check(err)

	var lastInsertId int
	err = db.QueryRow("INSERT into public.todo_list (user_id,todo,status) VALUES ($1,$2,$3) returning id;", userId, r.PostFormValue("note"), true).Scan(&lastInsertId)
	check(err)
}

func RemoveNote(r *http.Request, db *sql.DB) {
	err := r.ParseForm()
	check(err)

	stmt, err := db.Prepare("Delete from todo_list where id=$1")
	check(err)
	_, err = stmt.Exec(r.PostFormValue("id"))
	check(err)
}

func IsUniqueUsername(username string, db *sql.DB) bool {

	//strings.Contains(err, "account_username_key")
	rows, err := db.Query("SELECT id FROM account WHERE username = $1 limit 1;", username)
	check(err)

	return rows.Next() == false
}

func IsUniqueEmail(email string, db *sql.DB) bool {

	rows, err := db.Query("SELECT id FROM account WHERE email = $1 limit 1;", email)
	check(err)

	return rows.Next() == false
}
