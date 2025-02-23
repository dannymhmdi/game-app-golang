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

func (m *MysqlDB) NewConn() *sql.DB {
	return m.db
}

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Host     string `koanf:"host"`
	Port     uint   `koanf:"port"`
	DbName   string `koanf:"db_name"`
}

func New(cfg Config) *MysqlDB {
	db, oErr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName))
	if oErr != nil {
		log.Fatalf("failed to connect database: %v\n", oErr)
	}
	fmt.Println("connected to database")

	return &MysqlDB{db: db, config: cfg}

}
