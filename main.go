package main

import (
	"fmt"
	"time"
)

type Result struct {
	ErrorName          string
	NumberOfOccurances int64
}

func getErrorName(errorId string) (<-chan string, <-chan error) {
	names := map[string]string{
		"404": "Not Found",
		"400": "Bad Request",
		"500": "Internal Server Error",
	}

	out := make(chan string, 1)
	errs := make(chan error, 1)

	go func() {
		time.Sleep(time.Second)
		if name, ok := names[errorId]; ok {
			out <- name
		} else {
			errs <- fmt.Errorf("getErrorName: %s errorId not found", errorId)
		}

		close(out)
		close(errs)
	}()

	return out, errs
}

func getOccurances(errorId string) (<-chan int64, <-chan error) {
	occurances := map[string]int64{
		"404": 150, // Sample occurrence count
		"400": 75,  // Sample occurrence count
		"500": 30,
	}

	out := make(chan int64, 1)
	errs := make(chan error, 1)

	go func() {
		time.Sleep(time.Second)
		if occ, ok := occurances[errorId]; ok {
			out <- occ
		} else {
			errs <- fmt.Errorf("getOccurances: %s errorId not found", errorId)
		}

		close(out)
		close(errs)
	}()

	return out, errs
}

func getError(errorId string) (r *Result, err error) {

	nameOut, nameErr := getErrorName(errorId)
	occurancesOut, occurancesErr := getOccurances(errorId)

	var open bool

	if err, open = <-nameErr; open {
		return
	}
	if err, open = <-occurancesErr; open {
		return
	}
	r = &Result{ErrorName: <-nameOut, NumberOfOccurances: <-occurancesOut}

	return
}

func main() {
	fmt.Println("Using separate channels for error and result")
	errorIds := []string{
		"404",
		"400",
		"500",
	}
	for _, e := range errorIds {
		r, err := getError(e)
		if err != nil {
			fmt.Printf("Failed: %s\n", err.Error())
			continue
		}
		fmt.Printf("Name: \"%s\" has occurred \"%d\" times\n", r.ErrorName, r.NumberOfOccurances)
	}
}
