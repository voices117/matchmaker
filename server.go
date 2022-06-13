package main

import (
	"matchmaker/lobby"
	"net/http"
)

// Logf is the logging function type signature for the server.
type Logf func(f string, v ...interface{})

// Server is the type that implements the main server that
// handles all client HTTP requests.
type Server struct {
	// logf is the function pointer for the server to use as logger.
	logf Logf

	// mux routes the various endpoints to the appropriate handler.
	mux http.ServeMux

	// matchmaker is the match making service.
	matchmaker lobby.MatchServer
}

// NewServer creates a new server instance with the matchmaking,
// game and static file serving services.
func NewServer(logf Logf) *Server {
	server := Server{logf: logf}

	// adds a handler for static content served directly from the 'static' dir
	server.mux.Handle("/", http.FileServer(http.Dir("./static")))

	// adds a handler for matchmaking
	server.mux.HandleFunc("/matchmaker", server.matchmaker.AcceptClient)

	return &server
}

// ServeHTTP dispatches the request to the handlers defined by the server.
func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.mux.ServeHTTP(w, r)
}
