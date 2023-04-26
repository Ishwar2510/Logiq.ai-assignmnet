package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ishwar2510/Logiq.ai-assignmnet/cache"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type API struct {
	che cache.Cache
}


func NewAPI(che cache.Cache) *API {
	return &API{
		che : che,
	}
}

func writeResponse(w http.ResponseWriter, statusCode int, message string) {
	type data struct {
		StatusCode  int `json:"status_code""`
		Message   string    `json:"message""`
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data{
		StatusCode: statusCode,
		Message: message,
	})

}

func (a *API) Add(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Key        		string        `json:"key"`
			Value      		interface{}   `json:"value"`
			ExpirationInSecs int64		  `json:"expiration_in_secs"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			writeResponse(w,http.StatusBadRequest,err.Error())
			return
		}

		if err := a.che.Add(data.Key, data.Value, time.Duration(data.ExpirationInSecs)*time.Second); err != nil {
			if errors.Is(err,cache.ErrAlreadyExists) || errors.Is(err,cache.ErrMaxLimitReached){
				writeResponse(w,http.StatusBadRequest,err.Error())
				return
			}
			writeResponse(w,http.StatusInternalServerError,fmt.Sprintf("error adding value to cache: %v", err))
			return
		}

	writeResponse(w,http.StatusCreated,"added to cache")
}

func (a *API) Get(w http.ResponseWriter, r *http.Request) {
		key := mux.Vars(r)["key"]
		value, err := a.che.Get(key)
		if err != nil {
			if errors.Is(err,cache.ErrNotFound) {
				writeResponse(w,http.StatusNotFound,err.Error())
				return
			}
			writeResponse(w,http.StatusInternalServerError,fmt.Sprintf( "error getting value from cache: %v", err))
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(value)
}

func (a *API) Delete(w http.ResponseWriter, r *http.Request)  {
		key := mux.Vars(r)["key"]
		if err := a.che.Delete(key); err != nil {
			if errors.Is(err,cache.ErrNotFound) {
				writeResponse(w,http.StatusNotFound,err.Error())
				return
			}
			writeResponse(w,http.StatusInternalServerError,fmt.Sprintf( "error deleting value from cache: %v", err))
			return
		}
		writeResponse(w,http.StatusOK,"deleted from cache")
}

func (a *API) MaxSize(w http.ResponseWriter, r *http.Request)  {
	key := mux.Vars(r)["maxSize"]
	 size,err :=  strconv.Atoi(key)
	 if err != nil {
		 writeResponse(w,http.StatusBadRequest,err.Error())
		 return
	 }
	 a.che.MaxSize(size)
	writeResponse(w,http.StatusOK,"updated max size")
}

