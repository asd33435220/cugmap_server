package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math/rand"
	"strconv"
)

func run() (err error) {
	Db, err := sql.Open("mysql", "root:asd33435220@tcp/cug_map_db")
	if err != nil {
		fmt.Println(err.Error())
	}
	const insertUserStr = "insert into cug_map_users_tpl(student_id,username,password,lng,lat,signature) values(?,?,?,?,?,?);"
	stmt, err := Db.Prepare(insertUserStr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for i := 1; i < 4000000; i++ {
		studentId := 40001000000 + i
		student_id := strconv.Itoa(studentId)
		name := "测试账号" + strconv.Itoa(i)
		password := "asd33435220"
		lng := 75.000000 + rand.Float64()*70
		lat := 20.000000 + rand.Float64()*33
		signature := "测试签名" + strconv.Itoa(i)
		_, err = stmt.Exec(student_id, name, password, lng, lat, signature)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func main() {
	run()
}
