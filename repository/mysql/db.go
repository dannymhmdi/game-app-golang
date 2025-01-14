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
	db, oErr := sql.Open("mysql", "gameapp:gameappt0lk2o20@tcp(localhost:3308)/gameapp_db")
	if oErr != nil {
		log.Fatalf("failed to connect database: %v\n", oErr)
	}
	fmt.Println("connected to database")
	return &MysqlDB{db}

}
