package msg

// Client messages

// Login message expected from the client when the connection
// to the game is established.
type Login struct {
	// ClientId is the client's ID.
	ClientId string `json:"client_id"`
	// GameRoomId is the ID of the game the player will join.
	GameRoomId string `json:"game_room_id"`
}

type SelectPosition struct {
	// Position in the board where the player wants to place
	// the play. Must be an integer between 0 and 8.
	Position int `json:"position"`
}

// Server messages

// AssignPlay indicates the player which play (PLAY_X or PLAY_O)
// corresponds the player receiving this message.
type AssignPlay struct {
	Play rune `json:"play"`
}

// GameState represents the current game state.
type GameState struct {
	// Error indicates if the last play by the player was rejected
	// because of a problem (indicated in the string). If there was
	// no error, then this field is empty.
	Error string `json:"error"`

	// Board is the board state indicating which positions are
	// occupied (either an 'X' or a 'O') or empty (a ' ').
	Board string `json:"board"`

	// State if the game state indicating if it finished with a
	// winner, a tie or still ongoing.
	State string `json:"state"`

	// CurrentTurn is the player that should place the next move.
	CurrentTurn string `json:"current_turn"`
}
