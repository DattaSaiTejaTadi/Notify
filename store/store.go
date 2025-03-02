package store

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type store struct {
	DB *sql.DB
}

func New(username, password string) *store {
	dsn := username + ":" + password + "@tcp(127.0.0.1:3306)/user_management"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v \n issue with %v", err, dsn)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}
	return &store{
		DB: db,
	}
}
