package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlDB struct {
	config Config
	db     *sql.DB
}

type Config struct {
	Username string
	Password string
	Host     string
	Port     uint
	DbName   string
}

func New(cfg Config) *MysqlDB {
	db, oErr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName))
	if oErr != nil {
		log.Fatalf("failed to connect database: %v\n", oErr)
	}
	fmt.Println("connected to database")

	return &MysqlDB{db: db, config: cfg}

}
