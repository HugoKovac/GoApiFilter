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

func SubmitDomain(domain string, status chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	data := Domain{Domain: domain}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		log.Print(err)
		return
	}
	r, err := http.Post("http://server/v1/submit_domain", "application/json", &buffer)

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

	status <- string(body) + r.Status
}

func CheckDomain(domain string, status chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	data := Domain{Domain: domain}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		log.Print(err)
		return
	}
	r, err := http.Get("http://server/v1/domain_status?domain=" + domain)

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

func goRoutineRequests(urls [][]string, f func(domain string, status chan string, wg *sync.WaitGroup)) {
	rtnChan := make(chan string, 10)
	wg := sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)
		go f(url[0], rtnChan, &wg)
	}

	wg.Wait()
	close(rtnChan)

	for url := range urls {
		log.Printf("rtnChan: %s\n", url)
	}

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

	records := openCsv("./domain_name_1M.csv")

	goRoutineRequests(records, SubmitDomain)
	goRoutineRequests(records, CheckDomain)

}
