package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Domain struct {
	Domain string `json:"domain"`
}

func SubmitDomain(domain string, status chan string, client *http.Client) {
	data := Domain{Domain: domain}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		log.Print(err)
		return
	}

	r, err := client.Post("http://server/v1/submit_domain", "application/json", &buffer)
	
	if err != nil {
		log.Println(err)
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	status <- string(body) + r.Status

}

func CheckDomain(domain string, status chan string, client *http.Client) {
	data := Domain{Domain: domain}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		log.Print(err)
		return
	}
	r, err := client.Get("http://server/v1/domain_status?domain=" + domain)

	if err != nil {
		log.Print(err)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		return
	}

	status <- domain + ": " + string(body)
}

func goRoutineRequests(urls [][]string, f func(domain string, status chan string, client *http.Client), client_nb int) {
	rtnChan := make(chan string, 100)
	limiter := make(chan bool, client_nb) // limit the number of concurrent go routines
	// The limit is currently the number if port available to send requests
	wg := sync.WaitGroup{}
	var mutex = &sync.Mutex{}

	for i := 0; i < client_nb; i++ {
		wg.Add(1)
		limiter <- true
		client := &http.Client{}
		go func(url_part [][]string, f func(domain string, status chan string, client *http.Client), rtnChan chan string, client *http.Client) {
			defer wg.Done()
			for _, url := range url_part {
				f(url[0], rtnChan, client)
				mutex.Lock()
				log.Println(<-rtnChan)
				mutex.Unlock()
			}
			<- limiter
		}(urls[i*len(urls)/client_nb : (i+1)*len(urls)/client_nb], f, rtnChan, client)
	}



	wg.Wait()
	close(rtnChan)
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

func main() {

	records := openCsv("./alexa-domains-1M.csv")

	goRoutineRequests(records, SubmitDomain, 500)
	// goRoutineRequests(records, CheckDomain, 500)

}
