package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pion/webrtc/v3"
)

// Registers route and starts the server at specified port.
func CreateRoutes(port string) {

	router := mux.NewRouter()
	router.HandleFunc("/hello", sayHello).Methods("GET")
	router.HandleFunc("/feed", getFeeder).Methods("POST")

	http.Handle("/", router)
	log.Printf("Starting Server at localhost:8000")
	http.ListenAndServe(":8000", handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(router))

}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := []byte(`{"message": "hello"}`)
	_, err := w.Write(response)

	responseErrorHandling(err, "Some Exception has occured")

}

// SessionDescription as received from the frontend
type SessionInfo struct {
	SessionInformation webrtc.SessionDescription `json:"sdp"`
}

func getFeeder(w http.ResponseWriter, r *http.Request) {

	log.Printf("%s-%s", r.Method, r.RequestURI)

	w.Header().Set("Content-Type", "application/json")

	var sdp SessionInfo

	err := json.NewDecoder(r.Body).Decode(&sdp)

	if err != nil {
		log.Fatalf("Unable to decode the request body: %s\n", err)
	}

	log.Printf("SPD received")

	response := []byte(`{"message": "recieved sdp"}`)

	_, err = w.Write(response)
	responseErrorHandling(err, "Some Exception occured")
}

func responseErrorHandling(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
