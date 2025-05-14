package main

import (
	"exchange-rate-serivce/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/convert", handlers.Convert)
	r.Get("/rate", handlers.LatestRate)
	r.Get("/history", handlers.HistoricaRates)

	go handlers.StartRateRefresher()

	fmt.Println("server is listening on port: 5050")
	err := http.ListenAndServe(":5050", r)
	if err != nil {
		log.Fatal(err)
	}

}
