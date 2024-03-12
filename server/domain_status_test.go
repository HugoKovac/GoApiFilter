package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlerDomainStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerDomainStatus))
	defer server.Close()


	// valid domain
	resp, err := http.Get(server.URL + "/v1/domain_status?domain=example.com")
	if err != nil {
		t.Fatal(err)
	}
	
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status code should be 200, got %d", resp.StatusCode)
	}

	// no domain
	resp, err = http.Get(server.URL + "/v1/domain_status")
	if err != nil {
		t.Fatal(err)
	}
	
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code should be 400, got %d", resp.StatusCode)
	}

	// invalid domain
	resp, err = http.Get(server.URL + "/v1/domain_status?domain=example")
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status code should be 400, got %d", resp.StatusCode)
	}

}

// Function that will do all the requests in go routines called in TestHandlerSubmitDomainStress
func CheckDomain(domain string, client *http.Client, url string, t *testing.T) {
	data := Domain{Domain: domain}
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	r, err := client.Get(url + "/v1/domain_status?domain=" + domain)

	if err != nil {
		t.Fatal(err)
	}

	if !domainRegex.MatchString(domain) {
		if r.StatusCode != http.StatusBadRequest {
			t.Errorf("Status code should be 400, got %d for %s", r.StatusCode, domain)
		}
	} else if (r.StatusCode != http.StatusCreated) && (r.StatusCode != http.StatusOK) {
		t.Errorf("Status code should be 200 or 201, got %d for %s", r.StatusCode, domain)
	}

	defer r.Body.Close()
	validBody := []string{"blocked", "allowed", "unknown"}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !contains(validBody, string(body)) {
		t.Errorf("Invalid response body: %s", string(body))
	}

	// status <- domain + ": " + string(body)
}

func TestHandlerDomainStatusStress(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handlerDomainStatus))
	defer server.Close()

	records := openCsv("./domains/alexa-domains-1M.csv")
	goRoutineRequests(records, CheckDomain, 200, server.URL, t)

}


