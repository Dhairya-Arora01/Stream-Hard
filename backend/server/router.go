package server

import (
	"encoding/json"
	"fmt"
	"io"
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

	var client Client
	api := initialSetup()
	client.Peer, err = api.NewPeerConnection(webrtc.Configuration{
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

	var (
		ffmpegIn  io.WriteCloser
		ivfWriter *ivfwriter.IVFWriter
		candidate webrtc.ICECandidate
		offer     webrtc.SessionDescription
	)

	rtmpChan := make(chan string)
	// var oggWriter *oggwriter.OggWriter

	// Handling the incoming tracks
	client.Peer.OnTrack(func(tr *webrtc.TrackRemote, rtpr *webrtc.RTPReceiver) {

		if ffmpegIn == nil {
			client.RtmpLink.URL = <-rtmpChan
			ffmpegIn = ffmpegSetup(client.RtmpLink.URL)
		}
		if ivfWriter == nil {
			ivfWriter, err = ivfwriter.NewWith(ffmpegIn)
			if err != nil {
				panic(err)
			}
		}
		// if oggWriter == nil {
		// 	oggWriter, err = oggwriter.NewWith(ffmpegIn, 44100, 2)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for range ticker.C {
				rtcpSendErr := client.Peer.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(tr.SSRC())}})
				if rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		log.Printf("Track has started stream %s, id: %s, rid: %s, kind: %s", tr.StreamID(), tr.ID(), tr.RID(), tr.Kind())
		for {
			rtpPacket, _, err := tr.ReadRTP()
			if err != nil {
				panic(err)
			}
			// writing the video packet using ivfwriter.
			if rtpPacket.PayloadType == 96 {
				if err := ivfWriter.WriteRTP(rtpPacket); err != nil {
					panic(err)
				}
			}
			// if rtpPacket.PayloadType == 111 {
			// 	if err := oggWriter.WriteRTP(rtpPacket); err != nil {
			// 		panic(err)
			// 	}
			// }
		}
	})

	// Exchanging the ICECandidates
	client.Peer.OnICECandidate(func(c *webrtc.ICECandidate) {
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

	client.Peer.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	var rtmpLink RTMPLink

	for {

		_, p, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		switch {
		// if the message is an offer set it as remote description and send the answer.
		// Also set the answer as local description
		case json.Unmarshal(p, &offer) == nil && offer.SDP != "":
			if err := client.Peer.SetRemoteDescription(offer); err != nil {
				panic(err)
			}

			answer, answererr := client.Peer.CreateAnswer(nil)
			if answererr != nil {
				panic(answererr)
			}

			if err = client.Peer.SetLocalDescription(answer); err != nil {
				panic(err)
			}

			if err = conn.WriteJSON(answer); err != nil {
				panic(err)
			}
		// if the message is an RTMPLink.
		case json.Unmarshal(p, &rtmpLink) == nil && rtmpLink.URL != "":
			if err = rtmpLink.isValid(); err != nil {
				panic(err)
			}
			go func() {
				rtmpChan <- rtmpLink.URL
				close(rtmpChan)
			}()

		// if the message is an ICECandidate, add it.
		case json.Unmarshal(p, &candidate) == nil && candidate.ToJSON().Candidate != "":
			if err = client.Peer.AddICECandidate(candidate.ToJSON()); err != nil {
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
