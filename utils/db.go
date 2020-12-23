package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

var Db *sql.DB
var err error

const (
	insertUserSrt       = "insert into cug_map_users_tpl(student_id,username,password) values(?,?,?);"
	queryStudentIdSrt   = "select student_id from cug_map_users_tpl where student_id = ?;"
	queryStudentNameSrt   = "select username from cug_map_users_tpl where student_id = ?;"
	queryStudentSrt     = "select username,password from cug_map_users_tpl where student_id = ?;"
	queryStudentInfoSrt = "select student_id,username,password,position,signature from cug_map_users_tpl where student_id = ?;"
	updateUserPosition  = "update cug_map_users_tpl set position = ? where student_id = ?;"
	updateUserSignature = "update cug_map_users_tpl set signature = ? where student_id = ?;"
	queryUserPosition = "select position from cug_map_users_tpl where student_id = ?;"
	queryAllUserInfo  = "select student_id,username,position,signature from cug_map_users_tpl;"
	queryAllInfo  = "select student_id,username,position,signature from cug_map_users_tpl where student_id = ?;"

)

type User struct {
	StudentId string `form:"student_id" json:"student_id" binding:"required"`
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding"required"`
	Position  string `json:"position"`
	Signature string `json:"signature"`
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
func (newUser *User) QueryUserId() (id string) {
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
	result.Scan(&user.StudentId, &user.Username, &user.Password, &user.Position, &user.Signature)
	fmt.Println(user)
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
func (newUser *User) UpdateUserInfo() (err error) {
	fmt.Println(newUser)
	err = nil
	tx, err := Db.Begin()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	stmt, err := tx.Prepare(updateUserPosition)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	result, err := stmt.Exec(newUser.Position, newUser.StudentId)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	fmt.Println("id=", id)
	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	fmt.Println("rows=", rows)
	stmt, err = tx.Prepare(updateUserSignature)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	result, err = stmt.Exec(newUser.Signature, newUser.StudentId)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	id, err = result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	fmt.Println("id=", id)
	rows, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	fmt.Println("rows=", rows)
	tx.Commit()
	return
}
func (newUser *User) QueryUserPosition() {
	stmt, err := Db.Prepare(queryUserPosition)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	result := stmt.QueryRow(newUser.StudentId)
	result.Scan(&newUser.Position)
	return
}
func (newUser *User)QueryUserName()(name string){
	stmt, err := Db.Prepare(queryStudentNameSrt)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	result := stmt.QueryRow(newUser.StudentId)
	result.Scan(&name)
	return
}
func (newUser *User)QueryAllInfo()(err error){
	stmt,err := Db.Prepare(queryAllInfo)
	defer stmt.Close()
	if err != nil {
		return
	}
	rows := stmt.QueryRow(newUser.StudentId)
	err = rows.Scan(&newUser.StudentId,&newUser.Username,&newUser.Position,&newUser.Signature)
	return
}
func GetAllUserInfo(position string) (userList []*UserWithDistance, err error) {
	userList = make([]*UserWithDistance, 0)
	stmt, err := Db.Prepare(queryAllUserInfo)
	defer stmt.Close()
	if err != nil {
		return
	}
	rows, err := stmt.Query()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		rows.Scan(&user.StudentId, &user.Username, &user.Position,&user.Signature)
		if strings.TrimSpace(user.Position) != "" {
			userWithD := Getdistance(position, &user)
			userList = append(userList, userWithD)
		}

	}
	userList = getTargetList(userList)
	return
}
