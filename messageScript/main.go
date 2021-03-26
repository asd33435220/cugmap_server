package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"time"
)

func run() (err error) {
	Db, err := sql.Open("mysql", "root:asd33435220@tcp/cug_map_db")
	if err != nil {
		fmt.Println(err.Error())
	}
	const insertMessageStr = "insert into cug_map_message_tpl2(receiver_id,sender_id,message,send_time,send_time_str) values(?,?,?,?,?);"
	stmt, err := Db.Prepare(insertMessageStr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for i := 0;i<5 ;i++  {
		time.Sleep(time.Second)
		receiver_id := "30000000001"
		sender_id := "20181000664"
		send_time := strconv.Itoa(int(time.Now().UnixNano()/1e6))
		send_time_str := time.Now().Format("2006-01-02 15:04:05")
		message := "回复1号测试消息" + strconv.Itoa(i)
		_, err = stmt.Exec(receiver_id, sender_id, message, send_time, send_time_str)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}


	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//id, err := result.LastInsertId()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Println("id=", id)
	//rows, err := result.RowsAffected()
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//fmt.Println("rows=", rows)
	return nil
}

func main() {
	run()
}
