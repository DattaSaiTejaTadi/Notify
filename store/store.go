package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type store struct {
	DB *sql.DB
}

func New(password string) *store {
	dsn := fmt.Sprintf("postgresql://postgres.qehkntqehiukjqlbmygr:%s@aws-0-us-west-1.pooler.supabase.com:5432/postgres", password)
	db, err := sql.Open("postgres", dsn)
	db.SetMaxOpenConns(45)                  // Matches Supabase pool size
	db.SetMaxIdleConns(20)                  // Keep some idle connections to avoid frequent creation
	db.SetConnMaxLifetime(20 * time.Minute) // Refresh active connections every 20 mins
	db.SetConnMaxIdleTime(5 * time.Minute)  // Close idle connections after 5 mins
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
