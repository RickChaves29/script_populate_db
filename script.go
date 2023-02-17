package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

type Movie struct {
	ID    int32
	Title string
}

func main() {
	remountCSV("movies.csv")
}

func remountCSV(file string) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
	}

	rx, err := regexp.Compile(`[\\"\\"\/]`)
	if err != nil {
		log.Fatal(err.Error())
	}

	newFile := rx.ReplaceAll(body, []byte(""))
	readNewFile := csv.NewReader(bytes.NewBuffer(newFile))
	_, err = readNewFile.Read()
	if err != nil {
		log.Fatalln(err.Error())
	}
	var movies []Movie
	for {
		row, err := readNewFile.Read()
		if err == io.EOF {
			break
		}

		id, err := strconv.Atoi(row[0])
		if err != nil {
			log.Fatal(err.Error())
		}

		movies = append(movies, Movie{
			ID:    int32(id),
			Title: row[1],
		})

		fmt.Printf("id: %v, title: %v\n", id, row[1])
	}
}
