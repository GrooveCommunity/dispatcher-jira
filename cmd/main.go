package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GrooveCommunity/dispatcher-jira/entity"
	"google.golang.org/api/pubsub/v1"

	"github.com/gorilla/mux"
)

type pushRequest struct {
	Message      pubsub.PubsubMessage
	Subscription string
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/queue-dispatcher-jira", handleQueueDispatcher).Methods("POST")

	fmt.Println("Port: ", os.Getenv("APP_PORT"))

	panic(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entity.Healthy{Status: "Success!"})
}

func handleQueueDispatcher(w http.ResponseWriter, r *http.Request) {

	var pr pushRequest

	if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
		http.Error(w, fmt.Sprintf("Não foi possível decodificar o body"), http.StatusBadRequest)
		return
	}

	requestJson, err := base64.StdEncoding.DecodeString(pr.Message.Data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Não foi possível decodificar o json da fila"), http.StatusBadRequest)
		return
	}

	log.Println(string(requestJson))

	/*var dispatcherRequest interface{}

	if err := json.NewDecoder(r.Body).Decode(&dispatcherRequest); err != nil {
		http.Error(w, fmt.Sprintf("Não foi possível decodificar o body: %v", err), http.StatusBadRequest)
		return
	}

	fmt.Println(dispatcherRequest)*/
}
