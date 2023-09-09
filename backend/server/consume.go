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

	err = mediaEngine.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     webrtc.MimeTypeOpus,
			ClockRate:    48000,
			Channels:     2,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 111,
	},
		webrtc.RTPCodecTypeAudio)

	if err != nil {
		panic(err)
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))
	return api

}

// Setup and execution for the ffmpeg command.
func ffmpegSetup(rtmpLink string) (io.WriteCloser, *exec.Cmd, error) {
	ffmpegCmd := exec.Command("ffmpeg", "-re", "-i", "pipe:0", "-c:v", "libx264", "-preset", "veryfast", "-maxrate", "3000k", "-bufsize", "6000k", "-pix_fmt", "yuv420p", "-c:a", "aac", "-ar", "44100", "-f", "flv", rtmpLink)

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
		return nil, ffmpegCmd, err
	}
	return ffmpegIn, ffmpegCmd, nil
}
