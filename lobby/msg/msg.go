package msg

// Client messages

// Login message expected from the client when the connection
// to the lobby is established.
type Login struct {
	// ClientId is the client's ID.
	ClientId string `json:"client_id"`
}

// Server messages

// MatchReady is the message that the server sends to the
// client indicating it was matched and the game is ready.
type MatchReady struct {
	// GameId is the ID of the game created by the match maker.
	// The user must join it by using the game service and indicating
	// this ID.
	GameId string `json:"game_id"`
}
