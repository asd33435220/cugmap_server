package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Db *sql.DB
var err error

const (
	insertUserSrt       = "insert into cug_map_users_tpl(student_id,username,password) values(?,?,?);"
	queryStudentIdSrt   = "select student_id from cug_map_users_tpl where student_id = ?;"
	queryStudentSrt     = "select username,password from cug_map_users_tpl where student_id = ?;"
	queryStudentInfoSrt = "select student_id,username,password from cug_map_users_tpl where student_id = ?;"
	updateUserPosition  = "update cug_map_users_tpl set position = ? where student_id = ?;"
)

type User struct {
	StudentId string `form:"student_id" json:"student_id" binding:"required"`
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding"required"`
	Position  string `json:"position"`
}

func init() {
	Db, err = sql.Open("mysql", "root:asd33435220@tcp/cug_map_db")
	if err != nil {
		fmt.Println(err.Error())
	}
}
func (newUser *User) UpdateUser() (err error) {
	return nil
}
func (newUser *User) QueryUser() (id string) {
	stmt, err := Db.Prepare(queryStudentIdSrt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	result := stmt.QueryRow(newUser.StudentId)
	result.Scan(&id)
	return
}
func (newUser *User) QueryUserInfo() *User {
	stmt, err := Db.Prepare(queryStudentInfoSrt)
	user := &User{}
	if err != nil {
		log.Fatal(err)
		return user
	}
	defer stmt.Close()
	result := stmt.QueryRow(newUser.StudentId)
	result.Scan(&user.StudentId, &user.Username, &user.Password)
	return user
}
func (newUser *User) AddUser() (err error) {
	stmt, err := Db.Prepare(insertUserSrt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(newUser.StudentId, newUser.Username, newUser.Password)
	if err != nil {
		log.Fatal(err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("id=", id)
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("rows=", rows)
	return nil
}
func (newUser *User) GetUser() (username string, err error) {
	stmt, err := Db.Prepare(queryStudentSrt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	var password string
	result := stmt.QueryRow(newUser.StudentId)
	result.Scan(&username, &password)
	if password == newUser.Password {
		return
	} else {
		username = ""
		return
	}
}
func (newUser *User) UpdateUserPosition()(err error) {
	err = nil
	stmt, err := Db.Prepare(updateUserPosition)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(newUser.Position, newUser.StudentId)
	if err != nil {
		log.Fatal(err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("id=", id)
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("rows=", rows)
	return
}
