package main

import (
	"log"
	"net/http"
	"encoding/json"
	"math/rand"
	"time"
)

func domainStatus(_ string, statusChan chan string) {
	//* [DEBUG] To simulate low api call or processing of checking
	randomNumber := rand.Intn(100)
	time.Sleep(time.Second * 3)

    if randomNumber < 30 {
        statusChan <- "blocked"
    } else {
        statusChan <- "allowed"
    }
}

/*
	[POST Method] take domain name to process and add to db
*/
func handlerSubmitDomain(w http.ResponseWriter, r *http.Request) {
	//* [DEBUG] Time seed for the Debug Sleep
	rand.Seed(time.Now().UnixNano())

	if r.Method == "POST" {
		now := time.Now()
		domain_submit := Domain{}
		decoder := json.NewDecoder(r.Body).Decode(&domain_submit)
		if decoder != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}

		// Validate domain name using regular expression
		if !domainRegex.MatchString(domain_submit.Domain) {
			http.Error(w, "Invalid domain name", http.StatusBadRequest)
			return
		}

		_, err := redisClient.Get(ctx, domain_submit.Domain).Result()
		if err != nil { //* If not exist
			statusChan := make(chan string)
			go domainStatus(domain_submit.Domain, statusChan) //* concurrency to not be blocked by api call

			status := <- statusChan //* wait go routines and get the status of check domain

			err := redisClient.Set(ctx, domain_submit.Domain, status, 0).Err() //* value to domain
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			//* Logging
			log.Printf("%s have been created in: ", domain_submit.Domain)
			chrono := time.Since(now)
			if chrono > 4 * time.Second {
				log.Printf("\033[31m%s\033[0m", chrono)
			} else {
				log.Println(chrono)
			}
			//* End Logging
			
			w.WriteHeader(201) //* 201 CREATED
		} else {//* If already exist
			//* Logging
			log.Printf("%s already exist in db in: ", domain_submit.Domain)
			chrono := time.Since(now)
			if chrono > 1 * time.Second {
				log.Printf("\033[31m%s\033[0m", chrono)
			} else {
				log.Println(chrono)
			}
			//* End Logging
		}

		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

}
