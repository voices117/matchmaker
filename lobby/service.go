package lobby

import (
	"context"
	"fmt"
	"matchmaker/lobby/msg"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"math/rand"
)

var (
	test_user_id int = 0
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

	
	// id := PlayerId(login.ClientId)
	id := PlayerId(fmt.Sprintf("Test User %", test_user_id))
	test_user_id = test_user_id + 1
	
	// TODO: le should come from login
	// loe := PlayerId(login.ClientLoe)
	loe := rand.Intn(100)

	// create the channel from where the resulting created match
	// will be returned
	response := make(chan Game)

	wsjson.Write(r.Context(), conn, "Hi client "+id)
	s.Service.Add(r.Context(), id, loe, response)

	// await the match maker response when it's done
	select {
	case gameId := <-response:
		err = wsjson.Write(r.Context(), conn, gameId)
		if err != nil {
			fmt.Printf("Failed sending match creation message to %v", id)
		}
	case <-r.Context().Done():
	}
}
