package main

import (
	"net/http"

	routes "github.com/Vanodium/pricetracker/internal/routes"
	core_functions "github.com/Vanodium/pricetracker/internal/services"

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
