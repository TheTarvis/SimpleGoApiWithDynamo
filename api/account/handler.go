package account

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

func GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error invalid account id. %s", err)
		return
	}

	account := Account{}
	err = account.GetAccountById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching account id. %s", err)
		return
	}

	writeResponse(w, &account)
}

func UpdateAccountHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := uuid.Parse(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error invalid account id. %s", err)
		return
	}

	account, err := getAccountFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing body. %s", err)
		return
	}
	account.ID = id
	err = account.UpdateAccount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating account id. %s", err)
		return
	}

	writeResponse(w, account)
}

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	account, err := getAccountFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error parsing body. %s", err)
		return
	}

	err = account.NewAccount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating account. %s", err)
		return
	}

	err = account.SaveAccount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating account. %s", err)
		return
	}

	writeResponse(w, account)
}

func getAccountFromRequest(r *http.Request) (*Account, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	err = json.Unmarshal(reqBody, &account)
	if err != nil {
		log.Warnf("Error unmarshalling account from request. %s", err)
		return nil, err
	}
	return account, nil
}


func writeResponse(w http.ResponseWriter, account *Account) {
	jsonData, err := json.Marshal(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marhsalling account id. %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
