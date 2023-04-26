package main

import (
	"github.com/Ishwar2510/Logiq.ai-assignmnet/api"
	"log"
	"net/http"

	"github.com/Ishwar2510/Logiq.ai-assignmnet/cache"
	"github.com/gorilla/mux"
)

func main() {
	che := cache.NewCache(100)
	api := api.NewAPI(che)
	router := mux.NewRouter()
	router.HandleFunc("/cache/{maxSize}", api.MaxSize).Methods(http.MethodPut)
	router.HandleFunc("/cache", api.Add).Methods(http.MethodPost)
	router.HandleFunc("/cache/{key}", api.Get).Methods(http.MethodGet)
	router.HandleFunc("/cache/{key}", api.Delete).Methods(http.MethodDelete)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
	
}