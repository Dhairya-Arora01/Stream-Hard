package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/pion/webrtc/v3"
)

// initialSetup() returns a pointer to an instance of webrtc's API.
func initialSetup() *webrtc.API {

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
	return api

}

func ffmpegSetup() io.WriteCloser {
	ffmpegCmd := exec.Command("ffmpeg", "-re", "-i", "pipe:0", "-c:v", "libx264", "-preset", "veryfast", "-maxrate", "3000k", "-bufsize", "6000k", "-pix_fmt", "yuv420p", "-f", "flv", "rtmp://a.rtmp.youtube.com/live2/um28-897g-yfrq-zjr7-d559")

	ffmpegIn, err := ffmpegCmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	ffmpegOut, err := ffmpegCmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		scanner := bufio.NewScanner(ffmpegOut)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err := ffmpegCmd.Start(); err != nil {
		panic(err)
	}
	return ffmpegIn
}
