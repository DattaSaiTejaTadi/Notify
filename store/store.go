package store

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type store struct {
	DB *sql.DB
}

func New(password string) *store {
	dsn := fmt.Sprintf("postgresql://postgres.qehkntqehiukjqlbmygr:%s@aws-0-us-west-1.pooler.supabase.com:5432/postgres", password)
	db, err := sql.Open("postgres", dsn)
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
