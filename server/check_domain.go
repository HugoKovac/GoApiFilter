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
		domain_submit := Domain{}
		decoder := json.NewDecoder(r.Body).Decode(&domain_submit)
		if decoder != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}

		_, err := redisClient.Get(ctx, domain_submit.Domain).Result()
		if err != nil { //* If not exist
			statusChan := make(chan string)
			go domainStatus(domain_submit.Domain, statusChan) //* concurrency to not be blocked by api call

			status := <- statusChan //* wait go routines and get the status of check domain

			err := redisClient.Set(ctx, domain_submit.Domain, status, 0).Err() //* value to domain
			log.Printf("%s have been created\n", domain_submit.Domain)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.WriteHeader(201) //* 201 CREATED
		} else {//* If already exist
			log.Printf("%s already exist in db\n", domain_submit.Domain)
		}

		return
	}
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)

}
