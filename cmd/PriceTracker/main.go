package main

import (
	"net/http"

	core_functions "github.com/Vanodium/pricetracker/internal/services"
	routes "github.com/Vanodium/pricetracker/internal/transport/rest"

	"log"

	"github.com/jasonlvhit/gocron"
)

func main() {
	jobScheduler := gocron.NewScheduler()
	jobScheduler.Every(30).Minutes().Do(core_functions.CheckPrices)
	go jobScheduler.Start()

	router := routes.Router()
	log.Println("API launched")
	err := http.ListenAndServe(":8989", router)
	if err != nil {
		panic(err)
	}
}
