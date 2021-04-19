package main

import (
	"github.com/cyberlabsai/cybervox-client-go"

	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("package", "main")

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	var (
		ws    *websocket.Conn
		err   error
		delta int64
	)

	if ws, _, err = cybervox.Dial(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		ws.Close()
	}()

	// --- PING ---
	pingResponse := cybervox.Ping(ws)
	delta = time.Now().UnixNano() - pingResponse.Payload.Timestamp
	log.Printf("   PING: Round-trip: %8.2f ms, Success: %v\n", float64(delta)/1e6, pingResponse.Payload.Success)

	// --- TTS ---
	ttsResponse := cybervox.TTS(ws, "Olá Mundo!", "perola")
	delta = time.Now().UnixNano() - ttsResponse.Payload.Timestamp
	log.Printf("    TTS: Round-trip: %8.2f ms, Success: %v, Reason: %q, AudioURL: https://api.cybervox.ai%s\n",
		float64(delta)/1e6,
		ttsResponse.Payload.Success,
		ttsResponse.Payload.Reason,
		ttsResponse.Payload.AudioURL)

	// --- Upload ---
	uploadResponse := cybervox.Upload(ws, "ola-mundo.wav")
	delta = time.Now().UnixNano() - uploadResponse.Payload.Timestamp
	log.Printf(" Upload: Round-trip: %8.2f ms, Upload ID: %q\n", float64(delta)/1e6, uploadResponse.Payload.UploadID)

	// --- STT ---
	sttResponse := cybervox.STT(ws, uploadResponse.Payload.UploadID)
	delta = time.Now().UnixNano() - sttResponse.Payload.Timestamp
	log.Printf("    STT: Round-trip: %8.2f ms, Success: %v, Reason: %q, Text: %q\n",
		float64(delta)/1e6,
		sttResponse.Payload.Success,
		sttResponse.Payload.Reason,
		sttResponse.Payload.Text)
}
