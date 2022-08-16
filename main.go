package main

import (
	"fmt"
	"net/http"

	"simple-colly-example/crawler"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/{vin}", crawler.GetEngines)

	fmt.Println("We are up and running. localhost:8000/{vin}")
	http.ListenAndServe(":8000", router)
}
