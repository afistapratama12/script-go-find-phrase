package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Credential struct {
	Host         string
	Username     string
	Password     string
	DatabaseName string
	Port         int
	Schema       string
}

func ConnectDB(creds *Credential) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=require", creds.Host, creds.Username, creds.Password, creds.DatabaseName, creds.Port)

	// connect using database/sql + pq
	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
