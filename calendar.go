package icsreader

import (
	"net/http"
	"os"
)

// GetCalendarFromFile : automatic parsing from file
func GetCalendarFromFile(filename string) (Events, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, &CalendarError{}
	}
	return Parse(file), nil
}

// GetCalendarFromURL : automatic parsing from source (URL)
func GetCalendarFromURL(url string) (Events, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, &CalendarError{}
	}
	return Parse(resp.Body), nil
}
