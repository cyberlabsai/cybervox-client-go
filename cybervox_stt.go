package cybervox

import (
	"time"

	"github.com/gorilla/websocket"
)

type (
	// STTRequest contains the parameters to send a `stt` message to the platform.
	STTRequest struct {
		Emit    string            `json:"emit"`
		Payload STTRequestPayload `json:"payload"`
	}
	// STTRequestPayload most important field is the `UploadID` to be transcribed into a text block.
	STTRequestPayload struct {
		Timestamp int64  `json:"timestamp"`
		UploadID  string `json:"upload_id"`
	}
	// STTResponse contains the parameters sent in response to a `stt` message.
	STTResponse struct {
		Event   string             `json:"event"`
		Payload STTResponsePayload `json:"payload"`
	}
	STTResponsePayload struct {
		Success   bool   `json:"success"`   // successfully transcribed STTRequestPayload.UploadID into a block of text
		Reason    string `json:"reason"`    // if Success is false, the failure reason
		Text      string `json:"text"`      // if Success is true, the block of transcribed text
		Timestamp int64  `json:"timestamp"` // the given STTRequestPayload.Timestamp
	}
)

// STT sends a speech-to-text request on an established websocket connection.
// It receives the websocket.Conn and assumes it's properly connected and returns an STTResponse.
func STT(ws *websocket.Conn, uploadID string) (response STTResponse) {
	request := STTRequest{
		Emit: "stt",
		Payload: STTRequestPayload{
			Timestamp: time.Now().UnixNano(),
			UploadID:  uploadID,
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
