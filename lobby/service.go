package lobby

import (
	"context"
	"matchmaker/lobby/msg"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Logf is the logging function type signature for the server.
type Logf func(f string, v ...interface{})

// MatchService is the matchmaker server implementation that accepts
// client connections (websocket) and puts them in queue to make them
// eligible on the matchmaking algorithm to start a new game.
type MatchService struct {
	// logf controls where logs are sent.
	logf Logf

	Service MatchMaker
}

// NewMatchService creates and initializes a new MatchServer
// instance.
func NewMatchService(logf Logf) MatchService {
	return MatchService{
		logf:    logf,
		Service: NewMatchMaker(),
	}
}

// AcceptClient is the matchmaking server request handler for client
// connections. It will upgrade HTTP connections into a websocket and
// send/receive events.
func (s *MatchService) AcceptClient(w http.ResponseWriter, r *http.Request) {
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

	id := PlayerId(login.ClientId)


	player := NewPlayer(id)

	match := player.StartPlayer(r.Context(), &s.Service)

	wsjson.Write(r.Context(), conn, match)

}
