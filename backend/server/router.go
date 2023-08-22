package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

// Registers route and starts the server at specified port.
func CreateRoutes(port string) {

	router := mux.NewRouter()
	router.HandleFunc("/hello", sayHello).Methods("GET")
	// router.HandleFunc("/broadcast", broadcastHandler).Methods("POST")
	router.HandleFunc("/ws", ws)

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

type SendSdp struct {
	Sdp webrtc.SessionDescription `json:"sdp"`
}

// func broadcastHandler(w http.ResponseWriter, r *http.Request) {
// 	api := initialSetup()
// 	Peer, err := api.NewPeerConnection(webrtc.Configuration{
// 		ICEServers: []webrtc.ICEServer{
// 			{
// 				URLs: []string{
// 					"stun:stun1.l.google.com:19302",
// 				},
// 			},
// 		},
// 		ICECandidatePoolSize: 10,
// 	})
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	Peer.OnTrack(func(tr *webrtc.TrackRemote, rtpr *webrtc.RTPReceiver) {
// 		log.Println(tr.SSRC())
// 	})

// 	var sessionDescription SessionInfo
// 	json.NewDecoder(r.Body).Decode(&sessionDescription)
// 	err = Peer.SetRemoteDescription(sessionDescription.SessionInformation)
// 	if err != nil {
// 		log.Println("Error while setting remote description", err)
// 	}
// 	log.Println("Remote Description set sucessfully")

// 	answer, err := Peer.CreateAnswer(nil)
// 	if err != nil {
// 		log.Println("Error while creating answer", err)
// 	}
// 	Peer.SetLocalDescription(answer)
// 	log.Println("Local Description set sucessfully")

// 	sdpToSend, err := json.Marshal(SendSdp{
// 		Sdp: answer,
// 	})

// 	if err != nil {
// 		log.Println("Error while Marshalling", err)
// 	}
// 	_, err = w.Write(sdpToSend)
// 	if err != nil {
// 		log.Println("Some err occured", err)
// 	}

// 	go func() {

// 		for {
// 			log.Println(Peer.ConnectionState(), Peer.ICEConnectionState(), Peer.ICEGatheringState())
// 			time.Sleep(5 * time.Second)
// 		}
// 	}()
// }

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ws(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()

	api := initialSetup()
	Peer, err := api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{
					"stun:stun1.l.google.com:19302",
				},
			},
		},
		ICECandidatePoolSize: 10,
	})
	if err != nil {
		fmt.Println(err)
	}
	Peer.OnTrack(func(tr *webrtc.TrackRemote, rtpr *webrtc.RTPReceiver) {
		log.Println("SSRC", tr.SSRC())
		log.Println(tr.Codec().ClockRate)
	})

	Peer.OnICECandidate(func(c *webrtc.ICECandidate) {
		fmt.Println(c)
		if c == nil {
			return
		}
		outbound, marshalErr := json.Marshal(c.ToJSON())
		if marshalErr != nil {
			panic(marshalErr)
		}
		if err := conn.WriteMessage(websocket.TextMessage, outbound); err != nil {
			panic(err)
		}
	})

	Peer.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	var (
		candidate webrtc.ICECandidate
		offer     webrtc.SessionDescription
	)
	for {

		_, p, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		switch {
		case json.Unmarshal(p, &offer) == nil && offer.SDP != "":
			if err := Peer.SetRemoteDescription(offer); err != nil {
				panic(err)
			}

			answer, answererr := Peer.CreateAnswer(nil)
			if answererr != nil {
				panic(answererr)
			}

			if err = Peer.SetLocalDescription(answer); err != nil {
				panic(err)
			}

			// outbound, err := json.Marshal(answer)
			// if err != nil {
			// 	panic(err)
			// }

			if err = conn.WriteJSON(answer); err != nil {
				panic(err)
			}
		case json.Unmarshal(p, &candidate) == nil && candidate.ToJSON().Candidate != "":
			if err = Peer.AddICECandidate(candidate.ToJSON()); err != nil {
				panic(err)
			}
		}

	}

}

// func ws(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Print("upgrade failed: ", err)
// 		return
// 	}
// 	defer conn.Close()

// 	var IceCandidateRecieved webrtc.ICECandidateInit

// 	for {
// 		_, b, err := conn.ReadMessage()

// 		if err != nil {
// 			log.Println(err)
// 		}

// 		err = json.Unmarshal(b, &IceCandidateRecieved)
// 		if err != nil {
// 			log.Println("JSON unmarshal err", err)
// 		}

// 		err = Peer.AddICECandidate(IceCandidateRecieved)
// 		if err != nil {
// 			log.Println("Error adding Ice candidate", err)
// 		}

// 		Peer.OnICECandidate(func(i *webrtc.ICECandidate) {
// 			if i != nil {
// 				err = conn.WriteJSON(i.ToJSON())
// 				if err != nil {
// 					log.Println("Error while sending candidate", err)
// 				}
// 			}
// 		})
// 	}
// }

func responseErrorHandling(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
