package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"time"
)

type Mysql struct {
	USERNAME string
	PASSWORD string
	NETWORK  string
	SERVER   string
	PORT     int
	DATABASE string
	db 		 *sql.DB
}

func Getconn() *sql.DB {
	USERNAME := "root"
	PASSWORD := "123456"
	NETWORK  := "tcp"
	SERVER   := "localhost"
	PORT     := 3306
	DATABASE := "cicd"
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Errorf("Open mysql failed, err: %s", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}