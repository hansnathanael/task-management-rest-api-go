package app

import (
	"database/sql"
	"task-management/helper"
	"time"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:P@ssw0rd@tcp(localhost:3306)/task_management")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}