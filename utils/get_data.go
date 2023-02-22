package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func RemoveWhiteSpace(line string) string {
	rx, err := regexp.Compile(`^[\s]+|[\s]+$`)
	if err != nil {
		return ""
	}
	stringClean := rx.ReplaceAllString(line, "")
	return stringClean
}

func GetTitle(line string) string {
	rx, err := regexp.Compile(`\s*?[\(]{1}[0-9]{4}[\)]{1}`)
	if err != nil {
		return ""
	}
	newString := rx.ReplaceAllString(line, "")
	return newString
}

func GetYear(line string) (int, error) {
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
	stringClean := RemoveWhiteSpace(newString)
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

func GetAllGenres(line string) ([]string, error) {
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
