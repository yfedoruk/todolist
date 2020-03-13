package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yfedoruck/todolist/lang"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type postgres struct {
	db *sql.DB
}

func (p *postgres) Connect() {
	dbConf := Config()
	dbInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbConf.Host, dbConf.Port, dbConf.User, dbConf.Password, dbConf.Name)
	fmt.Println(dbInfo)
	var err error
	p.db, err = sql.Open("postgres", dbInfo)
	check(err)

	for i, connected := 0, false; connected == false && i < 4; i++ {
		err = p.db.Ping()
		if err == nil {
			connected = true
			return
		} else {
			log.Println("Error: Could not establish a connection with the database!", err, " but I still tried to connect...")
			time.Sleep(2 * time.Second)
		}
	}
	panic(err)
}

func (p *postgres) Close() {
	err := p.db.Close()
	check(err)
}

func (p *postgres) Tables() {
	files, err := filepath.Glob(BasePath() + "/sql/*.sql")
	check(err)

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		check(err)

		stmt, err := p.db.Prepare(string(data))
		check(err)

		_, err = stmt.Exec()
		check(err)
	}
}
type Todo struct {
	Id     int
	Todo   string
	Status bool
}

func (p *postgres) RegisterUser(username string, password string, email string) int {

	var lastInsertId int
	dbErr := p.db.QueryRow("INSERT into account (email,password,username) VALUES ($1,$2,$3) returning id;", email, password, username).Scan(&lastInsertId)
	check(dbErr)

	return lastInsertId
}

func (p *postgres) LoginUser(username string, password string) (int, error) {

	rows, err := p.db.Query("SELECT id, email FROM account WHERE username = $1 and password=$2 limit 1;", username, password)
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

func (p postgres) TodoListData(userId int) []Todo {
	rows, err := p.db.Query("SELECT id, todo, status FROM  public.todo_list where user_id = $1 ORDER BY id DESC", userId)
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

func (p *postgres) AddNote(userId int, note string) {

	var lastInsertId int
	err := p.db.QueryRow("INSERT into public.todo_list (user_id,todo,status) VALUES ($1,$2,$3) returning id;", userId, note, true).Scan(&lastInsertId)
	check(err)
}

func (p *postgres) RemoveNote(id int) {
	stmt, err := p.db.Prepare("Delete from todo_list where id=$1")
	check(err)
	_, err = stmt.Exec(id)
	check(err)
}

func (p postgres) IsUniqueUsername(username string) bool {

	// strings.Contains(err, "account_username_key")
	rows, err := p.db.Query("SELECT id FROM account WHERE username = $1 limit 1;", username)
	check(err)

	return rows.Next() == false
}

func (p postgres) IsUniqueEmail(email string) bool {

	rows, err := p.db.Query("SELECT id FROM account WHERE email = $1 limit 1;", email)
	check(err)

	return rows.Next() == false
}


type Conf struct {
	User     string `json:"User"`
	Password string `json:"Password"`
	Name     string `json:"Name"`
	Host     string `json:"Host"`
	Port     string `json:"Port"`
}

func Config() Conf {
	file, err := os.Open(BasePath() + filepath.FromSlash("/config/"+Env()+"/postgres.json"))
	check(err)

	dbConf := Conf{}
	err = json.NewDecoder(file).Decode(&dbConf)
	check(err)

	return dbConf
}