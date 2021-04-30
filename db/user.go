package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)


const (
	insertUserSrt       = "insert into cug_map_users_tpl(student_id,username,password) values(?,?,?);"
	queryStudentIdSrt   = "select student_id from cug_map_users_tpl where student_id = ?;"
	queryStudentNameSrt = "select username from cug_map_users_tpl where student_id = ?;"
	queryStudentSrt     = "select username,password from cug_map_users_tpl where student_id = ?;"
	queryStudentInfoSrt = "select student_id,username,password,lng,lat,signature from cug_map_users_tpl where student_id = ?;"
	updateUserPosition  = "update cug_map_users_tpl set lng = ?,lat = ? where student_id = ?;"
	updateUserSignature = "update cug_map_users_tpl set signature = ? where student_id = ?;"
	queryUserPosition   = "select lng,lat from cug_map_users_tpl where student_id = ?;"
	countUserNumber     = "select count(student_id) from cug_map_users_tpl where lng between ? and ? and lat between ? and ?;"
	queryAllUserInfo    = "select student_id,username,lng,lat,signature from cug_map_users_tpl where lng between ? and ? and lat between ? and ?;"
	queryAllInfo        = "select student_id,username,lng,lat,signature from cug_map_users_tpl where student_id = ?;"
)

type User struct {
	StudentId string  `form:"student_id" json:"student_id" binding:"required"`
	Username  string  `form:"username" json:"username" binding:"required"`
	Password  string  `form:"password" json:"password" binding"required"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Signature string  `json:"signature"`
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

	err = result.Scan(&user.StudentId, &user.Username, &user.Password, &user.Longitude, &user.Latitude, &user.Signature)
	if err != nil {
		fmt.Println(err.Error())
	}
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
	_, err = result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
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
	result, err := stmt.Exec(newUser.Longitude, newUser.Latitude, newUser.StudentId)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	_, err = result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
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
	_, err = result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return
	}
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
	result.Scan(&newUser.Longitude, &newUser.Latitude)
	return
}
func (newUser *User) QueryUserName() (name string) {
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
func (newUser *User) QueryAllInfo() (err error) {
	stmt, err := Db.Prepare(queryAllInfo)
	defer stmt.Close()
	if err != nil {
		return
	}
	rows := stmt.QueryRow(newUser.StudentId)
	err = rows.Scan(&newUser.StudentId, &newUser.Username, &newUser.Longitude, &newUser.Latitude, &newUser.Signature)
	return
}
func GetAllUserInfo(lng, lat float64, myId string) (userList []*User, err error) {
	number := 0
	Lrange := 0.6
	userList = make([]*User, 0)
	stmt2, err2 := Db.Prepare(queryAllUserInfo)
	if err2 != nil {
		fmt.Println("err2", err2.Error())
		return
	}
	defer stmt2.Close()
	stmt, err := Db.Prepare(countUserNumber)
	if err != nil {
		fmt.Println("err1", err.Error())
		return
	}
	defer stmt.Close()
	row := stmt.QueryRow(lng-Lrange, lng+Lrange, lat-Lrange, lat+Lrange)
	row.Scan(&number)
	Lrange = Lrange / 3
	i := 1
	for number > 100 {
		start := time.Now() // 获取当前时间
		row = stmt.QueryRow(lng-Lrange, lng+Lrange, lat-Lrange, lat+Lrange)
		elapsed := time.Since(start)
		row.Scan(&number)
		fmt.Printf("第%d次查询,耗时%s,共%d位用户符合要求\n", i, elapsed, number)
		i++
		Lrange = Lrange / 3
	}
	Lrange = Lrange * 4
	var rows *sql.Rows
	start := time.Now() // 获取当前时间
	if Lrange > 9 {
		rows, err = stmt2.Query(-180, 180, -90, 90)
	} else {
		rows, err = stmt2.Query(lng-Lrange, lng+Lrange, lat-Lrange, lat+Lrange)
	}
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		rows.Scan(&user.StudentId, &user.Username, &user.Longitude, &user.Latitude, &user.Signature)
		if user.Longitude != 0 && user.Latitude != 0 && user.StudentId != myId {
			userList = append(userList, &user)
		}
	}
	elapsed := time.Since(start)
	fmt.Println("提取符合要求的用户数据耗时：", elapsed)
	return
}
