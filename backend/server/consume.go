package server

import (
	"github.com/pion/webrtc/v3"
)

// initialSetup() returns a pointer to an instance of webrtc's API.
func initialSetup() *webrtc.API {

	mediaEngine := webrtc.MediaEngine{}

	err := mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeVP9,
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
	return api

}
