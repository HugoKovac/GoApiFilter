package main

import (
	"fmt"
	"net/http"
	"regexp"
)

// Regex for validating domain names
var domainRegex = regexp.MustCompile(`^(?:[-A-Za-z0-9]+\.)+[A-Za-z]{2,6}$`)

/*
	[GET Method] return if domain name is blocked, allowed, or unknown
*/
func handlerDomainStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		domain := r.URL.Query().Get("domain") // Get domain name from query parameter

		// Check if domain name is empty
		if domain == "" {
			http.Error(w, "Domain name is empty", http.StatusBadRequest)
			return
		}

		// Validate domain name using regular expression
		if !domainRegex.MatchString(domain) {
			http.Error(w, "Invalid domain name", http.StatusBadRequest)
			return
		}

		value, err := redisClient.Get(ctx, domain).Result() //* Request Redis 

		if err != nil {
			fmt.Fprintf(w, "unknown")
			return
		}

		fmt.Fprintf(w, "%s\n", value)

		
		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}
