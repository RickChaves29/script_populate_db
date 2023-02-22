package data

import (
	"database/sql"
	"log"
)

func ConnDatabase(strConnection string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", strConnection)
	if err != nil {
		return nil, err
	}
	log.Println("database connected")
	err = conn.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("connection is open")
	return conn, nil
}
