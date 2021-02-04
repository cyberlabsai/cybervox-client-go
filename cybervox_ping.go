package cybervox

import (
	"time"

	"github.com/gorilla/websocket"
)

type (
	// PingRequest contains the parameters to send a `ping` message to the platform.
	PingRequest struct {
		Emit    string             `json:"emit"`
		Payload PingRequestPayload `json:"payload"`
	}
	PingRequestPayload struct {
		Timestamp int64 `json:"timestamp"`
	}

	// PingResponse contains the parameters sent in response to a `ping` message.
	PingResponse struct {
		Event   string              `json:"event"`
		Payload PingResponsePayload `json:"payload"`
	}
	PingResponsePayload struct {
		Success   bool  `json:"success"`
		Timestamp int64 `json:"timestamp,omitempty"` // the given PingRequestPayload.Timestamp
	}
)

// Ping sends a ping request on an established websocket connection.
// It receives the websocket.Conn and assumes it's properly connected and returns a PingResponse with the given
// Timestamp.
func Ping(ws *websocket.Conn) (response PingResponse) {
	request := PingRequest{
		Emit: "ping",
		Payload: PingRequestPayload{
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
