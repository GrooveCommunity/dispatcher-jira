package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/GrooveCommunity/dispatcher-jira/entity"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/queue-dispatcher-jira", handleQueueDispatcher).Methods("POST")

	fmt.Println("Port: ", os.Getenv("APP_PORT"))

	panic(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
	//log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entity.Healthy{Status: "Success!"})
}

func handleQueueDispatcher(w http.ResponseWriter, r *http.Request) {
	var dispatcherRequest interface{}

	if err := json.NewDecoder(r.Body).Decode(&dispatcherRequest); err != nil {
		http.Error(w, fmt.Sprintf("Não foi possível decodificar o body: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Println(dispatcherRequest)
}
