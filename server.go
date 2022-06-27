package main

import (
	"context"
	"matchmaker/game"
	"matchmaker/lobby"
	"net"
	"net/http"
)

// Server is the type that implements the main server that
// handles all client HTTP requests.
type Server struct {
	// logf is the function pointer for the server to use as logger.
	logf lobby.Logf

	// mux routes the various endpoints to the appropriate handler.
	mux http.ServeMux

	// matchmaker is the match making service.
	matchmaker lobby.MatchService

	gameServer game.GameService

	// actual HTTP server.
	httpServer *http.Server
}

// NewServer creates a new server instance with the matchmaking,
// game and static file serving services.
func NewServer(logf lobby.Logf) *Server {
	server := Server{
		logf:       logf,
		matchmaker: lobby.NewMatchService(logf),
		gameServer: game.NewGameService(),
	}
	server.httpServer = &http.Server{Handler: &server}

	// adds a handler for static content served directly from the 'static' dir
	server.mux.Handle("/", http.FileServer(http.Dir("./static")))

	// adds handlers for matchmaking and game services
	server.mux.HandleFunc("/matchmaker", server.matchmaker.AcceptClient)
	server.mux.HandleFunc("/game", server.gameServer.JoinGame)

	return &server
}

// ServeHTTP dispatches the request to the handlers defined by the server.
func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.mux.ServeHTTP(w, r)
}

// Serve accepts incoming connections on the Listener l, creating a new service
// goroutine for each. The service goroutines read requests and then call
// srv.Handler to reply to them.
//
// HTTP/2 support is only enabled if the Listener returns *tls.Conn connections
// and they were configured with "h2" in the TLS Config.NextProtos.
//
// Serve always returns a non-nil error and closes l. After Shutdown or Close,
// the returned error is ErrServerClosed.
func (server *Server) Serve(address string, ctx context.Context) error {
	// start the matchmaking service
	server.logf("Starting matchmaking service")
	go server.matchmaker.Service.Start(ctx)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server.logf("listening on http://%v", listener.Addr())
	return server.httpServer.Serve(listener)
}

// Gracefully shutdown the server.
func (server *Server) Shutdown(ctx context.Context) error {
	return server.httpServer.Shutdown(ctx)
}
