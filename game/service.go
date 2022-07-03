package game

import (
	"context"
	"log"
	"matchmaker/game/msg"
	"matchmaker/game/room"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// GameService is the game server implementation that handles
// game rooms and game logic.
type GameService struct {
	Rooms room.RoomManager
}

// NewGameService creates and initializes a new GameServer
// instance.
func NewGameService() GameService {
	return GameService{
		Rooms: room.NewRoomManager(),
	}
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

	// expects the client to indicate the game room ID
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	login := msg.Login{}
	if err := wsjson.Read(ctx, conn, &login); err != nil {
		conn.Close(websocket.StatusPolicyViolation, "Expected 'Login' message")
		cancel()
		return
	}
	cancel()

	room := s.Rooms.GetOrCreate(login.GameRoomId)
	player, err := room.Join(login.ClientId)
	if err != nil {
		wsjson.Write(r.Context(), conn, msg.GameState{Error: err.Error()})
		return
	}
	// send the assigned play
	if wsjson.Write(r.Context(), conn, msg.AssignPlay{Play: player.Play}) != nil {
		log.Printf("Failed sending play assignation to player %v", login.ClientId)
		return
	}

	// game updates sender
	go func(ctx context.Context, conn *websocket.Conn) {
		for message := range player.GameEvents {
			wsjson.Write(ctx, conn, &message)
		}
	}(r.Context(), conn)

	// player events receiver
	for {
		playerMsg := msg.SelectPosition{}

		if err := wsjson.Read(r.Context(), conn, &playerMsg); err != nil {
			log.Printf("Failed reading player message: %v", err)
			close(player.PlayerEvents)
			return
		}
		player.PlayerEvents <- playerMsg
	}
}
