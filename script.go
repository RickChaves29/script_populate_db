package main

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Movie struct {
	ID    int32
	Title string
	Year  int
}

func main() {
	remountCSV("movies.csv")
}

func remountCSV(file string) {
	body, err := ioutil.ReadFile(file)
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
		idClean := removeWhiteSpace(row[0])
		titleClean := removeWhiteSpace(row[1])
		genderClean := removeWhiteSpace(row[2])

		title := getTitle(titleClean)
		year, err := getYear(titleClean)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Printf("id: %v", idClean)
		log.Printf("title: %v", title)
		log.Printf("year: %v", year)

		genders, err := getAllGenders(genderClean)
		if err != nil {
			log.Fatal(err.Error())
		}
		c := make(chan string)
		go channalOfGenders(genders, c)
		log.Println("genders: ", <-c)
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
	if line == "" {
		return 0, errors.New("string line is empty")
	}
	rx, err := regexp.Compile(`[^0-9]`)
	if err != nil {
		return 0, err
	}
	newString := rx.ReplaceAllString(line, "")
	if newString == "" {
		return 0, nil
	}
	stringToInt, err := strconv.Atoi(newString)
	if err != nil {
		return 0, err
	}
	return stringToInt, nil
}

func getAllGenders(line string) ([]string, error) {
	rx, err := regexp.Compile(`(?:[\(\)]|[\s][^a-zA-Z]+)`)
	if err != nil {
		return nil, err
	}
	stringRemove := rx.ReplaceAllString(line, "")
	genders := strings.Split(stringRemove, "|")
	if err != nil {
		return nil, err
	}
	return genders, nil
}

func channalOfGenders(genders []string, c chan string) {
	var gender string
	for i := range genders {
		gender += genders[i] + " "

	}
	c <- gender
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
