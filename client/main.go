package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Domain struct {
	Domain string `json:"domain"`
}


func SubmitDomain(domain string, status chan string) {
	data := Domain{Domain: domain}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.Post("http://server/v1/submit_domain", "application/json", &buffer)

	if err != nil {
		log.Fatal(err)
	}

	status <- r.Status
}

func CheckDomain(domain string, status chan string) {
	data := Domain{Domain: domain}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	r, err := http.Get("http://server/v1/domain_status?domain=" + domain)

	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	status <- domain + ": " + string(body)
}

func goRoutineRequests(url string, tld string, f func(domain string, status chan string)){
	rtnChan := make(chan string)

	for pos := range url {
		go f(url[:pos+1]+tld, rtnChan)
	}

	// Collect responses from f calls
	for range url {
		log.Printf("rtnChan: %s\n", <-rtnChan)
	}

	close(rtnChan) // Close the channel after all submissions are done	
}

func main() {
	goRoutineRequests("test", ".com", SubmitDomain) // submit new domain to check
	goRoutineRequests("test", ".com", CheckDomain) // get status of domain
	goRoutineRequests("test", "", CheckDomain) // get status of invalid domain
}
