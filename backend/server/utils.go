package server

import (
	"errors"
	"io"
	"log"
	"os/exec"

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

func closeAndExit(client Client, ffmpegIn io.WriteCloser, ffmpegCmd *exec.Cmd) {
	if client.Peer != nil {
		if err := client.Peer.Close(); err != nil {
			log.Println("Unable to close the connection")
		}
	}

	if ffmpegIn != nil {
		if err := ffmpegIn.Close(); err != nil {
			log.Println("Unable to close stdin pipe to ffmpegIn")
		} else {
			log.Println("ffmpeg stdin closed")
		}
	}

	if ffmpegCmd != nil {
		if err := ffmpegCmd.Process.Kill(); err != nil {
			log.Println("Unable to kill the ffmpeg cmd")
		} else {
			log.Println("ffmpeg cmd killed")
		}
	}
}

// Dummy error handing.
func responseErrorHandling(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
