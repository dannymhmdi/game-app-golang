package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:gameappRoo7t0lk2o20(127.0.0.1:3308)/gameapp_db"
	connetion, oErr := sql.Open("mysql", dsn)
	if oErr != nil {
		fmt.Println("failed to open db", oErr)

		return
	}
	fmt.Println("connected to db", connetion.Ping())
	
}
