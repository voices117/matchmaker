package game

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// GameService is the game server implementation that handles
// game rooms and game logic.
type GameService struct {
}

// NewGameService creates and initializes a new GameServer
// instance.
func NewGameService() GameService {
	return GameService{}
}

// JoinGame is the game server request handler for client
// connections. It will upgrade HTTP connections into a websocket and
// send/receive events.
func (s *GameService) JoinGame(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"game"},
	})
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "Server shutdown")

	wsjson.Write(r.Context(), conn, "Joined game!")
}
