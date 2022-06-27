package room

import "matchmaker/game/msg"

// Player represents a player in a particular game Room.
type Player struct {
	// Play is the type of player (X or O)
	Play rune

	// PlayerEvents is a channel where the player sends its events
	// for the game to be consumed.
	PlayerEvents chan msg.SelectPosition

	// GameEvents is the channel where the game sends events with
	// updates for the player.
	GameEvents chan msg.GameState
}

// NewPlayer creates and initializes a new Player instance.
func NewPlayer(play rune) Player {
	return Player{
		Play:         play,
		PlayerEvents: make(chan msg.SelectPosition),
		GameEvents:   make(chan msg.GameState),
	}
}
