package main

import (
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

	// todo read body

}


