package main

import (
	"os"

	"github.com/cyberlabsai/cybervox-client-go"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("package", "main")

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	var (
		ws  *websocket.Conn
		err error
	)

	if ws, _, err = cybervox.Dial(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		ws.Close()
	}()

	audioFile := "ola-mundo.wav"
	if len(os.Args) < 2 {
		log.Infof("no audio file was given, using default %q", audioFile)
	} else {
		audioFile = os.Args[1]
	}

	// --- Upload ---
	uploadResponse := cybervox.Upload(ws, audioFile)
	log.Printf("Upload ID: %q\n", uploadResponse.Payload.UploadID)

	// --- STT ---
	sttResponse := cybervox.STT(ws, uploadResponse.Payload.UploadID)
	if !sttResponse.Payload.Success {
		log.Fatalf("abort: reason: %q", sttResponse.Payload.Reason)
	}
	log.Printf("%q\n", sttResponse.Payload.Text)
}
