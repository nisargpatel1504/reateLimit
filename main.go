package main

import (
	"RateLimiter/fileReader"
	"RateLimiter/ratelimit"
	"net/http"

	"github.com/gorilla/mux"
)

func redirectHandler(w http.ResponseWriter, r *http.Request){

}
func main(){
	r := mux.NewRouter();
	fileReader.ReadingConfigFile();
	businessLogicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(" "));
    });
    r.Handle("/user/{id}", ratelimit.RateApiLimit(businessLogicHandler)).Methods("PATCH");
    r.HandleFunc("/userinfo/{id}", redirectHandler).Methods("PATCH")

    http.ListenAndServe(":8080", r)
}