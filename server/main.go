package main

import (
	"context"
	"log"
	"net/http"
	"github.com/redis/go-redis/v9"
)

// TODO if lot of json struct create "DMO" file
type Domain struct {
	Domain string `json:"domain"`
}

var ( // to use persistant map in all the server
	ctx = context.Background()
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB: 0,
	}) 
)

func main() {
	http.HandleFunc("/v1/submit_domain", handlerSubmitDomain)
	http.HandleFunc("/v1/domain_status", handlerDomainStatus)
	log.Fatal(http.ListenAndServe(":80", nil))
}
