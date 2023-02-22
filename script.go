package main

import (
	"bytes"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/RickChaves29/script_populate_db/internal/data"
	"github.com/RickChaves29/script_populate_db/utils"
	_ "github.com/lib/pq"
)

func init() {
	db, err := data.ConnDatabase(os.Getenv("CONNECT_DB"))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS movies (
		id BIGSERIAL PRIMARY KEY,
	  	title TEXT NOT NULL,
	  	year INT NULL, 
	  	genres TEXT NULL
	)
	`)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
func main() {
	remountCSV("movies.csv")
}

func remountCSV(file string) {
	db, err := data.ConnDatabase(os.Getenv("CONNECT_DB"))
	if err != nil {
		log.Fatalln("ERROR DB_CONN: ", err.Error())
	}
	defer db.Close()
	body, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln("ERROR READ_FILE: ", err.Error())
	}
	readNewFile := csv.NewReader(bytes.NewBuffer(body))
	readNewFile.LazyQuotes = true
	_, err = readNewFile.Read()
	if err != nil {
		log.Println("ERROR READ_ROW: ", err.Error())
	}
	for {
		row, err := readNewFile.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("ERROR READ_ROW: ", err.Error())
		}
		idClean, err := strconv.Atoi(utils.RemoveWhiteSpace(row[0]))
		if err != nil {
			log.Println("ERROR CONVERT_ID: ", err.Error())
		}
		titleClean := utils.RemoveWhiteSpace(row[1])
		genresClean := utils.RemoveWhiteSpace(row[2])

		id := int64(idClean)
		title := utils.GetTitle(titleClean)
		year, err := utils.GetYear(titleClean)
		if err != nil {
			log.Println("ERROR GET_YEAR: ", err.Error())
		}
		log.Printf("ID: %v", idClean)
		log.Printf("TITLE: %v", title)
		log.Printf("YEAR: %v", year)
		data.CreateMovie(id, title, year, genresClean, db)
		genres, err := utils.GetAllGenres(genresClean)

		if err != nil {
			log.Println("ERROR GET_ALL_GENRES: ", err.Error())
		}
		c := make(chan string)
		go channalOfGenres(genres, c)
		log.Println("GENRES: ", <-c)
	}
	log.Println("END: add all movies is success")
}

func channalOfGenres(genres []string, c chan string) {
	var genre string
	for i := range genres {
		genre += genres[i] + " "

	}
	c <- genre
}
