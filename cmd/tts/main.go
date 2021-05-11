package main

import (
	"os"
	"strings"

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

	text := strings.Join(os.Args[1:], " ")
	if text == "" {
		text = "Olá Mundo!"
		log.Infof("no text was given, using default %q", text)
	}

	// --- TTS ---
	ttsResponse := cybervox.TTS(ws, text, "perola")
	if !ttsResponse.Payload.Success {
		log.Fatalf("abort: reason: %q", ttsResponse.Payload.Reason)
	}
	log.Printf("https://api.cybervox.ai%s\n", ttsResponse.Payload.AudioURL)
}
