package lobby

import (
	"context"
	"matchmaker/lobby/msg"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// MatchServer is the matchmaker server implementation that accepts
// client connections (websocket) and puts them in queue to make them
// eligible on the matchmaking algorithm to start a new game.
type MatchServer struct {
	// logf controls where logs are sent.
	logf func(f string, v ...interface{})
}

// AcceptClient is the matchmaking server request handler for client
// connections. It will upgrade HTTP connections into a websocket and
// send/receive events.
func (s *MatchServer) AcceptClient(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		Subprotocols: []string{"matchmaker"},
	})
	if err != nil {
		s.logf("%v", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "Server shutdown")

	// expects the client Login message or else closes the connection
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	login := msg.Login{}
	if err := wsjson.Read(ctx, conn, &login); err != nil {
		conn.Close(websocket.StatusPolicyViolation, "Expected 'Login' message")
		cancel()
		return
	}
	cancel()

	wsjson.Write(r.Context(), conn, "Hi client "+login.ClientId)
}
