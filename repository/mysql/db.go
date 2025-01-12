package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlDB struct {
	db *sql.DB
}

func New() *MysqlDB {
	dsn := "root:gameappRoo7t0lk2o20(127.0.0.1:3308)/gameapp_db"
	db, oErr := sql.Open("mysql", dsn)
	if oErr != nil {
		log.Fatalf("failed to connect database: %v\n", oErr)
	}
	fmt.Println("connected to database")
	return &MysqlDB{db}
}
