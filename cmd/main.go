package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/GrooveCommunity/dispatcher-jira/internal"
	"github.com/GrooveCommunity/glib-noc-event-structs/entity"
	"github.com/gorilla/mux"
	"google.golang.org/api/pubsub/v1"
)

type pushRequest struct {
	Message      pubsub.PubsubMessage
	Subscription string
}

var (
	username, token, endpoint, appPort string
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/queue-dispatcher-jira", handleQueueDispatcher).Methods("POST")
	router.HandleFunc("/put-rule", handlePutRule).Methods("POST")

	username = os.Getenv("JIRA_USERNAME")
	token = os.Getenv("JIRA_TOKENAPI")
	endpoint = os.Getenv("JIRA_ENDPOINT")
	appPort = os.Getenv("APP_PORT")

	if username == "" || token == "" || endpoint == "" || appPort == "" {
		log.Fatal("Nem todas as variáveis de ambiente requeridas foram fornecidas. ")
	}

	panic(http.ListenAndServe(":"+appPort, router))
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

	var jiraIssue entity.JiraIssue

	if err := json.Unmarshal(requestJson, &jiraIssue); err != nil {
		panic(err)
	}

	go internal.ForwardIssue(jiraIssue, username, token, endpoint)
}

func handlePutRule(w http.ResponseWriter, r *http.Request) {
	var rule entity.Rule

	err := json.NewDecoder(r.Body).Decode(&rule)
	if err != nil {
		panic(err)
	}

	internal.WriteRule(rule)
}
