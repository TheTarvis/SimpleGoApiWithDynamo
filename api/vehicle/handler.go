package vehicle

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var log = logrus.New()

func GetVehicleHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accountId := params["accountId"]
	id := params["id"]

	_, err := uuid.Parse(accountId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error invalid accountId. %s", err)
		return
	}

	_, err = uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error invalid vehicle id. %s", err)
		return
	}

	vehicle := Vehicle{}
	err = vehicle.GetVehicleById(accountId, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching vehicle id. %s", err)
		return
	}

	writeResponse(w, &vehicle)
}

func UpdateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accountId := params["accountId"]
	id := params["id"]


	_, err := uuid.Parse(accountId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error invalid accountId. %s", err)
		return
	}

	_, err = uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error invalid vehicle id. %s", err)
		return
	}

	vehicle, err := getVehicleFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing body. %s", err)
		return
	}

	vehicle.ID = id
	vehicle.AccountId = accountId
	err = vehicle.UpdateVehicle()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating vehicle id. %s", err)
		return
	}

	writeResponse(w, vehicle)
}

func CreateVehicleHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accountId := params["accountId"]

	vehicle, err := getVehicleFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing body. %s", err)
		return
	}

	err = vehicle.NewVehicle(accountId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating vehicle. %s", err)
		return
	}

	err = vehicle.SaveVehicle()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating vehicle. %s", err)
		return
	}

	writeResponse(w, vehicle)
}

func getVehicleFromRequest(r *http.Request) (*Vehicle, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	vehicle := &Vehicle{}
	err = json.Unmarshal(reqBody, &vehicle)
	if err != nil {
		log.Warnf("Error unmarshalling vehicle from request. %s", err)
		return nil, err
	}
	return vehicle, nil
}

func writeResponse(w http.ResponseWriter, vehicle *Vehicle) {
	jsonData, err := json.Marshal(vehicle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marhsalling vehicle id. %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
