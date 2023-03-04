package database

import (
	"database/sql"
	"fmt"
	"time"
)

func NewConnection(config DbConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", config.User, config.Password, config.Host, config.Dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db, nil
}
