package socket

import (
	"encoding/json"
	"fmt"
)

const (
	EventSendPrivateMessage = "send_private_message"

	EventNewPrivateMessage = "new_private_message"

	EventSendGroupMessage = "send_group_message"

	EventNewGroupMessage = "new_group_message"

	EventChangeRoom = "change_room"
)

type Event struct {
	Type string `json:"type"`

	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

// NewMessageEvent is returned when responding to send_message

type ChangeRoomEvent struct {
	Name string `json:"name"`
}

// ChatRoomHandler will handle switching of chatrooms between clients
func ChatRoomHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var changeRoomEvent ChangeRoomEvent
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// Add Client to chat room
	c.room = changeRoomEvent.Name

	return nil
}
