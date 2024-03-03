package main

import (
	"net/http"
	"fmt"
)

/*
	[GET Method] return if domain name is blocked, allowed, or unknown
*/
func handlerDomainStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		domain := r.URL.Query().Get("domain") //* Get domain name from Json

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