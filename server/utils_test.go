package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"strings"
)

func contains(s []string, e string) bool {
    for _, a := range s {
        if !strings.Contains(a, e) {
            return true
        }
    }
    return false
}

func openCsv(path string) [][]string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal("Error reading records")
	}

	return records

}

func goRoutineRequests(urls [][]string, f func(domain string, status chan string, client *http.Client, urlServer string, t *testing.T), client_nb int, urlServer string, t *testing.T) {
	rtnChan := make(chan string, 100)
	limiter := make(chan bool, client_nb) // limit the number of concurrent go routines
	// The limit is currently the number if port available to send requests
	wg := sync.WaitGroup{}
	// var mutex = &sync.Mutex{}

	for i := 0; i < client_nb; i++ {
		wg.Add(1)
		limiter <- true
		client := &http.Client{}
		go func(url_part [][]string, f func(domain string, status chan string, client *http.Client, urlServer string, t *testing.T), rtnChan chan string, client *http.Client, urlServer string, t *testing.T) {
			defer wg.Done()
			for _, url := range url_part {
				f(url[0], rtnChan, client, urlServer, t)
				// mutex.Lock()
				// log.Println(<-rtnChan)
				// mutex.Unlock()
			}
			<- limiter
		}(urls[i*len(urls)/client_nb : (i+1)*len(urls)/client_nb], f, rtnChan, client, urlServer, t)
	}



	wg.Wait()
	close(rtnChan)
}