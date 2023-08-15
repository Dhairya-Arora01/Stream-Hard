package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Registers route and starts the server at specified port
func CreateRoutes(port string) {

	router := mux.NewRouter()
	router.HandleFunc("/hello", sayHello).Methods("GET")
	router.HandleFunc("/feed", getFeeder).Methods("POST")

	http.Handle("/", router)
	log.Printf("Starting Server at localhost:8000")
	http.ListenAndServe(":8000", nil)

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := []byte(`{"message": "hello"}`)
	_, err := w.Write(response)

	responseErrorHandling(err, "Some Exception has occured")

}

type Sdp struct {
	Name string `json:"name"`
	Sdp  string `json:"sdp"`
}

func getFeeder(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var sessionDescription Sdp
	err := json.NewDecoder(r.Body).Decode(&sessionDescription)

	if err != nil {
		log.Fatalf("Unable to decode the request body: %s\n", err)
	}

	log.Printf("Session Description--name: %s, sdp: %s\n", sessionDescription.Name, sessionDescription.Sdp)

	response := []byte(`{"message": "recieved sdp"}`)

	_, err = w.Write(response)
	responseErrorHandling(err, "Some Exception occured")
}

func responseErrorHandling(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
