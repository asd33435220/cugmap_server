package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB
var err error

func init() {
	Db, err = sql.Open("mysql", "root:asd33435220@tcp/cug_map_db")
	if err != nil {
		fmt.Println(err.Error())
	}
}
