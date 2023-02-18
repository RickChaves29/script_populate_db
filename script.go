package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func init() {
	db, err := connDatabase(os.Getenv("CONNECT_DB"))
	if err != nil {
		log.Fatalln(err.Error())
	}
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
	db, err := connDatabase(os.Getenv("CONNECT_DB"))
	if err != nil {
		log.Fatal(err.Error())
	}
	body, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	readNewFile := csv.NewReader(bytes.NewBuffer(body))
	readNewFile.LazyQuotes = true
	_, err = readNewFile.Read()
	if err != nil {
		log.Fatalln(err.Error())
	}
	for {
		row, err := readNewFile.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		idClean, err := strconv.Atoi(removeWhiteSpace(row[0]))
		if err != nil {
			log.Fatal(err.Error())
		}
		titleClean := removeWhiteSpace(row[1])
		genresClean := removeWhiteSpace(row[2])

		id := int64(idClean)
		title := getTitle(titleClean)
		year, err := getYear(titleClean)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = db.Exec(`INSERT INTO movies (id, title, year, genres) VALUES ($1, $2, $3, $4);`, id, title, year, genresClean)
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("id: %v", idClean)
		log.Printf("title: %v", title)
		log.Printf("year: %v", year)

		genres, err := getAllGenres(genresClean)
		if err != nil {
			log.Fatal(err.Error())
		}
		c := make(chan string)
		go channalOfGenres(genres, c)
		log.Println("genres: ", <-c)
	}

}

func removeWhiteSpace(line string) string {
	rx, err := regexp.Compile(`^[\s]+|[\s]+$`)
	if err != nil {
		return ""
	}
	stringClean := rx.ReplaceAllString(line, "")
	return stringClean
}

func getTitle(line string) string {
	rx, err := regexp.Compile(`\s*?[\(]{1}[0-9]{4}[\)]{1}`)
	if err != nil {
		return ""
	}
	newString := rx.ReplaceAllString(line, "")
	return newString
}
func getYear(line string) (int, error) {
	var stringToInt int
	if line == "" {
		return 0, errors.New("string line is empty")
	}
	rx, err := regexp.Compile(`[^0-9]`)
	if err != nil {
		return 0, err
	}
	newString := rx.ReplaceAllString(line, " ")
	if newString == "" {
		return 0, nil
	}
	stringClean := removeWhiteSpace(newString)
	if len(newString) > 4 {
		stringSplited := strings.Split(stringClean, " ")

		lastString := stringSplited[len(stringSplited)-1]
		stringToInt, err = strconv.Atoi(lastString)
		return stringToInt, err
	}
	stringToInt, err = strconv.Atoi(newString)
	if err != nil {
		return 0, err
	}
	return stringToInt, nil
}

func getAllGenres(line string) ([]string, error) {
	rx, err := regexp.Compile(`(?:[\(\)]|[\s][^a-zA-Z]+)`)
	if err != nil {
		return nil, err
	}
	stringRemove := rx.ReplaceAllString(line, "")
	genres := strings.Split(stringRemove, "|")
	if err != nil {
		return nil, err
	}
	return genres, nil
}

func channalOfGenres(genres []string, c chan string) {
	var genre string
	for i := range genres {
		genre += genres[i] + " "

	}
	c <- genre
}

func connDatabase(strConnection string) (*sql.DB, error) {
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
