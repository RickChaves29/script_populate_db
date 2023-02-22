package data

import (
	"database/sql"
	"log"
)

func CreateMovie(id int64, title string, year int, gender string, db *sql.DB) {
	_, err := db.Exec(`INSERT INTO movies (id, title, year, genres) VALUES ($1, $2, $3, $4);`, id, title, year, gender)
	if err != nil {
		log.Println(err.Error())
	}
}
