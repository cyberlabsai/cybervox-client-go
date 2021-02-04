package cybervox

import (
	"io/ioutil"
	"time"

	"github.com/gorilla/websocket"
)

type (
	// UploadRequest contains the parameters to send a `upload` message to the platform.
	UploadRequest struct {
		Emit    string               `json:"emit"`
		Payload UploadRequestPayload `json:"payload"`
	}
	UploadRequestPayload struct {
		Timestamp  int64 `json:"timestamp"`
		MaxUploads int   `json:"max_uploads"`
	}
	// UploadResponse contains the parameters sent in response to an `upload` message.
	UploadResponse struct {
		Event   string                `json:"event"`
		Payload UploadResponsePayload `json:"payload"`
	}
	UploadResponsePayload struct {
		UploadID  string `json:"upload_id"` // the generated uuid valid for some time
		Timestamp int64  `json:"timestamp"` // the given UploadRequestPayload.Timestamp
	}
)

// Upload sends an upload request followed by a bytes stream on an established websocket connection.
// It receives the websocket.Conn and assumes it's properly connected and returns an UploadResponse.
func Upload(ws *websocket.Conn, filename string) (response UploadResponse) {
	request := UploadRequest{
		Emit: "upload",
		Payload: UploadRequestPayload{
			Timestamp:  time.Now().UnixNano(),
			MaxUploads: 1,
		},
	}
	if err := ws.WriteJSON(request); err != nil {
		log.Error(err)
		return
	}

	var err error
	var contents []byte
	if contents, err = ioutil.ReadFile(filename); err != nil {
		log.Error(err)
		return
	}
	if err = ws.WriteMessage(websocket.BinaryMessage, contents); err != nil {
		log.Error(err)
		return
	}
	if err = ws.ReadJSON(&response); err != nil {
		log.Error(err)
		return
	}
	return
}
