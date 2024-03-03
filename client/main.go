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


func main() {

	status := make(chan string)
	test := "xxxxxxxxxxxxx"
	
	for pos, _ := range(test) {
		go SubmitDomain(test[:pos+1] + ".com", status)	
	}
	
	<- status

	for x:=range status {
		log.Printf("Status: %s\n", x)
	}
	
	for pos, _ := range(test) {
		go CheckDomain(test[:pos] + ".com", status)	
	}

	<- status

	for x:=range status {
		log.Printf("%s", x)
	}
}
