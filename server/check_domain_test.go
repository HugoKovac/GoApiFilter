package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"encoding/json"
)

func TestHandlerSubmitDomain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerSubmitDomain))
	defer server.Close()
	
	// empty body
	resp, err := http.Post(server.URL + "/v1/submit_domain", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code should be 400, got %d", resp.StatusCode)
	}

	// valid domain
	//* create json for request body	
	data := Domain{Domain: "example.com"}
	var buffer bytes.Buffer
	err = json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.Post(server.URL + "/v1/submit_domain", "application/json", &buffer)
	if err != nil {
		t.Fatal(err)
	}

	if (resp.StatusCode != http.StatusCreated) && (resp.StatusCode != http.StatusOK) {
		t.Errorf("Status code should be 200 or 201, got %d", resp.StatusCode)
	}

	// invalid domain
	//* create json for request body	
	data = Domain{Domain: "example"}
	err = json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.Post(server.URL + "/v1/submit_domain", "application/json", &buffer)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code should be 400, got %d", resp.StatusCode)
	}
}

// Function that will do all the requests in go routines called in TestHandlerSubmitDomainStress
func SubmitDomain(domain string, _status chan string, client *http.Client, url string, t *testing.T) {
	data := Domain{Domain: domain}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		t.Fatal(err)
		return
	}

	resp, err := client.Post(url + "/v1/submit_domain", "application/json", &buffer)
	
	if err != nil {
		t.Fatal(err)
	}
	
	if !domainRegex.MatchString(domain) {
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Status code should be 400, got %d for %s", resp.StatusCode, domain)
		}
	} else if (resp.StatusCode != http.StatusCreated) && (resp.StatusCode != http.StatusOK) {
		t.Errorf("Status code should be 200 or 201, got %d for %s", resp.StatusCode, domain)
	}
	// status <- resp.Status

}

func TestHandlerSubmitDomainStress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerSubmitDomain))
	defer server.Close()

	records := openCsv("./domains/alexa-domains-1M.csv")
	goRoutineRequests(records, SubmitDomain, 200, server.URL, t)

}
