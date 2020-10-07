package subscription

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionACKMessage ...
type ConnectionACKMessage struct {
	OperationID string `json:"id,omitempty"`
	Type        string `json:"type"`
	Payload     struct {
		Query string `json:"query"`
	} `json:"payload,omitempty"`
}

// Subscriber ...
type Subscriber struct {
	ID            int
	Conn          *websocket.Conn
	RequestString string
	OperationID   string
}

// Upgrader ...
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"graphql-ws"},
}

// Subscribers ...
var Subscribers sync.Map

// Handler ...
func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to do websocket upgrade: %v", err)
		return
	}
	connectionACK, err := json.Marshal(map[string]string{
		"type": "connection_ack",
	})
	if err != nil {
		log.Printf("failed to marshal ws connection ack: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, connectionACK); err != nil {
		log.Printf("failed to write to ws connection: %v", err)
		return
	}
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				return
			}
			if err != nil {
				log.Printf("failed to read websocket message: %s\n", err)
				return
			}
			var msg ConnectionACKMessage
			if err := json.Unmarshal(p, &msg); err != nil {
				log.Printf("failed to unmarshal: %s", err)
				return
			}
			if msg.Type == "start" {
				length := 0
				Subscribers.Range(func(key, value interface{}) bool {
					length++
					return true
				})
				var subscriber = Subscriber{
					ID:            length + 1,
					Conn:          conn,
					RequestString: msg.Payload.Query,
					OperationID:   msg.OperationID,
				}
				Subscribers.Store(subscriber.ID, &subscriber)
			}
		}
	}()
}
