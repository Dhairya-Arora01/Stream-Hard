package server

import (
	"errors"

	"github.com/pion/webrtc/v3"
)

// Definition of the Client connecting to the server.
type Client struct {
	Peer     *webrtc.PeerConnection
	RtmpLink RTMPLink
}

type RTMPLink struct {
	URL string `json:"rtmp"`
}

func (rl *RTMPLink) isValid() error {
	pref := rl.URL[:7]
	if pref != "rtmp://" {
		return errors.New("invalid RTMP link")
	}
	return nil
}
