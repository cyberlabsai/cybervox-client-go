package cybervox

import (
	"time"

	"github.com/gorilla/websocket"
)

type (
	// TTSRequest contains the parameters to send a `tts` message to the platform.
	TTSRequest struct {
		Emit    string            `json:"emit"`
		Payload TTSRequestPayload `json:"payload"`
	}
	// TTSRequestPayload most important field is the `Text` to be converted into a WAVE file.
	TTSRequestPayload struct {
		Timestamp int64  `json:"timestamp"`
		Text      string `json:"text"`
	}
	// TTSResponse contains the parameters sent in response to a `tts` message.
	TTSResponse struct {
		Event   string             `json:"event"`
		Payload TTSResponsePayload `json:"payload"`
	}
	TTSResponsePayload struct {
		Success   bool   `json:"success"`             // successfully converted TTSRequestPayload.Text into a WAVE file
		Reason    string `json:"reason,omitempty"`    // if Success is false, the failure reason
		AudioURL  string `json:"audio_url,omitempty"` // if Success is true, the WAVE file url to be downloaded
		Timestamp int64  `json:"timestamp,omitempty"` // the given TTSRequestPayload.Timestamp
	}
)

// TTS sends a text-to-speech request on an established websocket connection.
// It receives the websocket.Conn and assumes it's properly connected and returns a TTSResponse.
func TTS(ws *websocket.Conn, text string) (response TTSResponse) {
	request := TTSRequest{
		Emit: "tts",
		Payload: TTSRequestPayload{
			Text:      text,
			Timestamp: time.Now().UnixNano(),
		},
	}
	if err := ws.WriteJSON(request); err != nil {
		log.Error(err)
		return
	}
	if err := ws.ReadJSON(&response); err != nil {
		log.Error(err)
		return
	}
	return
}
