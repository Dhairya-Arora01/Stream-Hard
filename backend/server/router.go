package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media/ivfwriter"
)

// Registers route and starts the server at specified port.
func CreateRoutes(port string) {

	router := mux.NewRouter()
	router.HandleFunc("/hello", sayHello).Methods("GET")
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

		ffmpegIn := ffmpegSetup()
		ivfwriter, err := ivfwriter.NewWith(ffmpegIn)
		if err != nil {
			panic(err)
		}
		// oggwriter, err := oggwriter.NewWith(ffmpegIn, 44100, 2)
		// if err != nil {
		// 	panic(err)
		// }

		go func() {
			ticker := time.NewTicker(time.Second * 1)
			for range ticker.C {
				rtcpSendErr := Peer.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(tr.SSRC())}})
				if rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		log.Printf("Track has started stream %s, id: %s, rid: %s, kind: %s", tr.StreamID(), tr.ID(), tr.RID(), tr.Kind())
		if tr.Kind().String() == "video" {
			for {
				rtpPacket, _, err := tr.ReadRTP()
				if err != nil {
					panic(err)
				}
				if rtpPacket.PayloadType == 96 {
					if err := ivfwriter.WriteRTP(rtpPacket); err != nil {
						panic(err)
					}
				}
				// if rtpPacket.PayloadType == 111 {
				// 	if err := oggwriter.WriteRTP(rtpPacket); err != nil {
				// 		panic(err)
				// 	}
				// }
			}
		}
		// } else {
		// 	for {
		// 		rtpPacket, _, err := tr.ReadRTP()
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		if err := oggwriter.WriteRTP(rtpPacket); err != nil {
		// 			panic(err)
		// 		}
		// 	}
		// }
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
		// if the message is an offer set it as remote description and send the answer.
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

			if err = conn.WriteJSON(answer); err != nil {
				panic(err)
			}
		// if the message is an ICECandidate, add it.
		case json.Unmarshal(p, &candidate) == nil && candidate.ToJSON().Candidate != "":
			if err = Peer.AddICECandidate(candidate.ToJSON()); err != nil {
				panic(err)
			}
		}

	}

}

func responseErrorHandling(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
