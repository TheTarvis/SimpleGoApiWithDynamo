package main

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"herbie/api/account"
	"herbie/api/vehicle"
	"net/http"
)

var log = logrus.New()

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/account", account.CreateAccountHandler).Methods("POST")
	router.HandleFunc("/account/{id}", account.UpdateAccountHandler).Methods("POST")
	router.HandleFunc("/account/{id}", account.GetAccountHandler).Methods("GET")

	router.HandleFunc("/account/{accountId}/vehicle", vehicle.CreateVehicleHandler).Methods("POST")
	router.HandleFunc("/account/{accountId}/vehicle/{id}", vehicle.UpdateVehicleHandler).Methods("POST")
	router.HandleFunc("/account/{accountId}/vehicle/{id}", vehicle.GetVehicleHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":7070", router))
}
