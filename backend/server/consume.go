package server

import (
	"log"

	"github.com/pion/webrtc/v3"
)

var peer *webrtc.PeerConnection

func initialSetup(sdp webrtc.SessionDescription) *webrtc.SessionDescription {

	mediaEngine := webrtc.MediaEngine{}

	err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeVP8,
			ClockRate:    90000,
			Channels:     0,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 96,
	},
		webrtc.RTPCodecTypeVideo)

	if err != nil {
		panic(err)
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peer, err = api.NewPeerConnection(config)

	if err != nil {
		panic(err)
	}

	err = peer.SetRemoteDescription(sdp)
	if err != nil {
		log.Fatal("this is it")
		panic(err)
	}

	localsdp, err := peer.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	err = peer.SetLocalDescription(localsdp)
	if err != nil {
		panic(err)
	}
	consumeMedia()
	return peer.LocalDescription()
}

func consumeMedia() {

	peer.OnTrack(func(track *webrtc.TrackRemote, reciever *webrtc.RTPReceiver) {
		log.Println("SSRC: ", track.SSRC())
	})

}
