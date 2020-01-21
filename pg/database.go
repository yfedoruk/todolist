package pg

import (
	"database/sql"
	"errors"
	"github.com/yfedoruck/todolist/lang"
	"log"
)

type Todo struct {
	Id     int
	Todo   string
	Status bool
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterUser(db *sql.DB, username string, password string, email string) int {

	var lastInsertId int
	dbErr := db.QueryRow("INSERT into account (email,password,username) VALUES ($1,$2,$3) returning id;", email, password, username).Scan(&lastInsertId)
	check(dbErr)

	return lastInsertId
}

func LoginUser(db *sql.DB, username string, password string) (int, error) {

	rows, err := db.Query("SELECT id, email FROM account WHERE username = $1 and password=$2 limit 1;", username, password)
	check(err)

	if rows.Next() == false {
		return 0, errors.New(lang.LoginErr)
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

func AddNote(db *sql.DB, userId int, note string) {

	var lastInsertId int
	err := db.QueryRow("INSERT into public.todo_list (user_id,todo,status) VALUES ($1,$2,$3) returning id;", userId, note, true).Scan(&lastInsertId)
	check(err)
}

func RemoveNote(id int, db *sql.DB) {
	stmt, err := db.Prepare("Delete from todo_list where id=$1")
	check(err)
	_, err = stmt.Exec(id)
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
