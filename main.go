package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Ishwar2510/Logiq.ai-assignmnet/cache"
)

func main() {
	cache := cache.NewCache(100)

	http.HandleFunc("/cache", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var data struct {
				Key        string        `json:"key"`
				Value      interface{}   `json:"value"`
				Expiration int  		 `json:"expiration_in_sec"`
			}

			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "error decoding request body: %v", err)
				return
			}

			if err := cache.Add(data.Key, data.Value, data.Expiration); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "error adding value to cache: %v", err)
				return
			}

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "value added to cache")

		case http.MethodGet:
			key := r.URL.Path[len("/cache/"):]
			value, err := cache.Get(key)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "error getting value from cache: %v", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(value)

		case http.MethodDelete:
			key := r.URL.Path[len("/cache/"):]
			if err := cache.Delete(key); err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "error deleting value from cache: %v", err)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "value deleted from cache")

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "unsupported method: %v", r.Method)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("error starting server: %v", err)
	}
}
